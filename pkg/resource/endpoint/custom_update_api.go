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

// Use this file to add custom implementation for any operation of intercept
// the autogenerated code that trigger an update on an endpoint

package endpoint

import (
	"context"
	"errors"
	"strings"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
	"github.com/aws/aws-sdk-go/aws/awserr"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	FailUpdateError = awserr.New("EndpointUpdateError", "unable to update endpoint. check FailureReason", nil)

	FailureReasonInternalServiceErrorPrefix = "Request to service failed"
)

// customUpdateEndpoint adds specialized logic to check if controller should
// proceeed with updateEndpoint call.
// Update is blocked in the following cases:
//  1. while EndpointStatus != InService
//  2. EndpointStatus == Failed
//  3. A previous update to the Endpoint with same endpointConfigName failed
// Method returns nil if endpoint can be updated, otherwise error depending on above cases
func (rm *resourceManager) customUpdateEndpoint(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	latestStatus := latest.ko.Status.EndpointStatus
	if latestStatus == nil {
		return nil, nil
	}

	if *latestStatus != svcsdk.EndpointStatusFailed {
		// Case 1 - requeueAfter until endpoint is in InService state
		err := rm.endpointStatusAllowUpdates(ctx, latest)
		if err != nil {
			return nil, err
		}
	}

	failureReason := latest.ko.Status.FailureReason
	latestEndpointConfig := latest.ko.Spec.EndpointConfigName
	desiredEndpointConfig := desired.ko.Spec.EndpointConfigName
	lastEndpointConfigForUpdate := desired.ko.Status.LastEndpointConfigNameForUpdate

	// Case 2 - EndpointStatus == Failed
	if *latestStatus == svcsdk.EndpointStatusFailed ||
		// Case 3 - A previous update to the Endpoint with same endpointConfigName failed
		// Following checks indicate FailureReason is related to a failed update
		// Note: Internal service error is an exception for this case
		// "Request to service failed" means update failed because of ISE and can be retried
		(failureReason != nil && lastEndpointConfigForUpdate != nil &&
			!strings.HasPrefix(*failureReason, FailureReasonInternalServiceErrorPrefix) &&
			*desiredEndpointConfig != *latestEndpointConfig &&
			*desiredEndpointConfig == *lastEndpointConfigForUpdate) {
		// 1. FailureReason alone does mean an update failed it can appear because of other reasons(patching/scaling failed)
		// 2. *desiredEndpointConfig == *lastEndpointConfigForUpdate only tells us an update was tried with lastEndpointConfigForUpdate
		// but does not tell us anything if the update was successful or not in the past because it is set if updateEndpoint returns 200 (aync operation).
		// 3. Now, sdkUpdate can execute because of change in any field in Spec (like tags/deploymentConfig in future)

		// 1 & 2 does not guarantee an update Failed. Hence we need to look at `*latestEndpointConfigName` to determine if the update was unsuccessful
		// `*desiredEndpointConfig != *latestEndpointConfig` + `*desiredEndpointConfig == *lastEndpointConfigForUpdate`+ `FailureReason != nil` indicate that an update is needed,
		// has already been tried and failed.
		return nil, FailUpdateError
	}

	return nil, nil
}

// customDeleteEndpoint adds specialized logic to requeueAfter until endpoint is in
// InService or Failed state before a deleteEndpoint can be called
func (rm *resourceManager) customDeleteEndpoint(
	ctx context.Context,
	latest *resource,
) error {
	latestStatus := latest.ko.Status.EndpointStatus
	if latestStatus != nil && *latestStatus == svcsdk.EndpointStatusFailed {
		return nil
	}
	return rm.endpointStatusAllowUpdates(ctx, latest)
}

// endpointStatusAllowUpdates is a helper method to determine if endpoint allows modification
func (rm *resourceManager) endpointStatusAllowUpdates(
	ctx context.Context,
	r *resource,
) error {
	latestStatus := r.ko.Status.EndpointStatus
	if latestStatus != nil && *latestStatus != svcsdk.EndpointStatusInService {
		return requeue.NeededAfter(
			errors.New("endpoint status does not allow modification, it is not in 'InService' state"),
			requeue.DefaultRequeueAfterDuration)
	}

	return nil
}
