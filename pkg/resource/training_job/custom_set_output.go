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

// Use this file if the Status/Spec of the CR needs to be modified after
// create/describe/update operation

package training_job

import (
	"context"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

// customCreateTrainingJobSetOutput sets the resource in TempOutofSync if TrainingJob is
// in creating state. At this stage we know call to createTrainingJob was successful.
func (rm *resourceManager) customCreateTrainingJobSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.CreateTrainingJobOutput,
	ko *svcapitypes.TrainingJob,
) (*svcapitypes.TrainingJob, error) {
	rm.customSetOutput(r, aws.String(svcsdk.TrainingJobStatusInProgress), ko)
	return ko, nil
}

// customDescribeTrainingJobSetOutput sets the resource in TempOutofSync if
// TrainingJob is being modified by AWS. It has an additional check on the debugger status.
func (rm *resourceManager) customDescribeTrainingJobSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.DescribeTrainingJobOutput,
	ko *svcapitypes.TrainingJob,
) (*svcapitypes.TrainingJob, error) {
	trainingJobStatus := resp.TrainingJobStatus
	if resp.DebugRuleEvaluationStatuses != nil && resp.DebugRuleEvaluationStatuses[0].RuleEvaluationStatus != nil {
		debuggerStatus := resp.DebugRuleEvaluationStatuses[0].RuleEvaluationStatus
		rm.customSetOutput(r, debuggerStatus, ko)
	} else {
		rm.customSetOutput(r, trainingJobStatus, ko)
	}
	return ko, nil
}

// customStopTrainingJobSetOutput sets the resource in TempOutofSync if TrainingJob is
// in stopping state. At this stage we know call to deleteTrainingJob was successful.
func (rm *resourceManager) customStopTrainingJobSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.StopTrainingJobOutput,
	ko *svcapitypes.TrainingJob,
) (*svcapitypes.TrainingJob, error) {
	rm.customSetOutput(r, aws.String(svcsdk.TrainingJobStatusStopping), ko)
	return ko, nil
}

// customSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the trainingJobStatus on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customSetOutput(
	r *resource,
	trainingJobStatus *string,
	ko *svcapitypes.TrainingJob,
) {
	if trainingJobStatus == nil {
		return
	}

	// TODO: Re-check debugger statuses that shouldn't be requeued.
	syncConditionStatus := corev1.ConditionUnknown
	if *trainingJobStatus == svcsdk.TrainingJobStatusCompleted || *trainingJobStatus == svcsdk.TrainingJobStatusStopped || *trainingJobStatus == svcsdk.RuleEvaluationStatusNoIssuesFound {
		syncConditionStatus = corev1.ConditionTrue
	} else {
		syncConditionStatus = corev1.ConditionFalse
	}

	var resourceSyncedCondition *ackv1alpha1.Condition = nil
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	} else {
		for _, condition := range ko.Status.Conditions {
			if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
				resourceSyncedCondition = condition
				break
			}
		}
	}

	if resourceSyncedCondition == nil {
		resourceSyncedCondition = &ackv1alpha1.Condition{
			Type:   ackv1alpha1.ConditionTypeResourceSynced,
			Status: syncConditionStatus,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
	} else {
		resourceSyncedCondition.Status = syncConditionStatus
	}

}
