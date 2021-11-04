// Manually set the DomainID as Create only return the ARN
if resp.DomainArn != nil && ko.Status.DomainID == nil {
  ko.Status.DomainID = &strings.Split(*resp.DomainArn, "/")[1] 
}
