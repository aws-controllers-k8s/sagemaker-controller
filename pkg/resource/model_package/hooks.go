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

package model_package

import (
	"errors"
	"time"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	"github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var (
	ErrModelPackageDeleting   = errors.New("ModelPackage in 'DELETING' state, cannot be modified or deleted")
	ErrModelPackagePending    = errors.New("ModelPackage in 'PENDING' state, cannot be modified or deleted")
	ErrModelPackageInProgress = errors.New("ModelPackage in 'INPROGRESS' state, cannot be modified or deleted")
)

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// SageMaker ModelPackage
	TerminalStatuses = []v1alpha1.ModelPackageStatus_SDK{
		v1alpha1.ModelPackageStatus_SDK_Deleting,
	}
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		ErrModelPackageDeleting,
		5*time.Second,
	)
	requeueWaitWhilePending = ackrequeue.NeededAfter(
		ErrModelPackagePending,
		5*time.Second,
	)
	requeueWaitWhileInProgress = ackrequeue.NeededAfter(
		ErrModelPackageInProgress,
		5*time.Second,
	)
)

// ModelPackageHasTerminalStatus returns whether the supplied SageMaker ModelPackage is in a
// terminal state
func ModelPackageHasTerminalStatus(r *resource) bool {
	if r.ko.Status.ModelPackageStatus == nil {
		return false
	}
	ts := *r.ko.Status.ModelPackageStatus
	for _, s := range TerminalStatuses {
		if ts == string(s) {
			return true
		}
	}
	return false
}

// isModelPackagePending returns true if the supplied SageMaker ModelPackage is in the process
// of pending
func isModelPackagePending(r *resource) bool {
	if r.ko.Status.ModelPackageStatus == nil {
		return false
	}
	dbis := *r.ko.Status.ModelPackageStatus
	return dbis == string(v1alpha1.ModelPackageStatus_SDK_Pending)
}

// isModelPackageDeleting returns true if the supplied SageMaker ModelPackage is in the process
// of being deleted
func isModelPackageDeleting(r *resource) bool {
	if r.ko.Status.ModelPackageStatus == nil {
		return false
	}
	dbis := *r.ko.Status.ModelPackageStatus
	return dbis == string(v1alpha1.ModelPackageStatus_SDK_Deleting)
}

// isModelPackageInProgress returns true if the supplied SageMaker ModelPackage is in the process
// of being in progress
func isModelPackageInProgress(r *resource) bool {
	if r.ko.Status.ModelPackageStatus == nil {
		return false
	}
	dbis := *r.ko.Status.ModelPackageStatus
	return dbis == string(v1alpha1.ModelPackageStatus_SDK_InProgress)
}

func ModelPackageCustomSetOutput(r *resource) {
	if r.ko.Status.ModelPackageStatus == nil {
		return
	}
	sagemaker_status := *r.ko.Status.ModelPackageStatus
	if sagemaker_status == string(v1alpha1.ModelPackageStatus_SDK_Completed) || sagemaker_status == string(v1alpha1.ModelPackageStatus_SDK_Failed) {
		setSyncedCondition(r, corev1.ConditionTrue, nil, nil)
	} else {
		setSyncedCondition(r, corev1.ConditionFalse, nil, nil)
	}
}
