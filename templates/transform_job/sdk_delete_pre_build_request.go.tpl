	// Call StopTranformJob only if the job is InProgress, otherwise just return nil to mark the
	// resource Unmanaged
	latestStatus := r.ko.Status.TransformJobStatus
	if latestStatus != nil && *latestStatus != svcsdk.TransformJobStatusInProgress {
		return nil, nil
	}