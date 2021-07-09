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

package model_package

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func customSetDefaults(
	a *resource,
	b *resource,
) {
	if ackcompare.IsNil(a.ko.Spec.CertifyForMarketplace) && ackcompare.IsNotNil(b.ko.Spec.CertifyForMarketplace) {
		a.ko.Spec.CertifyForMarketplace = b.ko.Spec.CertifyForMarketplace
	}
	if ackcompare.IsNil(a.ko.Spec.ModelApprovalStatus) && ackcompare.IsNotNil(b.ko.Spec.ModelApprovalStatus) {
		a.ko.Spec.ModelApprovalStatus = b.ko.Spec.ModelApprovalStatus
	}
}
