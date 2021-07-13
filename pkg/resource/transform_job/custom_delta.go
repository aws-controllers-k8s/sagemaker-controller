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

package transform_job

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// TransformInput is a required field.
	if ackcompare.IsNotNil(a.ko.Spec.TransformInput) && ackcompare.IsNotNil(b.ko.Spec.TransformInput) {
	   	if ackcompare.IsNil(a.ko.Spec.TransformInput.CompressionType) && ackcompare.IsNotNil(b.ko.Spec.TransformInput.CompressionType) {
		   a.ko.Spec.TransformInput.CompressionType = b.ko.Spec.TransformInput.CompressionType
		}
	}

	if ackcompare.IsNotNil(a.ko.Spec.TransformInput) && ackcompare.IsNotNil(b.ko.Spec.TransformInput) {
	   	if ackcompare.IsNil(a.ko.Spec.TransformInput.SplitType) && ackcompare.IsNotNil(b.ko.Spec.TransformInput.SplitType) {
		   a.ko.Spec.TransformInput.SplitType = b.ko.Spec.TransformInput.SplitType
		}
	}

	// DataProcessing is not a required field, so first create it.
	if ackcompare.IsNil(a.ko.Spec.DataProcessing) && ackcompare.IsNotNil(b.ko.Spec.DataProcessing) {
		a.ko.Spec.DataProcessing = &svcapitypes.DataProcessing{}
	}

	if ackcompare.IsNotNil(a.ko.Spec.DataProcessing) && ackcompare.IsNotNil(b.ko.Spec.DataProcessing) {
		if ackcompare.IsNil(a.ko.Spec.DataProcessing.InputFilter) && ackcompare.IsNotNil(b.ko.Spec.DataProcessing.InputFilter) {
			a.ko.Spec.DataProcessing.InputFilter = b.ko.Spec.DataProcessing.InputFilter
		}
		if ackcompare.IsNil(a.ko.Spec.DataProcessing.JoinSource) && ackcompare.IsNotNil(b.ko.Spec.DataProcessing.JoinSource) {
			a.ko.Spec.DataProcessing.JoinSource = b.ko.Spec.DataProcessing.JoinSource
		}
		if ackcompare.IsNil(a.ko.Spec.DataProcessing.OutputFilter) && ackcompare.IsNotNil(b.ko.Spec.DataProcessing.OutputFilter) {
			a.ko.Spec.DataProcessing.OutputFilter = b.ko.Spec.DataProcessing.OutputFilter
		}
	}

	// TODO: TransformOutput is a required field, so this check should not be required
	if ackcompare.IsNil(a.ko.Spec.TransformOutput) && ackcompare.IsNotNil(b.ko.Spec.TransformOutput) {
		a.ko.Spec.TransformOutput = &svcapitypes.TransformOutput{}
	}

	if ackcompare.IsNotNil(a.ko.Spec.TransformOutput) && ackcompare.IsNotNil(b.ko.Spec.TransformOutput) {
		if ackcompare.IsNil(a.ko.Spec.TransformOutput.AssembleWith) && ackcompare.IsNotNil(b.ko.Spec.TransformOutput.AssembleWith) {
			a.ko.Spec.TransformOutput.AssembleWith = b.ko.Spec.TransformOutput.AssembleWith
		}
		if ackcompare.IsNil(a.ko.Spec.TransformOutput.KMSKeyID) && ackcompare.IsNotNil(b.ko.Spec.TransformOutput.KMSKeyID) {
			a.ko.Spec.TransformOutput.KMSKeyID = b.ko.Spec.TransformOutput.KMSKeyID
		}
	}

}
