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
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/sagemaker"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

var (
	trainingJobModifyingStatuses = []string{
		string(svcsdktypes.TrainingJobStatusInProgress),
		string(svcsdktypes.TrainingJobStatusStopping),
	}
	ruleModifyingStatuses = []string{
		string(svcsdktypes.RuleEvaluationStatusInProgress),
		string(svcsdktypes.RuleEvaluationStatusStopping),
	}
	WarmPoolModifyingStatuses = []string{
		string(svcsdktypes.WarmPoolResourceStatusAvailable),
		string(svcsdktypes.WarmPoolResourceStatusInuse),
	}
	TrainingJobTerminalProfiler = []string{
		string(svcsdktypes.TrainingJobStatusCompleted),
		string(svcsdktypes.TrainingJobStatusFailed),
		string(svcsdktypes.TrainingJobStatusStopping),
		string(svcsdktypes.TrainingJobStatusStopped),
	}
	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Stopping."),
		ackrequeue.DefaultRequeueAfterDuration,
	)

	requeueBeforeUpdate = ackrequeue.NeededAfter(
		errors.New("warm pool cannot be updated while TrainingJob status is InProgress, requeuing until TrainingJob completes."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// customSetOutput sets the resource ResourceSynced condition to False if
// TrainingJob is being modified by AWS. It checks for debug and profiler rule status in addition to TrainingJobStatus
func (rm *resourceManager) customSetOutput(r *resource) {
	trainingJobStatus := r.ko.Status.TrainingJobStatus
	// early exit if training job is InProgress
	if trainingJobStatus != nil && *trainingJobStatus == string(svcsdktypes.TrainingJobStatusInProgress) {
		svccommon.SetSyncedCondition(r, trainingJobStatus, &resourceName, &trainingJobModifyingStatuses)
		return
	}

	for _, rule := range r.ko.Status.DebugRuleEvaluationStatuses {
		if rule.RuleEvaluationStatus != nil && svccommon.IsModifyingStatus(rule.RuleEvaluationStatus, &ruleModifyingStatuses) {
			svccommon.SetSyncedCondition(r, rule.RuleEvaluationStatus, aws.String("DebugRule"), &ruleModifyingStatuses)
			return
		}
	}

	// Sometimes rule evaluation status will stay in InProgress state.
	if ackcompare.IsNotNil(r.ko.Status.ProfilingStatus) && *r.ko.Status.ProfilingStatus != "Disabled" {
		for _, rule := range r.ko.Status.ProfilerRuleEvaluationStatuses {
			if rule.RuleEvaluationStatus != nil && svccommon.IsModifyingStatus(rule.RuleEvaluationStatus, &ruleModifyingStatuses) {
				svccommon.SetSyncedCondition(r, rule.RuleEvaluationStatus, aws.String("ProfilerRule"), &ruleModifyingStatuses)
				return
			}
		}
	}

	svccommon.SetSyncedCondition(r, trainingJobStatus, &resourceName, &trainingJobModifyingStatuses)

}

// isWarmPoolUpdateable returns a requeue or terminal error depending on the warmpool/training job state
func (rm *resourceManager) isWarmPoolUpdatable(latest *resource) error {
	trainingJobStatus := latest.ko.Status.TrainingJobStatus
	if ackcompare.IsNil(latest.ko.Spec.ResourceConfig.KeepAlivePeriodInSeconds) {
		return ackerr.NewTerminalError(errors.New("warm pool does not exist and can only be configured at creation time"))
	}
	if ackcompare.IsNotNil(trainingJobStatus) {
		if *trainingJobStatus == string(svcsdktypes.TrainingJobStatusInProgress) {
			return requeueBeforeUpdate
		}
		if *trainingJobStatus == string(svcsdktypes.TrainingJobStatusCompleted) {
			if ackcompare.IsNotNil(latest.ko.Status.WarmPoolStatus) {
				if svccommon.IsModifyingStatus(latest.ko.Status.WarmPoolStatus.Status, &WarmPoolModifyingStatuses) {
					return nil
				} else {
					return ackerr.NewTerminalError(errors.New("warm pool cannot be updated if has been terminated or reused"))
				}
			} else {
				// Sometimes the API (briefly) does not return the WP status even if it completes.
				// This only occurs for a short time after training job has reached Completed state.
				return requeueBeforeUpdate
			}
		} else {
			// Training Job is in 'Failed'|'Stopping'|'Stopped' (Terminal)
			return ackerr.NewTerminalError(errors.New("warm pool can only be updated if TrainingJob status is Completed. Warm pool will be terminated automatically if trainingjob has not completed successfully"))
		}

	}
	return nil

}

// isProfilerUpdatable decides whether the training job is ready/eligible for update
// depending on the status.
func (rm *resourceManager) isProfilerUpdatable(r *resource) error {
	trainingJobStatus := r.ko.Status.TrainingJobStatus
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
// The profiler gets removed when ProfilerConfig or ProfilerRuleConfig (or both) are not present in the spec but were present before.
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
//  1. Customer updates both profiler parameters: Recreate the input for profiler Rule.
//  2. Customer only updates Profiler Config: Set the profiler rule configuration to nil to avoid validation error.
//  3. Customer only updates Rule Configurations: Recreate the input for profiler Rule and set Profiler config to nil.
//     safer to do this because the "only add" behavior might reappear.
func (rm *resourceManager) customSetUpdateInput(desired *resource, latest *resource, delta *ackcompare.Delta, input *svcsdk.UpdateTrainingJobInput) error {
	if !delta.DifferentAt("Spec.ProfilerConfig") {
		input.ProfilerConfig = nil
	}
	if !delta.DifferentAt("Spec.ProfilerRuleConfigurations") {
		input.ProfilerRuleConfigurations = nil
	} else {
		err := rm.buildProfilerRuleConfigUpdateInput(desired, latest, input)
		return err
	}

	return nil
}
