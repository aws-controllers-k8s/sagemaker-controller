	latestStatus := r.ko.Status.LabelingJobStatus
	if latestStatus != nil {
		if *latestStatus == string(svcsdktypes.LabelingJobStatusStopping) {
			return r, requeueWaitWhileDeleting
		}

		// Call StopLabelingJob only if the job is InProgress, otherwise just
		// return nil to mark the resource Unmanaged
		if *latestStatus != string(svcsdktypes.LabelingJobStatusInProgress) {
			return r, err
		}
	}
