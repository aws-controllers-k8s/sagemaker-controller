	curr := ko.GetAnnotations()
	if curr == nil {
		curr = make(map[string]string)
	}
	curr["done_updating"] = "true"
	ko.SetAnnotations(curr)
	rm.customSetOutputCreateUpdate(ko)