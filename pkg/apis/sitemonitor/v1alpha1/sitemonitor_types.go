package v1alpha1

import (
	"wen/site-monitor-operator/pkg/monitor"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SiteMonitorSpec defines the desired state of SiteMonitor
// +k8s:openapi-gen=true
type SiteMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	monitor.SiteMonitor
}

// SiteMonitorStatus defines the observed state of SiteMonitor
// +k8s:openapi-gen=true
type SiteMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Disabled bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SiteMonitor is the Schema for the sitemonitors API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type SiteMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SiteMonitorSpec   `json:"spec,omitempty"`
	Status SiteMonitorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SiteMonitorList contains a list of SiteMonitor
type SiteMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SiteMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SiteMonitor{}, &SiteMonitorList{})
}
