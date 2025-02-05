if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

latestStatus := r.ko.Status.NotebookInstanceStatus

if latestStatus != nil &&
 *latestStatus == string(svcsdktypes.NotebookInstanceStatusInService) {
	if err := rm.stopNotebookInstance(ctx, r); err != nil {
		return nil, err
	} else {
		return r, requeueWaitWhileStopping
	}
}