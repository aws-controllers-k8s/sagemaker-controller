err = rm.customSetOutputDescribe(r, ko)
if err != nil{
  return &resource{ko}, err
}