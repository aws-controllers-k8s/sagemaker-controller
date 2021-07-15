    if err == nil {
            if foundResource, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
                if isDeleting(foundResource) {
                    return requeueWaitWhileDeleting
                }
                return err
            }
    }