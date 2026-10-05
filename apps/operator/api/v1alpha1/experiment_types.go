package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ModelParameters struct {
	HiddenSize        int32 `json:"hiddenSize,omitempty"`
	NumLayers         int32 `json:"numLayers,omitempty"`
	NumAttentionHeads int32 `json:"numAttentionHeads,omitempty"`
	MaxSequenceLength int32 `json:"maxSequenceLength,omitempty"`
}

type ModelSpec struct {
	Family       string          `json:"family"`
	Architecture string          `json:"architecture"`
	Parameters   ModelParameters `json:"parameters,omitempty"`
}

type TrainingSpec struct {
	Backend                   string  `json:"backend,omitempty"`
	Seed                      int64   `json:"seed,omitempty"`
	MaxSteps                  int64   `json:"maxSteps,omitempty"`
	BatchSize                 int32   `json:"batchSize,omitempty"`
	GradientAccumulationSteps int32   `json:"gradientAccumulationSteps,omitempty"`
	LearningRate              float64 `json:"learningRate,omitempty"`
}

type ExperimentSpec struct {
	ProjectRef         LocalReference     `json:"projectRef"`
	Datasets           ExperimentDatasets `json:"datasets"`
	Model              ModelSpec          `json:"model"`
	Training           TrainingSpec       `json:"training"`
	Runtime            RuntimeSpec        `json:"runtime,omitempty"`
	EvaluationSuiteRef *LocalReference    `json:"evaluationSuiteRef,omitempty"`
}

type ExperimentStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Project",type="string",JSONPath=".spec.projectRef.name"
// +kubebuilder:printcolumn:name="Backend",type="string",JSONPath=".spec.training.backend"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ExperimentSpec   `json:"spec,omitempty"`
	Status            ExperimentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Experiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
}
