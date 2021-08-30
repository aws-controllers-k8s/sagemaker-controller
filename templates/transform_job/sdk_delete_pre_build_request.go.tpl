	latestStatus := r.ko.Status.TransformJobStatus
	if latestStatus != nil {
		if *latestStatus == svcsdk.TransformJobStatusStopping {
			return r, requeueWaitWhileDeleting
		}

		// Call StopTranformJob only if the job is InProgress, otherwise just 
		// return nil to mark the resource Unmanaged
		if *latestStatus != svcsdk.TransformJobStatusInProgress {
			return r, err
		}
	}