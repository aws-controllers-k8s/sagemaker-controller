	// Manually set the RuntimeConfig.CopyCount from read response RuntimeConfig.DesiredCopyCount
	if resp.RuntimeConfig != nil && ko.Spec.RuntimeConfig != nil {
		ko.Spec.RuntimeConfig.CopyCount = resp.RuntimeConfig.DesiredCopyCount
	}

	rm.customDescribeInferenceComponentSetOutput(ko)
