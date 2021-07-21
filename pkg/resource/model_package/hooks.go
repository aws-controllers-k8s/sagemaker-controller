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
	"context"
	"errors"

	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New("ModelPackage is deleting"),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

func (rm *resourceManager) customSetOutput(r *resource) {
	if r.ko.Status.ModelPackageStatus == nil {
		return
	}
	ModelPackageStatus := *r.ko.Status.ModelPackageStatus
	msg := "ModelPackage is in" + ModelPackageStatus + "status"
	if !isModifiying(r) {
		ackcondition.SetSynced(r, corev1.ConditionTrue, &msg, nil)
	} else {
		ackcondition.SetSynced(r, corev1.ConditionFalse, &msg, nil)
	}
}

func (rm *resourceManager) requeueUntilCanModify(ctx context.Context,
	latest *resource,
) error {
	if latest.ko.Status.ModelPackageStatus == nil {
		return nil
	}
	ModelPackageStatus := *latest.ko.Status.ModelPackageStatus
	if isModifiying(latest) {
		errMsg := "ModelPackage in" + ModelPackageStatus + "state cannot be modified or deleted"
		requeueWaitWhileModifying := ackrequeue.NeededAfter(
			errors.New(errMsg),
			ackrequeue.DefaultRequeueAfterDuration,
		)
		return requeueWaitWhileModifying
	}
	return nil
}

func isModifiying(r *resource) bool {
	if r == nil || r.ko.Status.ModelPackageStatus == nil {
		return false
	}
	ModelPackageStatus := *r.ko.Status.ModelPackageStatus

	if ModelPackageStatus == string(svcsdk.ModelPackageStatusInProgress) || ModelPackageStatus == string(svcsdk.ModelPackageStatusPending) || ModelPackageStatus == string(svcsdk.ModelPackageStatusDeleting) {
		return true
	}
	return false
}
