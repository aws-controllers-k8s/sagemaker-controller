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

package processing_job

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// StoppingCondition is not a required field, so first create it
	if ackcompare.IsNil(a.ko.Spec.StoppingCondition) && ackcompare.IsNotNil(b.ko.Spec.StoppingCondition) {
		a.ko.Spec.StoppingCondition = &svcapitypes.ProcessingStoppingCondition{}
	}

	if ackcompare.IsNotNil(a.ko.Spec.StoppingCondition) && ackcompare.IsNotNil(b.ko.Spec.StoppingCondition) {
		if ackcompare.IsNil(a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds) && ackcompare.IsNotNil(b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds) {
			a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds = b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds
		}
	}

	// Default value of AppManaged is false
	// Default value of S3DataDistributionType is FullyReplicated
	if ackcompare.IsNotNil(a.ko.Spec.ProcessingInputs) && ackcompare.IsNotNil(b.ko.Spec.ProcessingInputs) {
		for index := range a.ko.Spec.ProcessingInputs {
			if ackcompare.IsNil(a.ko.Spec.ProcessingInputs[index].AppManaged) && ackcompare.IsNotNil(b.ko.Spec.ProcessingInputs[index].AppManaged) {
				a.ko.Spec.ProcessingInputs[index].AppManaged =
					b.ko.Spec.ProcessingInputs[index].AppManaged
			}
			if ackcompare.IsNil(a.ko.Spec.ProcessingInputs[index].S3Input.S3DataDistributionType) && ackcompare.IsNotNil(b.ko.Spec.ProcessingInputs[index].S3Input.S3DataDistributionType) {
				a.ko.Spec.ProcessingInputs[index].S3Input.S3DataDistributionType =
					b.ko.Spec.ProcessingInputs[index].S3Input.S3DataDistributionType
			}
		}
	}

	// Default value of AppManaged is false
	if ackcompare.IsNotNil(a.ko.Spec.ProcessingOutputConfig) && ackcompare.IsNotNil(b.ko.Spec.ProcessingOutputConfig) {
		if ackcompare.IsNotNil(a.ko.Spec.ProcessingOutputConfig.Outputs) && ackcompare.IsNotNil(b.ko.Spec.ProcessingOutputConfig.Outputs) {
			for index := range a.ko.Spec.ProcessingOutputConfig.Outputs {
				if ackcompare.IsNil(a.ko.Spec.ProcessingOutputConfig.Outputs[index].AppManaged) && ackcompare.IsNotNil(b.ko.Spec.ProcessingOutputConfig.Outputs[index].AppManaged) {
					a.ko.Spec.ProcessingOutputConfig.Outputs[index].AppManaged =
						b.ko.Spec.ProcessingOutputConfig.Outputs[index].AppManaged
				}
			}
		}
	}
}
