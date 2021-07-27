started_notebook := rm.customSetOutputDescribe(r, ko)
if !started_notebook{
  ko.Status.StoppedByAck = &tmp //covers a scenario where the code generator sets r.ko.Status.StoppedByAck
}
