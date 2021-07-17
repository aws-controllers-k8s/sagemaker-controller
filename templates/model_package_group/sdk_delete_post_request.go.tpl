    if err == nil {
            if _, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
                return requeueWaitWhileDeleting
            }
    }