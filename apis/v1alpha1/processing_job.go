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

// ProcessingJobSpec defines the desired state of ProcessingJob.
//
// An Amazon SageMaker processing job that is used to analyze data and evaluate
// models. For more information, see Process Data and Evaluate Models (https://docs.aws.amazon.com/sagemaker/latest/dg/processing-job.html).
type ProcessingJobSpec struct {
	// Configures the processing job to run a specified Docker container image.
	// +kubebuilder:validation:Required
	AppSpecification *AppSpecification `json:"appSpecification"`
	// The environment variables to set in the Docker container. Up to 100 key and
	// values entries in the map are supported.
	Environment map[string]*string `json:"environment,omitempty"`

	ExperimentConfig *ExperimentConfig `json:"experimentConfig,omitempty"`
	// Networking options for a processing job, such as whether to allow inbound
	// and outbound network calls to and from processing containers, and the VPC
	// subnets and security groups to use for VPC-enabled processing jobs.
	NetworkConfig *NetworkConfig `json:"networkConfig,omitempty"`
	// An array of inputs configuring the data to download into the processing container.
	ProcessingInputs []*ProcessingInput `json:"processingInputs,omitempty"`
	// The name of the processing job. The name must be unique within an AWS Region
	// in the AWS account.
	// +kubebuilder:validation:Required
	ProcessingJobName *string `json:"processingJobName"`
	// Output configuration for the processing job.
	ProcessingOutputConfig *ProcessingOutputConfig `json:"processingOutputConfig,omitempty"`
	// Identifies the resources, ML compute instances, and ML storage volumes to
	// deploy for a processing job. In distributed training, you specify more than
	// one instance.
	// +kubebuilder:validation:Required
	ProcessingResources *ProcessingResources `json:"processingResources"`
	// The Amazon Resource Name (ARN) of an IAM role that Amazon SageMaker can assume
	// to perform tasks on your behalf.
	// +kubebuilder:validation:Required
	RoleARN *string `json:"roleARN"`
	// The time limit for how long the processing job is allowed to run.
	StoppingCondition *ProcessingStoppingCondition `json:"stoppingCondition,omitempty"`
}

// ProcessingJobStatus defines the observed state of ProcessingJob
type ProcessingJobStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// A string, up to one KB in size, that contains the reason a processing job
	// failed, if it failed.
	FailureReason *string `json:"failureReason,omitempty"`
	// Provides the status of a processing job.
	ProcessingJobStatus *string `json:"processingJobStatus,omitempty"`
}

// ProcessingJob is the Schema for the ProcessingJobs API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type ProcessingJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProcessingJobSpec   `json:"spec,omitempty"`
	Status            ProcessingJobStatus `json:"status,omitempty"`
}

// ProcessingJobList contains a list of ProcessingJob
// +kubebuilder:object:root=true
type ProcessingJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProcessingJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProcessingJob{}, &ProcessingJobList{})
}
