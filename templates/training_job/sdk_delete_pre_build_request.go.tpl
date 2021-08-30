	latestStatus := r.ko.Status.TrainingJobStatus
	if latestStatus != nil {
		if *latestStatus == svcsdk.TrainingJobStatusStopping {
			return r, requeueWaitWhileDeleting
		}

		// Call StopTrainingJob only if the job is InProgress, otherwise just 
		// return nil to mark the resource Unmanaged
		if *latestStatus != svcsdk.TrainingJobStatusInProgress {
			return r, err
		}
	}