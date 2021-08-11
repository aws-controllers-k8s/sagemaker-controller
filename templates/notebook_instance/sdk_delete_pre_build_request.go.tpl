if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

latestStatus := r.ko.Status.NotebookInstanceStatus

if latestStatus != nil &&
 *latestStatus == svcsdk.NotebookInstanceStatusInService {
	if err := rm.stopNotebookInstance(r); err != nil {
		return r, err
	} else {
		return r, requeueWaitWhileStopping
	}
}