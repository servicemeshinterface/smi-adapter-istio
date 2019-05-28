package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceRoleSpec defines the desired state of ServiceRole
// +k8s:openapi-gen=true
type ServiceRoleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Required. The set of access rules (permissions) that the role has.
	Rules []*AccessRule `json:"rules,omitempty"`
}

// ServiceRoleStatus defines the observed state of ServiceRole
// +k8s:openapi-gen=true
type ServiceRoleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRole is the Schema for the serviceroles API
// +k8s:openapi-gen=true
type ServiceRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceRoleSpec   `json:"spec,omitempty"`
	Status ServiceRoleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleList contains a list of ServiceRole
type ServiceRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceRole{}, &ServiceRoleList{})
}
