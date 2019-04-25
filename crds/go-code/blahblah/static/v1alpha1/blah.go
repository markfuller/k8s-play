package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Blah struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlahSpec   `json:"spec"`
	Status BlahStatus `json:"status"`
}
type BlahList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Blah `json:"items"`
}

type BlahSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
}

type BlahStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}
