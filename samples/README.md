# Job Sample Overview

This sample demonstrates how to start jobs using your own script, packaged in a SageMaker-compatible container, using the Amazon AWS Controllers for Kubernetes (ACK) service controller for Amazon SageMaker.                     

## Prerequisites    

This sample assumes that you have already configured an Kubernetes cluster with the ACK operator. It also assumes that you have installed `kubectl` - you can find a link on our [installation page](To do).

You will also need an IAM role which has permissions to access your S3 resources and SageMaker. If you have not yet created a role with these permissions, you can find an example policy at [Amazon SageMaker Roles](https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-roles.html#sagemaker-roles-createtrainingjob-perms).

### Creating your first Job

The easiest way to start is taking a look at the sample training jobs and its corresponding [README](/samples/training/README.md)