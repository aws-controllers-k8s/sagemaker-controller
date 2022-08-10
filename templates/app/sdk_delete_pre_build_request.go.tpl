	latestStatus := r.ko.Status.Status
	if latestStatus != nil && *latestStatus == svcsdk.AppStatusDeleted {
		return nil, nil
	}

	if err = rm.requeueUntilCanModify(ctx, r); err != nil {
		return r, err
	}
