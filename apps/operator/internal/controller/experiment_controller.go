package controller

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	aiv1alpha1 "github.com/teaglebuilt/model_factory/operator/api/v1alpha1"
)

// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=experiments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=experiments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=experiments/finalizers,verbs=update
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelprojects,verbs=get;list;watch

type ExperimentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var experiment aiv1alpha1.Experiment
	if err := r.Get(ctx, req.NamespacedName, &experiment); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	before := experiment.DeepCopy()
	experiment.Status.ObservedGeneration = experiment.Generation

	var project aiv1alpha1.ModelProject
	err := r.Get(ctx, types.NamespacedName{
		Namespace: experiment.Namespace,
		Name:      experiment.Spec.ProjectRef.Name,
	}, &project)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			setCondition(
				&experiment.Status.Conditions,
				"Ready",
				metav1.ConditionFalse,
				"ProjectNotFound",
				"Referenced ModelProject does not exist.",
				experiment.Generation,
			)
			if patchErr := r.Status().Patch(ctx, &experiment, client.MergeFrom(before)); patchErr != nil {
				return ctrl.Result{}, patchErr
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	setCondition(
		&experiment.Status.Conditions,
		"Ready",
		metav1.ConditionTrue,
		"ReferencesResolved",
		"Experiment references a valid ModelProject.",
		experiment.Generation,
	)

	if err := r.Status().Patch(ctx, &experiment, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1alpha1.Experiment{}).
		Named("experiment").
		Complete(r)
}
