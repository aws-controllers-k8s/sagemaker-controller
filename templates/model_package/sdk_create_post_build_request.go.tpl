    // If ModelPackageGroupName set after newRequestPayload, return error
    // This is because versioned modelpackage is not supported.
	if input.ModelPackageGroupName != nil {
		return nil, VersionedModelPackageNotSupported
	} 