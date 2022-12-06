	latestStatus := r.ko.Status.PipelineExecutionStatus
	if latestStatus != nil {
		if *latestStatus == svcsdk.PipelineExecutionStatusStopping {
			return r, requeueWaitWhileDeleting
		}

		// Call StopPipelineExecution only if the job is Executing, otherwise just 
		// return nil to mark the resource Unmanaged
		if *latestStatus != svcsdk.PipelineExecutionStatusExecuting {
			return r, err
		}
	}