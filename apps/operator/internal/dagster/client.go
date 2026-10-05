package dagster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RunStatus string

const (
	RunStatusQueued     RunStatus = "QUEUED"
	RunStatusNotStarted RunStatus = "NOT_STARTED"
	RunStatusStarting   RunStatus = "STARTING"
	RunStatusStarted    RunStatus = "STARTED"
	RunStatusSuccess    RunStatus = "SUCCESS"
	RunStatusFailure    RunStatus = "FAILURE"
	RunStatusCanceled   RunStatus = "CANCELED"
	RunStatusCanceling  RunStatus = "CANCELING"
)

type LaunchRequest struct {
	JobName       string
	RunConfigData map[string]any
}

type Client interface {
	LaunchRun(context.Context, LaunchRequest) (string, error)
	GetRunStatus(context.Context, string) (RunStatus, error)
	TerminateRun(context.Context, string) error
}

type HTTPClient struct {
	endpoint           string
	repositoryLocation string
	repositoryName     string
	httpClient         *http.Client
}

func NewHTTPClient(endpoint, repositoryLocation, repositoryName string) *HTTPClient {
	return &HTTPClient{
		endpoint:           endpoint,
		repositoryLocation: repositoryLocation,
		repositoryName:     repositoryName,
		httpClient:         &http.Client{Timeout: 30 * time.Second},
	}
}

type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *HTTPClient) execute(
	ctx context.Context,
	query string,
	variables map[string]any,
	result any,
) error {
	payload, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dagster graphql returned HTTP %d", resp.StatusCode)
	}

	var gql graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gql); err != nil {
		return err
	}
	if len(gql.Errors) > 0 {
		return fmt.Errorf("dagster graphql: %s", gql.Errors[0].Message)
	}

	return json.Unmarshal(gql.Data, result)
}

func (c *HTTPClient) LaunchRun(ctx context.Context, request LaunchRequest) (string, error) {
	const query = `mutation LaunchRunMutation(
	  $repositoryLocationName: String!,
	  $repositoryName: String!,
	  $jobName: String!,
	  $runConfigData: RunConfigData!
	) {
	  launchRun(executionParams: {
	    selector: {
	      repositoryLocationName: $repositoryLocationName,
	      repositoryName: $repositoryName,
	      jobName: $jobName
	    },
	    runConfigData: $runConfigData
	  }) {
	    __typename
	    ... on LaunchRunSuccess { run { runId } }
	    ... on RunConfigValidationInvalid { errors { message reason } }
	    ... on PythonError { message }
	  }
	}`

	variables := map[string]any{
		"repositoryLocationName": c.repositoryLocation,
		"repositoryName":         c.repositoryName,
		"jobName":                request.JobName,
		"runConfigData":          request.RunConfigData,
	}

	var result struct {
		LaunchRun struct {
			Typename string `json:"__typename"`
			Run      *struct {
				RunID string `json:"runId"`
			} `json:"run"`
			Message string `json:"message"`
			Errors  []struct {
				Message string `json:"message"`
			} `json:"errors"`
		} `json:"launchRun"`
	}

	if err := c.execute(ctx, query, variables, &result); err != nil {
		return "", err
	}

	if result.LaunchRun.Typename != "LaunchRunSuccess" || result.LaunchRun.Run == nil {
		if result.LaunchRun.Message != "" {
			return "", fmt.Errorf("dagster launch failed: %s", result.LaunchRun.Message)
		}
		if len(result.LaunchRun.Errors) > 0 {
			return "", fmt.Errorf("dagster launch failed: %s", result.LaunchRun.Errors[0].Message)
		}
		return "", fmt.Errorf("dagster launch failed: %s", result.LaunchRun.Typename)
	}

	return result.LaunchRun.Run.RunID, nil
}

func (c *HTTPClient) GetRunStatus(ctx context.Context, runID string) (RunStatus, error) {
	const query = `query RunQuery($runId: ID!) {
	  runOrError(runId: $runId) {
	    __typename
	    ... on Run { runId status }
	    ... on RunNotFoundError { runId }
	    ... on PythonError { message }
	  }
	}`

	var result struct {
		RunOrError struct {
			Typename string    `json:"__typename"`
			Status   RunStatus `json:"status"`
			Message  string    `json:"message"`
		} `json:"runOrError"`
	}

	if err := c.execute(ctx, query, map[string]any{"runId": runID}, &result); err != nil {
		return "", err
	}
	if result.RunOrError.Typename != "Run" {
		return "", fmt.Errorf(
			"dagster run lookup failed: %s %s",
			result.RunOrError.Typename,
			result.RunOrError.Message,
		)
	}
	return result.RunOrError.Status, nil
}

func (c *HTTPClient) TerminateRun(ctx context.Context, runID string) error {
	const query = `mutation TerminateRun($runId: String!) {
	  terminateRun(runId: $runId) {
	    __typename
	    ... on TerminateRunSuccess { run { runId } }
	    ... on TerminateRunFailure { message }
	    ... on RunNotFoundError { runId }
	    ... on PythonError { message }
	  }
	}`

	var result struct {
		TerminateRun struct {
			Typename string `json:"__typename"`
			Message  string `json:"message"`
		} `json:"terminateRun"`
	}

	if err := c.execute(ctx, query, map[string]any{"runId": runID}, &result); err != nil {
		return err
	}

	switch result.TerminateRun.Typename {
	case "TerminateRunSuccess", "RunNotFoundError":
		return nil
	default:
		return fmt.Errorf(
			"dagster terminate failed: %s %s",
			result.TerminateRun.Typename,
			result.TerminateRun.Message,
		)
	}
}
