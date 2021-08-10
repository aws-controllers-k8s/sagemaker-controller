//TODO: Replace the is_updating status with an annotation if the runtime can update annotations after a readOne call.
ko.Status.IsUpdating = aws.String("true")
//Making the controller requeue after calling update.
rm.customSetOutputUpdate(ko)