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

package endpoint

import (
	"context"
	"errors"
	"fmt"
	"strings"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/smithy-go"
)

var (
	modifyingStatuses = []string{
		string(svcsdktypes.EndpointStatusCreating),
		string(svcsdktypes.EndpointStatusUpdating),
		string(svcsdktypes.EndpointStatusSystemUpdating),
		string(svcsdktypes.EndpointStatusRollingBack),
		string(svcsdktypes.EndpointStatusDeleting),
	}

	resourceName = GroupKind.Kind

	lastEndpointConfigForUpdateAnnotation = fmt.Sprintf("%s/last-endpoint-config-for-update", GroupKind.Group)

	FailureReasonInternalServiceErrorPrefix = "Request to service failed"

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Deleting."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customDescribeEndpointSetOutput sets the resource ResourceSynced condition to False if endpoint is
// being modified by AWS
func (rm *resourceManager) customDescribeEndpointSetOutput(ko *svcapitypes.Endpoint) {
	svccommon.SetSyncedCondition(&resource{ko}, ko.Status.EndpointStatus, &resourceName, &modifyingStatuses)
}

// customUpdateEndpointSetOutput sets the resource ResourceSynced condition to False if endpoint is
// being updated. At this stage we know call to updateEndpoint was successful.
func (rm *resourceManager) customUpdateEndpointSetOutput(ko *svcapitypes.Endpoint) {
	// set last endpoint config name used for udapte in annotations
	// no nil check present since Spec.EndpointConfigName is a required field
	annotations := ko.ObjectMeta.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[lastEndpointConfigForUpdateAnnotation] = *ko.Spec.EndpointConfigName
	ko.ObjectMeta.SetAnnotations(annotations)

	// injecting Updating status to keep the Sync condition message and status.endpointStatus in sync
	ko.Status.EndpointStatus = aws.String(string(svcsdktypes.EndpointStatusUpdating))

	svccommon.SetSyncedCondition(&resource{ko}, ko.Status.EndpointStatus, &resourceName, &modifyingStatuses)
}

// customUpdateEndpointPreChecks adds specialized logic to check if controller should
// proceeed with updateEndpoint call.
// Update is blocked in the following cases:
//  1. while EndpointStatus != InService (handled by requeueUntilCanModify method)
//  2. EndpointStatus == Failed
//  3. A previous update to the Endpoint with same endpointConfigName failed
//
// Method returns nil if endpoint can be updated, otherwise error depending on above cases
func (rm *resourceManager) customUpdateEndpointPreChecks(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) error {
	latestStatus := latest.ko.Status.EndpointStatus
	if latestStatus == nil {
		return nil
	}

	failureReason := latest.ko.Status.FailureReason
	desiredEndpointConfig := desired.ko.Spec.EndpointConfigName

	var lastEndpointConfigForUpdate *string = nil
	// get last endpoint config name used for update from annotations
	annotations := desired.ko.ObjectMeta.GetAnnotations()
	for k, v := range annotations {
		if k == lastEndpointConfigForUpdateAnnotation {
			lastEndpointConfigForUpdate = &v
		}
	}

	// Case 2 - EndpointStatus == Failed
	if *latestStatus == string(svcsdktypes.EndpointStatusFailed) ||
		// Case 3 - A previous update to the Endpoint with same endpointConfigName failed
		// Following checks indicate FailureReason is related to a failed update
		// Note: Internal service error is an exception for this case
		// "Request to service failed" means update failed because of ISE and can be retried
		(failureReason != nil && lastEndpointConfigForUpdate != nil &&
			!strings.HasPrefix(*failureReason, FailureReasonInternalServiceErrorPrefix) &&
			delta.DifferentAt("Spec.EndpointConfigName") &&
			*desiredEndpointConfig == *lastEndpointConfigForUpdate) {
		// 1. FailureReason alone does mean an update failed it can appear because of other reasons(patching/scaling failed)
		// 2. *desiredEndpointConfig == *lastEndpointConfigForUpdate only tells us an update was tried with lastEndpointConfigForUpdate
		// but does not tell us anything if the update was successful or not in the past because it is set if updateEndpoint returns 200 (aync operation).
		// 3. Now, sdkUpdate can execute because of change in any field in Spec (like tags/deploymentConfig in future)

		// 1 & 2 does not guarantee an update Failed. Hence we need to look at `*latestEndpointConfigName` to determine if the update was unsuccessful
		// `*desiredEndpointConfig != *latestEndpointConfig` + `*desiredEndpointConfig == *lastEndpointConfigForUpdate`+ `FailureReason != nil` indicate that an update is needed,
		// has already been tried and failed.
		return &smithy.GenericAPIError{
			Code:    "EndpointUpdateError",
			Message: fmt.Sprintf("unable to update endpoint. check FailureReason. latest EndpointConfigName is %s", *latest.ko.Spec.EndpointConfigName),
			Fault:   0,
		}
	}

	return nil
}

// requeueUntilCanModify creates and returns an ackrequeue error
// if a resource's latest status matches any of the defined modifying statuses.
// This is so the controller requeues until the resource can be modifed
func (rm *resourceManager) requeueUntilCanModify(
	ctx context.Context,
	r *resource,
) error {
	latestStatus := r.ko.Status.EndpointStatus
	return svccommon.RequeueIfModifying(latestStatus, &resourceName, &modifyingStatuses)
}
