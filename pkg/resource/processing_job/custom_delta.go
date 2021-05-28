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
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	if a.ko.Spec.StoppingCondition == nil && b.ko.Spec.StoppingCondition != nil {
		a.ko.Spec.StoppingCondition = &svcapitypes.ProcessingStoppingCondition{}
	}

	if a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds == nil && b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds != nil{
		a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds = b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds
	}
}
