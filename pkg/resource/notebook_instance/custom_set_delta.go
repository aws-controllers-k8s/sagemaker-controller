package notebook_instance

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// customSetDefaults sets the default fields for DirectInternetAccess and RootAccess.
func customSetDefaults(
	a *resource,
	b *resource,
) {
	// Direct Internet Access describes whether Amazon SageMaker provides internet access to the notebook instance.
	// The default value is Enabled.
	if ackcompare.IsNil(a.ko.Spec.DirectInternetAccess) && ackcompare.IsNotNil(b.ko.Spec.DirectInternetAccess) {
		a.ko.Spec.DirectInternetAccess = b.ko.Spec.DirectInternetAccess
	}
	// Root Access describes whether root access is enabled or disabled for users of the notebook instance.
	// The default value is Enabled
	if ackcompare.IsNil(a.ko.Spec.RootAccess) && ackcompare.IsNotNil(b.ko.Spec.RootAccess) {
		a.ko.Spec.RootAccess = b.ko.Spec.RootAccess
	}
}
