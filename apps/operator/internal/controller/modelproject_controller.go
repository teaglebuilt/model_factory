package controller

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	aiv1alpha1 "github.com/teaglebuilt/model_factory/operator/api/v1alpha1"
)

// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ai.teaglebuilt.io,resources=modelprojects/finalizers,verbs=update

type ModelProjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ModelProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var project aiv1alpha1.ModelProject
	if err := r.Get(ctx, req.NamespacedName, &project); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	before := project.DeepCopy()
	project.Status.ObservedGeneration = project.Generation
	setCondition(
		&project.Status.Conditions,
		"Ready",
		metav1.ConditionTrue,
		"ConfigurationAccepted",
		"ModelProject configuration is valid.",
		project.Generation,
	)

	if err := r.Status().Patch(ctx, &project, client.MergeFrom(before)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ModelProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aiv1alpha1.ModelProject{}).
		Named("modelproject").
		Complete(r)
}
