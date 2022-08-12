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
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go/aws"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {

	// Default value of VolumeSizeInGB is 0
	defaultVolumeSizeInGB := aws.Int64(0)

	if ackcompare.IsNotNil(a.ko.Spec.ProfilerRuleConfigurations) && ackcompare.IsNotNil(b.ko.Spec.ProfilerRuleConfigurations) {
		for index := range a.ko.Spec.ProfilerRuleConfigurations {
			if ackcompare.IsNil(a.ko.Spec.ProfilerRuleConfigurations[index].VolumeSizeInGB) && ackcompare.IsNotNil(b.ko.Spec.ProfilerRuleConfigurations[index].VolumeSizeInGB) {
				a.ko.Spec.ProfilerRuleConfigurations[index].VolumeSizeInGB = defaultVolumeSizeInGB
			}
		}
	}

	// Default value of VolumeSizeInGB is 0
	if ackcompare.IsNotNil(a.ko.Spec.DebugRuleConfigurations) && ackcompare.IsNotNil(b.ko.Spec.DebugRuleConfigurations) {
		for index := range a.ko.Spec.DebugRuleConfigurations {
			if ackcompare.IsNil(a.ko.Spec.DebugRuleConfigurations[index].VolumeSizeInGB) && ackcompare.IsNotNil(b.ko.Spec.DebugRuleConfigurations[index].VolumeSizeInGB) {
				a.ko.Spec.DebugRuleConfigurations[index].VolumeSizeInGB = defaultVolumeSizeInGB
			}
		}
	}

	// Default value of RecordWrapperType is None
	defaultRecordWrapperType := aws.String("None")

	if ackcompare.IsNotNil(a.ko.Spec.InputDataConfig) && ackcompare.IsNotNil(b.ko.Spec.InputDataConfig) {
		for index := range a.ko.Spec.InputDataConfig {
			if ackcompare.IsNil(a.ko.Spec.InputDataConfig[index].RecordWrapperType) && ackcompare.IsNotNil(b.ko.Spec.InputDataConfig[index].RecordWrapperType) {
				a.ko.Spec.InputDataConfig[index].RecordWrapperType = defaultRecordWrapperType
			}
		}
	}
}
