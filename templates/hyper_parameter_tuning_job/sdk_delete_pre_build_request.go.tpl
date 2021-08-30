	latestStatus := r.ko.Status.HyperParameterTuningJobStatus
	if latestStatus != nil {
		if *latestStatus == svcsdk.HyperParameterTuningJobStatusStopping {
			return r, requeueWaitWhileDeleting
		}

		// Call StopHyperParameterTuningJob only if the job is InProgress, otherwise just 
		// return nil to mark the resource Unmanaged
		if *latestStatus != svcsdk.HyperParameterTuningJobStatusInProgress {
			return r, err
		}
	}