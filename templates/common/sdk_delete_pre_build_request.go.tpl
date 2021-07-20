    if err = rm.requeueUntilCanModify(ctx, r); err != nil {
            return err
        }