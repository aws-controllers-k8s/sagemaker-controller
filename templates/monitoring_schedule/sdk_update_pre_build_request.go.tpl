    // specialized logic to check if modification is allowed
    err = rm.statusAllowUpdates(ctx, latest)
    if err != nil {
        return nil, err
    }