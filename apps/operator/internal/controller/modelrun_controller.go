package controller

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	aiv1alpha1 "github.com/teaglebuilt/model_factory/operator/api/v1alpha1"
	"github.com/teaglebuilt/model_factory/operator/internal/dagster"
)

const modelRunFinalizer = "ai.teaglebuilt.io/dagster-run"

// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelruns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelruns/finalizers,verbs=update
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelprojects;experiments,verbs=get;list;watch

type ModelRunReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	Dagster        dagster.Client
	DagsterJobName string
}

func (r *ModelRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var run aiv1alpha1.ModelRun
	if err := r.Get(ctx, req.NamespacedName, &run); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !run.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, &run)
	}

	if !controllerutil.ContainsFinalizer(&run, modelRunFinalizer) {
		before := run.DeepCopy()
		controllerutil.AddFinalizer(&run, modelRunFinalizer)
		if err := r.Patch(ctx, &run, client.MergeFrom(before)); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	project, experiment, err := r.resolveReferences(ctx, &run)
	if err != nil {
		return r.setFailedCondition(ctx, &run, "ReferenceError", err.Error())
	}

	if experiment.Spec.ProjectRef.Name != project.Name {
		return r.setFailedCondition(
			ctx,
			&run,
			"ProjectMismatch",
			"Experiment belongs to a different ModelProject.",
		)
	}

	if run.Status.DagsterRunID == "" {
		return r.launchDagsterRun(ctx, &run, experiment)
	}

	status, err := r.Dagster.GetRunStatus(ctx, run.Status.DagsterRunID)
	if err != nil {
		return ctrl.Result{RequeueAfter: 15 * time.Second}, err
	}

	return r.syncDagsterStatus(ctx, &run, status)
}

func (r *ModelRunReconciler) resolveReferences(
	ctx context.Context,
	run *aiv1alpha1.ModelRun,
) (*aiv1alpha1.ModelProject, *aiv1alpha1.Experiment, error) {
	var project aiv1alpha1.ModelProject
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: run.Namespace,
		Name:      run.Spec.ProjectRef.Name,
	}, &project); err != nil {
		return nil, nil, fmt.Errorf("resolve ModelProject: %w", err)
	}

	var experiment aiv1alpha1.Experiment
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: run.Namespace,
		Name:      run.Spec.ExperimentRef.Name,
	}, &experiment); err != nil {
		return nil, nil, fmt.Errorf("resolve Experiment: %w", err)
	}

	return &project, &experiment, nil
}

func (r *ModelRunReconciler) launchDagsterRun(
	ctx context.Context,
	run *aiv1alpha1.ModelRun,
	experiment *aiv1alpha1.Experiment,
) (ctrl.Result, error) {
	runConfig := map[string]any{
		"ops": map[string]any{
			"resolve_dataset": map[string]any{
				"config": map[string]any{
					"dataset_uri": experiment.Spec.Datasets.Pretraining.URI,
				},
			},
		},
	}

	dagsterRunID, err := r.Dagster.LaunchRun(ctx, dagster.LaunchRequest{
		JobName:       r.DagsterJobName,
		RunConfigData: runConfig,
	})
	if err != nil {
		return r.setFailedCondition(ctx, run, "DagsterLaunchFailed", err.Error())
	}

	before := run.DeepCopy()
	run.Status.ObservedGeneration = run.Generation
	run.Status.DagsterRunID = dagsterRunID
	run.Status.Phase = aiv1alpha1.ModelRunPhaseQueued
	run.Status.MaxSteps = experiment.Spec.Training.MaxSteps

	setCondition(
		&run.Status.Conditions,
		"WorkflowSubmitted",
		metav1.ConditionTrue,
		"DagsterRunCreated",
		fmt.Sprintf("Dagster run %s was created.", dagsterRunID),
		run.Generation,
	)

	if err := r.Status().Patch(ctx, run, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

func (r *ModelRunReconciler) syncDagsterStatus(
	ctx context.Context,
	run *aiv1alpha1.ModelRun,
	status dagster.RunStatus,
) (ctrl.Result, error) {
	before := run.DeepCopy()
	run.Status.ObservedGeneration = run.Generation

	switch status {
	case dagster.RunStatusQueued, dagster.RunStatusNotStarted, dagster.RunStatusStarting:
		run.Status.Phase = aiv1alpha1.ModelRunPhaseQueued
		setCondition(&run.Status.Conditions, "Running", metav1.ConditionFalse, "Queued", "Dagster run is queued.", run.Generation)

	case dagster.RunStatusStarted:
		run.Status.Phase = aiv1alpha1.ModelRunPhaseTraining
		setCondition(&run.Status.Conditions, "Running", metav1.ConditionTrue, "DagsterRunning", "Dagster run is executing.", run.Generation)

	case dagster.RunStatusSuccess:
		run.Status.Phase = aiv1alpha1.ModelRunPhaseSucceeded
		setCondition(&run.Status.Conditions, "Completed", metav1.ConditionTrue, "DagsterSucceeded", "Dagster run succeeded.", run.Generation)

	case dagster.RunStatusFailure:
		run.Status.Phase = aiv1alpha1.ModelRunPhaseFailed
		setCondition(&run.Status.Conditions, "Completed", metav1.ConditionFalse, "DagsterFailed", "Dagster run failed.", run.Generation)

	case dagster.RunStatusCanceled, dagster.RunStatusCanceling:
		run.Status.Phase = aiv1alpha1.ModelRunPhaseCancelled
		setCondition(&run.Status.Conditions, "Completed", metav1.ConditionFalse, "DagsterCancelled", "Dagster run was cancelled.", run.Generation)

	default:
		return ctrl.Result{RequeueAfter: 15 * time.Second}, fmt.Errorf("unknown Dagster run status %q", status)
	}

	if err := r.Status().Patch(ctx, run, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}

	switch run.Status.Phase {
	case aiv1alpha1.ModelRunPhaseSucceeded,
		aiv1alpha1.ModelRunPhaseFailed,
		aiv1alpha1.ModelRunPhaseCancelled:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

func (r *ModelRunReconciler) setFailedCondition(
	ctx context.Context,
	run *aiv1alpha1.ModelRun,
	reason string,
	message string,
) (ctrl.Result, error) {
	before := run.DeepCopy()
	run.Status.ObservedGeneration = run.Generation
	run.Status.Phase = aiv1alpha1.ModelRunPhaseFailed
	setCondition(
		&run.Status.Conditions,
		"Ready",
		metav1.ConditionFalse,
		reason,
		message,
		run.Generation,
	)
	if err := r.Status().Patch(ctx, run, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ModelRunReconciler) reconcileDelete(
	ctx context.Context,
	run *aiv1alpha1.ModelRun,
) (ctrl.Result, error) {
	if !controllerutil.ContainsFinalizer(run, modelRunFinalizer) {
		return ctrl.Result{}, nil
	}

	if run.Status.DagsterRunID != "" &&
		run.Status.Phase != aiv1alpha1.ModelRunPhaseSucceeded &&
		run.Status.Phase != aiv1alpha1.ModelRunPhaseFailed &&
		run.Status.Phase != aiv1alpha1.ModelRunPhaseCancelled {
		if err := r.Dagster.TerminateRun(ctx, run.Status.DagsterRunID); err != nil {
			return ctrl.Result{RequeueAfter: 10 * time.Second}, err
		}
	}

	before := run.DeepCopy()
	controllerutil.RemoveFinalizer(run, modelRunFinalizer)
	if err := r.Patch(ctx, run, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ModelRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1alpha1.ModelRun{}).
		Named("modelrun").
		Complete(r)
}
