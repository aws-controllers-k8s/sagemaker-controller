started_notebook := rm.customSetOutputDescribe(r, ko)
//TODO: Take this out if runtime supports updating annotations when an error is returned.
if !started_notebook{
  //covers a scenario where the code generator generates code to set StoppedByAck
  ko.Status.StoppedByAck = &tmp 
}
