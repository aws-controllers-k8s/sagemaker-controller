//This will avoid exponential backoff
if err = rm.requeueUntilCanModify(ctx, r); err != nil {
	return r, err
}

//Stops the Notebook Instance
rm.customPreDelete(r)