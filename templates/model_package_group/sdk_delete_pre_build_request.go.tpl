    if isModelPackageGroupDeleting(r) {
        msg := "ModelPackageGroup is currently being deleted"
        setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
        return requeueWaitWhileDeleting
    }
    if isModelPackageGroupInProgress(r) {
        msg := "ModelPackageGroup is currently in progress"
        setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
        return requeueWaitWhileInProgress
    }
    if isModelPackageGroupPending(r) {
        msg := "ModelPackageGroup is currently pending"
        setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
        return requeueWaitWhilePending
    }
    if isModelPackageGroupDeleteFailed(r) {
        // TODO Implement exponential backoff for DeleteFailed
        msg := "ModelPackageGroup delete failed"
        setSyncedCondition(r, corev1.ConditionFalse, &msg, nil)
        return requeueWaitWhileDeleteFailed
    }