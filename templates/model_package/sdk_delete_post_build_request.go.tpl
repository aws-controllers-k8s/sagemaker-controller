    if input.ModelPackageName == nil {  
        if r.ko.Status.ACKResourceMetadata != nil && r.ko.Status.ACKResourceMetadata.ARN != nil {
            input.SetModelPackageName(string(*r.ko.Status.ACKResourceMetadata.ARN))
        } else {
            return ackerr.NotFound
        }
    }