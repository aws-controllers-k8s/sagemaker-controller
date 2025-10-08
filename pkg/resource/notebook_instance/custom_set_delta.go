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

package notebook_instance

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// customSetDefaults sets the default fields for DirectInternetAccess and RootAccess.
func customSetDefaults(
	a *resource,
	b *resource,
) {
	// Direct Internet Access describes whether Amazon SageMaker provides internet access to the notebook instance.
	// The default value is Enabled.
	DefaultDirectInternetAccess := aws.String("Enabled")
	if ackcompare.IsNil(a.ko.Spec.DirectInternetAccess) && ackcompare.IsNotNil(b.ko.Spec.DirectInternetAccess) {
		a.ko.Spec.DirectInternetAccess = DefaultDirectInternetAccess
	}

	// Root Access describes whether root access is enabled or disabled for users of the notebook instance.
	// The default value is Enabled
	DefaultRootAccess := aws.String("Enabled")
	if ackcompare.IsNil(a.ko.Spec.RootAccess) && ackcompare.IsNotNil(b.ko.Spec.RootAccess) {
		a.ko.Spec.RootAccess = DefaultRootAccess
	}

	// IPAddressType specify either dualstack ipv4/ipv6 or just ipv4
	// The default value is ipv4
	DefaultIPAddressType := aws.String("ipv4")
	if ackcompare.IsNil(a.ko.Spec.IPAddressType) && ackcompare.IsNotNil(b.ko.Spec.IPAddressType) {
		a.ko.Spec.RootAccess = DefaultIPAddressType
	}
}
