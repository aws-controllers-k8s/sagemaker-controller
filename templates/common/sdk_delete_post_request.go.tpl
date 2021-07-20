    if err == nil {
            if _, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
                if err != nil {
                    return err
                }
                return requeueWaitWhileDeleting
            }
    }