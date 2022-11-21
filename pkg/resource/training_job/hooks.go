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

package training_job

import (
	"errors"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svccommon "github.com/aws-controllers-k8s/sagemaker-controller/pkg/common"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	trainingJobModifyingStatuses = []string{
		svcsdk.TrainingJobStatusInProgress,
		svcsdk.TrainingJobStatusStopping,
	}
	ruleModifyingStatuses = []string{
		svcsdk.RuleEvaluationStatusInProgress,
		svcsdk.RuleEvaluationStatusStopping,
	}
	WarmPoolModifyingStatuses = []string{
		svcsdk.WarmPoolResourceStatusAvailable,
		svcsdk.WarmPoolResourceStatusInUse,
	}
	TrainingJobTerminalProfiler = []string{
		svcsdk.TrainingJobStatusCompleted,
		svcsdk.TrainingJobStatusFailed,
		svcsdk.TrainingJobStatusStopping,
		svcsdk.TrainingJobStatusStopped,
	}
	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Stopping."),
		ackrequeue.DefaultRequeueAfterDuration,
	)

	requeueBeforeUpdate = ackrequeue.NeededAfter(
		errors.New("Warm pool cannot be updated in InProgress state requeuing until TrainingJob reaches completed state."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueBeforeUpdateStarting = ackrequeue.NeededAfter(
		errors.New("Controller cannot update while secondary status is in Starting state."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customSetOutput sets the resource ResourceSynced condition to False if
// TrainingJob is being modified by AWS. It checks for debug and profiler rule status in addition to TrainingJobStatus
func (rm *resourceManager) customSetOutput(r *resource) {
	trainingJobStatus := r.ko.Status.TrainingJobStatus
	// early exit if training job is InProgress
	if trainingJobStatus != nil && *trainingJobStatus == svcsdk.TrainingJobStatusInProgress {
		svccommon.SetSyncedCondition(r, trainingJobStatus, &resourceName, &trainingJobModifyingStatuses)
		return
	}

	for _, rule := range r.ko.Status.DebugRuleEvaluationStatuses {
		if rule.RuleEvaluationStatus != nil && svccommon.IsModifyingStatus(rule.RuleEvaluationStatus, &ruleModifyingStatuses) {
			svccommon.SetSyncedCondition(r, rule.RuleEvaluationStatus, aws.String("DebugRule"), &ruleModifyingStatuses)
			return
		}
	}

	for _, rule := range r.ko.Status.ProfilerRuleEvaluationStatuses {
		if ackcompare.IsNotNil(r.ko.Status.ProfilingStatus) {
			// Sometimes rule evaluation status will stay in InProgress state.
			if *r.ko.Status.ProfilingStatus == "Disabled" {
				break
			}
		}
		if rule.RuleEvaluationStatus != nil && svccommon.IsModifyingStatus(rule.RuleEvaluationStatus, &ruleModifyingStatuses) {
			svccommon.SetSyncedCondition(r, rule.RuleEvaluationStatus, aws.String("ProfilerRule"), &ruleModifyingStatuses)
			return
		}
	}

	svccommon.SetSyncedCondition(r, trainingJobStatus, &resourceName, &trainingJobModifyingStatuses)

	warmpoolUsed := ackcompare.IsNotNil(r.ko.Spec.ResourceConfig) && ackcompare.IsNotNil(r.ko.Spec.ResourceConfig.KeepAlivePeriodInSeconds)

	// Only requeue when warm pool is being used and when training job is in the completed state.
	// WP will always have terminated status on error(Training Job or Warmpool).
	if ackcompare.IsNotNil(trainingJobStatus) && *trainingJobStatus == svcsdk.TrainingJobStatusCompleted &&
		warmpoolUsed {

		// Sometimes DescribeTrainingJob does not contain the warm pool status
		// In this condition the only possible status is Available or Terminated.
		if ackcompare.IsNotNil(trainingJobStatus) && ackcompare.IsNil(r.ko.Status.WarmPoolStatus) {
			svccommon.SetSyncedCondition(r, aws.String("Available"), aws.String("Warm Pool Infrastructure"), &WarmPoolModifyingStatuses)
		}

		if ackcompare.IsNotNil(r.ko.Status.WarmPoolStatus) && svccommon.IsModifyingStatus(r.ko.Status.WarmPoolStatus.Status, &WarmPoolModifyingStatuses) {
			svccommon.SetSyncedCondition(r, r.ko.Status.WarmPoolStatus.Status, aws.String("Warm Pool Infrastructure"), &WarmPoolModifyingStatuses)
		}
	}

}

// This function makes the controller requeue if there is an update and
// the training job is still in InProgress
func customSetOutputUpdateWarmpool(r *resource) error {
	trainingJobStatus := r.ko.Status.TrainingJobStatus
	if ackcompare.IsNotNil(trainingJobStatus) && *trainingJobStatus == svcsdk.TrainingJobStatusInProgress {
		return requeueBeforeUpdate
	}
	return nil
}

// Check if warm pool has reached a state where it is not updateable
func warmPoolTerminalCheck(latest *resource) bool {
	trainingJobStatus := latest.ko.Status.TrainingJobStatus
	if ackcompare.IsNotNil(latest.ko.Spec.ResourceConfig) {
		if ackcompare.IsNil(latest.ko.Spec.ResourceConfig.KeepAlivePeriodInSeconds) {
			return true // Warm pool can only be updated iff there is a provisioned cluster.
		}
	} else {
		return false
	}

	if ackcompare.IsNotNil(trainingJobStatus) {
		if *trainingJobStatus == svcsdk.TrainingJobStatusInProgress {
			return false
		}
		if *trainingJobStatus == svcsdk.TrainingJobStatusCompleted {
			if ackcompare.IsNotNil(latest.ko.Status.WarmPoolStatus) {
				wp_modifying := svccommon.IsModifyingStatus(latest.ko.Status.WarmPoolStatus.Status, &WarmPoolModifyingStatuses)
				return !wp_modifying
			} else {
				return false // Sometimes the API (briefly) does not return the WP status even if it completes.
			}
		} else {
			// Training Job is in 'Failed'|'Stopping'|'Stopped' (Terminal)
			return true
		}
	}

	// ACK OIDC is misconfigured (Terminal)
	return true
}

// Profiler cannot be updated at certain statuses.
func customSetOutputUpdateProfiler(r *resource) error {
	trainingSecondaryStatus := r.ko.Status.SecondaryStatus
	trainingJobStatus := r.ko.Status.TrainingJobStatus
	if ackcompare.IsNotNil(trainingSecondaryStatus) && *trainingSecondaryStatus == svcsdk.SecondaryStatusStarting {
		return requeueBeforeUpdateStarting
	}
	if ackcompare.IsNotNil(trainingJobStatus) {
		for _, terminalStatus := range TrainingJobTerminalProfiler {
			if terminalStatus == *trainingJobStatus {
				return errors.New("[ACK_SM] Profiler can only be updated when Training Job is in InProgress state")
			}
		}
	}
	return nil
}

// Checks if the profiler was removed.
func profilerRemovalCheck(desired *resource, latest *resource) bool {
	if ackcompare.IsNotNil(desired.ko.Spec) && ackcompare.IsNotNil(latest.ko.Spec) {
		if ackcompare.IsNil(desired.ko.Spec.ProfilerRuleConfigurations) && ackcompare.IsNotNil(latest.ko.Spec.ProfilerRuleConfigurations) {
			return true
		}
		if ackcompare.IsNil(desired.ko.Spec.ProfilerConfig) && ackcompare.IsNotNil(latest.ko.Spec.ProfilerConfig) {
			return true
		}
	}
	return false
}

// The statuses in ko object in the end of update are empty, using customSetOutput wont work.
func customSetOutputPostUpdate(ko *svcapitypes.TrainingJob, delta *ackcompare.Delta) {
	warmpool_diff := delta.DifferentAt("Spec.ResourceConfig.KeepAlivePeriodInSeconds")
	profiler_diff := delta.DifferentAt("Spec.ProfilerConfig") || delta.DifferentAt("Spec.ProfilerRuleConfigurations")
	if profiler_diff {
		svccommon.SetSyncedCondition(&resource{ko}, aws.string("InProgress"), &resourceName, &trainingJobModifyingStatuses)
	}
	if warmpool_diff {
		svccommon.SetSyncedCondition(&resource{ko}, aws.string("Availible"), aws.String("Warm Pool Infrastructure"), &WarmPoolModifyingStatuses)
	}

}
