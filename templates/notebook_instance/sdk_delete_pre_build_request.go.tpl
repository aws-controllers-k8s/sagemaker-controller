//This will avoid exponential backoff
if isNotebookStopping(r){
    return requeueWaitWhileStopping
}
//This will avoid exponential backoff
if isNotebookPending(r){
    return requeueWaitWhilePending
}

rm.customDelete(r)