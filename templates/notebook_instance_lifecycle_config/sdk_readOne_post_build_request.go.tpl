awErr, ok := ackerr.AWSError(err)

		if ok {
			errMsg := awErr.Message()
			if strings.Contains(errMsg, "Notebook Instance Lifecycle Config does not exist") {
				return nil, ackerr.NotFound
			}

		}