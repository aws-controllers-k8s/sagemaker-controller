//TODO: Take this out if runtime supports updating annotations during ReadOne
tmp := ""
is_updating_tmp := ""
if r != nil && r.ko != nil && r.ko.Status.StoppedByAck != nil{
  tmp = *r.ko.Status.StoppedByAck
}
if r != nil && r.ko != nil && r.ko.Status.IsUpdating != nil{
  is_updating_tmp = *r.ko.Status.IsUpdating
}