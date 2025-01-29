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
	svcsdk "github.com/aws/aws-sdk-go-v2/service/sagemaker"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

// buildProfilerRuleConfigUpdateInput sets the input of the ProfilerRuleConfiguration so that
// it is compatible with the sagemaker API.
// Update training job is post operation wrt to the profiler parameters.
// Because of this only NEW rules can be specified.
// In this function we check to see if any new profiler configurstions have been added.
// Four cases:
// 1. Rule gets added (handled normally)
// 2. Rule gets removed (error is returned)
// 3. Rule gets removed but others get added (error is returned)
// 4. Rule gets changed (error gets returned)
// 5. One/more rule gets changed and one/more rules get added : successful update and error in the next reconcilation loop
func (rm *resourceManager) buildProfilerRuleConfigUpdateInput(desired *resource, latest *resource, input *svcsdk.UpdateTrainingJobInput) error {
	profilerRuleDesired := desired.ko.Spec.ProfilerRuleConfigurations
	profilerRuleLatest := latest.ko.Spec.ProfilerRuleConfigurations

	if ackcompare.IsNil(profilerRuleLatest) {
		return nil
	}
	if len(profilerRuleDesired) <= len(profilerRuleLatest) {
		return ackerr.NewTerminalError(errors.New("cannot remove/modify existing profiler rules."))
	}

	latestRules, err := rm.markNonUpdatableRules(profilerRuleDesired, profilerRuleLatest)
	if err != nil {
		return err
	}
	profilerRuleInput := []svcsdktypes.ProfilerRuleConfiguration{}

	for _, rule := range profilerRuleDesired {
		if ackcompare.IsNotNil(rule) && ackcompare.IsNotNil(rule.RuleConfigurationName) {
			_, present := latestRules[*rule.RuleConfigurationName]
			if !present {
				profilerRuleInput = append(profilerRuleInput, rm.convertProfileRuleType(rule))
			}
		}
	}
	// If the length of this slice is zero that only the contents of the profile rule have changed
	if len(profilerRuleInput) == 0 {
		return ackerr.NewTerminalError(errors.New("cannot modify an existing profiler rule."))
	}
	input.ProfilerRuleConfigurations = profilerRuleInput
	return nil
}

// markNonUpdatableRules returns a map containing the rules that are not eligible for update.
// In addition it returns an error if a rule gets removed.
func (rm *resourceManager) markNonUpdatableRules(profilerRuleDesired []*svcapitypes.ProfilerRuleConfiguration, profilerRuleLatest []*svcapitypes.ProfilerRuleConfiguration) (map[string]int, error) {
	latestRules := map[string]int{}
	for _, rule := range profilerRuleLatest {
		latestRules[*rule.RuleConfigurationName] = 0
	}
	// If a Rule Configuration is present in both latest and desired, set it to one.
	for _, rule := range profilerRuleDesired {
		_, present := latestRules[*rule.RuleConfigurationName]
		if present {
			latestRules[*rule.RuleConfigurationName] = 1
		}
	}
	// If a value in the map is equal to 0, the user must have removed the rule because
	// added rules would not be present in the map.
	for _, val := range latestRules {
		// This means that there exists a rule in latest that is not present in desired
		// which means that the input is invalid.
		if val == 0 {
			return nil, ackerr.NewTerminalError(errors.New("cannot remove an existing profiler rule"))
		}
	}

	return latestRules, nil
}

// handleProfilerRemoval sets the input parameters to disable the profiler.
func (rm *resourceManager) handleProfilerRemoval(input *svcsdk.UpdateTrainingJobInput) {
	input.ProfilerRuleConfigurations = nil
	profilerConfig := svcsdktypes.ProfilerConfigForUpdate{}
	disableProfilerCopy := true
	profilerConfig.DisableProfiler = &disableProfilerCopy
	input.ProfilerConfig = &profilerConfig
}

// convertProfileRuleType converts the kubernetes object ProfilerRuleConfiguration into
// a type that is compatible with the AWS API.
// Sagemaker and kubernetes types are not the same so the input has to be reconstructed.
func (rm *resourceManager) convertProfileRuleType(kubernetesObjectRule *svcapitypes.ProfilerRuleConfiguration) svcsdktypes.ProfilerRuleConfiguration {
	sagemakerAPIRule := svcsdktypes.ProfilerRuleConfiguration{}
	if kubernetesObjectRule.InstanceType != nil {
		sagemakerAPIRule.InstanceType = svcsdktypes.ProcessingInstanceType(*kubernetesObjectRule.InstanceType)
	}
	if kubernetesObjectRule.LocalPath != nil {
		sagemakerAPIRule.LocalPath = kubernetesObjectRule.LocalPath
	}
	if kubernetesObjectRule.RuleConfigurationName != nil {
		sagemakerAPIRule.RuleConfigurationName = kubernetesObjectRule.RuleConfigurationName
	}
	if kubernetesObjectRule.RuleEvaluatorImage != nil {
		sagemakerAPIRule.RuleEvaluatorImage = kubernetesObjectRule.RuleEvaluatorImage
	}
	if kubernetesObjectRule.RuleParameters != nil {
		f1elemf4 := map[string]string{}
		for f1elemf4key, f1elemf4valiter := range kubernetesObjectRule.RuleParameters {
			var f1elemf4val string
			f1elemf4val = *f1elemf4valiter
			f1elemf4[f1elemf4key] = f1elemf4val
		}
		sagemakerAPIRule.RuleParameters = f1elemf4
	}
	if kubernetesObjectRule.S3OutputPath != nil {
		sagemakerAPIRule.S3OutputPath = kubernetesObjectRule.S3OutputPath
	}
	if kubernetesObjectRule.VolumeSizeInGB != nil {
		volumeSizeInGBCopy0 := int32(*kubernetesObjectRule.VolumeSizeInGB)
		sagemakerAPIRule.VolumeSizeInGB = &volumeSizeInGBCopy0
	}
	return sagemakerAPIRule
}
