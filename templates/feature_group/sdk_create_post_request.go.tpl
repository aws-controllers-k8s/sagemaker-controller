
// If creating, requeue with wait untill status becomes created.
if err == nil {
        if foundResource, err := rm.sdkFind(ctx, desired); err != ackerr.NotFound {
                if isCreating(foundResource) {
                        return nil, requeueWaitWhileCreating
                }
                return nil, err
        }
}
