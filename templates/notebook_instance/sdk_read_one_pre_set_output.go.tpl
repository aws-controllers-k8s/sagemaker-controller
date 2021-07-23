tmp := ""
if r != nil && r.ko != nil && r.ko.Status.StoppedByAck != nil{
  tmp = *r.ko.Status.StoppedByAck
}