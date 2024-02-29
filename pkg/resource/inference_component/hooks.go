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

package inference_component

import (
	"context"
	"errors"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	modifyingStatuses = []string{
		svcsdk.InferenceComponentStatusCreating,
		svcsdk.InferenceComponentStatusUpdating,
		svcsdk.InferenceComponentStatusDeleting,
	}

	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Deleting."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customDescribeInferenceComponentSetOutput sets the resource ResourceSynced condition to False if
// InferenceComponent is being modified by AWS
func (rm *resourceManager) customDescribeInferenceComponentSetOutput(ko *svcapitypes.InferenceComponent) {
	latestStatus := ko.Status.InferenceComponentStatus
	svccommon.SetSyncedCondition(&resource{ko}, latestStatus, &resourceName, &modifyingStatuses)
}

// customUpdateInferenceComponentSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the InferenceComponentStatus on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customUpdateInferenceComponentSetOutput(ko *svcapitypes.InferenceComponent) {

	// injecting Updating status to keep the Sync condition message and status.InferenceComponentStatus in sync
	//ko.Status.InferenceComponentStatus = aws.String(svcsdk.InferenceComponentStatusUpdating)

	latestStatus := ko.Status.InferenceComponentStatus
	svccommon.SetSyncedCondition(&resource{ko}, latestStatus, &resourceName, &modifyingStatuses)
}

// requeueUntilCanModify creates and returns an ackrequeue error
// if a resource's latest status matches any of the defined modifying statuses.
// This is so the controller requeues until the resource can be modifed
func (rm *resourceManager) requeueUntilCanModify(
	ctx context.Context,
	r *resource,
) error {
	latestStatus := r.ko.Status.InferenceComponentStatus
	return svccommon.RequeueIfModifying(latestStatus, &resourceName, &modifyingStatuses)
}
