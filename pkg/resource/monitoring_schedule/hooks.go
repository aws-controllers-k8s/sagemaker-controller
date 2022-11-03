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

package monitoring_schedule

import (
	"context"
	"errors"
	"fmt"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
)

var (
	modifyingStatuses = []string{
		"Pending",
	}

	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Deleting."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the monitoringScheduleStatus on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customSetOutput(
	r *resource,
	latestStatus *string,
) {
	svccommon.SetSyncedCondition(r, latestStatus, &resourceName, &modifyingStatuses)
}

// requeueUntilCanModify is a helper method to determine if monitoring schedule status allows modification
// Modifications to monitoring schedule are only allowed if:
//  1. The schedule is in a terminal state i.e. Status != Pending
//  2. There are no pending or in-progress jobs/executions
func (rm *resourceManager) requeueUntilCanModify(
	ctx context.Context,
	r *resource,
) error {
	scheduleStatus := r.ko.Status.MonitoringScheduleStatus
	if scheduleStatus == nil {
		return nil
	}

	errMsg := fmt.Sprintf("%s in Pending state cannot be modified or deleted.", resourceName)
	inProgressExecutions := false
	// It is possible that schedule doesn't have any execution yet
	if r.ko.Status.LastMonitoringExecutionSummary != nil {
		executionStatus := r.ko.Status.LastMonitoringExecutionSummary.MonitoringExecutionStatus
		if *executionStatus == "Pending" || *executionStatus == "InProgress" || *executionStatus == "Stopping" {
			inProgressExecutions = true
			errMsg = fmt.Sprintf("Monitoring Job in %s state, %s cannot be modified or deleted.", *executionStatus, resourceName)
		}
	}
	if *scheduleStatus == "Pending" || inProgressExecutions {
		return ackrequeue.NeededAfter(
			errors.New(errMsg),
			ackrequeue.DefaultRequeueAfterDuration)
	}

	return nil
}
