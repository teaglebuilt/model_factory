# Model Factory Operator

Kubernetes-native control-plane for the Model Factory.

This operator **does not replace** the existing Python Dagster orchestrator or
`src/model_factory` package.

Architecture:

```text
ModelProject / Experiment / ModelRun / ModelArtifact / ModelRelease
                              |
                              v
                    Go controller-runtime
                              |
                              v
                        Dagster GraphQL
                              |
                              v
                  apps/orchestrator (Python)
                              |
                              v
                    src/model_factory (Python)
                              |
                              v
                       Kubernetes Jobs
```

## API group

```text
ai.teaglebuilt.io/v1alpha1
```

Resources:

- `ModelProject`
- `Experiment`
- `ModelRun`
- `ModelArtifact`
- `ModelRelease`

## Prerequisites

- Go 1.26+
- kubectl
- a Kubernetes cluster
- Dagster GraphQL endpoint reachable from the operator
- Kubebuilder/controller-gen tooling (the Makefile installs local tools)

## Generate CRDs and deepcopy code

```bash
cd operator
make generate
make manifests
```

## Run locally

```bash
export DAGSTER_GRAPHQL_URL=http://localhost:3000/graphql
export DAGSTER_REPOSITORY_LOCATION=model_factory
export DAGSTER_REPOSITORY_NAME=__repository__
export DAGSTER_JOB_NAME=foundation_model_job

make install
make run
```

## Deploy

```bash
make manifests
make generate

make docker-build IMG=ghcr.io/teaglebuilt/model-factory-operator:dev
make docker-push IMG=ghcr.io/teaglebuilt/model-factory-operator:dev

make deploy IMG=ghcr.io/teaglebuilt/model-factory-operator:dev
```

## Create a project and experiment

```bash
kubectl apply -f config/samples/ai_v1alpha1_modelproject.yaml
kubectl apply -f config/samples/ai_v1alpha1_experiment.yaml
kubectl apply -f config/samples/ai_v1alpha1_modelrun.yaml
```

Watch:

```bash
kubectl get modelprojects,experiments,modelruns,modelartifacts,modelreleases
kubectl describe modelrun coder-smoke-001
```

## Current ModelRun behavior

The `ModelRun` reconciler:

1. Resolves `ModelProject`.
2. Resolves `Experiment`.
3. Verifies that the experiment belongs to the project.
4. Launches the existing Python `foundation_model_job` through Dagster GraphQL.
5. Stores the Dagster run ID in `ModelRun.status.dagsterRunId`.
6. Polls Dagster and maps Dagster status into `ModelRun.status.phase`.
7. Uses Kubernetes `metav1.Condition` for observable state.
8. Terminates the Dagster run when a `ModelRun` with an active run is deleted.

Dataset resolution, model artifact publication, and serving reconciliation are deliberately
left behind interfaces/boundaries for the next implementation pass.
