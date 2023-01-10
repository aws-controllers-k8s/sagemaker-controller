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

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FeatureGroupSpec defines the desired state of FeatureGroup.
//
// Amazon SageMaker Feature Store stores features in a collection called Feature
// Group. A Feature Group can be visualized as a table which has rows, with
// a unique identifier for each row where each column in the table is a feature.
// In principle, a Feature Group is composed of features and values per features.
type FeatureGroupSpec struct {

	// A free-form description of a FeatureGroup.
	Description *string `json:"description,omitempty"`
	// The name of the feature that stores the EventTime of a Record in a FeatureGroup.
	//
	// An EventTime is a point in time when a new event occurs that corresponds
	// to the creation or update of a Record in a FeatureGroup. All Records in the
	// FeatureGroup must have a corresponding EventTime.
	//
	// An EventTime can be a String or Fractional.
	//
	//    * Fractional: EventTime feature values must be a Unix timestamp in seconds.
	//
	//    * String: EventTime feature values must be an ISO-8601 string in the format.
	//    The following formats are supported yyyy-MM-dd'T'HH:mm:ssZ and yyyy-MM-dd'T'HH:mm:ss.SSSZ
	//    where yyyy, MM, and dd represent the year, month, and day respectively
	//    and HH, mm, ss, and if applicable, SSS represent the hour, month, second
	//    and milliseconds respsectively. 'T' and Z are constants.
	// +kubebuilder:validation:Required
	EventTimeFeatureName *string `json:"eventTimeFeatureName"`
	// A list of Feature names and types. Name and Type is compulsory per Feature.
	//
	// Valid feature FeatureTypes are Integral, Fractional and String.
	//
	// FeatureNames cannot be any of the following: is_deleted, write_time, api_invocation_time
	//
	// You can create up to 2,500 FeatureDefinitions per FeatureGroup.
	// +kubebuilder:validation:Required
	FeatureDefinitions []*FeatureDefinition `json:"featureDefinitions"`
	// The name of the FeatureGroup. The name must be unique within an Amazon Web
	// Services Region in an Amazon Web Services account. The name:
	//
	//    * Must start and end with an alphanumeric character.
	//
	//    * Can only contain alphanumeric character and hyphens. Spaces are not
	//    allowed.
	// +kubebuilder:validation:Required
	FeatureGroupName *string `json:"featureGroupName"`
	// Use this to configure an OfflineFeatureStore. This parameter allows you to
	// specify:
	//
	//    * The Amazon Simple Storage Service (Amazon S3) location of an OfflineStore.
	//
	//    * A configuration for an Amazon Web Services Glue or Amazon Web Services
	//    Hive data catalog.
	//
	//    * An KMS encryption key to encrypt the Amazon S3 location used for OfflineStore.
	//    If KMS encryption key is not specified, by default we encrypt all data
	//    at rest using Amazon Web Services KMS key. By defining your bucket-level
	//    key (https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucket-key.html)
	//    for SSE, you can reduce Amazon Web Services KMS requests costs by up to
	//    99 percent.
	//
	// To learn more about this parameter, see OfflineStoreConfig.
	OfflineStoreConfig *OfflineStoreConfig `json:"offlineStoreConfig,omitempty"`
	// You can turn the OnlineStore on or off by specifying True for the EnableOnlineStore
	// flag in OnlineStoreConfig; the default value is False.
	//
	// You can also include an Amazon Web Services KMS key ID (KMSKeyId) for at-rest
	// encryption of the OnlineStore.
	OnlineStoreConfig *OnlineStoreConfig `json:"onlineStoreConfig,omitempty"`
	// The name of the Feature whose value uniquely identifies a Record defined
	// in the FeatureStore. Only the latest record per identifier value will be
	// stored in the OnlineStore. RecordIdentifierFeatureName must be one of feature
	// definitions' names.
	//
	// You use the RecordIdentifierFeatureName to access data in a FeatureStore.
	//
	// This name:
	//
	//    * Must start and end with an alphanumeric character.
	//
	//    * Can only contains alphanumeric characters, hyphens, underscores. Spaces
	//    are not allowed.
	// +kubebuilder:validation:Required
	RecordIdentifierFeatureName *string `json:"recordIdentifierFeatureName"`
	// The Amazon Resource Name (ARN) of the IAM execution role used to persist
	// data into the OfflineStore if an OfflineStoreConfig is provided.
	RoleARN *string `json:"roleARN,omitempty"`
	// Tags used to identify Features in each FeatureGroup.
	Tags []*Tag `json:"tags,omitempty"`
}

// FeatureGroupStatus defines the observed state of FeatureGroup
type FeatureGroupStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// The reason that the FeatureGroup failed to be replicated in the OfflineStore.
	// This is failure can occur because:
	//
	//    * The FeatureGroup could not be created in the OfflineStore.
	//
	//    * The FeatureGroup could not be deleted from the OfflineStore.
	// +kubebuilder:validation:Optional
	FailureReason *string `json:"failureReason,omitempty"`
	// The status of the feature group.
	// +kubebuilder:validation:Optional
	FeatureGroupStatus *string `json:"featureGroupStatus,omitempty"`
}

// FeatureGroup is the Schema for the FeatureGroups API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FAILURE-REASON",type=string,priority=1,JSONPath=`.status.failureReason`
// +kubebuilder:printcolumn:name="STATUS",type=string,priority=0,JSONPath=`.status.featureGroupStatus`
type FeatureGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FeatureGroupSpec   `json:"spec,omitempty"`
	Status            FeatureGroupStatus `json:"status,omitempty"`
}

// FeatureGroupList contains a list of FeatureGroup
// +kubebuilder:object:root=true
type FeatureGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FeatureGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FeatureGroup{}, &FeatureGroupList{})
}
