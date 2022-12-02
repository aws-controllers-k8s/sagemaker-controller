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
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
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
		errors.New("warm pool cannot be updated in InProgress state, requeuing until TrainingJob reaches completed state."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueBeforeUpdateStarting = ackrequeue.NeededAfter(
		errors.New("training job cannot be updated while secondary status is in Starting state."),
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
		if ackcompare.IsNotNil(r.ko.Status.ProfilingStatus) && *r.ko.Status.ProfilingStatus == "Disabled" {
			// Sometimes rule evaluation status will stay in InProgress state.
			break
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

// isWarmPoolUpdateable returns a requeue or terminal error depending on the warmpool/training job state
func (rm *resourceManager) isWarmPoolUpdatable(latest *resource) error {
	trainingJobStatus := latest.ko.Status.TrainingJobStatus
	if ackcompare.IsNil(latest.ko.Spec.ResourceConfig.KeepAlivePeriodInSeconds) {
		return ackerr.NewTerminalError(errors.New("warm pool does not exist"))
	}
	if ackcompare.IsNotNil(trainingJobStatus) {
		if *trainingJobStatus == svcsdk.TrainingJobStatusInProgress {
			return requeueBeforeUpdate
		}
		if *trainingJobStatus == svcsdk.TrainingJobStatusCompleted {
			if ackcompare.IsNotNil(latest.ko.Status.WarmPoolStatus) {
				wp_modifying := svccommon.IsModifyingStatus(latest.ko.Status.WarmPoolStatus.Status, &WarmPoolModifyingStatuses)
				if wp_modifying {
					return nil
				} else {
					return ackerr.NewTerminalError(errors.New("warm pool is in a non updateable state"))
				}
			} else {
				// Sometimes the API (briefly) does not return the WP status even if it completes.
				// This only occurs for a short time after training job has reached Completed state.
				return requeueBeforeUpdate
			}
		} else {
			// Training Job is in 'Failed'|'Stopping'|'Stopped' (Terminal)
			return ackerr.NewTerminalError(errors.New("warm pool is in a non updateable state"))
		}

	}
	return nil

}

// customSetOutputUpdateProfiler decides whether the training job is ready/eligible for update
// depending on the status.
func (rm *resourceManager) customSetOutputUpdateProfiler(r *resource) error {
	trainingSecondaryStatus := r.ko.Status.SecondaryStatus
	trainingJobStatus := r.ko.Status.TrainingJobStatus
	if ackcompare.IsNotNil(trainingSecondaryStatus) && *trainingSecondaryStatus == svcsdk.SecondaryStatusStarting {
		return requeueBeforeUpdateStarting
	}
	if ackcompare.IsNotNil(trainingJobStatus) {
		for _, terminalStatus := range TrainingJobTerminalProfiler {
			if terminalStatus == *trainingJobStatus {
				return ackerr.NewTerminalError(errors.New("profiler can only be updated when Training Job is in InProgress state"))
			}
		}
	}
	return nil
}

// isProfilerRemoved checks if the profiler was removed.
func (rm *resourceManager) isProfilerRemoved(desired *resource, latest *resource) bool {
	if ackcompare.IsNil(desired.ko.Spec.ProfilerRuleConfigurations) && ackcompare.IsNotNil(latest.ko.Spec.ProfilerRuleConfigurations) {
		return true
	}
	if ackcompare.IsNil(desired.ko.Spec.ProfilerConfig) && ackcompare.IsNotNil(latest.ko.Spec.ProfilerConfig) {
		return true
	}
	return false
}

// customSetUpdateInput modifies the input of UpdateTrainingJob.
// Three conditions:
// 1. Customer updates both profiler parameters: Recreate the input for profiler Rule.
// 2. Customer only updates Profiler Config: Set the profiler rule configuration to nil to avoid validation error.
// 3. Customer only updates Rule Configurations: Recreate the input for profiler Rule and set Profiler config to nil.
//	  safer to do this because the "only add" behavior might reappear.
func (rm *resourceManager) customSetUpdateInput(desired *resource, latest *resource, delta *ackcompare.Delta, input *svcsdk.UpdateTrainingJobInput) error {
	if !delta.DifferentAt("Spec.ProfilerConfig") {
		input.SetProfilerConfig(nil)
	}
	if !delta.DifferentAt("Spec.ProfilerRuleConfigurations") {
		input.SetProfilerRuleConfigurations(nil)
	} else {
		err := buildProfilerRuleConfigUpdateInput(desired, latest, input)
		return err
	}

	return nil
}
