    // If ModelPackageName not set after newDeleteRequestPayload attempt to use ARN
    // This is because versioned modelpackage uses ARN not name
    if input.ModelPackageName == nil {  
        if r.ko.Status.ACKResourceMetadata != nil && r.ko.Status.ACKResourceMetadata.ARN != nil {
            input.SetModelPackageName(string(*r.ko.Status.ACKResourceMetadata.ARN))
        } else {
            return nil, ackerr.NotFound
        }
    }