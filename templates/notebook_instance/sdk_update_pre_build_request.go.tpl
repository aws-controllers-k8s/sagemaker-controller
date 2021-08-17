if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
	return latest, err
}

latestStatus := latest.ko.Status.NotebookInstanceStatus
// The Notebook Instance is stopped and the StoppedByControllerMetadata status is
// set to UpdatePending.
if latestStatus != nil &&
		*latestStatus == svcsdk.NotebookInstanceStatusInService {
		if err := rm.stopNotebookInstance(latest); err != nil {
			return nil, err
		} else {
			//TODO: Replace with annotations once rutime supports it.
			latest.ko.Status.StoppedByControllerMetadata = aws.String("UpdatePending")
			return latest, requeueWaitWhileStopping
		}
}