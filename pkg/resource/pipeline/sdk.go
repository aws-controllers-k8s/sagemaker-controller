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

package pipeline

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
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
	_ = &svcapitypes.Pipeline{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
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

	var resp *svcsdk.DescribePipelineOutput
	resp, err = rm.sdkapi.DescribePipelineWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribePipeline", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "ResourceNotFound" && strings.HasSuffix(awsErr.Message(), "does not exist.") {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.CreationTime != nil {
		ko.Status.CreationTime = &metav1.Time{*resp.CreationTime}
	} else {
		ko.Status.CreationTime = nil
	}
	if resp.LastModifiedTime != nil {
		ko.Status.LastModifiedTime = &metav1.Time{*resp.LastModifiedTime}
	} else {
		ko.Status.LastModifiedTime = nil
	}
	if resp.ParallelismConfiguration != nil {
		f5 := &svcapitypes.ParallelismConfiguration{}
		if resp.ParallelismConfiguration.MaxParallelExecutionSteps != nil {
			f5.MaxParallelExecutionSteps = resp.ParallelismConfiguration.MaxParallelExecutionSteps
		}
		ko.Spec.ParallelismConfiguration = f5
	} else {
		ko.Spec.ParallelismConfiguration = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.PipelineArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.PipelineArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.PipelineDefinition != nil {
		ko.Spec.PipelineDefinition = resp.PipelineDefinition
	} else {
		ko.Spec.PipelineDefinition = nil
	}
	if resp.PipelineDescription != nil {
		ko.Spec.PipelineDescription = resp.PipelineDescription
	} else {
		ko.Spec.PipelineDescription = nil
	}
	if resp.PipelineDisplayName != nil {
		ko.Spec.PipelineDisplayName = resp.PipelineDisplayName
	} else {
		ko.Spec.PipelineDisplayName = nil
	}
	if resp.PipelineName != nil {
		ko.Spec.PipelineName = resp.PipelineName
	} else {
		ko.Spec.PipelineName = nil
	}
	if resp.PipelineStatus != nil {
		ko.Status.PipelineStatus = resp.PipelineStatus
	} else {
		ko.Status.PipelineStatus = nil
	}
	if resp.RoleArn != nil {
		ko.Spec.RoleARN = resp.RoleArn
	} else {
		ko.Spec.RoleARN = nil
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
	return r.ko.Spec.PipelineName == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribePipelineInput, error) {
	res := &svcsdk.DescribePipelineInput{}

	if r.ko.Spec.PipelineName != nil {
		res.SetPipelineName(*r.ko.Spec.PipelineName)
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
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreatePipelineOutput
	_ = resp
	resp, err = rm.sdkapi.CreatePipelineWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreatePipeline", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.PipelineArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.PipelineArn)
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
) (*svcsdk.CreatePipelineInput, error) {
	res := &svcsdk.CreatePipelineInput{}

	if r.ko.Spec.ParallelismConfiguration != nil {
		f0 := &svcsdk.ParallelismConfiguration{}
		if r.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps != nil {
			f0.SetMaxParallelExecutionSteps(*r.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps)
		}
		res.SetParallelismConfiguration(f0)
	}
	if r.ko.Spec.PipelineDefinition != nil {
		res.SetPipelineDefinition(*r.ko.Spec.PipelineDefinition)
	}
	if r.ko.Spec.PipelineDefinitionS3Location != nil {
		f2 := &svcsdk.PipelineDefinitionS3Location{}
		if r.ko.Spec.PipelineDefinitionS3Location.Bucket != nil {
			f2.SetBucket(*r.ko.Spec.PipelineDefinitionS3Location.Bucket)
		}
		if r.ko.Spec.PipelineDefinitionS3Location.ObjectKey != nil {
			f2.SetObjectKey(*r.ko.Spec.PipelineDefinitionS3Location.ObjectKey)
		}
		if r.ko.Spec.PipelineDefinitionS3Location.VersionID != nil {
			f2.SetVersionId(*r.ko.Spec.PipelineDefinitionS3Location.VersionID)
		}
		res.SetPipelineDefinitionS3Location(f2)
	}
	if r.ko.Spec.PipelineDescription != nil {
		res.SetPipelineDescription(*r.ko.Spec.PipelineDescription)
	}
	if r.ko.Spec.PipelineDisplayName != nil {
		res.SetPipelineDisplayName(*r.ko.Spec.PipelineDisplayName)
	}
	if r.ko.Spec.PipelineName != nil {
		res.SetPipelineName(*r.ko.Spec.PipelineName)
	}
	if r.ko.Spec.RoleARN != nil {
		res.SetRoleArn(*r.ko.Spec.RoleARN)
	}
	if r.ko.Spec.Tags != nil {
		f7 := []*svcsdk.Tag{}
		for _, f7iter := range r.ko.Spec.Tags {
			f7elem := &svcsdk.Tag{}
			if f7iter.Key != nil {
				f7elem.SetKey(*f7iter.Key)
			}
			if f7iter.Value != nil {
				f7elem.SetValue(*f7iter.Value)
			}
			f7 = append(f7, f7elem)
		}
		res.SetTags(f7)
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
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newUpdateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdatePipelineOutput
	_ = resp
	resp, err = rm.sdkapi.UpdatePipelineWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdatePipeline", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.PipelineArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.PipelineArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.UpdatePipelineInput, error) {
	res := &svcsdk.UpdatePipelineInput{}

	if r.ko.Spec.ParallelismConfiguration != nil {
		f0 := &svcsdk.ParallelismConfiguration{}
		if r.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps != nil {
			f0.SetMaxParallelExecutionSteps(*r.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps)
		}
		res.SetParallelismConfiguration(f0)
	}
	if r.ko.Spec.PipelineDefinition != nil {
		res.SetPipelineDefinition(*r.ko.Spec.PipelineDefinition)
	}
	if r.ko.Spec.PipelineDefinitionS3Location != nil {
		f2 := &svcsdk.PipelineDefinitionS3Location{}
		if r.ko.Spec.PipelineDefinitionS3Location.Bucket != nil {
			f2.SetBucket(*r.ko.Spec.PipelineDefinitionS3Location.Bucket)
		}
		if r.ko.Spec.PipelineDefinitionS3Location.ObjectKey != nil {
			f2.SetObjectKey(*r.ko.Spec.PipelineDefinitionS3Location.ObjectKey)
		}
		if r.ko.Spec.PipelineDefinitionS3Location.VersionID != nil {
			f2.SetVersionId(*r.ko.Spec.PipelineDefinitionS3Location.VersionID)
		}
		res.SetPipelineDefinitionS3Location(f2)
	}
	if r.ko.Spec.PipelineDescription != nil {
		res.SetPipelineDescription(*r.ko.Spec.PipelineDescription)
	}
	if r.ko.Spec.PipelineDisplayName != nil {
		res.SetPipelineDisplayName(*r.ko.Spec.PipelineDisplayName)
	}
	if r.ko.Spec.PipelineName != nil {
		res.SetPipelineName(*r.ko.Spec.PipelineName)
	}
	if r.ko.Spec.RoleARN != nil {
		res.SetRoleArn(*r.ko.Spec.RoleARN)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeletePipelineOutput
	_ = resp
	resp, err = rm.sdkapi.DeletePipelineWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeletePipeline", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeletePipelineInput, error) {
	res := &svcsdk.DeletePipelineInput{}

	if r.ko.Spec.PipelineName != nil {
		res.SetPipelineName(*r.ko.Spec.PipelineName)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Pipeline,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
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
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
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
				errorMessage = awsErr.Error()
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
	case "InvalidParameterCombination",
		"InvalidParameterValue",
		"MissingParameter",
		"ResourceNotFound":
		return true
	default:
		return false
	}
}