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
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

var (
	modifyingStatuses = []string{svcsdk.ModelPackageGroupStatusInProgress,
		svcsdk.ModelPackageGroupStatusPending,
		svcsdk.ModelPackageGroupStatusDeleting}

	resourceName = resourceGK.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is deleting."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

func (rm *resourceManager) customSetOutput(r *resource) {
	if r.ko.Status.ModelPackageGroupStatus == nil {
		return
	}
	ModelPackageGroupStatus := *r.ko.Status.ModelPackageGroupStatus
	msg := "ModelPackageGroup is in" + ModelPackageGroupStatus + "status"
	if !svccommon.IsModifyingStatus(&ModelPackageGroupStatus, &modifyingStatuses) {
		ackcondition.SetSynced(r, corev1.ConditionTrue, &msg, nil)
	} else {
		ackcondition.SetSynced(r, corev1.ConditionFalse, &msg, nil)
	}
}

// requeueUntilCanModify creates and returns an
// ackrequeue error if a resource's latest status matches
// any of the defined modifying statuses below.
func (rm *resourceManager) requeueUntilCanModify(
	ctx context.Context,
	r *resource,
) error {
	latestStatus := r.ko.Status.ModelPackageGroupStatus
	return svccommon.RequeueIfModifying(latestStatus, &resourceName, &modifyingStatuses)
}
