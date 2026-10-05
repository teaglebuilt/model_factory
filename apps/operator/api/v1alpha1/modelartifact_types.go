package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ModelArtifactPhase string

const (
	ModelArtifactPhaseCandidate ModelArtifactPhase = "Candidate"
	ModelArtifactPhaseApproved  ModelArtifactPhase = "Approved"
	ModelArtifactPhaseRejected  ModelArtifactPhase = "Rejected"
)

type DatasetProvenance struct {
	LogicalURI string `json:"logicalUri"`
	Revision   string `json:"revision"`
	Digest     string `json:"digest,omitempty"`
}

type ModelArtifactSpec struct {
	ProjectRef LocalReference      `json:"projectRef"`
	RunRef     LocalReference      `json:"runRef"`
	Artifact   ArtifactLocation    `json:"artifact"`
	Format     string              `json:"format,omitempty"`
	Family     string              `json:"family,omitempty"`
	Tokenizer  *ArtifactLocation   `json:"tokenizer,omitempty"`
	GitSHA     string              `json:"gitSha,omitempty"`
	Datasets   []DatasetProvenance `json:"datasets,omitempty"`
}

type EvaluationStatus struct {
	Passed  bool               `json:"passed"`
	Metrics map[string]float64 `json:"metrics,omitempty"`
}

type ModelArtifactStatus struct {
	ObservedGeneration int64                       `json:"observedGeneration,omitempty"`
	Phase              ModelArtifactPhase          `json:"phase,omitempty"`
	Evaluations        map[string]EvaluationStatus `json:"evaluations,omitempty"`
	Conditions         []metav1.Condition          `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Project",type="string",JSONPath=".spec.projectRef.name"
// +kubebuilder:printcolumn:name="Run",type="string",JSONPath=".spec.runRef.name"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ModelArtifact struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelArtifactSpec   `json:"spec,omitempty"`
	Status            ModelArtifactStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ModelArtifactList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelArtifact `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelArtifact{}, &ModelArtifactList{})
}
