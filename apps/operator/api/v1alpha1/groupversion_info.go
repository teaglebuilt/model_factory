// Package v1alpha1 contains API Schema definitions for ai.teaglebuilt.io/v1alpha1.
// +kubebuilder:object:generate=true
// +groupName=ai.teaglebuilt.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var GroupVersion = schema.GroupVersion{Group: "ai.teaglebuilt.io", Version: "v1alpha1"}
var SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
var AddToScheme = SchemeBuilder.AddToScheme
