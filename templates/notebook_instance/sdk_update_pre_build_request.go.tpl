if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
	return latest, err
}
stopped_by_ack := false
latestStatus := latest.ko.Status.NotebookInstanceStatus
if latestStatus != nil && *latestStatus == svcsdk.NotebookInstanceStatusInService {
		err := rm.stopNotebookInstance(latest)
		if err == nil {
			stopped_by_ack = true
		} else {
			return latest, err
		}
	}

//TODO: Take this out if the runtime supports updating annotations if an error is returned and use annotations for this.
if stopped_by_ack {
	latest.ko.Status.StoppedByAck = aws.String("true")
	return latest, requeueWaitWhileStopping
}