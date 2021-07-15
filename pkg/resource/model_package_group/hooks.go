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

	condition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

var (
	ErrModelPackageGroupDeleting     = errors.New("ModelPackageGroup in 'DELETING' state, cannot be modified or deleted")
	ErrModelPackageGroupInProgress   = errors.New("ModelPackageGroup in 'INPROGRESS' state, cannot be modified or deleted")
	ErrModelPackageGroupPending      = errors.New("ModelPackageGroup in 'PENDING' state, cannot be modified or deleted")
	ErrModelPackageGroupDeleteFailed = errors.New("ModelPackageGroup in 'DELETEFAILED' state")
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		ErrModelPackageGroupDeleting,
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileInProgress = ackrequeue.NeededAfter(
		ErrModelPackageGroupInProgress,
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhilePending = ackrequeue.NeededAfter(
		ErrModelPackageGroupPending,
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileDeleteFailed = ackrequeue.NeededAfter(
		ErrModelPackageGroupDeleteFailed,
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

func CustomSetOutput(r *resource) (err error) {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return nil
	}
	ModelPackageGroupStatus := *r.ko.Status.ModelPackageGroupStatus
	msg := "ModelPackageGroup is in" + ModelPackageGroupStatus + "status"
	if ModelPackageGroupStatus == string(svcsdk.ModelPackageGroupStatusCompleted) || ModelPackageGroupStatus == string(svcsdk.ModelPackageGroupStatusFailed) {
		condition.SetSynced(r, corev1.ConditionTrue, &msg, nil)
		return nil
	}
	requeue := &ackrequeue.RequeueNeededAfter{}
	switch ModelPackageGroupStatus {
	case string(svcsdk.ModelPackageGroupStatusInProgress):
		requeue = requeueWaitWhileInProgress
	case string(svcsdk.ModelPackageGroupStatusDeleting):
		requeue = requeueWaitWhileDeleting
	case string(svcsdk.ModelPackageGroupStatusDeleteFailed):
		requeue = requeueWaitWhileDeleteFailed
	case string(svcsdk.ModelPackageGroupStatusPending):
		requeue = requeueWaitWhilePending
	}
	condition.SetSynced(r, corev1.ConditionFalse, &msg, nil)
	return requeue
}
