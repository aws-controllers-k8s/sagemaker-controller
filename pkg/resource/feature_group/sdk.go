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

package feature_group

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/sagemaker"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	smithy "github.com/aws/smithy-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &svcsdk.Client{}
	_ = &svcapitypes.FeatureGroup{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
	_ = &aws.Config{}
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

	var resp *svcsdk.DescribeFeatureGroupOutput
	resp, err = rm.sdkapi.DescribeFeatureGroup(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeFeatureGroup", err)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "ResourceNotFound" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if resp.Description != nil {
		ko.Spec.Description = resp.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.EventTimeFeatureName != nil {
		ko.Spec.EventTimeFeatureName = resp.EventTimeFeatureName
	} else {
		ko.Spec.EventTimeFeatureName = nil
	}
	if resp.FailureReason != nil {
		ko.Status.FailureReason = resp.FailureReason
	} else {
		ko.Status.FailureReason = nil
	}
	if resp.FeatureDefinitions != nil {
		f4 := []*svcapitypes.FeatureDefinition{}
		for _, f4iter := range resp.FeatureDefinitions {
			f4elem := &svcapitypes.FeatureDefinition{}
			if f4iter.CollectionConfig != nil {
				f4elemf0 := &svcapitypes.CollectionConfig{}
				switch f4iter.CollectionConfig.(type) {
				case *svcsdktypes.CollectionConfigMemberVectorConfig:
					f4elemf0f0 := f4iter.CollectionConfig.(*svcsdktypes.CollectionConfigMemberVectorConfig)
					if f4elemf0f0 != nil {
						f4elemf0f0f0 := &svcapitypes.VectorConfig{}
						if f4elemf0f0.Value.Dimension != nil {
							dimensionCopy := int64(*f4elemf0f0.Value.Dimension)
							f4elemf0f0f0.Dimension = &dimensionCopy
						}
						f4elemf0.VectorConfig = f4elemf0f0f0
					}
				}
				f4elem.CollectionConfig = f4elemf0
			}
			if f4iter.CollectionType != "" {
				f4elem.CollectionType = aws.String(string(f4iter.CollectionType))
			}
			if f4iter.FeatureName != nil {
				f4elem.FeatureName = f4iter.FeatureName
			}
			if f4iter.FeatureType != "" {
				f4elem.FeatureType = aws.String(string(f4iter.FeatureType))
			}
			f4 = append(f4, f4elem)
		}
		ko.Spec.FeatureDefinitions = f4
	} else {
		ko.Spec.FeatureDefinitions = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.FeatureGroupArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.FeatureGroupArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.FeatureGroupName != nil {
		ko.Spec.FeatureGroupName = resp.FeatureGroupName
	} else {
		ko.Spec.FeatureGroupName = nil
	}
	if resp.FeatureGroupStatus != "" {
		ko.Status.FeatureGroupStatus = aws.String(string(resp.FeatureGroupStatus))
	} else {
		ko.Status.FeatureGroupStatus = nil
	}
	if resp.OfflineStoreConfig != nil {
		f11 := &svcapitypes.OfflineStoreConfig{}
		if resp.OfflineStoreConfig.DataCatalogConfig != nil {
			f11f0 := &svcapitypes.DataCatalogConfig{}
			if resp.OfflineStoreConfig.DataCatalogConfig.Catalog != nil {
				f11f0.Catalog = resp.OfflineStoreConfig.DataCatalogConfig.Catalog
			}
			if resp.OfflineStoreConfig.DataCatalogConfig.Database != nil {
				f11f0.Database = resp.OfflineStoreConfig.DataCatalogConfig.Database
			}
			if resp.OfflineStoreConfig.DataCatalogConfig.TableName != nil {
				f11f0.TableName = resp.OfflineStoreConfig.DataCatalogConfig.TableName
			}
			f11.DataCatalogConfig = f11f0
		}
		if resp.OfflineStoreConfig.DisableGlueTableCreation != nil {
			f11.DisableGlueTableCreation = resp.OfflineStoreConfig.DisableGlueTableCreation
		}
		if resp.OfflineStoreConfig.S3StorageConfig != nil {
			f11f2 := &svcapitypes.S3StorageConfig{}
			if resp.OfflineStoreConfig.S3StorageConfig.KmsKeyId != nil {
				f11f2.KMSKeyID = resp.OfflineStoreConfig.S3StorageConfig.KmsKeyId
			}
			if resp.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3Uri != nil {
				f11f2.ResolvedOutputS3URI = resp.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3Uri
			}
			if resp.OfflineStoreConfig.S3StorageConfig.S3Uri != nil {
				f11f2.S3URI = resp.OfflineStoreConfig.S3StorageConfig.S3Uri
			}
			f11.S3StorageConfig = f11f2
		}
		ko.Spec.OfflineStoreConfig = f11
	} else {
		ko.Spec.OfflineStoreConfig = nil
	}
	if resp.OnlineStoreConfig != nil {
		f13 := &svcapitypes.OnlineStoreConfig{}
		if resp.OnlineStoreConfig.EnableOnlineStore != nil {
			f13.EnableOnlineStore = resp.OnlineStoreConfig.EnableOnlineStore
		}
		if resp.OnlineStoreConfig.SecurityConfig != nil {
			f13f1 := &svcapitypes.OnlineStoreSecurityConfig{}
			if resp.OnlineStoreConfig.SecurityConfig.KmsKeyId != nil {
				f13f1.KMSKeyID = resp.OnlineStoreConfig.SecurityConfig.KmsKeyId
			}
			f13.SecurityConfig = f13f1
		}
		if resp.OnlineStoreConfig.StorageType != "" {
			f13.StorageType = aws.String(string(resp.OnlineStoreConfig.StorageType))
		}
		if resp.OnlineStoreConfig.TtlDuration != nil {
			f13f3 := &svcapitypes.TTLDuration{}
			if resp.OnlineStoreConfig.TtlDuration.Unit != "" {
				f13f3.Unit = aws.String(string(resp.OnlineStoreConfig.TtlDuration.Unit))
			}
			if resp.OnlineStoreConfig.TtlDuration.Value != nil {
				valueCopy := int64(*resp.OnlineStoreConfig.TtlDuration.Value)
				f13f3.Value = &valueCopy
			}
			f13.TTLDuration = f13f3
		}
		ko.Spec.OnlineStoreConfig = f13
	} else {
		ko.Spec.OnlineStoreConfig = nil
	}
	if resp.RecordIdentifierFeatureName != nil {
		ko.Spec.RecordIdentifierFeatureName = resp.RecordIdentifierFeatureName
	} else {
		ko.Spec.RecordIdentifierFeatureName = nil
	}
	if resp.RoleArn != nil {
		ko.Spec.RoleARN = resp.RoleArn
	} else {
		ko.Spec.RoleARN = nil
	}
	if resp.ThroughputConfig != nil {
		f17 := &svcapitypes.ThroughputConfig{}
		if resp.ThroughputConfig.ProvisionedReadCapacityUnits != nil {
			provisionedReadCapacityUnitsCopy := int64(*resp.ThroughputConfig.ProvisionedReadCapacityUnits)
			f17.ProvisionedReadCapacityUnits = &provisionedReadCapacityUnitsCopy
		}
		if resp.ThroughputConfig.ProvisionedWriteCapacityUnits != nil {
			provisionedWriteCapacityUnitsCopy := int64(*resp.ThroughputConfig.ProvisionedWriteCapacityUnits)
			f17.ProvisionedWriteCapacityUnits = &provisionedWriteCapacityUnitsCopy
		}
		if resp.ThroughputConfig.ThroughputMode != "" {
			f17.ThroughputMode = aws.String(string(resp.ThroughputConfig.ThroughputMode))
		}
		ko.Spec.ThroughputConfig = f17
	} else {
		ko.Spec.ThroughputConfig = nil
	}

	rm.setStatusDefaults(ko)
	rm.customSetOutput(&resource{ko})
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.FeatureGroupName == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeFeatureGroupInput, error) {
	res := &svcsdk.DescribeFeatureGroupInput{}

	if r.ko.Spec.FeatureGroupName != nil {
		res.FeatureGroupName = r.ko.Spec.FeatureGroupName
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

	var resp *svcsdk.CreateFeatureGroupOutput
	_ = resp
	resp, err = rm.sdkapi.CreateFeatureGroup(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateFeatureGroup", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.FeatureGroupArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.FeatureGroupArn)
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
) (*svcsdk.CreateFeatureGroupInput, error) {
	res := &svcsdk.CreateFeatureGroupInput{}

	if r.ko.Spec.Description != nil {
		res.Description = r.ko.Spec.Description
	}
	if r.ko.Spec.EventTimeFeatureName != nil {
		res.EventTimeFeatureName = r.ko.Spec.EventTimeFeatureName
	}
	if r.ko.Spec.FeatureDefinitions != nil {
		f2 := []svcsdktypes.FeatureDefinition{}
		for _, f2iter := range r.ko.Spec.FeatureDefinitions {
			f2elem := &svcsdktypes.FeatureDefinition{}
			if f2iter.CollectionConfig != nil {
				var f2elemf0 svcsdktypes.CollectionConfig
				isInterfaceSet := false
				if f2iter.CollectionConfig.VectorConfig != nil {
					if isInterfaceSet {
						return nil, ackerr.NewTerminalError(fmt.Errorf("can only set one of the members for VectorConfig"))
					}
					f2elemf0f0Parent := &svcsdktypes.CollectionConfigMemberVectorConfig{}
					f2elemf0f0 := &svcsdktypes.VectorConfig{}
					if f2iter.CollectionConfig.VectorConfig.Dimension != nil {
						dimensionCopy0 := *f2iter.CollectionConfig.VectorConfig.Dimension
						if dimensionCopy0 > math.MaxInt32 || dimensionCopy0 < math.MinInt32 {
							return nil, fmt.Errorf("error: field Dimension is of type int32")
						}
						dimensionCopy := int32(dimensionCopy0)
						f2elemf0f0.Dimension = &dimensionCopy
					}
					f2elemf0f0Parent.Value = *f2elemf0f0
				}
				f2elem.CollectionConfig = f2elemf0
			}
			if f2iter.CollectionType != nil {
				f2elem.CollectionType = svcsdktypes.CollectionType(*f2iter.CollectionType)
			}
			if f2iter.FeatureName != nil {
				f2elem.FeatureName = f2iter.FeatureName
			}
			if f2iter.FeatureType != nil {
				f2elem.FeatureType = svcsdktypes.FeatureType(*f2iter.FeatureType)
			}
			f2 = append(f2, *f2elem)
		}
		res.FeatureDefinitions = f2
	}
	if r.ko.Spec.FeatureGroupName != nil {
		res.FeatureGroupName = r.ko.Spec.FeatureGroupName
	}
	if r.ko.Spec.OfflineStoreConfig != nil {
		f4 := &svcsdktypes.OfflineStoreConfig{}
		if r.ko.Spec.OfflineStoreConfig.DataCatalogConfig != nil {
			f4f0 := &svcsdktypes.DataCatalogConfig{}
			if r.ko.Spec.OfflineStoreConfig.DataCatalogConfig.Catalog != nil {
				f4f0.Catalog = r.ko.Spec.OfflineStoreConfig.DataCatalogConfig.Catalog
			}
			if r.ko.Spec.OfflineStoreConfig.DataCatalogConfig.Database != nil {
				f4f0.Database = r.ko.Spec.OfflineStoreConfig.DataCatalogConfig.Database
			}
			if r.ko.Spec.OfflineStoreConfig.DataCatalogConfig.TableName != nil {
				f4f0.TableName = r.ko.Spec.OfflineStoreConfig.DataCatalogConfig.TableName
			}
			f4.DataCatalogConfig = f4f0
		}
		if r.ko.Spec.OfflineStoreConfig.DisableGlueTableCreation != nil {
			f4.DisableGlueTableCreation = r.ko.Spec.OfflineStoreConfig.DisableGlueTableCreation
		}
		if r.ko.Spec.OfflineStoreConfig.S3StorageConfig != nil {
			f4f2 := &svcsdktypes.S3StorageConfig{}
			if r.ko.Spec.OfflineStoreConfig.S3StorageConfig.KMSKeyID != nil {
				f4f2.KmsKeyId = r.ko.Spec.OfflineStoreConfig.S3StorageConfig.KMSKeyID
			}
			if r.ko.Spec.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3URI != nil {
				f4f2.ResolvedOutputS3Uri = r.ko.Spec.OfflineStoreConfig.S3StorageConfig.ResolvedOutputS3URI
			}
			if r.ko.Spec.OfflineStoreConfig.S3StorageConfig.S3URI != nil {
				f4f2.S3Uri = r.ko.Spec.OfflineStoreConfig.S3StorageConfig.S3URI
			}
			f4.S3StorageConfig = f4f2
		}
		res.OfflineStoreConfig = f4
	}
	if r.ko.Spec.OnlineStoreConfig != nil {
		f5 := &svcsdktypes.OnlineStoreConfig{}
		if r.ko.Spec.OnlineStoreConfig.EnableOnlineStore != nil {
			f5.EnableOnlineStore = r.ko.Spec.OnlineStoreConfig.EnableOnlineStore
		}
		if r.ko.Spec.OnlineStoreConfig.SecurityConfig != nil {
			f5f1 := &svcsdktypes.OnlineStoreSecurityConfig{}
			if r.ko.Spec.OnlineStoreConfig.SecurityConfig.KMSKeyID != nil {
				f5f1.KmsKeyId = r.ko.Spec.OnlineStoreConfig.SecurityConfig.KMSKeyID
			}
			f5.SecurityConfig = f5f1
		}
		if r.ko.Spec.OnlineStoreConfig.StorageType != nil {
			f5.StorageType = svcsdktypes.StorageType(*r.ko.Spec.OnlineStoreConfig.StorageType)
		}
		if r.ko.Spec.OnlineStoreConfig.TTLDuration != nil {
			f5f3 := &svcsdktypes.TtlDuration{}
			if r.ko.Spec.OnlineStoreConfig.TTLDuration.Unit != nil {
				f5f3.Unit = svcsdktypes.TtlDurationUnit(*r.ko.Spec.OnlineStoreConfig.TTLDuration.Unit)
			}
			if r.ko.Spec.OnlineStoreConfig.TTLDuration.Value != nil {
				valueCopy0 := *r.ko.Spec.OnlineStoreConfig.TTLDuration.Value
				if valueCopy0 > math.MaxInt32 || valueCopy0 < math.MinInt32 {
					return nil, fmt.Errorf("error: field Value is of type int32")
				}
				valueCopy := int32(valueCopy0)
				f5f3.Value = &valueCopy
			}
			f5.TtlDuration = f5f3
		}
		res.OnlineStoreConfig = f5
	}
	if r.ko.Spec.RecordIdentifierFeatureName != nil {
		res.RecordIdentifierFeatureName = r.ko.Spec.RecordIdentifierFeatureName
	}
	if r.ko.Spec.RoleARN != nil {
		res.RoleArn = r.ko.Spec.RoleARN
	}
	if r.ko.Spec.Tags != nil {
		f8 := []svcsdktypes.Tag{}
		for _, f8iter := range r.ko.Spec.Tags {
			f8elem := &svcsdktypes.Tag{}
			if f8iter.Key != nil {
				f8elem.Key = f8iter.Key
			}
			if f8iter.Value != nil {
				f8elem.Value = f8iter.Value
			}
			f8 = append(f8, *f8elem)
		}
		res.Tags = f8
	}
	if r.ko.Spec.ThroughputConfig != nil {
		f9 := &svcsdktypes.ThroughputConfig{}
		if r.ko.Spec.ThroughputConfig.ProvisionedReadCapacityUnits != nil {
			provisionedReadCapacityUnitsCopy0 := *r.ko.Spec.ThroughputConfig.ProvisionedReadCapacityUnits
			if provisionedReadCapacityUnitsCopy0 > math.MaxInt32 || provisionedReadCapacityUnitsCopy0 < math.MinInt32 {
				return nil, fmt.Errorf("error: field ProvisionedReadCapacityUnits is of type int32")
			}
			provisionedReadCapacityUnitsCopy := int32(provisionedReadCapacityUnitsCopy0)
			f9.ProvisionedReadCapacityUnits = &provisionedReadCapacityUnitsCopy
		}
		if r.ko.Spec.ThroughputConfig.ProvisionedWriteCapacityUnits != nil {
			provisionedWriteCapacityUnitsCopy0 := *r.ko.Spec.ThroughputConfig.ProvisionedWriteCapacityUnits
			if provisionedWriteCapacityUnitsCopy0 > math.MaxInt32 || provisionedWriteCapacityUnitsCopy0 < math.MinInt32 {
				return nil, fmt.Errorf("error: field ProvisionedWriteCapacityUnits is of type int32")
			}
			provisionedWriteCapacityUnitsCopy := int32(provisionedWriteCapacityUnitsCopy0)
			f9.ProvisionedWriteCapacityUnits = &provisionedWriteCapacityUnitsCopy
		}
		if r.ko.Spec.ThroughputConfig.ThroughputMode != nil {
			f9.ThroughputMode = svcsdktypes.ThroughputMode(*r.ko.Spec.ThroughputConfig.ThroughputMode)
		}
		res.ThroughputConfig = f9
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
	return nil, ackerr.NewTerminalError(ackerr.NotImplemented)
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
	if err = rm.requeueUntilCanModify(ctx, r); err != nil {
		return r, err
	}

	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteFeatureGroupOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteFeatureGroup(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteFeatureGroup", err)

	if err == nil {
		if observed, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
			if err != nil {
				return nil, err
			}
			r.SetStatus(observed)
			return r, requeueWaitWhileDeleting
		}
	}

	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteFeatureGroupInput, error) {
	res := &svcsdk.DeleteFeatureGroupInput{}

	if r.ko.Spec.FeatureGroupName != nil {
		res.FeatureGroupName = r.ko.Spec.FeatureGroupName
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.FeatureGroup,
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
	// custom update conditions
	customUpdate := rm.CustomUpdateConditions(ko, r, err)
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil || customUpdate {
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

	var terminalErr smithy.APIError
	if !errors.As(err, &terminalErr) {
		return false
	}
	switch terminalErr.ErrorCode() {
	case "ResourceNotFound",
		"ResourceInUse",
		"InvalidParameterCombination",
		"InvalidParameterValue",
		"MissingParameter":
		return true
	default:
		return false
	}
}
