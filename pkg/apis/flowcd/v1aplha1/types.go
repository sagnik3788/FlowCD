package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FlowCD struct {
    // Kubernetes object metadata
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    
    // custom fields
    Spec   FlowCDSpec   `json:"spec"`
    Status FlowCDStatus `json:"status,omitempty"`
}

type FlowCDSpec struct {
	Source      ApplicationSource      `json:"source"`
	Destination ApplicationDestination `json:"destination"`
	DeploymentStrategy DeploymentStrategy `json:"deploymentStrategy,omitempty"`
}

type ApplicationSource struct {
	// url of the source repo
	RepoURL string `json:"repoURL"`

	// specific branch you want to deploy
	Branch  string `json:"branch"`

	// path of the manifest files
	Path    string `json:"path,omitempty"`
}

type ApplicationDestination struct {
	// cluster API server
	Server    string `json:"server,omitempty"`
	// target namespace
	Namespace string `json:"namespace"`
}

type DeploymentStrategy struct {
	QuickSync *QuickSyncStrategy `json:"quickSync,omitempty"`
    Pipeline  *PipelineStrategy  `json:"pipeline,omitempty"`
    Custom    *CustomStrategy    `json:"custom,omitempty"`
}
type FlowCDStatus struct {
	Sync       SyncStatus       `json:"sync,omitempty"`
}