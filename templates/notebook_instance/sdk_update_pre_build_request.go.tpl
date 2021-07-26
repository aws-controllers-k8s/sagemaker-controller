if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
	return latest, err
}

stopped_by_ack := rm.customPreUpdate(ctx,desired,latest)
if stopped_by_ack {
		latest.ko.Status.StoppedByAck = aws.String("true")
		return latest, requeueWaitWhileStopping
}