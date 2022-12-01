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

package hyper_parameter_tuning_job

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	if ackcompare.IsNotNil(a.ko.Spec.TrainingJobDefinition) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinition) {
		// SageMaker adds StaticHyperParameters prefixed with an underscore. We must ignore these when comparing.
		latestStaticHyperParameters := b.ko.Spec.TrainingJobDefinition.StaticHyperParameters
		if ackcompare.IsNotNil(latestStaticHyperParameters) {
			for key, _ := range latestStaticHyperParameters {
				if key[0:1] == "_" {
					delete(b.ko.Spec.TrainingJobDefinition.StaticHyperParameters, key)
				}
			}
		}
	}

	// TODO: Use late initialize instead once code generator supports late initializing slices.
	if ackcompare.IsNotNil(a.ko.Spec.TrainingJobDefinitions) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions) {
		if len(a.ko.Spec.TrainingJobDefinitions) == len(b.ko.Spec.TrainingJobDefinitions) {
			for index := range a.ko.Spec.TrainingJobDefinitions {
				latestStaticHyperParameters := b.ko.Spec.TrainingJobDefinitions[index].StaticHyperParameters
				if ackcompare.IsNotNil(latestStaticHyperParameters) {
					for key, _ := range latestStaticHyperParameters {
						if key[0:1] == "_" {
							delete(b.ko.Spec.TrainingJobDefinitions[index].StaticHyperParameters, key)
						}
					}
				}
				if ackcompare.IsNotNil(a.ko.Spec.TrainingJobDefinitions[index].AlgorithmSpecification) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].AlgorithmSpecification) {
					if ackcompare.IsNil(a.ko.Spec.TrainingJobDefinitions[index].AlgorithmSpecification.MetricDefinitions) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].AlgorithmSpecification.MetricDefinitions) {
						a.ko.Spec.TrainingJobDefinitions[index].AlgorithmSpecification.MetricDefinitions = b.ko.Spec.TrainingJobDefinitions[index].AlgorithmSpecification.MetricDefinitions
					}
				}
				if ackcompare.IsNil(a.ko.Spec.TrainingJobDefinitions[index].EnableInterContainerTrafficEncryption) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].EnableInterContainerTrafficEncryption) {
					a.ko.Spec.TrainingJobDefinitions[index].EnableInterContainerTrafficEncryption = b.ko.Spec.TrainingJobDefinitions[index].EnableInterContainerTrafficEncryption
				}
				if ackcompare.IsNil(a.ko.Spec.TrainingJobDefinitions[index].EnableManagedSpotTraining) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].EnableManagedSpotTraining) {
					a.ko.Spec.TrainingJobDefinitions[index].EnableManagedSpotTraining = b.ko.Spec.TrainingJobDefinitions[index].EnableManagedSpotTraining
				}
				if ackcompare.IsNil(a.ko.Spec.TrainingJobDefinitions[index].EnableNetworkIsolation) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].EnableNetworkIsolation) {
					a.ko.Spec.TrainingJobDefinitions[index].EnableNetworkIsolation = b.ko.Spec.TrainingJobDefinitions[index].EnableNetworkIsolation
				}
				if ackcompare.IsNotNil(a.ko.Spec.TrainingJobDefinitions[index].ResourceConfig) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].ResourceConfig) {
					if ackcompare.IsNil(a.ko.Spec.TrainingJobDefinitions[index].ResourceConfig.InstanceCount) && ackcompare.IsNotNil(b.ko.Spec.TrainingJobDefinitions[index].ResourceConfig.InstanceCount) {
						a.ko.Spec.TrainingJobDefinitions[index].ResourceConfig.InstanceCount = b.ko.Spec.TrainingJobDefinitions[index].ResourceConfig.InstanceCount
					}
				}
			}
		}
	}

}
