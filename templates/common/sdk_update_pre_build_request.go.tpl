    if err = rm.requeueUntilCanModify(ctx, latest); err != nil {
        return nil, err
    }