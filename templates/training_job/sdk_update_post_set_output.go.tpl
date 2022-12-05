return desired, ackrequeue.NeededAfter(
	errors.New("training job is updating"),
	ackrequeue.DefaultRequeueAfterDuration,
)