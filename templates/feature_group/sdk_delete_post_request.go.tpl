// If a delete failed, requeue on delete.
if err == nil {
        if foundResource, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
                if isDeleteFailed(foundResource) {
                        return requeueWaitWhileDeleteFailed
                }
                return err
        }
}
