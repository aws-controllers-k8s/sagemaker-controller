	if isModelPackageDeleting(r) {
        msg := "ModelPackage is currently being deleted"
		setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
		return requeueWaitWhileDeleting
	}
    if isModelPackageInProgress(r) {
        msg := "ModelPackage is currently in progress"
		setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
		return requeueWaitWhileInProgress
	}