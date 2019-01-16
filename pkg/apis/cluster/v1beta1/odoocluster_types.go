/*
 * This file is part of the Odoo-Operator (R) project.
 * Copyright (c) 2018-2018 XOE Corp. SAS
 * Authors: David Arnold, et al.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * ALTERNATIVE LICENCING OPTION
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license.  Buying such a license is mandatory as soon as you
 * develop commercial activities involving the Odoo-Operator software without
 * disclosing the source code of your own applications. These activities
 * include: Offering paid services to a customer as an ASP, shipping Odoo-
 * Operator with a closed source product.
 *
 */

package v1beta1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OdooImageRepoSpec defines the cluster-level image configuration
type OdooImageRepoSpec struct {
	Registry string `json:"registry"`
	Repo     string `json:"repo"`
	// +optional
	Secret string `json:"secret,omitempty"`
}

// OdooResource defines the desired resource spec
type OdooResource struct {
	Tiers   []OdooTierSpec   `json:"tiers,omitempty"`
	Volumes []OdooVolumeSpec `json:"volumes,omitempty"`
}

// OdooTierSpec specs a tier
type OdooTierSpec struct {
	Name string `json:"name"`
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
	// +optional
	QOS *v1.PodQOSClass `json:"qos,omitempty"`
	// +optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`
}

// OdooVolumeSpec specs a volume
type OdooVolumeSpec struct {
	Name string `json:"name"`
	// +optional
	PVCSpec v1.PersistentVolumeClaimSpec `json:"pvcspec"`
}

// OdooClusterStatus defines the observed state of OdooCluster
type OdooClusterStatus struct {
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []OdooClusterStatusCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// OdooClusterStatusCondition defines an observable OdooClusterStatus condition
type OdooClusterStatusCondition struct {
	// Type of the OdooClusterStatus condition.
	Type OdooClusterStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=OdooClusterStatusConditionType"`
	StatusCondition
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OdooCluster is the Schema for the odooclusters API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OdooCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OdooClusterSpec   `json:"spec,omitempty"`
	Status OdooClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OdooClusterList contains a list of OdooCluster
type OdooClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OdooCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OdooCluster{}, &OdooClusterList{})
}
