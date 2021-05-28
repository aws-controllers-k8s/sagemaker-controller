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

func customSetDefaults(
	a *resource,
	b *resource,
) {
	// TODO: This throws a nil pointer error
	// if a.ko.Spec.PrimaryContainer.Mode == nil && b.ko.Spec.PrimaryContainer.Mode != nil {
	// 	a.ko.Spec.PrimaryContainer.Mode = b.ko.Spec.PrimaryContainer.Mode
	// }

	if a.ko.Spec.EnableNetworkIsolation == nil && b.ko.Spec.EnableNetworkIsolation != nil {
		a.ko.Spec.EnableNetworkIsolation = b.ko.Spec.EnableNetworkIsolation
	}
}
