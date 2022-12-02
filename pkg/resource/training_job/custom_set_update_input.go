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
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// buildProfilerRuleConfigUpdateInput sets the input of the ProfilerRuleConfiguration so that
// it is compatible with the sagemaker API.
// Update training job is post operation wrt to the profiler parameters.
// Because of this only NEW rules can be specified.
// In this function we check to see if any new profiler configurstions have been added.
func buildProfilerRuleConfigUpdateInput(desired *resource, latest *resource, input *svcsdk.UpdateTrainingJobInput) error {
	profilerRuleDesired := desired.ko.Spec.ProfilerRuleConfigurations
	profilerRuleLatest := latest.ko.Spec.ProfilerRuleConfigurations

	if ackcompare.IsNil(profilerRuleLatest) {
		return nil
	}
	if len(profilerRuleDesired) < len(profilerRuleLatest) {
		return ackerr.NewTerminalError(errors.New("cannot remove a profiler rule."))
	}

	ruleMap := map[string]int{}
	profilerRuleInput := []*svcsdk.ProfilerRuleConfiguration{}
	for _, rule := range profilerRuleLatest {
		if ackcompare.IsNotNil(rule.RuleConfigurationName) {
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

// handleProfilerRemoval sets the input parameters to disable the profiler.
func handleProfilerRemoval(input *svcsdk.UpdateTrainingJobInput) {
	input.SetProfilerRuleConfigurations(nil)
	profilerConfig := svcsdk.ProfilerConfigForUpdate{}
	profilerConfig.SetDisableProfiler(true)
	input.SetProfilerConfig(&profilerConfig)
}

// convertProfileRuleType converts the kubernetes object ProfilerRuleConfiguration into
// a type that is compatible with the AWS API.
// Sagemaker and kubernetes types are not the same so the input has to be reconstructed.
func convertProfileRuleType(rule *svcapitypes.ProfilerRuleConfiguration) *svcsdk.ProfilerRuleConfiguration {
	rule := &svcsdk.ProfilerRuleConfiguration{}
	if rule.InstanceType != nil {
		rule.SetInstanceType(*rule.InstanceType)
	}
	if rule.LocalPath != nil {
		rule.SetLocalPath(*rule.LocalPath)
	}
	if rule.RuleConfigurationName != nil {
		rule.SetRuleConfigurationName(*rule.RuleConfigurationName)
	}
	if rule.RuleEvaluatorImage != nil {
		rule.SetRuleEvaluatorImage(*rule.RuleEvaluatorImage)
	}
	if rule.RuleParameters != nil {
		f1elemf4 := map[string]*string{}
		for f1elemf4key, f1elemf4valiter := range rule.RuleParameters {
			var f1elemf4val string
			f1elemf4val = *f1elemf4valiter
			f1elemf4[f1elemf4key] = &f1elemf4val
		}
		rule.SetRuleParameters(f1elemf4)
	}
	if rule.S3OutputPath != nil {
		rule.SetS3OutputPath(*rule.S3OutputPath)
	}
	if rule.VolumeSizeInGB != nil {
		rule.SetVolumeSizeInGB(*rule.VolumeSizeInGB)
	}
	return rule
}
