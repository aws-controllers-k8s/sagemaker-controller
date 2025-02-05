	latestStatus := r.ko.Status.TrainingJobStatus
	if latestStatus != nil {
		if *latestStatus == string(svcsdktypes.TrainingJobStatusStopping) {
			return r, requeueWaitWhileDeleting
		}

		// Call StopTrainingJob only if the job is InProgress, otherwise just
		// return nil to mark the resource Unmanaged
		if *latestStatus != string(svcsdktypes.TrainingJobStatusInProgress) {
			return r, err
		}
	}