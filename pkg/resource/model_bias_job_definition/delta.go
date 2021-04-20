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

// Code generated by ack-generate. DO NOT EDIT.

package model_bias_job_definition

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}

	if ackcompare.HasNilDifference(a.ko.Spec.JobDefinitionName, b.ko.Spec.JobDefinitionName) {
		delta.Add("Spec.JobDefinitionName", a.ko.Spec.JobDefinitionName, b.ko.Spec.JobDefinitionName)
	} else if a.ko.Spec.JobDefinitionName != nil && b.ko.Spec.JobDefinitionName != nil {
		if *a.ko.Spec.JobDefinitionName != *b.ko.Spec.JobDefinitionName {
			delta.Add("Spec.JobDefinitionName", a.ko.Spec.JobDefinitionName, b.ko.Spec.JobDefinitionName)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.JobResources, b.ko.Spec.JobResources) {
		delta.Add("Spec.JobResources", a.ko.Spec.JobResources, b.ko.Spec.JobResources)
	} else if a.ko.Spec.JobResources != nil && b.ko.Spec.JobResources != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.JobResources.ClusterConfig, b.ko.Spec.JobResources.ClusterConfig) {
			delta.Add("Spec.JobResources.ClusterConfig", a.ko.Spec.JobResources.ClusterConfig, b.ko.Spec.JobResources.ClusterConfig)
		} else if a.ko.Spec.JobResources.ClusterConfig != nil && b.ko.Spec.JobResources.ClusterConfig != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.JobResources.ClusterConfig.InstanceCount, b.ko.Spec.JobResources.ClusterConfig.InstanceCount) {
				delta.Add("Spec.JobResources.ClusterConfig.InstanceCount", a.ko.Spec.JobResources.ClusterConfig.InstanceCount, b.ko.Spec.JobResources.ClusterConfig.InstanceCount)
			} else if a.ko.Spec.JobResources.ClusterConfig.InstanceCount != nil && b.ko.Spec.JobResources.ClusterConfig.InstanceCount != nil {
				if *a.ko.Spec.JobResources.ClusterConfig.InstanceCount != *b.ko.Spec.JobResources.ClusterConfig.InstanceCount {
					delta.Add("Spec.JobResources.ClusterConfig.InstanceCount", a.ko.Spec.JobResources.ClusterConfig.InstanceCount, b.ko.Spec.JobResources.ClusterConfig.InstanceCount)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.JobResources.ClusterConfig.InstanceType, b.ko.Spec.JobResources.ClusterConfig.InstanceType) {
				delta.Add("Spec.JobResources.ClusterConfig.InstanceType", a.ko.Spec.JobResources.ClusterConfig.InstanceType, b.ko.Spec.JobResources.ClusterConfig.InstanceType)
			} else if a.ko.Spec.JobResources.ClusterConfig.InstanceType != nil && b.ko.Spec.JobResources.ClusterConfig.InstanceType != nil {
				if *a.ko.Spec.JobResources.ClusterConfig.InstanceType != *b.ko.Spec.JobResources.ClusterConfig.InstanceType {
					delta.Add("Spec.JobResources.ClusterConfig.InstanceType", a.ko.Spec.JobResources.ClusterConfig.InstanceType, b.ko.Spec.JobResources.ClusterConfig.InstanceType)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID, b.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID) {
				delta.Add("Spec.JobResources.ClusterConfig.VolumeKMSKeyID", a.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID, b.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID)
			} else if a.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID != nil && b.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID != nil {
				if *a.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID != *b.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID {
					delta.Add("Spec.JobResources.ClusterConfig.VolumeKMSKeyID", a.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID, b.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB, b.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB) {
				delta.Add("Spec.JobResources.ClusterConfig.VolumeSizeInGB", a.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB, b.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB)
			} else if a.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB != nil && b.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB != nil {
				if *a.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB != *b.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB {
					delta.Add("Spec.JobResources.ClusterConfig.VolumeSizeInGB", a.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB, b.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasAppSpecification, b.ko.Spec.ModelBiasAppSpecification) {
		delta.Add("Spec.ModelBiasAppSpecification", a.ko.Spec.ModelBiasAppSpecification, b.ko.Spec.ModelBiasAppSpecification)
	} else if a.ko.Spec.ModelBiasAppSpecification != nil && b.ko.Spec.ModelBiasAppSpecification != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasAppSpecification.ConfigURI, b.ko.Spec.ModelBiasAppSpecification.ConfigURI) {
			delta.Add("Spec.ModelBiasAppSpecification.ConfigURI", a.ko.Spec.ModelBiasAppSpecification.ConfigURI, b.ko.Spec.ModelBiasAppSpecification.ConfigURI)
		} else if a.ko.Spec.ModelBiasAppSpecification.ConfigURI != nil && b.ko.Spec.ModelBiasAppSpecification.ConfigURI != nil {
			if *a.ko.Spec.ModelBiasAppSpecification.ConfigURI != *b.ko.Spec.ModelBiasAppSpecification.ConfigURI {
				delta.Add("Spec.ModelBiasAppSpecification.ConfigURI", a.ko.Spec.ModelBiasAppSpecification.ConfigURI, b.ko.Spec.ModelBiasAppSpecification.ConfigURI)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasAppSpecification.Environment, b.ko.Spec.ModelBiasAppSpecification.Environment) {
			delta.Add("Spec.ModelBiasAppSpecification.Environment", a.ko.Spec.ModelBiasAppSpecification.Environment, b.ko.Spec.ModelBiasAppSpecification.Environment)
		} else if a.ko.Spec.ModelBiasAppSpecification.Environment != nil && b.ko.Spec.ModelBiasAppSpecification.Environment != nil {
			if !ackcompare.MapStringStringPEqual(a.ko.Spec.ModelBiasAppSpecification.Environment, b.ko.Spec.ModelBiasAppSpecification.Environment) {
				delta.Add("Spec.ModelBiasAppSpecification.Environment", a.ko.Spec.ModelBiasAppSpecification.Environment, b.ko.Spec.ModelBiasAppSpecification.Environment)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasAppSpecification.ImageURI, b.ko.Spec.ModelBiasAppSpecification.ImageURI) {
			delta.Add("Spec.ModelBiasAppSpecification.ImageURI", a.ko.Spec.ModelBiasAppSpecification.ImageURI, b.ko.Spec.ModelBiasAppSpecification.ImageURI)
		} else if a.ko.Spec.ModelBiasAppSpecification.ImageURI != nil && b.ko.Spec.ModelBiasAppSpecification.ImageURI != nil {
			if *a.ko.Spec.ModelBiasAppSpecification.ImageURI != *b.ko.Spec.ModelBiasAppSpecification.ImageURI {
				delta.Add("Spec.ModelBiasAppSpecification.ImageURI", a.ko.Spec.ModelBiasAppSpecification.ImageURI, b.ko.Spec.ModelBiasAppSpecification.ImageURI)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasBaselineConfig, b.ko.Spec.ModelBiasBaselineConfig) {
		delta.Add("Spec.ModelBiasBaselineConfig", a.ko.Spec.ModelBiasBaselineConfig, b.ko.Spec.ModelBiasBaselineConfig)
	} else if a.ko.Spec.ModelBiasBaselineConfig != nil && b.ko.Spec.ModelBiasBaselineConfig != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName, b.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName) {
			delta.Add("Spec.ModelBiasBaselineConfig.BaseliningJobName", a.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName, b.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName)
		} else if a.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName != nil && b.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName != nil {
			if *a.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName != *b.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName {
				delta.Add("Spec.ModelBiasBaselineConfig.BaseliningJobName", a.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName, b.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource, b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource) {
			delta.Add("Spec.ModelBiasBaselineConfig.ConstraintsResource", a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource, b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource)
		} else if a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource != nil && b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI, b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI) {
				delta.Add("Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI", a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI, b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI)
			} else if a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI != nil && b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI != nil {
				if *a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI != *b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI {
					delta.Add("Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI", a.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI, b.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput, b.ko.Spec.ModelBiasJobInput) {
		delta.Add("Spec.ModelBiasJobInput", a.ko.Spec.ModelBiasJobInput, b.ko.Spec.ModelBiasJobInput)
	} else if a.ko.Spec.ModelBiasJobInput != nil && b.ko.Spec.ModelBiasJobInput != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput, b.ko.Spec.ModelBiasJobInput.EndpointInput) {
			delta.Add("Spec.ModelBiasJobInput.EndpointInput", a.ko.Spec.ModelBiasJobInput.EndpointInput, b.ko.Spec.ModelBiasJobInput.EndpointInput)
		} else if a.ko.Spec.ModelBiasJobInput.EndpointInput != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset, b.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset", a.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset, b.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset != *b.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset", a.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset, b.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName, b.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.EndpointName", a.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName, b.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName != *b.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.EndpointName", a.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName, b.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute != *b.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute != *b.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath, b.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.LocalPath", a.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath, b.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath != *b.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.LocalPath", a.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath, b.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute != *b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute != *b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute", a.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute, b.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType, b.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType", a.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType, b.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType != *b.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType", a.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType, b.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode, b.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.S3InputMode", a.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode, b.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode != *b.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.S3InputMode", a.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode, b.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset, b.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset) {
				delta.Add("Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset", a.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset, b.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset)
			} else if a.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset != nil && b.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset != nil {
				if *a.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset != *b.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset {
					delta.Add("Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset", a.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset, b.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset)
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input, b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input) {
			delta.Add("Spec.ModelBiasJobInput.GroundTruthS3Input", a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input, b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input)
		} else if a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input != nil && b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI, b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI) {
				delta.Add("Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI", a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI, b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI)
			} else if a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI != nil && b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI != nil {
				if *a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI != *b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI {
					delta.Add("Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI", a.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI, b.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobOutputConfig, b.ko.Spec.ModelBiasJobOutputConfig) {
		delta.Add("Spec.ModelBiasJobOutputConfig", a.ko.Spec.ModelBiasJobOutputConfig, b.ko.Spec.ModelBiasJobOutputConfig)
	} else if a.ko.Spec.ModelBiasJobOutputConfig != nil && b.ko.Spec.ModelBiasJobOutputConfig != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID, b.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID) {
			delta.Add("Spec.ModelBiasJobOutputConfig.KMSKeyID", a.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID, b.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID)
		} else if a.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID != nil && b.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID != nil {
			if *a.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID != *b.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID {
				delta.Add("Spec.ModelBiasJobOutputConfig.KMSKeyID", a.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID, b.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID)
			}
		}

	}
	if ackcompare.HasNilDifference(a.ko.Spec.NetworkConfig, b.ko.Spec.NetworkConfig) {
		delta.Add("Spec.NetworkConfig", a.ko.Spec.NetworkConfig, b.ko.Spec.NetworkConfig)
	} else if a.ko.Spec.NetworkConfig != nil && b.ko.Spec.NetworkConfig != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption, b.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption) {
			delta.Add("Spec.NetworkConfig.EnableInterContainerTrafficEncryption", a.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption, b.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption)
		} else if a.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption != nil && b.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption != nil {
			if *a.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption != *b.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption {
				delta.Add("Spec.NetworkConfig.EnableInterContainerTrafficEncryption", a.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption, b.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.NetworkConfig.EnableNetworkIsolation, b.ko.Spec.NetworkConfig.EnableNetworkIsolation) {
			delta.Add("Spec.NetworkConfig.EnableNetworkIsolation", a.ko.Spec.NetworkConfig.EnableNetworkIsolation, b.ko.Spec.NetworkConfig.EnableNetworkIsolation)
		} else if a.ko.Spec.NetworkConfig.EnableNetworkIsolation != nil && b.ko.Spec.NetworkConfig.EnableNetworkIsolation != nil {
			if *a.ko.Spec.NetworkConfig.EnableNetworkIsolation != *b.ko.Spec.NetworkConfig.EnableNetworkIsolation {
				delta.Add("Spec.NetworkConfig.EnableNetworkIsolation", a.ko.Spec.NetworkConfig.EnableNetworkIsolation, b.ko.Spec.NetworkConfig.EnableNetworkIsolation)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.NetworkConfig.VPCConfig, b.ko.Spec.NetworkConfig.VPCConfig) {
			delta.Add("Spec.NetworkConfig.VPCConfig", a.ko.Spec.NetworkConfig.VPCConfig, b.ko.Spec.NetworkConfig.VPCConfig)
		} else if a.ko.Spec.NetworkConfig.VPCConfig != nil && b.ko.Spec.NetworkConfig.VPCConfig != nil {

			if !ackcompare.SliceStringPEqual(a.ko.Spec.NetworkConfig.VPCConfig.SecurityGroupIDs, b.ko.Spec.NetworkConfig.VPCConfig.SecurityGroupIDs) {
				delta.Add("Spec.NetworkConfig.VPCConfig.SecurityGroupIDs", a.ko.Spec.NetworkConfig.VPCConfig.SecurityGroupIDs, b.ko.Spec.NetworkConfig.VPCConfig.SecurityGroupIDs)
			}

			if !ackcompare.SliceStringPEqual(a.ko.Spec.NetworkConfig.VPCConfig.Subnets, b.ko.Spec.NetworkConfig.VPCConfig.Subnets) {
				delta.Add("Spec.NetworkConfig.VPCConfig.Subnets", a.ko.Spec.NetworkConfig.VPCConfig.Subnets, b.ko.Spec.NetworkConfig.VPCConfig.Subnets)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.RoleARN, b.ko.Spec.RoleARN) {
		delta.Add("Spec.RoleARN", a.ko.Spec.RoleARN, b.ko.Spec.RoleARN)
	} else if a.ko.Spec.RoleARN != nil && b.ko.Spec.RoleARN != nil {
		if *a.ko.Spec.RoleARN != *b.ko.Spec.RoleARN {
			delta.Add("Spec.RoleARN", a.ko.Spec.RoleARN, b.ko.Spec.RoleARN)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.StoppingCondition, b.ko.Spec.StoppingCondition) {
		delta.Add("Spec.StoppingCondition", a.ko.Spec.StoppingCondition, b.ko.Spec.StoppingCondition)
	} else if a.ko.Spec.StoppingCondition != nil && b.ko.Spec.StoppingCondition != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds, b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds) {
			delta.Add("Spec.StoppingCondition.MaxRuntimeInSeconds", a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds, b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds)
		} else if a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds != nil && b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds != nil {
			if *a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds != *b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds {
				delta.Add("Spec.StoppingCondition.MaxRuntimeInSeconds", a.ko.Spec.StoppingCondition.MaxRuntimeInSeconds, b.ko.Spec.StoppingCondition.MaxRuntimeInSeconds)
			}
		}
	}

	return delta
}
