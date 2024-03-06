// checking if any delta is found other than tags
if len(delta.Differences) > 0 {
	if delta.DifferentExcept("Spec.Tags") {	
		for _, parts := range delta.Differences {
			if !parts.Path.Contains("Tags") {
				return nil, fmt.Errorf("cannot update the following fields: %s , Allowed field to change: Spec.Tags", parts.Path)
			}
		}
	}
}
// this to handle delete/remove tags 
_ , err = rm.deleteTags(ctx,desired,latest)
if err != nil {
	return nil, err
}