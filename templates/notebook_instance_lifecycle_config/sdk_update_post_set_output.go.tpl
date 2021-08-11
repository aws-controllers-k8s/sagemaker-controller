//Done because controller finishes reconciling after update.
updated, err = rm.sdkFind(ctx, latest)
if err != nil {
	return latest, err
}
ko.Status.LastModifiedTime = updated.ko.Status.LastModifiedTime