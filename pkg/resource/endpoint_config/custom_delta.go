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

package endpoint_config

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// The Sagemaker Service returns a value for VolumeSizeInGB for a Serverless endpoint, despite the user not being able
	// to specify it.
	// TODO: Use Late Initialization instead whenever the code generator supports it for slices/arrays.

	if ackcompare.IsNotNil(a.ko.Spec.ProductionVariants) && ackcompare.IsNotNil(a.ko.Spec.ProductionVariants) {
		if len(a.ko.Spec.ProductionVariants) == len(b.ko.Spec.ProductionVariants) {
			for i, _ := range a.ko.Spec.ProductionVariants {
				if a.ko.Spec.ProductionVariants[i].ServerlessConfig != nil && a.ko.Spec.ProductionVariants[i].VolumeSizeInGB == nil &&
					b.ko.Spec.ProductionVariants[i].VolumeSizeInGB != nil {
					a.ko.Spec.ProductionVariants[i].VolumeSizeInGB = b.ko.Spec.ProductionVariants[i].VolumeSizeInGB
				}
			}
		}
	}
}
