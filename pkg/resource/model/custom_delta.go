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

package model

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// PrimaryContainer is not a required field, so this check is required
	if ackcompare.IsNil(a.ko.Spec.PrimaryContainer) && ackcompare.IsNotNil(b.ko.Spec.PrimaryContainer) {
		a.ko.Spec.PrimaryContainer = &svcapitypes.ContainerDefinition{}
	}

	if ackcompare.IsNotNil(a.ko.Spec.PrimaryContainer) && ackcompare.IsNotNil(b.ko.Spec.PrimaryContainer) {
		if ackcompare.IsNil(a.ko.Spec.PrimaryContainer.Mode) && ackcompare.IsNotNil(b.ko.Spec.PrimaryContainer.Mode) {
			a.ko.Spec.PrimaryContainer.Mode = b.ko.Spec.PrimaryContainer.Mode
		}
	}

	// Default value of Mode is SingleModel
	mode := aws.String("SingleModel")

	if ackcompare.IsNotNil(a.ko.Spec.Containers) && ackcompare.IsNotNil(b.ko.Spec.Containers) {
		for index := range a.ko.Spec.Containers {
			if ackcompare.IsNil(a.ko.Spec.Containers[index].Mode) && ackcompare.IsNotNil(b.ko.Spec.Containers[index].Mode) {
				a.ko.Spec.Containers[index].Mode = mode
			}
			if ackcompare.IsNotNil(a.ko.Spec.Containers[index].ImageConfig) && ackcompare.IsNotNil(b.ko.Spec.Containers[index].ImageConfig) {
				if ackcompare.IsNil(a.ko.Spec.Containers[index].ImageConfig.RepositoryAuthConfig) && ackcompare.IsNotNil(b.ko.Spec.Containers[index].ImageConfig.RepositoryAuthConfig) {
					a.ko.Spec.Containers[index].ImageConfig.RepositoryAuthConfig = &svcapitypes.RepositoryAuthConfig{}
				}
			}

		}
	}
}
