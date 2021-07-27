package notebook_instance

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

/* We just set the defualts here in case the user does not specify them. This avoids the controller trying to update itself right after creation.
Direct Internet Access does ...
Root Acess does .....
*/
func customSetDefaults(
	a *resource,
	b *resource,
) {
	if ackcompare.IsNil(a.ko.Spec.DirectInternetAccess) && ackcompare.IsNotNil(b.ko.Spec.DirectInternetAccess) {
		a.ko.Spec.DirectInternetAccess = b.ko.Spec.DirectInternetAccess
	}
	if ackcompare.IsNil(a.ko.Spec.RootAccess) && ackcompare.IsNotNil(b.ko.Spec.RootAccess) {
		a.ko.Spec.RootAccess = b.ko.Spec.RootAccess
	}
}
