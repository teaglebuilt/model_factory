package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ServingSpec struct {
	Backend   string       `json:"backend,omitempty"`
	Replicas  int32        `json:"replicas,omitempty"`
	Resources ResourceSpec `json:"resources,omitempty"`
	ModelName string       `json:"modelName,omitempty"`
}

type ModelReleaseSpec struct {
	ProjectRef  LocalReference `json:"projectRef"`
	ModelRef    LocalReference `json:"modelRef"`
	Environment string         `json:"environment"`
	Serving     ServingSpec    `json:"serving,omitempty"`
}

type ModelReleaseStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Phase              string             `json:"phase,omitempty"`
	Endpoint           string             `json:"endpoint,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Project",type="string",JSONPath=".spec.projectRef.name"
// +kubebuilder:printcolumn:name="Model",type="string",JSONPath=".spec.modelRef.name"
// +kubebuilder:printcolumn:name="Environment",type="string",JSONPath=".spec.environment"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ModelRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelReleaseSpec   `json:"spec,omitempty"`
	Status            ModelReleaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ModelReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelRelease{}, &ModelReleaseList{})
}
