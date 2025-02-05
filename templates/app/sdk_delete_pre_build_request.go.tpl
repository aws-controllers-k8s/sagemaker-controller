	latestStatus := r.ko.Status.Status
	if latestStatus != nil && *latestStatus == string(svcsdktypes.AppStatusDeleted) {
		return nil, nil
	}

	if err = rm.requeueUntilCanModify(ctx, r); err != nil {
		return r, err
	}
