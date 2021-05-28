import urllib.request
from urllib.parse import urlparse
import sys

# Gets your bucket name from command line
try:
    bucket = str(sys.argv[1])
except Exception as error:
    print("Please pass your bucket name as a commandline argument")
    sys.exit(1)

# Download dataset from pinned commit
url = "https://github.com/aws/amazon-sagemaker-examples/raw/af6667bd0be3c9cdec23fecda7f0be6d0e3fa3ea/sagemaker-debugger/xgboost_realtime_analysis/data_utils.py"
urllib.request.urlretrieve(url, "data_utils.py")

from data_utils import load_mnist, upload_to_s3

prefix = "sagemaker/xgboost"
train_file, validation_file = load_mnist()
upload_to_s3(train_file, bucket, f"{prefix}/train/mnist.train.libsvm")
upload_to_s3(validation_file, bucket, f"{prefix}/validation/mnist.validation.libsvm")

# Remove downloaded file
import os

os.remove("data_utils.py")
