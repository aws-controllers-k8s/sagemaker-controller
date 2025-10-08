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
	"github.com/aws/aws-sdk-go-v2/aws"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// Default is for ImageDigest and ModelDataETag to be generated automatically by Sagemaker if not specified
	if ackcompare.IsNotNil(a.ko.Spec.InferenceSpecification) && ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification) {
		if ackcompare.IsNotNil(a.ko.Spec.InferenceSpecification.Containers) && ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification.Containers) {
			for index := range a.ko.Spec.InferenceSpecification.Containers {

				// Set default for ImageDigest
				if ackcompare.IsNil(a.ko.Spec.InferenceSpecification.Containers[index].ImageDigest) &&
					ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification.Containers[index].ImageDigest) {
					a.ko.Spec.InferenceSpecification.Containers[index].ImageDigest =
						b.ko.Spec.InferenceSpecification.Containers[index].ImageDigest
				}

				// Set default for ModelDataETag
				if ackcompare.IsNil(a.ko.Spec.InferenceSpecification.Containers[index].ModelDataETag) &&
					ackcompare.IsNotNil(b.ko.Spec.InferenceSpecification.Containers[index].ModelDataETag) {
					a.ko.Spec.InferenceSpecification.Containers[index].ModelDataETag =
						b.ko.Spec.InferenceSpecification.Containers[index].ModelDataETag
				}
			}
		}
	}
	// Default is for KMSKeyID to be ""
	defaultKMSKeyID := aws.String("")

	if ackcompare.IsNotNil(a.ko.Spec.ValidationSpecification) && ackcompare.IsNotNil(b.ko.Spec.ValidationSpecification) {
		for index := range a.ko.Spec.ValidationSpecification.ValidationProfiles {
			if ackcompare.IsNil(a.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID) &&
				ackcompare.IsNotNil(b.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID) {
				a.ko.Spec.ValidationSpecification.ValidationProfiles[index].TransformJobDefinition.TransformOutput.KMSKeyID = defaultKMSKeyID
			}
		}
	}
}
