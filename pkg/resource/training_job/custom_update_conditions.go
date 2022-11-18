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
// of training job which is not evident from API response

package training_job

import (
	"strings"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var (
	terminalCode string = "[ACK_SM]"
)

// If the controller runs into an error that contains "[ACK_SM]"
// it will set the resource to a terminal state because it is an unrecoverable error.
func (rm *resourceManager) CustomUpdateConditions(
	ko *svcapitypes.TrainingJob,
	r *resource,
	err error,
) bool {

	if ackcompare.IsNil(err) {
		return false
	}

	if strings.Contains(err.Error(), terminalCode) {
		conditionManager := &resource{ko}
		exception := err.Error()
		ackcondition.SetTerminal(conditionManager, corev1.ConditionTrue, &exception, nil)
		return true
	}

	return false

}
