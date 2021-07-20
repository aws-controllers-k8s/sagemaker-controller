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

package data_quality_job_definition

import (
	"context"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
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
	_ = &svcapitypes.DataQualityJobDefinition{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer exit(err)
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

	var resp *svcsdk.DescribeDataQualityJobDefinitionOutput
	resp, err = rm.sdkapi.DescribeDataQualityJobDefinitionWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeDataQualityJobDefinition", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "ResourceNotFound" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.DataQualityAppSpecification != nil {
		f1 := &svcapitypes.DataQualityAppSpecification{}
		if resp.DataQualityAppSpecification.ContainerArguments != nil {
			f1f0 := []*string{}
			for _, f1f0iter := range resp.DataQualityAppSpecification.ContainerArguments {
				var f1f0elem string
				f1f0elem = *f1f0iter
				f1f0 = append(f1f0, &f1f0elem)
			}
			f1.ContainerArguments = f1f0
		}
		if resp.DataQualityAppSpecification.ContainerEntrypoint != nil {
			f1f1 := []*string{}
			for _, f1f1iter := range resp.DataQualityAppSpecification.ContainerEntrypoint {
				var f1f1elem string
				f1f1elem = *f1f1iter
				f1f1 = append(f1f1, &f1f1elem)
			}
			f1.ContainerEntrypoint = f1f1
		}
		if resp.DataQualityAppSpecification.Environment != nil {
			f1f2 := map[string]*string{}
			for f1f2key, f1f2valiter := range resp.DataQualityAppSpecification.Environment {
				var f1f2val string
				f1f2val = *f1f2valiter
				f1f2[f1f2key] = &f1f2val
			}
			f1.Environment = f1f2
		}
		if resp.DataQualityAppSpecification.ImageUri != nil {
			f1.ImageURI = resp.DataQualityAppSpecification.ImageUri
		}
		if resp.DataQualityAppSpecification.PostAnalyticsProcessorSourceUri != nil {
			f1.PostAnalyticsProcessorSourceURI = resp.DataQualityAppSpecification.PostAnalyticsProcessorSourceUri
		}
		if resp.DataQualityAppSpecification.RecordPreprocessorSourceUri != nil {
			f1.RecordPreprocessorSourceURI = resp.DataQualityAppSpecification.RecordPreprocessorSourceUri
		}
		ko.Spec.DataQualityAppSpecification = f1
	} else {
		ko.Spec.DataQualityAppSpecification = nil
	}
	if resp.DataQualityBaselineConfig != nil {
		f2 := &svcapitypes.DataQualityBaselineConfig{}
		if resp.DataQualityBaselineConfig.BaseliningJobName != nil {
			f2.BaseliningJobName = resp.DataQualityBaselineConfig.BaseliningJobName
		}
		if resp.DataQualityBaselineConfig.ConstraintsResource != nil {
			f2f1 := &svcapitypes.MonitoringConstraintsResource{}
			if resp.DataQualityBaselineConfig.ConstraintsResource.S3Uri != nil {
				f2f1.S3URI = resp.DataQualityBaselineConfig.ConstraintsResource.S3Uri
			}
			f2.ConstraintsResource = f2f1
		}
		if resp.DataQualityBaselineConfig.StatisticsResource != nil {
			f2f2 := &svcapitypes.MonitoringStatisticsResource{}
			if resp.DataQualityBaselineConfig.StatisticsResource.S3Uri != nil {
				f2f2.S3URI = resp.DataQualityBaselineConfig.StatisticsResource.S3Uri
			}
			f2.StatisticsResource = f2f2
		}
		ko.Spec.DataQualityBaselineConfig = f2
	} else {
		ko.Spec.DataQualityBaselineConfig = nil
	}
	if resp.DataQualityJobInput != nil {
		f3 := &svcapitypes.DataQualityJobInput{}
		if resp.DataQualityJobInput.EndpointInput != nil {
			f3f0 := &svcapitypes.EndpointInput{}
			if resp.DataQualityJobInput.EndpointInput.EndTimeOffset != nil {
				f3f0.EndTimeOffset = resp.DataQualityJobInput.EndpointInput.EndTimeOffset
			}
			if resp.DataQualityJobInput.EndpointInput.EndpointName != nil {
				f3f0.EndpointName = resp.DataQualityJobInput.EndpointInput.EndpointName
			}
			if resp.DataQualityJobInput.EndpointInput.FeaturesAttribute != nil {
				f3f0.FeaturesAttribute = resp.DataQualityJobInput.EndpointInput.FeaturesAttribute
			}
			if resp.DataQualityJobInput.EndpointInput.InferenceAttribute != nil {
				f3f0.InferenceAttribute = resp.DataQualityJobInput.EndpointInput.InferenceAttribute
			}
			if resp.DataQualityJobInput.EndpointInput.LocalPath != nil {
				f3f0.LocalPath = resp.DataQualityJobInput.EndpointInput.LocalPath
			}
			if resp.DataQualityJobInput.EndpointInput.ProbabilityAttribute != nil {
				f3f0.ProbabilityAttribute = resp.DataQualityJobInput.EndpointInput.ProbabilityAttribute
			}
			if resp.DataQualityJobInput.EndpointInput.ProbabilityThresholdAttribute != nil {
				f3f0.ProbabilityThresholdAttribute = resp.DataQualityJobInput.EndpointInput.ProbabilityThresholdAttribute
			}
			if resp.DataQualityJobInput.EndpointInput.S3DataDistributionType != nil {
				f3f0.S3DataDistributionType = resp.DataQualityJobInput.EndpointInput.S3DataDistributionType
			}
			if resp.DataQualityJobInput.EndpointInput.S3InputMode != nil {
				f3f0.S3InputMode = resp.DataQualityJobInput.EndpointInput.S3InputMode
			}
			if resp.DataQualityJobInput.EndpointInput.StartTimeOffset != nil {
				f3f0.StartTimeOffset = resp.DataQualityJobInput.EndpointInput.StartTimeOffset
			}
			f3.EndpointInput = f3f0
		}
		ko.Spec.DataQualityJobInput = f3
	} else {
		ko.Spec.DataQualityJobInput = nil
	}
	if resp.DataQualityJobOutputConfig != nil {
		f4 := &svcapitypes.MonitoringOutputConfig{}
		if resp.DataQualityJobOutputConfig.KmsKeyId != nil {
			f4.KMSKeyID = resp.DataQualityJobOutputConfig.KmsKeyId
		}
		if resp.DataQualityJobOutputConfig.MonitoringOutputs != nil {
			f4f1 := []*svcapitypes.MonitoringOutput{}
			for _, f4f1iter := range resp.DataQualityJobOutputConfig.MonitoringOutputs {
				f4f1elem := &svcapitypes.MonitoringOutput{}
				if f4f1iter.S3Output != nil {
					f4f1elemf0 := &svcapitypes.MonitoringS3Output{}
					if f4f1iter.S3Output.LocalPath != nil {
						f4f1elemf0.LocalPath = f4f1iter.S3Output.LocalPath
					}
					if f4f1iter.S3Output.S3UploadMode != nil {
						f4f1elemf0.S3UploadMode = f4f1iter.S3Output.S3UploadMode
					}
					if f4f1iter.S3Output.S3Uri != nil {
						f4f1elemf0.S3URI = f4f1iter.S3Output.S3Uri
					}
					f4f1elem.S3Output = f4f1elemf0
				}
				f4f1 = append(f4f1, f4f1elem)
			}
			f4.MonitoringOutputs = f4f1
		}
		ko.Spec.DataQualityJobOutputConfig = f4
	} else {
		ko.Spec.DataQualityJobOutputConfig = nil
	}
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
		f7 := &svcapitypes.MonitoringResources{}
		if resp.JobResources.ClusterConfig != nil {
			f7f0 := &svcapitypes.MonitoringClusterConfig{}
			if resp.JobResources.ClusterConfig.InstanceCount != nil {
				f7f0.InstanceCount = resp.JobResources.ClusterConfig.InstanceCount
			}
			if resp.JobResources.ClusterConfig.InstanceType != nil {
				f7f0.InstanceType = resp.JobResources.ClusterConfig.InstanceType
			}
			if resp.JobResources.ClusterConfig.VolumeKmsKeyId != nil {
				f7f0.VolumeKMSKeyID = resp.JobResources.ClusterConfig.VolumeKmsKeyId
			}
			if resp.JobResources.ClusterConfig.VolumeSizeInGB != nil {
				f7f0.VolumeSizeInGB = resp.JobResources.ClusterConfig.VolumeSizeInGB
			}
			f7.ClusterConfig = f7f0
		}
		ko.Spec.JobResources = f7
	} else {
		ko.Spec.JobResources = nil
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
// for the ReadOne Input shape that are required but not present in the
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
) (*svcsdk.DescribeDataQualityJobDefinitionInput, error) {
	res := &svcsdk.DescribeDataQualityJobDefinitionInput{}

	if r.ko.Spec.JobDefinitionName != nil {
		res.SetJobDefinitionName(*r.ko.Spec.JobDefinitionName)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer exit(err)
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateDataQualityJobDefinitionOutput
	_ = resp
	resp, err = rm.sdkapi.CreateDataQualityJobDefinitionWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateDataQualityJobDefinition", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

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
) (*svcsdk.CreateDataQualityJobDefinitionInput, error) {
	res := &svcsdk.CreateDataQualityJobDefinitionInput{}

	if r.ko.Spec.DataQualityAppSpecification != nil {
		f0 := &svcsdk.DataQualityAppSpecification{}
		if r.ko.Spec.DataQualityAppSpecification.ContainerArguments != nil {
			f0f0 := []*string{}
			for _, f0f0iter := range r.ko.Spec.DataQualityAppSpecification.ContainerArguments {
				var f0f0elem string
				f0f0elem = *f0f0iter
				f0f0 = append(f0f0, &f0f0elem)
			}
			f0.SetContainerArguments(f0f0)
		}
		if r.ko.Spec.DataQualityAppSpecification.ContainerEntrypoint != nil {
			f0f1 := []*string{}
			for _, f0f1iter := range r.ko.Spec.DataQualityAppSpecification.ContainerEntrypoint {
				var f0f1elem string
				f0f1elem = *f0f1iter
				f0f1 = append(f0f1, &f0f1elem)
			}
			f0.SetContainerEntrypoint(f0f1)
		}
		if r.ko.Spec.DataQualityAppSpecification.Environment != nil {
			f0f2 := map[string]*string{}
			for f0f2key, f0f2valiter := range r.ko.Spec.DataQualityAppSpecification.Environment {
				var f0f2val string
				f0f2val = *f0f2valiter
				f0f2[f0f2key] = &f0f2val
			}
			f0.SetEnvironment(f0f2)
		}
		if r.ko.Spec.DataQualityAppSpecification.ImageURI != nil {
			f0.SetImageUri(*r.ko.Spec.DataQualityAppSpecification.ImageURI)
		}
		if r.ko.Spec.DataQualityAppSpecification.PostAnalyticsProcessorSourceURI != nil {
			f0.SetPostAnalyticsProcessorSourceUri(*r.ko.Spec.DataQualityAppSpecification.PostAnalyticsProcessorSourceURI)
		}
		if r.ko.Spec.DataQualityAppSpecification.RecordPreprocessorSourceURI != nil {
			f0.SetRecordPreprocessorSourceUri(*r.ko.Spec.DataQualityAppSpecification.RecordPreprocessorSourceURI)
		}
		res.SetDataQualityAppSpecification(f0)
	}
	if r.ko.Spec.DataQualityBaselineConfig != nil {
		f1 := &svcsdk.DataQualityBaselineConfig{}
		if r.ko.Spec.DataQualityBaselineConfig.BaseliningJobName != nil {
			f1.SetBaseliningJobName(*r.ko.Spec.DataQualityBaselineConfig.BaseliningJobName)
		}
		if r.ko.Spec.DataQualityBaselineConfig.ConstraintsResource != nil {
			f1f1 := &svcsdk.MonitoringConstraintsResource{}
			if r.ko.Spec.DataQualityBaselineConfig.ConstraintsResource.S3URI != nil {
				f1f1.SetS3Uri(*r.ko.Spec.DataQualityBaselineConfig.ConstraintsResource.S3URI)
			}
			f1.SetConstraintsResource(f1f1)
		}
		if r.ko.Spec.DataQualityBaselineConfig.StatisticsResource != nil {
			f1f2 := &svcsdk.MonitoringStatisticsResource{}
			if r.ko.Spec.DataQualityBaselineConfig.StatisticsResource.S3URI != nil {
				f1f2.SetS3Uri(*r.ko.Spec.DataQualityBaselineConfig.StatisticsResource.S3URI)
			}
			f1.SetStatisticsResource(f1f2)
		}
		res.SetDataQualityBaselineConfig(f1)
	}
	if r.ko.Spec.DataQualityJobInput != nil {
		f2 := &svcsdk.DataQualityJobInput{}
		if r.ko.Spec.DataQualityJobInput.EndpointInput != nil {
			f2f0 := &svcsdk.EndpointInput{}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.EndTimeOffset != nil {
				f2f0.SetEndTimeOffset(*r.ko.Spec.DataQualityJobInput.EndpointInput.EndTimeOffset)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.EndpointName != nil {
				f2f0.SetEndpointName(*r.ko.Spec.DataQualityJobInput.EndpointInput.EndpointName)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.FeaturesAttribute != nil {
				f2f0.SetFeaturesAttribute(*r.ko.Spec.DataQualityJobInput.EndpointInput.FeaturesAttribute)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.InferenceAttribute != nil {
				f2f0.SetInferenceAttribute(*r.ko.Spec.DataQualityJobInput.EndpointInput.InferenceAttribute)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.LocalPath != nil {
				f2f0.SetLocalPath(*r.ko.Spec.DataQualityJobInput.EndpointInput.LocalPath)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.ProbabilityAttribute != nil {
				f2f0.SetProbabilityAttribute(*r.ko.Spec.DataQualityJobInput.EndpointInput.ProbabilityAttribute)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.ProbabilityThresholdAttribute != nil {
				f2f0.SetProbabilityThresholdAttribute(*r.ko.Spec.DataQualityJobInput.EndpointInput.ProbabilityThresholdAttribute)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.S3DataDistributionType != nil {
				f2f0.SetS3DataDistributionType(*r.ko.Spec.DataQualityJobInput.EndpointInput.S3DataDistributionType)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.S3InputMode != nil {
				f2f0.SetS3InputMode(*r.ko.Spec.DataQualityJobInput.EndpointInput.S3InputMode)
			}
			if r.ko.Spec.DataQualityJobInput.EndpointInput.StartTimeOffset != nil {
				f2f0.SetStartTimeOffset(*r.ko.Spec.DataQualityJobInput.EndpointInput.StartTimeOffset)
			}
			f2.SetEndpointInput(f2f0)
		}
		res.SetDataQualityJobInput(f2)
	}
	if r.ko.Spec.DataQualityJobOutputConfig != nil {
		f3 := &svcsdk.MonitoringOutputConfig{}
		if r.ko.Spec.DataQualityJobOutputConfig.KMSKeyID != nil {
			f3.SetKmsKeyId(*r.ko.Spec.DataQualityJobOutputConfig.KMSKeyID)
		}
		if r.ko.Spec.DataQualityJobOutputConfig.MonitoringOutputs != nil {
			f3f1 := []*svcsdk.MonitoringOutput{}
			for _, f3f1iter := range r.ko.Spec.DataQualityJobOutputConfig.MonitoringOutputs {
				f3f1elem := &svcsdk.MonitoringOutput{}
				if f3f1iter.S3Output != nil {
					f3f1elemf0 := &svcsdk.MonitoringS3Output{}
					if f3f1iter.S3Output.LocalPath != nil {
						f3f1elemf0.SetLocalPath(*f3f1iter.S3Output.LocalPath)
					}
					if f3f1iter.S3Output.S3UploadMode != nil {
						f3f1elemf0.SetS3UploadMode(*f3f1iter.S3Output.S3UploadMode)
					}
					if f3f1iter.S3Output.S3URI != nil {
						f3f1elemf0.SetS3Uri(*f3f1iter.S3Output.S3URI)
					}
					f3f1elem.SetS3Output(f3f1elemf0)
				}
				f3f1 = append(f3f1, f3f1elem)
			}
			f3.SetMonitoringOutputs(f3f1)
		}
		res.SetDataQualityJobOutputConfig(f3)
	}
	if r.ko.Spec.JobDefinitionName != nil {
		res.SetJobDefinitionName(*r.ko.Spec.JobDefinitionName)
	}
	if r.ko.Spec.JobResources != nil {
		f5 := &svcsdk.MonitoringResources{}
		if r.ko.Spec.JobResources.ClusterConfig != nil {
			f5f0 := &svcsdk.MonitoringClusterConfig{}
			if r.ko.Spec.JobResources.ClusterConfig.InstanceCount != nil {
				f5f0.SetInstanceCount(*r.ko.Spec.JobResources.ClusterConfig.InstanceCount)
			}
			if r.ko.Spec.JobResources.ClusterConfig.InstanceType != nil {
				f5f0.SetInstanceType(*r.ko.Spec.JobResources.ClusterConfig.InstanceType)
			}
			if r.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID != nil {
				f5f0.SetVolumeKmsKeyId(*r.ko.Spec.JobResources.ClusterConfig.VolumeKMSKeyID)
			}
			if r.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB != nil {
				f5f0.SetVolumeSizeInGB(*r.ko.Spec.JobResources.ClusterConfig.VolumeSizeInGB)
			}
			f5.SetClusterConfig(f5f0)
		}
		res.SetJobResources(f5)
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
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer exit(err)
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteDataQualityJobDefinitionOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteDataQualityJobDefinitionWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteDataQualityJobDefinition", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteDataQualityJobDefinitionInput, error) {
	res := &svcsdk.DeleteDataQualityJobDefinitionInput{}

	if r.ko.Spec.JobDefinitionName != nil {
		res.SetJobDefinitionName(*r.ko.Spec.JobDefinitionName)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.DataQualityJobDefinition,
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
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
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
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
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
