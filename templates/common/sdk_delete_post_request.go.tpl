    if err == nil {
            if foundResource, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
                if err != nil {
                    return foundResource, err
                }
                return foundResource, requeueWaitWhileDeleting
            }
    }