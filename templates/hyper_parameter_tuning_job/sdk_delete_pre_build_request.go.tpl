    // Call StopHyperparameterTuningJob only if the job is InProgress, otherwise just return nil to mark the
	// resource Unmanaged
	latestStatus := r.ko.Status.HyperParameterTuningJobStatus
	if latestStatus != nil && *latestStatus != svcsdk.HyperParameterTuningJobStatusInProgress {
		return nil, nil
	}