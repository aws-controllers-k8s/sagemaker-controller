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
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// SetSyncedCondition sets the ACK Synced Condition
// status to true if the resource status exists but
// is not in one of the given modifying statuses,
// or false if the resource status is one of the
// given modifying statuses.
func SetSyncedCondition(
	r acktypes.AWSResource,
	latestStatus *string,
	resourceName *string,
	modifyingStatuses *[]string,
) {
	if latestStatus == nil {
		return
	}

	msg := *resourceName + " is in " + *latestStatus + " status."
	conditionStatus := corev1.ConditionTrue
	if IsModifyingStatus(latestStatus, modifyingStatuses) {
		conditionStatus = corev1.ConditionFalse
	}

	ackcondition.SetSynced(r, conditionStatus, &msg, nil)
}

// SetTerminalState sets conditions (terminal) on
// a resource's supplied status conditions if the
// latest status matches the terminal status.
// It returns true if conditions are updated.
func SetTerminalState(
	r acktypes.AWSResource,
	latestStatus *string,
	resourceName *string,
	terminalStatus string,
) bool {
	if latestStatus == nil || *latestStatus != terminalStatus {
		return false
	}

	terminalCondition := ackcondition.Terminal(r)
	if terminalCondition != nil && terminalCondition.Status == corev1.ConditionTrue {
		// some other exception already put the resource in terminal condition
		return false
	}

	// setting terminal condition since controller can no longer recover by retrying
	terminalMessage := *resourceName + " status reached terminal state: " + terminalStatus + ". Check the FailureReason."
	ackcondition.SetTerminal(r, corev1.ConditionTrue, &terminalMessage, nil)

	return true
}
