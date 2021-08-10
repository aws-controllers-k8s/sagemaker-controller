if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

stopped_by_controller := rm.customPreDelete(r)
if stopped_by_controller{
	return r,requeueWaitWhileStopping
}