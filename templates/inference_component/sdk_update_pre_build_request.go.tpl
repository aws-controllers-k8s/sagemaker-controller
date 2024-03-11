	if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
		return nil, err
	}

	if err = rm.customUpdateInferenceComponentPreChecks(ctx, desired, latest, delta); err != nil {
		return nil, err
	}
