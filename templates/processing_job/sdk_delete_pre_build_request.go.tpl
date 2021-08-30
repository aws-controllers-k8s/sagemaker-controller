    latestStatus := r.ko.Status.ProcessingJobStatus
	if latestStatus != nil {
		if *latestStatus == svcsdk.ProcessingJobStatusStopping {
			return r, requeueWaitWhileDeleting
		}

		// Call StopProcessingJob only if the job is InProgress, otherwise just 
		// return nil to mark the resource Unmanaged
		if *latestStatus != svcsdk.ProcessingJobStatusInProgress {
			return r, err
		}
	}