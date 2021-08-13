if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
	return latest, err
}
stopped_by_ack := false
latestStatus := latest.ko.Status.NotebookInstanceStatus
if latestStatus != nil &&
 *latestStatus == svcsdk.NotebookInstanceStatusInService {
	if err := rm.stopNotebookInstance(latest); err != nil {
		return latest, err
	} else {
		stopped_by_ack = true
	}
}

//TODO: Take this out if the runtime supports updating annotations if an error is returned and use annotations for this.
if stopped_by_ack {
	latest.ko.Status.StoppedByController = aws.String("true")
	return latest, requeueWaitWhileStopping
}