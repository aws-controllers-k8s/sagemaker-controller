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
	if ackcompare.IsNil(a.ko.Spec.StoppingCondition) && ackcompare.IsNotNil(b.ko.Spec.StoppingCondition) {
		a.ko.Spec.StoppingCondition = &svcapitypes.ProcessingStoppingCondition{}
	}

	if ackcompare.IsNil(a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds) && ackcompare.IsNotNil(b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds) {
		a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds = b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds
	}
}
