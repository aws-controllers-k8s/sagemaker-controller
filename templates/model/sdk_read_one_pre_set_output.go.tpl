
// This is require if only spec.primarycontainer.modeldataurl is in use and not the ko.spec.primarycontainer.ModelDataSource
// in this case, during find "ko.spec.primarycontainer.ModelDataSource" gets updated as well , which creates a new k8s generation    

if ko.Spec.PrimaryContainer.ModelDataSource == nil && resp.PrimaryContainer.ModelDataSource != nil {
	resp.PrimaryContainer.ModelDataSource = nil	
}