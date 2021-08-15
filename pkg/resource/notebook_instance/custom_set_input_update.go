package notebook_instance

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// handleUpdateOnlyParameters sets Disassociate<field> to true if corresponding value in desired is nil(ie not included in the spec)
// and if the corresponding value in the latest spec is not nil (ie contains some value)
// Ex if LifecycleConfigName was "a" in the latest spec but nil in the desired spec, DisassociateLifecycleConfig is set to true
func handleUpdateOnlyParameters(
	desired *resource,
	latest *resource,
	update_input *svcsdk.UpdateNotebookInstanceInput) {
	if ackcompare.IsNil(desired.ko.Spec.AcceleratorTypes) && ackcompare.IsNotNil(latest.ko.Spec.AcceleratorTypes) {
		update_input.DisassociateAcceleratorTypes = aws.Bool(true)
	}
	if ackcompare.IsNil(desired.ko.Spec.AdditionalCodeRepositories) && ackcompare.IsNotNil(latest.ko.Spec.AdditionalCodeRepositories) {
		update_input.DisassociateAdditionalCodeRepositories = aws.Bool(true)
	}
	if ackcompare.IsNil(desired.ko.Spec.DefaultCodeRepository) && ackcompare.IsNotNil(latest.ko.Spec.DefaultCodeRepository) {
		update_input.DisassociateDefaultCodeRepository = aws.Bool(true)
	}
	if ackcompare.IsNil(desired.ko.Spec.LifecycleConfigName) && ackcompare.IsNotNil(latest.ko.Spec.LifecycleConfigName) {
		update_input.DisassociateLifecycleConfig = aws.Bool(true)
	}
}
