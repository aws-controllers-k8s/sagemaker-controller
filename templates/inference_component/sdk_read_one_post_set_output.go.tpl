	// Manually set the RuntimeConfig.CopyCount from read response RuntimeConfig.DesiredCopyCount
	if resp.RuntimeConfig != nil && ko.Spec.RuntimeConfig != nil {
		desiredCountCopy := int64(*resp.RuntimeConfig.DesiredCopyCount)
		ko.Spec.RuntimeConfig.CopyCount = &desiredCountCopy
	}

	rm.customDescribeInferenceComponentSetOutput(ko)
