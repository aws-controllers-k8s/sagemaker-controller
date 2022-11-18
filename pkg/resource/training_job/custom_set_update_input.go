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

// Use this file if conditions need to be updated based on the latest status
// of training job which is not evident from API response

package training_job

import (
	"errors"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	smv1alpha "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// Three conditions:
// 1. Customer updates both profiler parameters: Recreate the input for profiler Rule.
// 2. Customer only updates Profiler Config: Set the profiler rule configuration to nil to avoid validation error.
// 3. Customer only updates Rule Configurations: Recreate the input for profiler Rule and set Profiler config to nil.
//	  safer to do this because the "only add" behavior might reappear.

func customSetUpdateInput(desired *resource, latest *resource, delta *ackcompare.Delta, input *svcsdk.UpdateTrainingJobInput) error {
	if delta.DifferentAt("Spec.ProfilerConfig") && delta.DifferentAt("Spec.ProfilerRuleConfigurations") {
		err := handleProfilerRuleConfig(desired, latest, input)
		return err
	}
	if delta.DifferentAt("Spec.ProfilerConfig") && !delta.DifferentAt("Spec.ProfilerRuleConfigurations") {
		input.SetProfilerRuleConfigurations(nil)
		return nil
	}
	if delta.DifferentAt("Spec.ProfilerRuleConfigurations") && !delta.DifferentAt("Spec.ProfilerConfig") {
		err := handleProfilerRuleConfig(desired, latest, input)
		input.SetProfilerConfig(nil) // SM still assumes the profiler config is the same.
		return err
	}
	return nil
}

// Update training job is post operation wrt to the profiler parameters.
// Because of this only NEW rules can be specified.
// In this function we check to see if any new profiler configurstions have been added.
func handleProfilerRuleConfig(desired *resource, latest *resource, input *svcsdk.UpdateTrainingJobInput) error {
	profilerRuleDesired := desired.ko.Spec.ProfilerRuleConfigurations
	profilerRuleLatest := latest.ko.Spec.ProfilerRuleConfigurations

	if ackcompare.IsNil(profilerRuleDesired) {
		return errors.New("[ACK_SM] Cannot remove a profiler rule.")
	}
	if ackcompare.IsNil(profilerRuleLatest) {
		return nil
	}
	if len(profilerRuleDesired) < len(profilerRuleLatest) {
		return errors.New("[ACK_SM] Cannot remove a profiler rule.")
	}

	ruleMap := map[string]int{}
	profilerRuleInput := []*svcsdk.ProfilerRuleConfiguration{}
	for _, rule := range profilerRuleLatest {
		if ackcompare.IsNotNil(rule) && ackcompare.IsNotNil(rule.RuleConfigurationName) {
			ruleMap[*rule.RuleConfigurationName] = 1
		}
	}
	for _, rule := range profilerRuleDesired {
		if ackcompare.IsNotNil(rule) && ackcompare.IsNotNil(rule.RuleConfigurationName) {
			_, present := ruleMap[*rule.RuleConfigurationName]
			if !present {
				profilerRuleInput = append(profilerRuleInput, convertProfileRuleType(rule))
			}
		}
	}
	input.SetProfilerRuleConfigurations(profilerRuleInput)
	return nil
}

// Recreates input and sets disable profiler to true
func handleProfilerRemoval(input *svcsdk.UpdateTrainingJobInput) {
	input.SetProfilerRuleConfigurations(nil)
	profilerConfig := svcsdk.ProfilerConfigForUpdate{}
	profilerConfig.SetDisableProfiler(true)
	input.SetProfilerConfig(&profilerConfig)
}

// Sagemaker and kubernetes types are not the same so the input has to be reconstructed.
func convertProfileRuleType(rule *smv1alpha.ProfilerRuleConfiguration) *svcsdk.ProfilerRuleConfiguration {
	smRule := &svcsdk.ProfilerRuleConfiguration{}
	if rule.InstanceType != nil {
		smRule.SetInstanceType(*rule.InstanceType)
	}
	if rule.LocalPath != nil {
		smRule.SetLocalPath(*rule.LocalPath)
	}
	if rule.RuleConfigurationName != nil {
		smRule.SetRuleConfigurationName(*rule.RuleConfigurationName)
	}
	if rule.RuleEvaluatorImage != nil {
		smRule.SetRuleEvaluatorImage(*rule.RuleEvaluatorImage)
	}
	if rule.RuleParameters != nil {
		f1elemf4 := map[string]*string{}
		for f1elemf4key, f1elemf4valiter := range rule.RuleParameters {
			var f1elemf4val string
			f1elemf4val = *f1elemf4valiter
			f1elemf4[f1elemf4key] = &f1elemf4val
		}
		smRule.SetRuleParameters(f1elemf4)
	}
	if rule.S3OutputPath != nil {
		smRule.SetS3OutputPath(*rule.S3OutputPath)
	}
	if rule.VolumeSizeInGB != nil {
		smRule.SetVolumeSizeInGB(*rule.VolumeSizeInGB)
	}
	return smRule
}
