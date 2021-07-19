    // Call StopTrainingJob only if the job is InProgress, otherwise just return nil to mark the
	// resource Unmanaged
	latestStatus := r.ko.Status.TrainingJobStatus
	if latestStatus != nil && *latestStatus != svcsdk.TrainingJobStatusInProgress {
		return nil, err
	}