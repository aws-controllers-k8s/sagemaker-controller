    if isModelPackageDeleting(latest) {
		msg := "ModelPackage is currently being deleted"
		setSyncedCondition(desired, corev1.ConditionFalse, &msg, nil)
		return desired, requeueWaitWhileDeleting
	}
	if isModelPackagePending(latest) {
		msg := "ModelPackage is currently pending"
		setSyncedCondition(desired, corev1.ConditionFalse, &msg, nil)
		return desired, requeueWaitWhilePending
	}
	if isModelPackageInProgress(latest) {
		msg := "ModelPackage is currently in progress"
		setSyncedCondition(desired, corev1.ConditionFalse, &msg, nil)
		return desired, requeueWaitWhileInProgress
	}
	if ModelPackageHasTerminalStatus(latest) {
		msg := "ModelPackage is in '"+*latest.ko.Status.ModelPackageStatus+"' status"
		setTerminalCondition(desired, corev1.ConditionTrue, &msg, nil)
		setSyncedCondition(desired, corev1.ConditionTrue, nil, nil)
		return desired, nil
	}