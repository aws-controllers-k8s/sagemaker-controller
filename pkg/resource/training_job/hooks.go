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
	resourceName = GroupKind.Kind

	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New(resourceName+" is Stopping."),
		ackrequeue.DefaultRequeueAfterDuration,
	)

	requeueWaitWhileWarmPoolInUse = ackrequeue.NeededAfter(
		errors.New("Warm Pool Cluster is still active."),
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
		if rule.RuleEvaluationStatus != nil && svccommon.IsModifyingStatus(rule.RuleEvaluationStatus, &ruleModifyingStatuses) {
			svccommon.SetSyncedCondition(r, rule.RuleEvaluationStatus, aws.String("ProfilerRule"), &ruleModifyingStatuses)
			return
		}
	}

	svccommon.SetSyncedCondition(r, trainingJobStatus, &resourceName, &trainingJobModifyingStatuses)
}

func (rm *resourceManager) customSetWarmPoolOutput(r *resource) error {
	if ackcompare.IsNil(r.ko.Status.WarmPoolStatus) {
		return nil
	}
	if svccommon.IsModifyingStatus(r.ko.Status.WarmPoolStatus.Status, &WarmPoolModifyingStatuses) {
		return requeueWaitWhileWarmPoolInUse
	}
	return nil
}
