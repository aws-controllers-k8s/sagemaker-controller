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

// Use this file if the Status/Spec of the CR needs to be modified after
// create/describe/update operation

package endpoint

import (
	"context"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

// customCreateEndpointSetOutput sets the resource in TempOutofSync if endpoint is
// in creating state. At this stage we know call to createEndpoint was successful.
func (rm *resourceManager) customCreateEndpointSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.CreateEndpointOutput,
	ko *svcapitypes.Endpoint,
) (*svcapitypes.Endpoint, error) {
	rm.customSetOutput(r, aws.String(svcsdk.EndpointStatusCreating), ko)
	return ko, nil
}

// customDescribeEndpointSetOutput sets the resource in TempOutofSync if endpoint is
// being modified by AWS
func (rm *resourceManager) customDescribeEndpointSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.DescribeEndpointOutput,
	ko *svcapitypes.Endpoint,
) (*svcapitypes.Endpoint, error) {
	rm.customSetOutput(r, resp.EndpointStatus, ko)
	// Workaround: Field config for LatestEndpointConfigName of generator config
	// does not code generate this correctly since this field is part of Spec
	// SageMaker users will need following information:
	// 	 - latestEndpointConfig
	// 	 - desiredEndpointConfig
	// 	 - LastEndpointConfigNameForUpdate
	// 	 - FailureReason
	// to determine the correct course of action in case of update to Endpoint fails
	if resp.EndpointConfigName != nil {
		ko.Status.LatestEndpointConfigName = resp.EndpointConfigName
	} else {
		ko.Status.LatestEndpointConfigName = nil
	}
	return ko, nil
}

// customUpdateEndpointSetOutput sets the resource in TempOutofSync if endpoint is
// being updated. At this stage we know call to updateEndpoint was successful.
func (rm *resourceManager) customUpdateEndpointSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.UpdateEndpointOutput,
	ko *svcapitypes.Endpoint,
) (*svcapitypes.Endpoint, error) {
	rm.customSetOutput(r, aws.String(svcsdk.EndpointStatusUpdating), ko)
	// no nil check present here since Spec.EndpointConfigName is a required field
	ko.Status.LastEndpointConfigNameForUpdate = r.ko.Spec.EndpointConfigName
	return ko, nil
}

// customSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the endpoint status on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customSetOutput(
	r *resource,
	endpointStatus *string,
	ko *svcapitypes.Endpoint,
) {
	if endpointStatus == nil {
		return
	}

	syncConditionStatus := corev1.ConditionUnknown
	if *endpointStatus == svcsdk.EndpointStatusInService || *endpointStatus == svcsdk.EndpointStatusFailed {
		syncConditionStatus = corev1.ConditionTrue
	} else {
		syncConditionStatus = corev1.ConditionFalse
	}

	var resourceSyncedCondition *ackv1alpha1.Condition = nil
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	} else {
		for _, condition := range ko.Status.Conditions {
			if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
				resourceSyncedCondition = condition
				break
			}
		}
	}

	if resourceSyncedCondition == nil {
		resourceSyncedCondition = &ackv1alpha1.Condition{
			Type:   ackv1alpha1.ConditionTypeResourceSynced,
			Status: syncConditionStatus,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
	} else {
		resourceSyncedCondition.Status = syncConditionStatus
	}

}
