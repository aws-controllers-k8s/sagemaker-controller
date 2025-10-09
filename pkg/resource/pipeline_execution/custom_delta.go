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

package pipeline_execution

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// customSetDefaults sets the default fields for DirectInternetAccess and RootAccess.
func customSetDefaults(
	a *resource,
	b *resource,
) {
	// PipelineVersion ID starts at 1.
	DefaultPipelineVersionID := aws.Int64(1)
	if ackcompare.IsNil(a.ko.Spec.PipelineVersionID) && ackcompare.IsNotNil(b.ko.Spec.PipelineVersionID) {
		a.ko.Spec.PipelineVersionID = DefaultPipelineVersionID
	}
}
