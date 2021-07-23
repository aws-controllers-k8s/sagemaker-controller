
if err == nil {
        if _, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
                if err != nil {
                    return nil, err
                }
                return r, requeueWaitWhileDeleting
            }
    }
