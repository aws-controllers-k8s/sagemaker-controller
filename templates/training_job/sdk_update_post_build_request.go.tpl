warmpool_diff := delta.DifferentAt("Spec.ResourceConfig.KeepAlivePeriodInSeconds")
profiler_diff := delta.DifferentAt("Spec.ProfilerConfig") || delta.DifferentAt("Spec.ProfilerRuleConfigurations")
if warmpool_diff && profiler_diff {
	return latest, ackerr.NewTerminalError(errors.New("cannot update Warm pool and Profiler at the same time"))
}
if !warmpool_diff && !profiler_diff {
	return latest, ackerr.NewTerminalError(errors.New("only Warm Pool or Profiler can be updated"))
}
if warmpool_diff {
	input.SetProfilerConfig(nil)
	input.SetProfilerRuleConfigurations(nil)
	warmpool_terminal := warmPoolTerminalCheck(latest)
	if warmpool_terminal {
		return latest, ackerr.NewTerminalError(errors.New("warm pool either does not exist or has reached a non updatable state"))
	}
	//Requeue if TrainingJob is in InProgress state
	if err := customSetOutputUpdateWarmpool(latest); err != nil {
		return nil,err
	}
}
if profiler_diff {
	if up_err := customSetOutputUpdateProfiler(latest); up_err != nil {
		return nil, up_err
	}
	input.SetResourceConfig(nil)
	if profilerRemovalCheck(desired, latest) {
		handleProfilerRemoval(input)
	} else{
		inp_err := customSetUpdateInput(desired, latest, delta, input)
		if inp_err != nil {
			return nil, err	
		} 
	} 
}
