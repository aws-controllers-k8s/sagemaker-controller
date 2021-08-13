if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
	return latest, err
}
latestStatus := latest.ko.Status.NotebookInstanceStatus
// The Notebook Instance is stopped and the StoppedByController status is
// set to UpdatePending.
if latestStatus != nil &&
 *latestStatus == svcsdk.NotebookInstanceStatusInService {
	if err := rm.stopNotebookInstance(latest); err != nil {
		return latest, err
	} else {
		//TODO: Take this out if the runtime supports updating annotations if an error is returned and use annotations for this.
		latest.ko.Status.StoppedByControllerMETA = aws.String("UpdatePending")
	  return latest, requeueWaitWhileStopping
	}
}