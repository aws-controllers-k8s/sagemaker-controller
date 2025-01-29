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
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	corev1 "k8s.io/api/core/v1"
)

// CustomUpdateConditions sets conditions (terminal) on supplied inference component.
// it examines supplied resource to determine conditions.
// It returns true if conditions are updated.
func (rm *resourceManager) CustomUpdateConditions(
	ko *svcapitypes.InferenceComponent,
	r *resource,
	err error,
) bool {
	latestStatus := r.ko.Status.InferenceComponentStatus
	terminalStatus := string(svcsdktypes.InferenceComponentStatusFailed)
	conditionManager := &resource{ko}
	resourceName := GroupKind.Kind
	// If the latestStatus == terminalStatus we will set
	// the terminal condition and terminal message.
	updated := svccommon.SetTerminalState(conditionManager, latestStatus, &resourceName, terminalStatus)

	// Continue setting ResourceSynced condition to false in case of failed update
	// since desired and latest will be different until the issue is fixed.
	// Customer can use this condition state and FailureReason to determine
	// the correct course of action in case the update to InferenceComponent fails.
	if err != nil {
		awsErr, ok := ackerr.AWSError(err)
		if ok && awsErr.ErrorCode() == "InferenceComponentUpdateError" {
			ackcondition.SetSynced(conditionManager, corev1.ConditionFalse, aws.String(awsErr.Error()), nil)
			return true
		}
	}

	return updated
}
