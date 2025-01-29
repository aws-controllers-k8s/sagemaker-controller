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

package user_profile

import (
	"context"
	"errors"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

var (
	modifyingStatuses = []string{
		string(svcsdktypes.UserProfileStatusPending),
		string(svcsdktypes.UserProfileStatusUpdating),
		string(svcsdktypes.UserProfileStatusDeleting),
	}

	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Deleting."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// requeueUntilCanModify creates and returns an
// ackrequeue error if a resource's latest status matches
// any of the defined modifying statuses below.
func (rm *resourceManager) requeueUntilCanModify(ctx context.Context, r *resource) error {
	latestStatus := r.ko.Status.Status
	return svccommon.RequeueIfModifying(latestStatus, &resourceName, &modifyingStatuses)
}

// Sets the ResourceSynced condition to False if resource is being modified by AWS
func (rm *resourceManager) customDescribeUserProfileSetOutput(ko *svcapitypes.UserProfile) {
	svccommon.SetSyncedCondition(&resource{ko}, ko.Status.Status, &resourceName, &modifyingStatuses)
}
