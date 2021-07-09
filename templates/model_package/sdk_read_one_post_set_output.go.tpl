    if isModelPackagePending(&resource{ko}) {
            msg := "ModelPackage is currently pending"
		    setSyncedCondition(&resource{ko}, corev1.ConditionFalse, &msg, nil)
            return &resource{ko}, requeueWaitWhilePending
    }
    if isModelPackageInProgress(&resource{ko}) {
            msg := "ModelPackage is currently in progress"
		    setSyncedCondition(&resource{ko}, corev1.ConditionFalse, &msg, nil)
            return &resource{ko}, requeueWaitWhileInProgress
    }
    ModelPackageCustomSetOutput(&resource{ko})