if isNotebookStopping(latest){
    return latest,requeueWaitWhileStopping
}
if isNotebookPending(latest){
    return latest,requeueWaitWhilePending
}
if isNotebookUpdating(latest) && latest.ko.Status.FailureReason == nil {
	return latest, requeueWaitWhileUpdating
}
stopped_by_ack := rm.customPreUpdate(ctx,desired,latest)
if stopped_by_ack {
		stopped_by_ack_str := "true"
		latest.ko.Status.StoppedByAck = &stopped_by_ack_str
		return latest, requeueWaitWhileStopping
}