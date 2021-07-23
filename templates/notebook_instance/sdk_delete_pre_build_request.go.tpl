//This will avoid exponential backoff
if isNotebookStopping(r){
    return r,requeueWaitWhileStopping
}
//This will avoid exponential backoff
if isNotebookPending(r){
    return r,requeueWaitWhilePending
}
if isNotebookDeleting(r){
    return nil,requeueWaitWhileDeleting
}

rm.customPreDelete(r)