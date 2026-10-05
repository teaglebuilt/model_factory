package v1alpha1

import corev1 "k8s.io/api/core/v1"

type LocalReference struct {
	Name string `json:"name"`
}

type DatasetRef struct {
	URI string `json:"uri"`
}

type ExperimentDatasets struct {
	Pretraining DatasetRef  `json:"pretraining"`
	Instruction *DatasetRef `json:"instruction,omitempty"`
	Evaluation  *DatasetRef `json:"evaluation,omitempty"`
}

type ResourceSpec struct {
	GPU    int32               `json:"gpu,omitempty"`
	CPU    string              `json:"cpu,omitempty"`
	Memory string              `json:"memory,omitempty"`
	Extra  corev1.ResourceList `json:"extra,omitempty"`
}

type RuntimeSpec struct {
	Queue        string            `json:"queue,omitempty"`
	Resources    ResourceSpec      `json:"resources,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

type ArtifactLocation struct {
	URI    string `json:"uri"`
	Digest string `json:"digest,omitempty"`
}

type ResolvedDatasetRef struct {
	LogicalURI  string `json:"logicalUri"`
	SourceType  string `json:"sourceType"`
	PhysicalURI string `json:"physicalUri"`
	Revision    string `json:"revision"`
	Digest      string `json:"digest,omitempty"`
}

type CheckpointStatus struct {
	URI  string `json:"uri"`
	Step int64  `json:"step"`
}
