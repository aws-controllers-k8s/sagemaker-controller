// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ModelBiasJobDefinitionSpec defines the desired state of ModelBiasJobDefinition
type ModelBiasJobDefinitionSpec struct {
	// +kubebuilder:validation:Required
	JobDefinitionName *string `json:"jobDefinitionName"`
	// +kubebuilder:validation:Required
	JobResources *MonitoringResources `json:"jobResources"`
	// +kubebuilder:validation:Required
	ModelBiasAppSpecification *ModelBiasAppSpecification `json:"modelBiasAppSpecification"`
	ModelBiasBaselineConfig   *ModelBiasBaselineConfig   `json:"modelBiasBaselineConfig,omitempty"`
	// +kubebuilder:validation:Required
	ModelBiasJobInput *ModelBiasJobInput `json:"modelBiasJobInput"`
	// +kubebuilder:validation:Required
	ModelBiasJobOutputConfig *MonitoringOutputConfig  `json:"modelBiasJobOutputConfig"`
	NetworkConfig            *MonitoringNetworkConfig `json:"networkConfig,omitempty"`
	// +kubebuilder:validation:Required
	RoleARN           *string                      `json:"roleARN"`
	StoppingCondition *MonitoringStoppingCondition `json:"stoppingCondition,omitempty"`
	Tags              []*Tag                       `json:"tags,omitempty"`
}

// ModelBiasJobDefinitionStatus defines the observed state of ModelBiasJobDefinition
type ModelBiasJobDefinitionStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	Conditions       []*ackv1alpha1.Condition `json:"conditions"`
	JobDefinitionARN *string                  `json:"jobDefinitionARN,omitempty"`
}

// ModelBiasJobDefinition is the Schema for the ModelBiasJobDefinitions API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type ModelBiasJobDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelBiasJobDefinitionSpec   `json:"spec,omitempty"`
	Status            ModelBiasJobDefinitionStatus `json:"status,omitempty"`
}

// ModelBiasJobDefinitionList contains a list of ModelBiasJobDefinition
// +kubebuilder:object:root=true
type ModelBiasJobDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelBiasJobDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelBiasJobDefinition{}, &ModelBiasJobDefinitionList{})
}
