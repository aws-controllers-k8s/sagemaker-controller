warmpool_diff := delta.DifferentAt("Spec.ResourceConfig.KeepAlivePeriodInSeconds")
profiler_diff := delta.DifferentAt("Spec.ProfilerConfig") || delta.DifferentAt("Spec.ProfilerRuleConfigurations")
if warmpool_diff && profiler_diff {
	return latest, ackerr.NewTerminalError(errors.New("cannot update Warm pool and Profiler at the same time"))
}
if !warmpool_diff && !profiler_diff {
	return latest, ackerr.NewTerminalError(errors.New("only Warm Pool or Profiler can be updated"))
}
trainingSecondaryStatus := latest.ko.Status.SecondaryStatus
if ackcompare.IsNotNil(trainingSecondaryStatus) && *trainingSecondaryStatus == svcsdk.SecondaryStatusStarting {
	return nil, ackrequeue.NeededAfter(
		errors.New("training job cannot be updated while secondary status is in Starting state."),
		ackrequeue.DefaultRequeueAfterDuration,
	)
}
if warmpool_diff {
	input.SetProfilerConfig(nil)
	input.SetProfilerRuleConfigurations(nil)
	if err := rm.isWarmPoolUpdatable(latest); err != nil {
		return nil, err
	}
}
if profiler_diff {
	if err := rm.isProfilerUpdatable(latest); err != nil {
		return nil, err
	}
	input.SetResourceConfig(nil)
	if rm.isProfilerRemoved(desired, latest) {
		rm.handleProfilerRemoval(input)
	} else {
		if err := rm.customSetUpdateInput(desired, latest, delta, input); err != nil {
			return nil, err
		} 
	} 
}
