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

// Use this file if conditions need to be updated based on the latest status
// of endpoint which is not evident from API response

package endpoint

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

// CustomUpdateConditions sets conditions (terminal) on supplied endpoint
// it examines supplied resource to determine conditions.
// It returns true if conditions are updated
func (rm *resourceManager) customUpdateConditions(
	ko *svcapitypes.Endpoint,
	r *resource,
	err error,
) bool {
	latestStatus := r.ko.Status.EndpointStatus
	failureReason := r.ko.Status.FailureReason

	if latestStatus == nil || failureReason == nil {
		return false
	}
	var terminalCondition *ackv1alpha1.Condition = nil
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	} else {
		for _, condition := range ko.Status.Conditions {
			if condition.Type == ackv1alpha1.ConditionTypeTerminal {
				terminalCondition = condition
				break
			}
		}
		if terminalCondition != nil && terminalCondition.Status == corev1.ConditionTrue {
			// some other exception already put the resource in terminal condition
			return false
		}
	}

	if (err != nil && err == FailUpdateError) || (latestStatus != nil && *latestStatus == svcsdk.EndpointStatusFailed) {
		// setting terminal condition since controller can no longer recover by retrying
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		terminalCondition.Status = corev1.ConditionTrue
		if *latestStatus == svcsdk.EndpointStatusFailed {
			terminalCondition.Message = aws.String("Cannot update endpoint with Failed status")
		} else {
			terminalCondition.Message = aws.String(FailUpdateError.Error())
		}
		return true
	}

	return false
}
