	latestStatus := r.ko.Status.PipelineExecutionStatus
	if latestStatus != nil {
		if *latestStatus == string(svcsdktypes.PipelineExecutionStatusStopping) {
			return r, requeueWaitWhileDeleting
		}

		// Call StopPipelineExecution only if the job is Executing, otherwise just
		// return nil to mark the resource Unmanaged
		if *latestStatus != string(svcsdktypes.PipelineExecutionStatusExecuting) {
			return r, err
		}
	}