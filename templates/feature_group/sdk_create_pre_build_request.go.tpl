
// If creating, requeue with wait until status becomes created.
if isCreating(desired) {
        return nil, requeueWaitWhileCreating
}
