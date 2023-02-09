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

package pipeline_execution

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
	_ = &acktags.Tags{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}

	if ackcompare.HasNilDifference(a.ko.Spec.ParallelismConfiguration, b.ko.Spec.ParallelismConfiguration) {
		delta.Add("Spec.ParallelismConfiguration", a.ko.Spec.ParallelismConfiguration, b.ko.Spec.ParallelismConfiguration)
	} else if a.ko.Spec.ParallelismConfiguration != nil && b.ko.Spec.ParallelismConfiguration != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps, b.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps) {
			delta.Add("Spec.ParallelismConfiguration.MaxParallelExecutionSteps", a.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps, b.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps)
		} else if a.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps != nil && b.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps != nil {
			if *a.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps != *b.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps {
				delta.Add("Spec.ParallelismConfiguration.MaxParallelExecutionSteps", a.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps, b.ko.Spec.ParallelismConfiguration.MaxParallelExecutionSteps)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.PipelineExecutionDescription, b.ko.Spec.PipelineExecutionDescription) {
		delta.Add("Spec.PipelineExecutionDescription", a.ko.Spec.PipelineExecutionDescription, b.ko.Spec.PipelineExecutionDescription)
	} else if a.ko.Spec.PipelineExecutionDescription != nil && b.ko.Spec.PipelineExecutionDescription != nil {
		if *a.ko.Spec.PipelineExecutionDescription != *b.ko.Spec.PipelineExecutionDescription {
			delta.Add("Spec.PipelineExecutionDescription", a.ko.Spec.PipelineExecutionDescription, b.ko.Spec.PipelineExecutionDescription)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.PipelineExecutionDisplayName, b.ko.Spec.PipelineExecutionDisplayName) {
		delta.Add("Spec.PipelineExecutionDisplayName", a.ko.Spec.PipelineExecutionDisplayName, b.ko.Spec.PipelineExecutionDisplayName)
	} else if a.ko.Spec.PipelineExecutionDisplayName != nil && b.ko.Spec.PipelineExecutionDisplayName != nil {
		if *a.ko.Spec.PipelineExecutionDisplayName != *b.ko.Spec.PipelineExecutionDisplayName {
			delta.Add("Spec.PipelineExecutionDisplayName", a.ko.Spec.PipelineExecutionDisplayName, b.ko.Spec.PipelineExecutionDisplayName)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.PipelineName, b.ko.Spec.PipelineName) {
		delta.Add("Spec.PipelineName", a.ko.Spec.PipelineName, b.ko.Spec.PipelineName)
	} else if a.ko.Spec.PipelineName != nil && b.ko.Spec.PipelineName != nil {
		if *a.ko.Spec.PipelineName != *b.ko.Spec.PipelineName {
			delta.Add("Spec.PipelineName", a.ko.Spec.PipelineName, b.ko.Spec.PipelineName)
		}
	}
	if !reflect.DeepEqual(a.ko.Spec.PipelineParameters, b.ko.Spec.PipelineParameters) {
		delta.Add("Spec.PipelineParameters", a.ko.Spec.PipelineParameters, b.ko.Spec.PipelineParameters)
	}

	return delta
}
