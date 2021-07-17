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
	"context"
	"errors"

	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New("ModelPackageGroup is deleting"),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

func (rm *resourceManager) customSetOutput(r *resource) {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return
	}
	ModelPackageGroupStatus := *r.ko.Status.ModelPackageGroupStatus
	msg := "ModelPackageGroup is in" + ModelPackageGroupStatus + "status"
	if !isModifiying(r) {
		ackcondition.SetSynced(r, corev1.ConditionTrue, &msg, nil)
	} else {
		ackcondition.SetSynced(r, corev1.ConditionFalse, &msg, nil)
	}
}

func (rm *resourceManager) customDeleteModelPackageGroup(ctx context.Context,
	latest *resource,
) error {
	if latest.ko.Status.ModelPackageGroupStatus == nil {
		return nil
	}
	ModelPackageGroupStatus := *latest.ko.Status.ModelPackageGroupStatus
	if isModifiying(latest) {
		errMsg := "ModelPackageGroup in" + ModelPackageGroupStatus + "state cannot be modified or deleted"
		requeueWaitWhileModifying := ackrequeue.NeededAfter(
			errors.New(errMsg),
			ackrequeue.DefaultRequeueAfterDuration,
		)
		return requeueWaitWhileModifying
	}
	return nil
}

func isModifiying(r *resource) bool {
	if r == nil || r.ko.Status.ModelPackageGroupStatus == nil {
		return false
	}
	ModelPackageGroupStatus := *r.ko.Status.ModelPackageGroupStatus

	if ModelPackageGroupStatus == string(svcsdk.ModelPackageGroupStatusInProgress) || ModelPackageGroupStatus == string(svcsdk.ModelPackageGroupStatusPending) || ModelPackageGroupStatus == string(svcsdk.ModelPackageGroupStatusDeleting) {
		return true
	}
	return false
}
