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

package model_package_group

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
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

	if ackcompare.HasNilDifference(a.ko.Spec.ModelPackageGroupDescription, b.ko.Spec.ModelPackageGroupDescription) {
		delta.Add("Spec.ModelPackageGroupDescription", a.ko.Spec.ModelPackageGroupDescription, b.ko.Spec.ModelPackageGroupDescription)
	} else if a.ko.Spec.ModelPackageGroupDescription != nil && b.ko.Spec.ModelPackageGroupDescription != nil {
		if *a.ko.Spec.ModelPackageGroupDescription != *b.ko.Spec.ModelPackageGroupDescription {
			delta.Add("Spec.ModelPackageGroupDescription", a.ko.Spec.ModelPackageGroupDescription, b.ko.Spec.ModelPackageGroupDescription)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelPackageGroupName, b.ko.Spec.ModelPackageGroupName) {
		delta.Add("Spec.ModelPackageGroupName", a.ko.Spec.ModelPackageGroupName, b.ko.Spec.ModelPackageGroupName)
	} else if a.ko.Spec.ModelPackageGroupName != nil && b.ko.Spec.ModelPackageGroupName != nil {
		if *a.ko.Spec.ModelPackageGroupName != *b.ko.Spec.ModelPackageGroupName {
			delta.Add("Spec.ModelPackageGroupName", a.ko.Spec.ModelPackageGroupName, b.ko.Spec.ModelPackageGroupName)
		}
	}
	if !reflect.DeepEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}

	return delta
}
