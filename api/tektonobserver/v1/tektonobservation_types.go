/*
Copyright 2024 kcloutie.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TektonObservationSpec defines the desired state of TektonObservation
type TektonObservationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// PubSubTopics is a list of PubSub topics to which the controller will publish events
	PubSubTopics []PubSubTopic `json:"pubSubTopics" yaml:"pubSubTopics"`
}
type PubSubTopic struct {
	// ProjectID is the GCP project ID where the PubSub topic is located
	PubSubProjectID string `json:"pubSubProjectID" yaml:"pubSubProjectID"`
	// PubSubTopicID is the ID of the PubSub topic
	PubSubTopicID string `json:"pubSubTopicID" yaml:"pubSubTopicID"`
}

// TektonObservationStatus defines the observed state of TektonObservation
type TektonObservationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TektonObservation is the Schema for the tektonobservations API
type TektonObservation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TektonObservationSpec   `json:"spec,omitempty"`
	Status TektonObservationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TektonObservationList contains a list of TektonObservation
type TektonObservationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TektonObservation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TektonObservation{}, &TektonObservationList{})
}
