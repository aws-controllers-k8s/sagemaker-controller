if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

stopped_by_controller,err := rm.customPreDelete(r)
if err != nil{
	return latest,err
}
if stopped_by_controller{
	return r,requeueWaitWhileStopping
}