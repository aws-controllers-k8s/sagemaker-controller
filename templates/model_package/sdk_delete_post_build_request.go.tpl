    // If ModelPackageName not set after newDeleteRequestPayload attempt to use ARN
    // This is because versioned modelpackage uses ARN not name
	if input.ModelPackageName == nil {
		arn := r.Identifiers().ARN()
		if arn == nil {
			return nil, ackerr.NotFound
		}
		input.SetModelPackageName(string(*arn))
	}