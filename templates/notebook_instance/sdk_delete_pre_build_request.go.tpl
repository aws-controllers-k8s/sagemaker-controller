//This will avoid exponential backoff
if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

rm.customPreDelete(r)