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
// of resource which is not evident from API response

package monitoring_schedule

import (
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// CustomUpdateConditions sets conditions (terminal) on supplied monitoring schedule.
// it examines supplied resource to determine conditions.
// It returns true if conditions are updated.
func (rm *resourceManager) CustomUpdateConditions(
	ko *svcapitypes.MonitoringSchedule,
	r *resource,
	err error,
) bool {
	latestStatus := r.ko.Status.MonitoringScheduleStatus
	terminalStatus := svcsdk.ScheduleStatusFailed
	conditionManager := &resource{ko}
	resourceName := resourceGK.Kind
	// If the latestStatus == terminalStatus we will set
	// the terminal condition and terminal message.
	return svccommon.SetTerminalState(conditionManager, latestStatus, &resourceName, terminalStatus)
}
