observed, err := rm.sdkFind(ctx, latest)
if err != nil {
    return observed, err
}
tmp_resource := &resource{ko}
tmp_resource.SetStatus(observed)
