//Todo Take this if statement out if code generator can generate this field.
if resp.Url != nil{
  ko.Status.NotebookInstanceURL = resp.Url
}