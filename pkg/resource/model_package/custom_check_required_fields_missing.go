package model_package

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) customCheckRequiredFieldsMissingMethod(
	r *resource,
) bool {
	return r.ko.Spec.ModelPackageName == nil && r.ko.Spec.ModelPackageGroupName == nil

}
