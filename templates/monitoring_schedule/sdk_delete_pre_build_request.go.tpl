    // specialized logic to check if modification is allowed
    err = rm.statusAllowUpdates(ctx, r)
    if err != nil {
        return err
    }