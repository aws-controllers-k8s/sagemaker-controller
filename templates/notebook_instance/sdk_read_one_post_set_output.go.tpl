modified_status := rm.customSetOutputDescribe(r, ko)
//TODO: Take this out if runtime supports updating annotations when an error is returned.
if !modified_status{
  //covers a scenario where the code generator generates code to set StoppedByAck
  ko.Status.StoppedByAck = &tmp 
  ko.Status.IsUpdating = &is_updating_tmp
}
