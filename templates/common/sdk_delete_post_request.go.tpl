
if err == nil {
	if observed, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
		if err != nil {
			return nil, err
		}
		return observed, requeueWaitWhileDeleting
	}
}
    