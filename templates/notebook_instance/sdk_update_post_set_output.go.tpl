	curr := ko.GetAnnotations()
	if curr == nil {
		curr = make(map[string]string)
	}
	curr["done_updating"] = "true"
	ko.SetAnnotations(curr)
	//Making the controller requeue after calling update.
	rm.customSetOutputUpdate(ko)