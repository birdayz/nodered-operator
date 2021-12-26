/*
Copyright 2021.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NoderedSpec defines the desired state of Nodered
type NoderedSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Nodered. Edit nodered_types.go to remove/update
	Image    string          `json:"image,omitempty"`
	Settings NoderedSettings `json:"settings,omitempty"`
}

type NoderedSettings struct {
	EditorTheme NoderedSettingsEditorTheme `json:"editorTheme,omitempty"`
}

type NoderedSettingsEditorTheme struct {
	Page NoderedSettingsEditorThemePage `json:"page,omitempty"`
}

type NoderedSettingsEditorThemePage struct {
	Title string `json:"title,omitempty"`
}

// NoderedStatus defines the observed state of Nodered
type NoderedStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Nodered is the Schema for the nodereds API
type Nodered struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NoderedSpec   `json:"spec,omitempty"`
	Status NoderedStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NoderedList contains a list of Nodered
type NoderedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nodered `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nodered{}, &NoderedList{})
}
