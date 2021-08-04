//Todo Take this if statement out if code generator can generate this field.
if resp.Url != nil{
  ko.Status.NotebookInstanceURL = resp.Url
}
tmp := ""
if r != nil && r.ko != nil && r.ko.Status.StoppedByAck != nil{
  tmp = *r.ko.Status.StoppedByAck
}