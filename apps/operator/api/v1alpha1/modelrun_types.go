package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ModelRunPhase string

const (
	ModelRunPhasePending   ModelRunPhase = "Pending"
	ModelRunPhaseQueued    ModelRunPhase = "Queued"
	ModelRunPhaseTraining  ModelRunPhase = "Training"
	ModelRunPhaseSucceeded ModelRunPhase = "Succeeded"
	ModelRunPhaseFailed    ModelRunPhase = "Failed"
	ModelRunPhaseCancelled ModelRunPhase = "Cancelled"
)

type ModelRunSpec struct {
	ProjectRef    LocalReference `json:"projectRef"`
	ExperimentRef LocalReference `json:"experimentRef"`
	Runtime       RuntimeSpec    `json:"runtime,omitempty"`
}

type ModelRunStatus struct {
	ObservedGeneration int64                         `json:"observedGeneration,omitempty"`
	Phase              ModelRunPhase                 `json:"phase,omitempty"`
	DagsterRunID       string                        `json:"dagsterRunId,omitempty"`
	ResolvedDatasets   map[string]ResolvedDatasetRef `json:"resolvedDatasets,omitempty"`
	CurrentStep        int64                         `json:"currentStep,omitempty"`
	MaxSteps           int64                         `json:"maxSteps,omitempty"`
	LatestCheckpoint   *CheckpointStatus             `json:"latestCheckpoint,omitempty"`
	ArtifactRef        *LocalReference               `json:"artifactRef,omitempty"`
	Conditions         []metav1.Condition            `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Project",type="string",JSONPath=".spec.projectRef.name"
// +kubebuilder:printcolumn:name="Experiment",type="string",JSONPath=".spec.experimentRef.name"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Dagster Run",type="string",JSONPath=".status.dagsterRunId"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ModelRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelRunSpec   `json:"spec,omitempty"`
	Status            ModelRunStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ModelRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelRun{}, &ModelRunList{})
}
