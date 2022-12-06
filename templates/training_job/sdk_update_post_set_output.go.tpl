observed, err := rm.sdkFind(ctx, latest)
if err != nil {
    return observed, err
}
desired.SetStatus(observed)
return desired, ackrequeue.NeededAfter(
	errors.New("training job is updating"),
	ackrequeue.DefaultRequeueAfterDuration,
)