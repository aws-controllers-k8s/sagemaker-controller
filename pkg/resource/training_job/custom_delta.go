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

package training_job

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	if ackcompare.IsNil(a.ko.Spec.EnableInterContainerTrafficEncryption) && ackcompare.IsNotNil(b.ko.Spec.EnableInterContainerTrafficEncryption) {
		a.ko.Spec.EnableInterContainerTrafficEncryption = b.ko.Spec.EnableInterContainerTrafficEncryption
	}

	if ackcompare.IsNil(a.ko.Spec.EnableManagedSpotTraining) && ackcompare.IsNotNil(b.ko.Spec.EnableManagedSpotTraining) {
		a.ko.Spec.EnableManagedSpotTraining = b.ko.Spec.EnableManagedSpotTraining
	}

	if ackcompare.IsNil(a.ko.Spec.EnableNetworkIsolation) && ackcompare.IsNotNil(b.ko.Spec.EnableNetworkIsolation) {
		a.ko.Spec.EnableNetworkIsolation = b.ko.Spec.EnableNetworkIsolation
	}

	if ackcompare.IsNil(a.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries) && ackcompare.IsNotNil(b.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries) {
		a.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries = b.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries
	}

	// The KMS Key is an empty string by default but cannot be nil.
	if ackcompare.IsNil(a.ko.Spec.OutputDataConfig.KMSKeyID) && ackcompare.IsNotNil(b.ko.Spec.OutputDataConfig.KMSKeyID) {
		a.ko.Spec.OutputDataConfig.KMSKeyID = b.ko.Spec.OutputDataConfig.KMSKeyID
	}
}
