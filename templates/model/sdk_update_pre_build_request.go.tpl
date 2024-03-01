// this to handle delete/remove tags 
_ , err = rm.deleteTags(ctx,desired,latest)
if err != nil {
	return nil, err
}