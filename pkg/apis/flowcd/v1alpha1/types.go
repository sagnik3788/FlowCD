package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FlowCD is the Schema for the flowcds API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type FlowCD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlowCDSpec   `json:"spec,omitempty"`
	Status FlowCDStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type FlowCDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FlowCD `json:"items"`
}

type FlowCDSpec struct {
	Source             ApplicationSource      `json:"source"`
	Destination        ApplicationDestination `json:"destination"`
	DeploymentStrategy DeploymentStrategy     `json:"deploymentStrategy,omitempty"`
}

type ApplicationSource struct {
	RepoURL string `json:"repoURL"`
	Branch  string `json:"branch"`
	Path    string `json:"path,omitempty"`
}

type ApplicationDestination struct {
	Server    string `json:"server,omitempty"`
	Namespace string `json:"namespace"`
}

type DeploymentStrategy struct {
	Type      string             `json:"type"`
	QuickSync *QuickSyncStrategy `json:"quickSync,omitempty"`
	Pipeline  *PipelineStrategy  `json:"pipeline,omitempty"`
	Custom    *CustomStrategy    `json:"custom,omitempty"`
}

type QuickSyncStrategy struct {
	Prune bool `json:"prune,omitempty"`
}

type PipelineStrategy struct {
	Stages []PipelineStage `json:"stages,omitempty"`
}

type PipelineStage struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CustomStrategy struct {
	Script  string `json:"script"`
	Timeout string `json:"timeout,omitempty"`
}

type FlowCDStatus struct {
	Sync       SyncStatus         `json:"sync,omitempty"`
	Message    string             `json:"message,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type SyncStatus struct {
	Status       string       `json:"status,omitempty"`
	Revision     string       `json:"revision,omitempty"`
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`
}
