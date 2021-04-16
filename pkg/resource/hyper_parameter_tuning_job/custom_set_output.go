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

package hyper_parameter_tuning_job

import (
	"context"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
)

// customCreateHyperParameterTuningJobSetOutput sets the resource in TempOutofSync if HyperParameterTuningJob is
// in creating state. At this stage we know call to createHyperParameterTuningJob was successful.
func (rm *resourceManager) customCreateHyperParameterTuningJobSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.CreateHyperParameterTuningJobOutput,
	ko *svcapitypes.HyperParameterTuningJob,
) (*svcapitypes.HyperParameterTuningJob, error) {
	rm.customSetOutput(r, aws.String(svcsdk.HyperParameterTuningJobStatusInProgress), ko)
	return ko, nil
}

// customDescribeHyperParameterTuningJobSetOutput sets the resource in TempOutofSync if
// HyperParameterTuningJob is being modified by AWS.
func (rm *resourceManager) customDescribeHyperParameterTuningJobSetOutput(
	ctx context.Context,
	r *resource,
	resp *svcsdk.DescribeHyperParameterTuningJobOutput,
	ko *svcapitypes.HyperParameterTuningJob,
) (*svcapitypes.HyperParameterTuningJob, error) {
	rm.customSetOutput(r, resp.HyperParameterTuningJobStatus, ko)
	return ko, nil
}

// customSetOutput sets ConditionTypeResourceSynced condition to True or False
// based on the hyperParameterTuningJobStatus on AWS so the reconciler can determine if a
// requeue is needed
func (rm *resourceManager) customSetOutput(
	r *resource,
	hyperParameterTuningJobStatus *string,
	ko *svcapitypes.HyperParameterTuningJob,
) {
	if hyperParameterTuningJobStatus == nil {
		return
	}

	syncConditionStatus := corev1.ConditionUnknown
	if *hyperParameterTuningJobStatus == svcsdk.HyperParameterTuningJobStatusCompleted || *hyperParameterTuningJobStatus == svcsdk.HyperParameterTuningJobStatusStopped || *hyperParameterTuningJobStatus == svcsdk.HyperParameterTuningJobStatusFailed {
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
			Type: ackv1alpha1.ConditionTypeResourceSynced,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
	}
	resourceSyncedCondition.Status = syncConditionStatus

}
