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

package model_package_group

import (
	"errors"
	"time"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	"github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var (
	ErrModelPackageGroupDeleting     = errors.New("ModelPackageGroup in 'DELETING' state, cannot be modified or deleted")
	ErrModelPackageGroupInProgress   = errors.New("ModelPackageGroup in 'INPROGRESS' state, cannot be modified or deleted")
	ErrModelPackageGroupPending      = errors.New("ModelPackageGroup in 'PENDING' state, cannot be modified or deleted")
	ErrModelPackageGroupDeleteFailed = errors.New("ModelPackageGroup in 'DELETEFAILED' state")
)

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// ModelPackageGroup
	TerminalStatuses = []v1alpha1.ModelPackageGroupStatus_SDK{
		v1alpha1.ModelPackageGroupStatus_SDK_Deleting,
	}
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		ErrModelPackageGroupDeleting,
		5*time.Second,
	)
	requeueWaitWhileInProgress = ackrequeue.NeededAfter(
		ErrModelPackageGroupInProgress,
		5*time.Second,
	)
	requeueWaitWhilePending = ackrequeue.NeededAfter(
		ErrModelPackageGroupPending,
		5*time.Second,
	)
	requeueWaitWhileDeleteFailed = ackrequeue.NeededAfter(
		ErrModelPackageGroupDeleteFailed,
		5*time.Second,
	)
)

// modelPackageGroupHasTerminalStatus returns whether the supplied SageMaker ModelPackageGroup is in a
// terminal state
func modelPackageGroupHasTerminalStatus(r *resource) bool {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return false
	}
	ts := *r.ko.Status.ModelPackageGroupStatus
	for _, s := range TerminalStatuses {
		if ts == string(s) {
			return true
		}
	}
	return false
}

// isModelPackageGroupPending returns true if the supplied SageMaker ModelPackageGroup is in the process
// of pending
func isModelPackageGroupPending(r *resource) bool {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return false
	}
	sagemaker_status := *r.ko.Status.ModelPackageGroupStatus
	return sagemaker_status == string(v1alpha1.ModelPackageGroupStatus_SDK_Pending)
}

// isModelPackageGroupProgressreturns true if the supplied SageMaker ModelPackageGroup is in progress
func isModelPackageGroupInProgress(r *resource) bool {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return false
	}
	sagemaker_status := *r.ko.Status.ModelPackageGroupStatus
	return sagemaker_status == string(v1alpha1.ModelPackageGroupStatus_SDK_InProgress)
}

// isModelPackageGroupDeleting returns true if the supplied SageMaker ModelPackageGroup is in the process
// of being deleted
func isModelPackageGroupDeleting(r *resource) bool {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return false
	}
	sagemaker_status := *r.ko.Status.ModelPackageGroupStatus
	return sagemaker_status == string(v1alpha1.ModelPackageGroupStatus_SDK_Deleting)
}

// isModelPackageGroupDeleting returns true if the supplied SageMaker ModelPackageGroup delete failed
func isModelPackageGroupDeleteFailed(r *resource) bool {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return false
	}
	sagemaker_status := *r.ko.Status.ModelPackageGroupStatus
	return sagemaker_status == string(v1alpha1.ModelPackageGroupStatus_SDK_DeleteFailed)
}

func ModelPackageGroupCustomSetOutput(r *resource) {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return
	}
	sagemaker_status := *r.ko.Status.ModelPackageGroupStatus
	if sagemaker_status == string(v1alpha1.ModelPackageGroupStatus_SDK_Completed) || sagemaker_status == string(v1alpha1.ModelPackageGroupStatus_SDK_Failed) {
		setSyncedCondition(r, corev1.ConditionTrue, nil, nil)
	} else {
		setSyncedCondition(r, corev1.ConditionFalse, nil, nil)
	}
}
