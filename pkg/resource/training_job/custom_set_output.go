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

// customDescribeTrainingJobSetOutput sets the resource ResourceSynced condition to False if
// TrainingJob is being modified by AWS. It has an additional check on the debugger status.
func (rm *resourceManager) customDescribeTrainingJobSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.DescribeTrainingJobOutput,
	ko *svcapitypes.TrainingJob,
) (*svcapitypes.TrainingJob, error) {
	trainingJobStatus := resp.TrainingJobStatus
	debuggerRuleInProgress := false
	if resp.DebugRuleEvaluationStatuses != nil {
		for _, rule := range resp.DebugRuleEvaluationStatuses {
			if rule.RuleEvaluationStatus != nil && *rule.RuleEvaluationStatus == svcsdk.RuleEvaluationStatusInProgress {
				debuggerRuleInProgress = true
				rm.customSetOutput(r, aws.String(svcsdk.TrainingJobStatusInProgress), ko)
				break
			}
		}
	}

	if !debuggerRuleInProgress {
		rm.customSetOutput(r, trainingJobStatus, ko)
	}

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

	syncConditionStatus := corev1.ConditionUnknown
	if *trainingJobStatus == svcsdk.TrainingJobStatusCompleted || *trainingJobStatus == svcsdk.TrainingJobStatusStopped || *trainingJobStatus == svcsdk.TrainingJobStatusFailed {
		syncConditionStatus = corev1.ConditionTrue
	} else {
		syncConditionStatus = corev1.ConditionFalse
	}

	var resourceSyncedCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			resourceSyncedCondition = condition
			break
		}
	}

	if resourceSyncedCondition == nil {
		resourceSyncedCondition = &ackv1alpha1.Condition{
			Type: ackv1alpha1.ConditionTypeResourceSynced,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
	}
	resourceSyncedCondition.Status = syncConditionStatus

}
