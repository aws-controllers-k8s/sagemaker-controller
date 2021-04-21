    // Call StopProcessingJob only if the job is InProgress, otherwise just return nil to mark the
	// resource Unmanaged
	latestStatus := r.ko.Status.ProcessingJobStatus
	if latestStatus != nil && *latestStatus != svcsdk.ProcessingJobStatusInProgress {
		return nil
	}