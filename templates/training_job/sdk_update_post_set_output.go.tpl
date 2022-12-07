observed, err := rm.sdkFind(ctx, latest)
if err != nil {
    return observed, err
}
ko.Status = observed.ko.Status
return &resource{ko}, ackrequeue.NeededAfter(
	errors.New("training job is updating"),
	ackrequeue.DefaultRequeueAfterDuration,
)