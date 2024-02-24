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

// InferenceComponentSpec defines the desired state of InferenceComponent.
type InferenceComponentSpec struct {

	// The name of an existing endpoint where you host the inference component.
	// +kubebuilder:validation:Required
	EndpointName *string `json:"endpointName"`
	// A unique name to assign to the inference component.
	// +kubebuilder:validation:Required
	InferenceComponentName *string `json:"inferenceComponentName"`
	// Runtime settings for a model that is deployed with an inference component.
	// +kubebuilder:validation:Required
	RuntimeConfig *InferenceComponentRuntimeConfig `json:"runtimeConfig"`
	// Details about the resources to deploy with this inference component, including
	// the model, container, and compute resources.
	// +kubebuilder:validation:Required
	Specification *InferenceComponentSpecification `json:"specification"`
	// A list of key-value pairs associated with the model. For more information,
	// see Tagging Amazon Web Services resources (https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html)
	// in the Amazon Web Services General Reference.
	Tags []*Tag `json:"tags,omitempty"`
	// The name of an existing production variant where you host the inference component.
	// +kubebuilder:validation:Required
	VariantName *string `json:"variantName"`
}

// InferenceComponentStatus defines the observed state of InferenceComponent
type InferenceComponentStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
}

// InferenceComponent is the Schema for the InferenceComponents API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type InferenceComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              InferenceComponentSpec   `json:"spec,omitempty"`
	Status            InferenceComponentStatus `json:"status,omitempty"`
}

// InferenceComponentList contains a list of InferenceComponent
// +kubebuilder:object:root=true
type InferenceComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InferenceComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InferenceComponent{}, &InferenceComponentList{})
}
