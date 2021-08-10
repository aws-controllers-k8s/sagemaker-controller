if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}
latestStatus := r.ko.Status.NotebookInstanceStatus

if latestStatus != nil && *latestStatus == svcsdk.NotebookInstanceStatusInService {
		err := rm.stopNotebookInstance(r)
		if err == nil {
			return r,requeueWaitWhileStopping
		} else {
			return r, err
		}
	}
if err != nil{
	return r,err
}