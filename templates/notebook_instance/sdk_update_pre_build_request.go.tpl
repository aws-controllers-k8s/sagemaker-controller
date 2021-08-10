if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
	return latest, err
}
stopped_by_ack := rm.customStopNotebook(latest)
//TODO: Take this out if the runtime supports updating annotations if an error is returned and use annotations for this.
if stopped_by_ack {
	latest.ko.Status.StoppedByAck = aws.String("true")
	return latest, requeueWaitWhileStopping
}