    if isModelPackageGroupPending(&resource{ko}) {
        msg := "ModelPackageGroup is currently pending"
		setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
		return &resource{ko}, requeueWaitWhilePending
	}
    if isModelPackageGroupInProgress(&resource{ko}) {
        msg := "ModelPackageGroup is currently in progress"
		setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
		return &resource{ko}, requeueWaitWhileInProgress
	}
    ModelPackageGroupCustomSetOutput(&resource{ko})