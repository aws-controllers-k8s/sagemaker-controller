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

package labeling_job

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// StoppingConditions is not a required field, so first create it
	if ackcompare.IsNil(a.ko.Spec.StoppingConditions) && ackcompare.IsNotNil(b.ko.Spec.StoppingConditions) {
		a.ko.Spec.StoppingConditions = &svcapitypes.LabelingJobStoppingConditions{}
	}

	if ackcompare.IsNotNil(a.ko.Spec.StoppingConditions) && ackcompare.IsNotNil(b.ko.Spec.StoppingConditions) {
		if ackcompare.IsNil(a.ko.Spec.StoppingConditions.MaxHumanLabeledObjectCount) && ackcompare.IsNotNil(b.ko.Spec.StoppingConditions.MaxHumanLabeledObjectCount) {
			a.ko.Spec.StoppingConditions.MaxHumanLabeledObjectCount = b.ko.Spec.StoppingConditions.MaxHumanLabeledObjectCount
		}
		if ackcompare.IsNil(a.ko.Spec.StoppingConditions.MaxPercentageOfInputDatasetLabeled) && ackcompare.IsNotNil(b.ko.Spec.StoppingConditions.MaxPercentageOfInputDatasetLabeled) {
			a.ko.Spec.StoppingConditions.MaxPercentageOfInputDatasetLabeled = b.ko.Spec.StoppingConditions.MaxPercentageOfInputDatasetLabeled
		}
	}
}
