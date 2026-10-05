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

// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelartifacts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelartifacts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelartifacts/finalizers,verbs=update
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelruns;modelprojects,verbs=get;list;watch

type ModelArtifactReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ModelArtifactReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var artifact aiv1alpha1.ModelArtifact
	if err := r.Get(ctx, req.NamespacedName, &artifact); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	before := artifact.DeepCopy()
	artifact.Status.ObservedGeneration = artifact.Generation

	var run aiv1alpha1.ModelRun
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: artifact.Namespace,
		Name:      artifact.Spec.RunRef.Name,
	}, &run); err != nil {
		if client.IgnoreNotFound(err) == nil {
			setCondition(&artifact.Status.Conditions, "Ready", metav1.ConditionFalse, "RunNotFound", "Referenced ModelRun does not exist.", artifact.Generation)
			if patchErr := r.Status().Patch(ctx, &artifact, client.MergeFrom(before)); patchErr != nil {
				return ctrl.Result{}, patchErr
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if run.Status.Phase != aiv1alpha1.ModelRunPhaseSucceeded {
		setCondition(&artifact.Status.Conditions, "Ready", metav1.ConditionFalse, "RunNotSucceeded", "Referenced ModelRun has not succeeded.", artifact.Generation)
		if err := r.Status().Patch(ctx, &artifact, client.MergeFrom(before)); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if artifact.Status.Phase == "" {
		artifact.Status.Phase = aiv1alpha1.ModelArtifactPhaseCandidate
	}
	setCondition(&artifact.Status.Conditions, "Ready", metav1.ConditionTrue, "ArtifactRegistered", "Artifact references a successful ModelRun.", artifact.Generation)

	if err := r.Status().Patch(ctx, &artifact, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ModelArtifactReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1alpha1.ModelArtifact{}).
		Named("modelartifact").
		Complete(r)
}
