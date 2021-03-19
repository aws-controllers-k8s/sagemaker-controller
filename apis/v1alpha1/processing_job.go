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

// ProcessingJobSpec defines the desired state of ProcessingJob
type ProcessingJobSpec struct {
	// +kubebuilder:validation:Required
	AppSpecification *AppSpecification  `json:"appSpecification"`
	Environment      map[string]*string `json:"environment,omitempty"`
	ExperimentConfig *ExperimentConfig  `json:"experimentConfig,omitempty"`
	NetworkConfig    *NetworkConfig     `json:"networkConfig,omitempty"`
	ProcessingInputs []*ProcessingInput `json:"processingInputs,omitempty"`
	// +kubebuilder:validation:Required
	ProcessingJobName      *string                 `json:"processingJobName"`
	ProcessingOutputConfig *ProcessingOutputConfig `json:"processingOutputConfig,omitempty"`
	// +kubebuilder:validation:Required
	ProcessingResources *ProcessingResources `json:"processingResources"`
	// +kubebuilder:validation:Required
	RoleARN           *string                      `json:"roleARN"`
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
	Conditions          []*ackv1alpha1.Condition `json:"conditions"`
	FailureReason       *string                  `json:"failureReason,omitempty"`
	ProcessingJobStatus *string                  `json:"processingJobStatus,omitempty"`
}

// ProcessingJob is the Schema for the ProcessingJobs API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FailureReason",type=string,JSONPath=`.status.failureReason`
// +kubebuilder:printcolumn:name="ProcessingJobStatus",type=string,JSONPath=`.status.processingJobStatus`
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
