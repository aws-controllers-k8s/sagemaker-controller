import urllib.request
from urllib.parse import urlparse
import sys

# Gets your bucket name from command line
try:
    bucket = str(sys.argv[1])
except Exception as error:
    print("Please pass your bucket name as a commandline argument")
    sys.exit(1)

# Download model from pinned commit
url = "https://github.com/aws/amazon-sagemaker-examples/raw/af6667bd0be3c9cdec23fecda7f0be6d0e3fa3ea/sagemaker_model_monitor/introduction/model/xgb-churn-prediction-model.tar.gz"
urllib.request.urlretrieve(url, "xgb-churn-prediction-model.tar.gz")

import os
import boto3

# Upload model to S3
prefix = "sagemaker/xgboost/model"
model_file = open("xgb-churn-prediction-model.tar.gz", "rb")
s3_key = os.path.join(prefix, "xgb-churn-prediction-model.tar.gz")
boto3.Session().resource("s3").Bucket(bucket).Object(s3_key).upload_fileobj(model_file)
s3url = f"s3://{bucket}/{s3_key}"
print(f"Writing to {s3url}")

# Remove downloaded file
os.remove("xgb-churn-prediction-model.tar.gz")
