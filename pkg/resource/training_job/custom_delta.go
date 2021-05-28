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

func customSetDefaults(
	a *resource,
	b *resource,
) {
	if a.ko.Spec.EnableInterContainerTrafficEncryption == nil && b.ko.Spec.EnableInterContainerTrafficEncryption != nil {
		a.ko.Spec.EnableInterContainerTrafficEncryption = b.ko.Spec.EnableInterContainerTrafficEncryption
	}

	if a.ko.Spec.EnableManagedSpotTraining == nil && b.ko.Spec.EnableManagedSpotTraining != nil {
		a.ko.Spec.EnableManagedSpotTraining = b.ko.Spec.EnableManagedSpotTraining
	}

	if a.ko.Spec.EnableNetworkIsolation == nil && b.ko.Spec.EnableNetworkIsolation != nil {
		a.ko.Spec.EnableNetworkIsolation = b.ko.Spec.EnableNetworkIsolation
	}

	if a.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries == nil && b.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries != nil {
		a.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries = b.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries
	}

	if a.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries == nil && b.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries != nil {
		a.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries = b.ko.Spec.AlgorithmSpecification.EnableSageMakerMetricsTimeSeries
	}

	if a.ko.Spec.OutputDataConfig.KMSKeyID == nil && b.ko.Spec.OutputDataConfig.KMSKeyID != nil {
		a.ko.Spec.OutputDataConfig.KMSKeyID = b.ko.Spec.OutputDataConfig.KMSKeyID
	}
}
