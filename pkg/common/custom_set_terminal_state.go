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

package common

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// SetTerminalState sets conditions (terminal) on
// a resource's supplied status conditions if the
// latest status matches the terminal status.
// It returns true if conditions are updated.
func SetTerminalState(
	koStatusConditions *[]*ackv1alpha1.Condition,
	latestStatus *string,
	terminalStatus string,
	terminalMessage *string,
) bool {
	if latestStatus == nil || *latestStatus != terminalStatus {
		return false
	}

	var terminalCondition *ackv1alpha1.Condition = nil
	for _, condition := range *koStatusConditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
			break
		}
	}
	if terminalCondition != nil && terminalCondition.Status == corev1.ConditionTrue {
		// some other exception already put the resource in terminal condition
		return false
	}

	// setting terminal condition since controller can no longer recover by retrying
	if terminalCondition == nil {
		terminalCondition = &ackv1alpha1.Condition{
			Type: ackv1alpha1.ConditionTypeTerminal,
		}
		*koStatusConditions = append(*koStatusConditions, terminalCondition)
	}
	terminalCondition.Status = corev1.ConditionTrue
	terminalCondition.Message = terminalMessage
	return true
}
