if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

if latestStatus := r.ko.Status.NotebookInstanceStatus; latestStatus != nil &&
 *latestStatus == svcsdk.NotebookInstanceStatusInService {
	if err := rm.stopNotebookInstance(r); err == nil {
		return r,requeueWaitWhileStopping
	} else {
		return r, err
	}
}