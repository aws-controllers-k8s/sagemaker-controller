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

// ModelSpec defines the desired state of Model.
type ModelSpec struct {
	// Specifies the containers in the inference pipeline.
	Containers []*ContainerDefinition `json:"containers,omitempty"`
	// Isolates the model container. No inbound or outbound network calls can be
	// made to or from the model container.
	EnableNetworkIsolation *bool `json:"enableNetworkIsolation,omitempty"`
	// The Amazon Resource Name (ARN) of the IAM role that Amazon SageMaker can
	// assume to access model artifacts and docker image for deployment on ML compute
	// instances or for batch transform jobs. Deploying on ML compute instances
	// is part of model hosting. For more information, see Amazon SageMaker Roles
	// (https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html).
	//
	// To be able to pass this role to Amazon SageMaker, the caller of this API
	// must have the iam:PassRole permission.
	// +kubebuilder:validation:Required
	ExecutionRoleARN *string `json:"executionRoleARN"`
	// Specifies details of how containers in a multi-container endpoint are called.
	InferenceExecutionConfig *InferenceExecutionConfig `json:"inferenceExecutionConfig,omitempty"`
	// The name of the new model.
	// +kubebuilder:validation:Required
	ModelName *string `json:"modelName"`
	// The location of the primary docker image containing inference code, associated
	// artifacts, and custom environment map that the inference code uses when the
	// model is deployed for predictions.
	PrimaryContainer *ContainerDefinition `json:"primaryContainer,omitempty"`
	// An array of key-value pairs. You can use tags to categorize your AWS resources
	// in different ways, for example, by purpose, owner, or environment. For more
	// information, see Tagging AWS Resources (https://docs.aws.amazon.com/general/latest/gr/aws_tagging.html).
	Tags []*Tag `json:"tags,omitempty"`
	// A VpcConfig object that specifies the VPC that you want your model to connect
	// to. Control access to and from your model container by configuring the VPC.
	// VpcConfig is used in hosting services and in batch transform. For more information,
	// see Protect Endpoints by Using an Amazon Virtual Private Cloud (https://docs.aws.amazon.com/sagemaker/latest/dg/host-vpc.html)
	// and Protect Data in Batch Transform Jobs by Using an Amazon Virtual Private
	// Cloud (https://docs.aws.amazon.com/sagemaker/latest/dg/batch-vpc.html).
	VPCConfig *VPCConfig `json:"vpcConfig,omitempty"`
}

// ModelStatus defines the observed state of Model
type ModelStatus struct {
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

// Model is the Schema for the Models API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ModelSpec   `json:"spec,omitempty"`
	Status            ModelStatus `json:"status,omitempty"`
}

// ModelList contains a list of Model
// +kubebuilder:object:root=true
type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Model `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Model{}, &ModelList{})
}
