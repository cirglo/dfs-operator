/*
Copyright 2025.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DFSSpec defines the desired state of DFS
type DFSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Required
	NumDataNodeServers *int32 `json:"numDataNodeServers"`

	// +kubebuilder:default:=1
	// +kubebuilder:validation:Required
	NumNameNodeServers *int32 `json:"numNameNodeServers"`

	// +kubebuilder:default:=5135
	// +kubebuilder:validation:Required`
	DataNodePort *int32 `json:"dataNodePort"`

	// +kubebuilder:default:=50070
	// +kubebuilder:validation:Required`
	NameNodePort *int32 `json:"nameNodePort"`

	// +kubebuilder:default:=default
	// +kubebuilder:validation:Required`
	DataNodeStorageClassName *string `json:"dataNodetorageClassName"`

	// +kubebuilder:default:=10Gi
	// +kubebuilder:validation:Required`
	DataNodeStorageSize *string `json:"dataNodeStorageSize"`
}

// DFSStatus defines the observed state of DFS
type DFSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DFS is the Schema for the dfs API
type DFS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DFSSpec   `json:"spec,omitempty"`
	Status DFSStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DFSList contains a list of DFS
type DFSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DFS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DFS{}, &DFSList{})
}
