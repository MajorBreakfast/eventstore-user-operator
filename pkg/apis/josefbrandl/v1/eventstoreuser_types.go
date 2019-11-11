package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventStoreUserSpec defines the desired state of EventStoreUser
// +k8s:openapi-gen=true
type EventStoreUserSpec struct {
	// EventStore defines the Event Store on which the user should be created
	EventStore string `json:"eventStore"`

	// Groups defines the list of groups the Event Store user should belong to
	// +optional
	Groups []string `json:"groups,omitempty"`
}

// EventStoreUserStatus defines the observed state of EventStoreUser
// +k8s:openapi-gen=true
type EventStoreUserStatus struct{}

// EventStoreUser is the Schema for the eventstoreusers API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventstoreusers,scope=Namespaced,shortName=esuser;esu
type EventStoreUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventStoreUserSpec   `json:"spec,omitempty"`
	Status EventStoreUserStatus `json:"status,omitempty"`
}

// EventStoreUserList contains a list of EventStoreUser
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EventStoreUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventStoreUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventStoreUser{}, &EventStoreUserList{})
}
