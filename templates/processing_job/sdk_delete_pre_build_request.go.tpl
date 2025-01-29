    latestStatus := r.ko.Status.ProcessingJobStatus
	if latestStatus != nil {
		if *latestStatus == string(svcsdktypes.ProcessingJobStatusStopping) {
			return r, requeueWaitWhileDeleting
		}

		// Call StopProcessingJob only if the job is InProgress, otherwise just
		// return nil to mark the resource Unmanaged
		if *latestStatus != string(svcsdktypes.ProcessingJobStatusInProgress) {
			return r, err
		}
	}