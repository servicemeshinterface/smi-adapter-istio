package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceRoleBindingSpec defines the desired state of ServiceRoleBinding
// +k8s:openapi-gen=true
type ServiceRoleBindingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Required. List of subjects that are assigned the ServiceRole object.
	Subjects []*Subject `json:"subjects,omitempty"`
	// Required. Reference to the ServiceRole object.
	RoleRef *RoleRef `json:"roleRef,omitempty"`
}

// ServiceRoleBindingStatus defines the observed state of ServiceRoleBinding
// +k8s:openapi-gen=true
type ServiceRoleBindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleBinding is the Schema for the servicerolebindings API
// +k8s:openapi-gen=true
type ServiceRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceRoleBindingSpec   `json:"spec,omitempty"`
	Status ServiceRoleBindingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleBindingList contains a list of ServiceRoleBinding
type ServiceRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceRoleBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceRoleBinding{}, &ServiceRoleBindingList{})
}
