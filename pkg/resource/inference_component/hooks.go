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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

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
		string(svcsdktypes.InferenceComponentStatusCreating),
		string(svcsdktypes.InferenceComponentStatusUpdating),
		string(svcsdktypes.InferenceComponentStatusDeleting),
	}

	resourceName = GroupKind.Kind

	lastSpecForUpdateAnnotation = fmt.Sprintf("%s/last-spec-for-update", GroupKind.Group)

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
func (rm *resourceManager) customUpdateInferenceComponentSetOutput(ko *svcapitypes.InferenceComponent) error {
	//set last inference component spec used for update in annotations
	annotations := ko.ObjectMeta.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	spec := ko.Spec.DeepCopy()
	spec.Tags = nil
	serializedSpec, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	annotations[lastSpecForUpdateAnnotation] = string(serializedSpec)
	ko.ObjectMeta.SetAnnotations(annotations)

	// injecting Updating status to keep the Sync condition message and status.InferenceComponentStatus in sync
	ko.Status.InferenceComponentStatus = aws.String(string(svcsdktypes.InferenceComponentStatusUpdating))

	latestStatus := ko.Status.InferenceComponentStatus
	svccommon.SetSyncedCondition(&resource{ko}, latestStatus, &resourceName, &modifyingStatuses)

	return nil
}

// customUpdateInferenceComponentPreChecks adds specialized logic to check if controller should
// proceed with UpdateInferenceComponent call.
// Update is blocked in the following cases:
//  1. while InferenceComponentStatus != InService (handled by requeueUntilCanModify method).
//  2. InferenceComponentStatus == Failed.
//  3. A previous update to the InferenceComponent with same spec failed.
//
// Method returns nil if InferenceComponent can be updated,
// otherwise InferenceComponentUpdateError depending on above cases.
func (rm *resourceManager) customUpdateInferenceComponentPreChecks(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) error {
	latestStatus := latest.ko.Status.InferenceComponentStatus
	if latestStatus == nil {
		return nil
	}

	failureReason := latest.ko.Status.FailureReason

	desiredSpec := desired.ko.Spec.DeepCopy()
	desiredSpec.Tags = nil

	var lastSpecForUpdateString *string = nil
	// get last endpoint config name used for update from annotations
	annotations := desired.ko.ObjectMeta.GetAnnotations()
	for k, v := range annotations {
		if k == lastSpecForUpdateAnnotation {
			lastSpecForUpdateString = &v
		}
	}

	var lastSpecForUpdate *svcapitypes.InferenceComponentSpec

	if lastSpecForUpdateString != nil {
		err := json.Unmarshal([]byte(*lastSpecForUpdateString), &lastSpecForUpdate)
		if err != nil {
			return err
		}
	}

	// Case 2 - InferenceComponentStatus == Failed
	if *latestStatus == string(svcsdktypes.InferenceComponentStatusFailed) ||
		// Case 3 - A previous update to the InferenceComponent with same spec failed
		// Following checks indicate FailureReason is related to a failed update
		(failureReason != nil && lastSpecForUpdateString != nil &&
			EqualInferenceComponentSpec(desiredSpec, lastSpecForUpdate)) {
		// 1. FailureReason alone doesn't mean an update failed it can appear because of other
		// reasons(patching/scaling failed).
		// 2. desiredSpec == lastSpecForUpdate only tells us an update was tried with lastSpecForUpdate
		// but does not tell us anything if the update was successful or not in the past because
		// it is set if updateInferenceComponent returns 200 (async operation).
		// 3. Now, sdkUpdate can execute because of change in any field in Spec.

		// 1 & 2 does not guarantee an update Failed. Hence, we need to look at `lastSpecForUpdate` to determine if the update was unsuccessful
		// `desiredSpec != latestSpec` + `desiredSpec == lastSpecForUpdate
		//`+ `FailureReason != nil` indicate that an update is needed, has already been tried and failed.
		return &smithy.GenericAPIError{
			Code:    "InferenceComponentUpdateError",
			Message: "Unable to update inference component." + " Check FailureReason.",
			Fault:   1,
		}
	}

	return nil
}

// EqualInferenceComponentSpec checks if two InferenceComponentSpec instances are equal
func EqualInferenceComponentSpec(desiredSpec *svcapitypes.InferenceComponentSpec,
	lastSpec *svcapitypes.InferenceComponentSpec) bool {
	if desiredSpec == nil || lastSpec == nil {
		return desiredSpec == lastSpec
	}
	return reflect.DeepEqual(desiredSpec, lastSpec)
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
