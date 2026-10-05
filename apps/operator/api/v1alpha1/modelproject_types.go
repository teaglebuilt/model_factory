package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ModelProjectDefaults struct {
	Runtime         RuntimeSpec `json:"runtime,omitempty"`
	CheckpointStore string      `json:"checkpointStore,omitempty"`
	ArtifactStore   string      `json:"artifactStore,omitempty"`
}

type ModelProjectSpec struct {
	Description       string               `json:"description,omitempty"`
	DatasetNamespace  string               `json:"datasetNamespace,omitempty"`
	RegistryNamespace string               `json:"registryNamespace,omitempty"`
	Defaults          ModelProjectDefaults `json:"defaults,omitempty"`
}

type ModelProjectStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ModelProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelProjectSpec   `json:"spec,omitempty"`
	Status            ModelProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ModelProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelProject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelProject{}, &ModelProjectList{})
}
