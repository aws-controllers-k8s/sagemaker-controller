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
	"context"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.SageMaker{}
	_ = &svcapitypes.ModelBiasJobDefinition{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.DescribeModelBiasJobDefinitionWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeModelBiasJobDefinition", respErr)
	if respErr != nil {
		if awsErr, ok := ackerr.AWSError(respErr); ok && awsErr.Code() == "ResourceNotFound" {
			return nil, ackerr.NotFound
		}
		return nil, respErr
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.JobDefinitionArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.JobDefinitionArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.JobDefinitionName != nil {
		ko.Spec.JobDefinitionName = resp.JobDefinitionName
	} else {
		ko.Spec.JobDefinitionName = nil
	}
	if resp.JobResources != nil {
		f3 := &svcapitypes.MonitoringResources{}
		if resp.JobResources.ClusterConfig != nil {
			f3f0 := &svcapitypes.MonitoringClusterConfig{}
			if resp.JobResources.ClusterConfig.InstanceCount != nil {
				f3f0.InstanceCount = resp.JobResources.ClusterConfig.InstanceCount
			}
			if resp.JobResources.ClusterConfig.InstanceType != nil {
				f3f0.InstanceType = resp.JobResources.ClusterConfig.InstanceType
			}
			if resp.JobResources.ClusterConfig.VolumeKmsKeyId != nil {
				f3f0.VolumeKMSKeyID = resp.JobResources.ClusterConfig.VolumeKmsKeyId
			}
			if resp.JobResources.ClusterConfig.VolumeSizeInGB != nil {
				f3f0.VolumeSizeInGB = resp.JobResources.ClusterConfig.VolumeSizeInGB
			}
			f3.ClusterConfig = f3f0
		}
		ko.Spec.JobResources = f3
	} else {
		ko.Spec.JobResources = nil
	}
	if resp.ModelBiasAppSpecification != nil {
		f4 := &svcapitypes.ModelBiasAppSpecification{}
		if resp.ModelBiasAppSpecification.ConfigUri != nil {
			f4.ConfigURI = resp.ModelBiasAppSpecification.ConfigUri
		}
		if resp.ModelBiasAppSpecification.Environment != nil {
			f4f1 := map[string]*string{}
			for f4f1key, f4f1valiter := range resp.ModelBiasAppSpecification.Environment {
				var f4f1val string
				f4f1val = *f4f1valiter
				f4f1[f4f1key] = &f4f1val
			}
			f4.Environment = f4f1
		}
		if resp.ModelBiasAppSpecification.ImageUri != nil {
			f4.ImageURI = resp.ModelBiasAppSpecification.ImageUri
		}
		ko.Spec.ModelBiasAppSpecification = f4
	} else {
		ko.Spec.ModelBiasAppSpecification = nil
	}
	if resp.ModelBiasBaselineConfig != nil {
		f5 := &svcapitypes.ModelBiasBaselineConfig{}
		if resp.ModelBiasBaselineConfig.BaseliningJobName != nil {
			f5.BaseliningJobName = resp.ModelBiasBaselineConfig.BaseliningJobName
		}
		if resp.ModelBiasBaselineConfig.ConstraintsResource != nil {
			f5f1 := &svcapitypes.MonitoringConstraintsResource{}
			if resp.ModelBiasBaselineConfig.ConstraintsResource.S3Uri != nil {
				f5f1.S3URI = resp.ModelBiasBaselineConfig.ConstraintsResource.S3Uri
			}
			f5.ConstraintsResource = f5f1
		}
		ko.Spec.ModelBiasBaselineConfig = f5
	} else {
		ko.Spec.ModelBiasBaselineConfig = nil
	}
	if resp.ModelBiasJobInput != nil {
		f6 := &svcapitypes.ModelBiasJobInput{}
		if resp.ModelBiasJobInput.EndpointInput != nil {
			f6f0 := &svcapitypes.EndpointInput{}
			if resp.ModelBiasJobInput.EndpointInput.EndTimeOffset != nil {
				f6f0.EndTimeOffset = resp.ModelBiasJobInput.EndpointInput.EndTimeOffset
			}
			if resp.ModelBiasJobInput.EndpointInput.EndpointName != nil {
				f6f0.EndpointName = resp.ModelBiasJobInput.EndpointInput.EndpointName
			}
			if resp.ModelBiasJobInput.EndpointInput.FeaturesAttribute != nil {
				f6f0.FeaturesAttribute = resp.ModelBiasJobInput.EndpointInput.FeaturesAttribute
			}
			if resp.ModelBiasJobInput.EndpointInput.InferenceAttribute != nil {
				f6f0.InferenceAttribute = resp.ModelBiasJobInput.EndpointInput.InferenceAttribute
			}
			if resp.ModelBiasJobInput.EndpointInput.LocalPath != nil {
				f6f0.LocalPath = resp.ModelBiasJobInput.EndpointInput.LocalPath
			}
			if resp.ModelBiasJobInput.EndpointInput.ProbabilityAttribute != nil {
				f6f0.ProbabilityAttribute = resp.ModelBiasJobInput.EndpointInput.ProbabilityAttribute
			}
			if resp.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute != nil {
				f6f0.ProbabilityThresholdAttribute = resp.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute
			}
			if resp.ModelBiasJobInput.EndpointInput.S3DataDistributionType != nil {
				f6f0.S3DataDistributionType = resp.ModelBiasJobInput.EndpointInput.S3DataDistributionType
			}
			if resp.ModelBiasJobInput.EndpointInput.S3InputMode != nil {
				f6f0.S3InputMode = resp.ModelBiasJobInput.EndpointInput.S3InputMode
			}
			if resp.ModelBiasJobInput.EndpointInput.StartTimeOffset != nil {
				f6f0.StartTimeOffset = resp.ModelBiasJobInput.EndpointInput.StartTimeOffset
			}
			f6.EndpointInput = f6f0
		}
		if resp.ModelBiasJobInput.GroundTruthS3Input != nil {
			f6f1 := &svcapitypes.MonitoringGroundTruthS3Input{}
			if resp.ModelBiasJobInput.GroundTruthS3Input.S3Uri != nil {
				f6f1.S3URI = resp.ModelBiasJobInput.GroundTruthS3Input.S3Uri
			}
			f6.GroundTruthS3Input = f6f1
		}
		ko.Spec.ModelBiasJobInput = f6
	} else {
		ko.Spec.ModelBiasJobInput = nil
	}
	if resp.ModelBiasJobOutputConfig != nil {
		f7 := &svcapitypes.MonitoringOutputConfig{}
		if resp.ModelBiasJobOutputConfig.KmsKeyId != nil {
			f7.KMSKeyID = resp.ModelBiasJobOutputConfig.KmsKeyId
		}
		if resp.ModelBiasJobOutputConfig.MonitoringOutputs != nil {
			f7f1 := []*svcapitypes.MonitoringOutput{}
			for _, f7f1iter := range resp.ModelBiasJobOutputConfig.MonitoringOutputs {
				f7f1elem := &svcapitypes.MonitoringOutput{}
				if f7f1iter.S3Output != nil {
					f7f1elemf0 := &svcapitypes.MonitoringS3Output{}
					if f7f1iter.S3Output.LocalPath != nil {
						f7f1elemf0.LocalPath = f7f1iter.S3Output.LocalPath
					}
					if f7f1iter.S3Output.S3UploadMode != nil {
						f7f1elemf0.S3UploadMode = f7f1iter.S3Output.S3UploadMode
					}
					if f7f1iter.S3Output.S3Uri != nil {
						f7f1elemf0.S3URI = f7f1iter.S3Output.S3Uri
					}
					f7f1elem.S3Output = f7f1elemf0
				}
				f7f1 = append(f7f1, f7f1elem)
			}
			f7.MonitoringOutputs = f7f1
		}
		ko.Spec.ModelBiasJobOutputConfig = f7
	} else {
		ko.Spec.ModelBiasJobOutputConfig = nil
	}
	if resp.NetworkConfig != nil {
		f8 := &svcapitypes.MonitoringNetworkConfig{}
		if resp.NetworkConfig.EnableInterContainerTrafficEncryption != nil {
			f8.EnableInterContainerTrafficEncryption = resp.NetworkConfig.EnableInterContainerTrafficEncryption
		}
		if resp.NetworkConfig.EnableNetworkIsolation != nil {
			f8.EnableNetworkIsolation = resp.NetworkConfig.EnableNetworkIsolation
		}
		if resp.NetworkConfig.VpcConfig != nil {
			f8f2 := &svcapitypes.VPCConfig{}
			if resp.NetworkConfig.VpcConfig.SecurityGroupIds != nil {
				f8f2f0 := []*string{}
				for _, f8f2f0iter := range resp.NetworkConfig.VpcConfig.SecurityGroupIds {
					var f8f2f0elem string
					f8f2f0elem = *f8f2f0iter
					f8f2f0 = append(f8f2f0, &f8f2f0elem)
				}
				f8f2.SecurityGroupIDs = f8f2f0
			}
			if resp.NetworkConfig.VpcConfig.Subnets != nil {
				f8f2f1 := []*string{}
				for _, f8f2f1iter := range resp.NetworkConfig.VpcConfig.Subnets {
					var f8f2f1elem string
					f8f2f1elem = *f8f2f1iter
					f8f2f1 = append(f8f2f1, &f8f2f1elem)
				}
				f8f2.Subnets = f8f2f1
			}
			f8.VPCConfig = f8f2
		}
		ko.Spec.NetworkConfig = f8
	} else {
		ko.Spec.NetworkConfig = nil
	}
	if resp.RoleArn != nil {
		ko.Spec.RoleARN = resp.RoleArn
	} else {
		ko.Spec.RoleARN = nil
	}
	if resp.StoppingCondition != nil {
		f10 := &svcapitypes.MonitoringStoppingCondition{}
		if resp.StoppingCondition.MaxRuntimeInSeconds != nil {
			f10.MaxRuntimeInSeconds = resp.StoppingCondition.MaxRuntimeInSeconds
		}
		ko.Spec.StoppingCondition = f10
	} else {
		ko.Spec.StoppingCondition = nil
	}

	rm.setStatusDefaults(ko)

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required by not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.JobDefinitionName == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeModelBiasJobDefinitionInput, error) {
	res := &svcsdk.DescribeModelBiasJobDefinitionInput{}

	if r.ko.Spec.JobDefinitionName != nil {
		res.SetJobDefinitionName(*r.ko.Spec.JobDefinitionName)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a new resource with any fields in the Status field filled in
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	input, err := rm.newCreateRequestPayload(ctx, r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.CreateModelBiasJobDefinitionWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateModelBiasJobDefinition", respErr)
	if respErr != nil {
		return nil, respErr
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.JobDefinitionArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.JobDefinitionArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}

	rm.setStatusDefaults(ko)

	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateModelBiasJobDefinitionInput, error) {
	res := &svcsdk.CreateModelBiasJobDefinitionInput{}

	if r.ko.Spec.JobDefinitionName != nil {
		res.SetJobDefinitionName(*r.ko.Spec.JobDefinitionName)
	}
	if r.ko.Spec.JobResources != nil {
		f1 := &svcsdk.MonitoringResources{}
		if r.ko.Spec.JobResources.ClusterConfig != nil {
			f1f0 := &svcsdk.MonitoringClusterConfig{}
			if r.ko.Spec.JobResources.ClusterConfig.InstanceCount != nil {
				f1f0.SetInstanceCount(*r.ko.Spec.JobResources.ClusterConfig.InstanceCount)
			}
			if r.ko.Spec.JobResources.ClusterConfig.InstanceType != nil {
				f1f0.SetInstanceType(*r.ko.Spec.JobResources.ClusterConfig.InstanceType)
			}
			if r.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID != nil {
				f1f0.SetVolumeKmsKeyId(*r.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID)
			}
			if r.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB != nil {
				f1f0.SetVolumeSizeInGB(*r.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB)
			}
			f1.SetClusterConfig(f1f0)
		}
		res.SetJobResources(f1)
	}
	if r.ko.Spec.ModelBiasAppSpecification != nil {
		f2 := &svcsdk.ModelBiasAppSpecification{}
		if r.ko.Spec.ModelBiasAppSpecification.ConfigURI != nil {
			f2.SetConfigUri(*r.ko.Spec.ModelBiasAppSpecification.ConfigURI)
		}
		if r.ko.Spec.ModelBiasAppSpecification.Environment != nil {
			f2f1 := map[string]*string{}
			for f2f1key, f2f1valiter := range r.ko.Spec.ModelBiasAppSpecification.Environment {
				var f2f1val string
				f2f1val = *f2f1valiter
				f2f1[f2f1key] = &f2f1val
			}
			f2.SetEnvironment(f2f1)
		}
		if r.ko.Spec.ModelBiasAppSpecification.ImageURI != nil {
			f2.SetImageUri(*r.ko.Spec.ModelBiasAppSpecification.ImageURI)
		}
		res.SetModelBiasAppSpecification(f2)
	}
	if r.ko.Spec.ModelBiasBaselineConfig != nil {
		f3 := &svcsdk.ModelBiasBaselineConfig{}
		if r.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName != nil {
			f3.SetBaseliningJobName(*r.ko.Spec.ModelBiasBaselineConfig.BaseliningJobName)
		}
		if r.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource != nil {
			f3f1 := &svcsdk.MonitoringConstraintsResource{}
			if r.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI != nil {
				f3f1.SetS3Uri(*r.ko.Spec.ModelBiasBaselineConfig.ConstraintsResource.S3URI)
			}
			f3.SetConstraintsResource(f3f1)
		}
		res.SetModelBiasBaselineConfig(f3)
	}
	if r.ko.Spec.ModelBiasJobInput != nil {
		f4 := &svcsdk.ModelBiasJobInput{}
		if r.ko.Spec.ModelBiasJobInput.EndpointInput != nil {
			f4f0 := &svcsdk.EndpointInput{}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset != nil {
				f4f0.SetEndTimeOffset(*r.ko.Spec.ModelBiasJobInput.EndpointInput.EndTimeOffset)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName != nil {
				f4f0.SetEndpointName(*r.ko.Spec.ModelBiasJobInput.EndpointInput.EndpointName)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute != nil {
				f4f0.SetFeaturesAttribute(*r.ko.Spec.ModelBiasJobInput.EndpointInput.FeaturesAttribute)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute != nil {
				f4f0.SetInferenceAttribute(*r.ko.Spec.ModelBiasJobInput.EndpointInput.InferenceAttribute)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath != nil {
				f4f0.SetLocalPath(*r.ko.Spec.ModelBiasJobInput.EndpointInput.LocalPath)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute != nil {
				f4f0.SetProbabilityAttribute(*r.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityAttribute)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute != nil {
				f4f0.SetProbabilityThresholdAttribute(*r.ko.Spec.ModelBiasJobInput.EndpointInput.ProbabilityThresholdAttribute)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType != nil {
				f4f0.SetS3DataDistributionType(*r.ko.Spec.ModelBiasJobInput.EndpointInput.S3DataDistributionType)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode != nil {
				f4f0.SetS3InputMode(*r.ko.Spec.ModelBiasJobInput.EndpointInput.S3InputMode)
			}
			if r.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset != nil {
				f4f0.SetStartTimeOffset(*r.ko.Spec.ModelBiasJobInput.EndpointInput.StartTimeOffset)
			}
			f4.SetEndpointInput(f4f0)
		}
		if r.ko.Spec.ModelBiasJobInput.GroundTruthS3Input != nil {
			f4f1 := &svcsdk.MonitoringGroundTruthS3Input{}
			if r.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI != nil {
				f4f1.SetS3Uri(*r.ko.Spec.ModelBiasJobInput.GroundTruthS3Input.S3URI)
			}
			f4.SetGroundTruthS3Input(f4f1)
		}
		res.SetModelBiasJobInput(f4)
	}
	if r.ko.Spec.ModelBiasJobOutputConfig != nil {
		f5 := &svcsdk.MonitoringOutputConfig{}
		if r.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID != nil {
			f5.SetKmsKeyId(*r.ko.Spec.ModelBiasJobOutputConfig.KMSKeyID)
		}
		if r.ko.Spec.ModelBiasJobOutputConfig.MonitoringOutputs != nil {
			f5f1 := []*svcsdk.MonitoringOutput{}
			for _, f5f1iter := range r.ko.Spec.ModelBiasJobOutputConfig.MonitoringOutputs {
				f5f1elem := &svcsdk.MonitoringOutput{}
				if f5f1iter.S3Output != nil {
					f5f1elemf0 := &svcsdk.MonitoringS3Output{}
					if f5f1iter.S3Output.LocalPath != nil {
						f5f1elemf0.SetLocalPath(*f5f1iter.S3Output.LocalPath)
					}
					if f5f1iter.S3Output.S3UploadMode != nil {
						f5f1elemf0.SetS3UploadMode(*f5f1iter.S3Output.S3UploadMode)
					}
					if f5f1iter.S3Output.S3URI != nil {
						f5f1elemf0.SetS3Uri(*f5f1iter.S3Output.S3URI)
					}
					f5f1elem.SetS3Output(f5f1elemf0)
				}
				f5f1 = append(f5f1, f5f1elem)
			}
			f5.SetMonitoringOutputs(f5f1)
		}
		res.SetModelBiasJobOutputConfig(f5)
	}
	if r.ko.Spec.NetworkConfig != nil {
		f6 := &svcsdk.MonitoringNetworkConfig{}
		if r.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption != nil {
			f6.SetEnableInterContainerTrafficEncryption(*r.ko.Spec.NetworkConfig.EnableInterContainerTrafficEncryption)
		}
		if r.ko.Spec.NetworkConfig.EnableNetworkIsolation != nil {
			f6.SetEnableNetworkIsolation(*r.ko.Spec.NetworkConfig.EnableNetworkIsolation)
		}
		if r.ko.Spec.NetworkConfig.VPCConfig != nil {
			f6f2 := &svcsdk.VpcConfig{}
			if r.ko.Spec.NetworkConfig.VPCConfig.SecurityGroupIDs != nil {
				f6f2f0 := []*string{}
				for _, f6f2f0iter := range r.ko.Spec.NetworkConfig.VPCConfig.SecurityGroupIDs {
					var f6f2f0elem string
					f6f2f0elem = *f6f2f0iter
					f6f2f0 = append(f6f2f0, &f6f2f0elem)
				}
				f6f2.SetSecurityGroupIds(f6f2f0)
			}
			if r.ko.Spec.NetworkConfig.VPCConfig.Subnets != nil {
				f6f2f1 := []*string{}
				for _, f6f2f1iter := range r.ko.Spec.NetworkConfig.VPCConfig.Subnets {
					var f6f2f1elem string
					f6f2f1elem = *f6f2f1iter
					f6f2f1 = append(f6f2f1, &f6f2f1elem)
				}
				f6f2.SetSubnets(f6f2f1)
			}
			f6.SetVpcConfig(f6f2)
		}
		res.SetNetworkConfig(f6)
	}
	if r.ko.Spec.RoleARN != nil {
		res.SetRoleArn(*r.ko.Spec.RoleARN)
	}
	if r.ko.Spec.StoppingCondition != nil {
		f8 := &svcsdk.MonitoringStoppingCondition{}
		if r.ko.Spec.StoppingCondition.MaxRuntimeInSeconds != nil {
			f8.SetMaxRuntimeInSeconds(*r.ko.Spec.StoppingCondition.MaxRuntimeInSeconds)
		}
		res.SetStoppingCondition(f8)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// TODO(jaypipes): Figure this out...
	return nil, ackerr.NotImplemented
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) error {

	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return err
	}
	_, respErr := rm.sdkapi.DeleteModelBiasJobDefinitionWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteModelBiasJobDefinition", respErr)
	return respErr
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteModelBiasJobDefinitionInput, error) {
	res := &svcsdk.DeleteModelBiasJobDefinitionInput{}

	if r.ko.Spec.JobDefinitionName != nil {
		res.SetJobDefinitionName(*r.ko.Spec.JobDefinitionName)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.ModelBiasJobDefinition,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
	}

	if rm.terminalAWSError(err) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		terminalCondition.Status = corev1.ConditionTrue
		awsErr, _ := ackerr.AWSError(err)
		errorMessage := awsErr.Message()
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Message()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	if terminalCondition != nil || recoverableCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "ResourceLimitExceeded",
		"ResourceNotFound",
		"ResourceInUse",
		"OptInRequired",
		"InvalidParameterCombination",
		"InvalidParameterValue",
		"MissingParameter",
		"MissingAction",
		"InvalidQueryParameter",
		"MalformedQueryString",
		"InvalidAction",
		"UnrecognizedClientException":
		return true
	default:
		return false
	}
}
