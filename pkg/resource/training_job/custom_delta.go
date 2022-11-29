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
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
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
			// Prevent out of bounds panics.
			if index == len(a.ko.Spec.ProfilerRuleConfigurations) || index == len(b.ko.Spec.ProfilerRuleConfigurations) {
				break
			}
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

	// inputDataConfig is an optional field. Its default value is an empty list
	if ackcompare.IsNil(a.ko.Spec.InputDataConfig) && ackcompare.IsNotNil(b.ko.Spec.InputDataConfig) {
		a.ko.Spec.InputDataConfig = []*svcapitypes.Channel{}
	}

	if ackcompare.IsNotNil(a.ko.Spec.InputDataConfig) && ackcompare.IsNotNil(b.ko.Spec.InputDataConfig) {
		for index := range a.ko.Spec.InputDataConfig {
			if ackcompare.IsNil(a.ko.Spec.InputDataConfig[index].RecordWrapperType) && ackcompare.IsNotNil(b.ko.Spec.InputDataConfig[index].RecordWrapperType) {
				a.ko.Spec.InputDataConfig[index].RecordWrapperType = defaultRecordWrapperType
			}
		}
	}
}

// SM returns profiler related objects even if the user disables the profiler
// customPostCompare detects if there is a diff
func customPostCompare(latest *resource, desired *resource, delta *ackcompare.Delta) {
	profilerConfigDiff := delta.DifferentAt("Spec.ProfilerConfig")
	profilerRuleDiff := delta.DifferentAt("Spec.ProfilerRuleConfigurations")
	if !profilerConfigDiff && !profilerRuleDiff {
		return
	}
	profilerStatus := latest.ko.Status.ProfilingStatus
	profilerDisabled := false

	if ackcompare.IsNotNil(profilerStatus) {
		//Do not remove profiler if user wants to enable it
		if *profilerStatus == "Disabled" && !userInitiatesProfilerCheck(desired) {
			profilerDisabled = true
		} else {
			return
		}
	} else {
		return
	}
	// TODO: Replace remove delta with an ack version when its natively supported
	if profilerConfigDiff && profilerDisabled {
		removeDelta(delta, "Spec.ProfilerConfig")
	}
	if profilerRuleDiff && profilerDisabled {
		removeDelta(delta, "Spec.ProfilerRuleConfigurations")
	}
}

// userInitiatesProfilerCheck checks if the user enabled/re enabled the profiler.
func userInitiatesProfilerCheck(desired *resource) bool {
	profilerConfigPresent := ackcompare.IsNotNil(desired.ko.Spec.ProfilerConfig)
	profilerRuleConfigPresent := ackcompare.IsNotNil(desired.ko.Spec.ProfilerRuleConfigurations)
	return profilerConfigPresent && profilerRuleConfigPresent
}

// removeDelta Removes fieldName from the delta slice.
// TODO: Replace when ack runtime can do this.
func removeDelta(delta *ackcompare.Delta, fieldName string) {
	differences := delta.Differences
	for index, diff := range differences {
		if diff.Path.Contains(fieldName) {
			differences = append(differences[:index], differences[index+1:]...)
			delta.Differences = differences
			return
		}
	}
}
