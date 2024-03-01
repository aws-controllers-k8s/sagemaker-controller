var resp_tags *svcsdk.ListTagsOutput

resp_tags, err = rm.sdkapi.ListTagsWithContext(ctx,&svcsdk.ListTagsInput{ResourceArn: resp.ModelArn})
rm.metrics.RecordAPICall("READ_ONE", "DescribeTags", err)

if resp_tags != nil {
	f6 := []*svcapitypes.Tag{}
	for _, f6iter := range resp_tags.Tags {
		f6elem := &svcapitypes.Tag{}
		if f6iter.Key != nil {
			f6elem.Key = f6iter.Key
		}
		if f6iter.Value != nil {
			f6elem.Value = f6iter.Value
		}
		f6 = append(f6, f6elem)
	}
	ko.Spec.Tags = f6
	}