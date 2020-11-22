/*
Copyright 2020 LQ.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppTemplate struct {
	Name string  `json:"name"`
	Spec AppSpec `json:"spec,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FuwuSpec defines the desired state of Fuwu
type FuwuSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Selector string `json:"selector"`
	//Apps     []AppTemplate `json:"apps,omitempty"`
}

// FuwuStatus defines the observed state of Fuwu
type FuwuStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Status        string `json:"status"`
	Phase         string `json:"phase,omitempty"`
	AvailableApps int32  `json:"availableApps,omitempty"`
	TotalApps     int32  `json:"totalApps,omitempty"`
}

// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Selector",type="string",JSONPath=".spec.Selector",description="The selector of Fuwu"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of Fuwu"

// Fuwu is the Schema for the fuwus API
type Fuwu struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FuwuSpec   `json:"spec,omitempty"`
	Status FuwuStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FuwuList contains a list of Fuwu
type FuwuList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fuwu `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fuwu{}, &FuwuList{})
}
