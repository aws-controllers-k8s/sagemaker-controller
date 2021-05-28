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
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	if a.ko.Spec.TransformInput.CompressionType == nil && b.ko.Spec.TransformInput.CompressionType != nil {
		a.ko.Spec.TransformInput.CompressionType = b.ko.Spec.TransformInput.CompressionType
	}

	if a.ko.Spec.TransformInput.SplitType == nil && b.ko.Spec.TransformInput.SplitType != nil {
		a.ko.Spec.TransformInput.SplitType = b.ko.Spec.TransformInput.SplitType
	}

	if a.ko.Spec.DataProcessing == nil && b.ko.Spec.DataProcessing != nil {
		a.ko.Spec.DataProcessing = &svcapitypes.DataProcessing{}
	}
	if a.ko.Spec.DataProcessing.InputFilter == nil && b.ko.Spec.DataProcessing.InputFilter != nil {
		a.ko.Spec.DataProcessing.InputFilter = b.ko.Spec.DataProcessing.InputFilter
	}
	if a.ko.Spec.DataProcessing.JoinSource == nil && b.ko.Spec.DataProcessing.JoinSource != nil {
		a.ko.Spec.DataProcessing.JoinSource = b.ko.Spec.DataProcessing.JoinSource
	}
	if a.ko.Spec.DataProcessing.OutputFilter == nil && b.ko.Spec.DataProcessing.OutputFilter != nil {
		a.ko.Spec.DataProcessing.OutputFilter = b.ko.Spec.DataProcessing.OutputFilter
	}

	if a.ko.Spec.TransformOutput == nil && b.ko.Spec.TransformOutput != nil{
		a.ko.Spec.TransformOutput = &svcapitypes.TransformOutput{}
	}
	if a.ko.Spec.TransformOutput.AssembleWith == nil && b.ko.Spec.TransformOutput.AssembleWith != nil {
		a.ko.Spec.TransformOutput.AssembleWith = b.ko.Spec.TransformOutput.AssembleWith
	}
	if a.ko.Spec.TransformOutput.KMSKeyID == nil && b.ko.Spec.TransformOutput.KMSKeyID != nil {
		a.ko.Spec.TransformOutput.KMSKeyID = b.ko.Spec.TransformOutput.KMSKeyID
	}


}
