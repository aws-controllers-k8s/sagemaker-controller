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

package model_package

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// Default is for CertifyForMarketplace to be set to false
	if ackcompare.IsNil(a.ko.Spec.CertifyForMarketplace) && ackcompare.IsNotNil(b.ko.Spec.CertifyForMarketplace) {
		a.ko.Spec.CertifyForMarketplace = b.ko.Spec.CertifyForMarketplace
	}
	// Default is for ModelApprovalStatus to be set to pending manual approval
	if ackcompare.IsNil(a.ko.Spec.ModelApprovalStatus) && ackcompare.IsNotNil(b.ko.Spec.ModelApprovalStatus) {
		a.ko.Spec.ModelApprovalStatus = b.ko.Spec.ModelApprovalStatus
	}
	// Default is for ImageDigest to be generated automatically by Sagemaker if not specified
	if ackcompare.IsNotNil(a.ko.Spec.InferenceSpecification) && ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification) {
		if ackcompare.IsNotNil(a.ko.Spec.InferenceSpecification.Containers) && ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification.Containers) {
			for index := range a.ko.Spec.InferenceSpecification.Containers {
				if ackcompare.IsNil(a.ko.Spec.InferenceSpecification.Containers[index].ImageDigest) &&
					ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification.Containers[index].ImageDigest) {
					a.ko.Spec.InferenceSpecification.Containers[index].ImageDigest =
						b.ko.Spec.InferenceSpecification.Containers[index].ImageDigest
				}
			}
		}
	}
	// Default is for KMSKeyID to be ""
	if ackcompare.IsNotNil(a.ko.Spec.ValidationSpecification) && ackcompare.IsNotNil(b.ko.Spec.ValidationSpecification) {
		for index := range a.ko.Spec.ValidationSpecification.ValidationProfiles {
			if ackcompare.IsNil(a.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID) &&
				ackcompare.IsNotNil(b.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID) {
				a.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID =
					b.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID
			}
		}
	}
}
