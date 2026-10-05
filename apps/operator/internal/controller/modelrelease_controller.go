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

// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelreleases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelreleases/finalizers,verbs=update
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelartifacts;modelprojects,verbs=get;list;watch

type ModelReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ModelReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var release aiv1alpha1.ModelRelease
	if err := r.Get(ctx, req.NamespacedName, &release); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	before := release.DeepCopy()
	release.Status.ObservedGeneration = release.Generation

	var artifact aiv1alpha1.ModelArtifact
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: release.Namespace,
		Name:      release.Spec.ModelRef.Name,
	}, &artifact); err != nil {
		if client.IgnoreNotFound(err) == nil {
			release.Status.Phase = "Blocked"
			setCondition(&release.Status.Conditions, "Ready", metav1.ConditionFalse, "ArtifactNotFound", "Referenced ModelArtifact does not exist.", release.Generation)
			if patchErr := r.Status().Patch(ctx, &release, client.MergeFrom(before)); patchErr != nil {
				return ctrl.Result{}, patchErr
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if artifact.Status.Phase != aiv1alpha1.ModelArtifactPhaseApproved {
		release.Status.Phase = "Blocked"
		setCondition(&release.Status.Conditions, "Ready", metav1.ConditionFalse, "ArtifactNotApproved", "ModelArtifact must be Approved before release.", release.Generation)
	} else {
		release.Status.Phase = "Pending"
		setCondition(&release.Status.Conditions, "Ready", metav1.ConditionFalse, "ServingIntegrationPending", "Artifact is approved; serving reconciliation is the next implementation step.", release.Generation)
	}

	if err := r.Status().Patch(ctx, &release, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *ModelReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1alpha1.ModelRelease{}).
		Named("modelrelease").
		Complete(r)
}
