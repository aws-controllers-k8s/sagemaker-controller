
// If creating, requeue with wait untill status becomes created.
if isCreating(desired) {
        return nil, requeueWaitWhileCreating
}
