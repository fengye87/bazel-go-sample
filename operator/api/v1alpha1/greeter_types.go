package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GreeterSpec defines the desired state of Greeter
type GreeterSpec struct {
}

// GreeterStatus defines the observed state of Greeter
type GreeterStatus struct {
}

// +kubebuilder:object:root=true

// Greeter is the Schema for the greeters API
type Greeter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GreeterSpec   `json:"spec,omitempty"`
	Status GreeterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GreeterList contains a list of Greeter
type GreeterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Greeter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Greeter{}, &GreeterList{})
}
