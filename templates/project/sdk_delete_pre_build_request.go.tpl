	latestStatus := r.ko.Status.ProjectStatus
	if latestStatus != nil && *latestStatus == string(svcsdktypes.ProjectStatusDeleteCompleted) {
		return nil, nil
	}

	if err = rm.requeueUntilCanModify(ctx, r); err != nil {
		return r, err
	}
