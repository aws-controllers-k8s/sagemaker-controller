# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# 	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.
"""Stores the values used by each of the integration tests for replacing the
SageMaker-specific test variables.
"""

from acktest.aws.identity import get_region
from e2e.bootstrap_resources import get_bootstrap_resources

# Taken from the SageMaker Python SDK
# Rather than including the entire SDK
XGBOOST_IMAGE_URIS = {
    "us-west-1": "746614075791.dkr.ecr.us-west-1.amazonaws.com",
    "us-west-2": "246618743249.dkr.ecr.us-west-2.amazonaws.com",
    "us-east-1": "683313688378.dkr.ecr.us-east-1.amazonaws.com",
    "us-east-2": "257758044811.dkr.ecr.us-east-2.amazonaws.com",
    "ap-east-1": "651117190479.dkr.ecr.ap-east-1.amazonaws.com",
    "ap-northeast-1": "354813040037.dkr.ecr.ap-northeast-1.amazonaws.com",
    "ap-northeast-2": "366743142698.dkr.ecr.ap-northeast-2.amazonaws.com",
    "ap-south-1": "720646828776.dkr.ecr.ap-south-1.amazonaws.com",
    "ap-southeast-1": "121021644041.dkr.ecr.ap-southeast-1.amazonaws.com",
    "ap-southeast-2": "783357654285.dkr.ecr.ap-southeast-2.amazonaws.com",
    "ca-central-1": "341280168497.dkr.ecr.ca-central-1.amazonaws.com",
    "cn-north-1": "450853457545.dkr.ecr.cn-north-1.amazonaws.com.cn",
    "cn-northwest-1": "451049120500.dkr.ecr.cn-northwest-1.amazonaws.com.cn",
    "eu-central-1": "492215442770.dkr.ecr.eu-central-1.amazonaws.com",
    "eu-north-1": "662702820516.dkr.ecr.eu-north-1.amazonaws.com",
    "eu-west-1": "141502667606.dkr.ecr.eu-west-1.amazonaws.com",
    "eu-west-2": "764974769150.dkr.ecr.eu-west-2.amazonaws.com",
    "eu-west-3": "659782779980.dkr.ecr.eu-west-3.amazonaws.com",
    "me-south-1": "801668240914.dkr.ecr.me-south-1.amazonaws.com",
    "sa-east-1": "737474898029.dkr.ecr.sa-east-1.amazonaws.com",
}

DEBUGGER_IMAGE_URIS = {
    "us-west-1": "685455198987.dkr.ecr.us-west-1.amazonaws.com",
    "us-west-2": "895741380848.dkr.ecr.us-west-2.amazonaws.com",
    "us-east-1": "503895931360.dkr.ecr.us-east-1.amazonaws.com",
    "us-east-2": "915447279597.dkr.ecr.us-east-2.amazonaws.com",
    "ap-east-1": "199566480951.dkr.ecr.ap-east-1.amazonaws.com",
    "ap-northeast-1": "430734990657.dkr.ecr.ap-northeast-1.amazonaws.com",
    "ap-northeast-2": "578805364391.dkr.ecr.ap-northeast-2.amazonaws.com",
    "ap-south-1": "904829902805.dkr.ecr.ap-south-1.amazonaws.com",
    "ap-southeast-1": "972752614525.dkr.ecr.ap-southeast-1.amazonaws.com",
    "ap-southeast-2": "184798709955.dkr.ecr.ap-southeast-2.amazonaws.com",
    "ca-central-1": "519511493484.dkr.ecr.ca-central-1.amazonaws.com",
    "cn-north-1": "618459771430.dkr.ecr.cn-north-1.amazonaws.com.cn",
    "cn-northwest-1": "658757709296.dkr.ecr.cn-northwest-1.amazonaws.com.cn",
    "eu-central-1": "482524230118.dkr.ecr.eu-central-1.amazonaws.com",
    "eu-north-1": "314864569078.dkr.ecr.eu-north-1.amazonaws.com",
    "eu-west-1": "929884845733.dkr.ecr.eu-west-1.amazonaws.com",
    "eu-west-2": "250201462417.dkr.ecr.eu-west-2.amazonaws.com",
    "eu-west-3": "447278800020.dkr.ecr.eu-west-3.amazonaws.com",
    "me-south-1": "986000313247.dkr.ecr.me-south-1.amazonaws.com",
    "sa-east-1": "818342061345.dkr.ecr.sa-east-1.amazonaws.com",
}

# https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
XGBOOST_V1_IMAGE_URIS = {
    "us-west-1": "632365934929.dkr.ecr.us-west-1.amazonaws.com",
    "us-west-2": "433757028032.dkr.ecr.us-west-2.amazonaws.com",
    "us-east-1": "811284229777.dkr.ecr.us-east-1.amazonaws.com",
    "us-east-2": "825641698319.dkr.ecr.us-east-2.amazonaws.com",
    "ap-east-1": "286214385809.dkr.ecr.ap-east-1.amazonaws.com",
    "ap-northeast-1": "501404015308.dkr.ecr.ap-northeast-1.amazonaws.com",
    "ap-northeast-2": "306986355934.dkr.ecr.ap-northeast-2.amazonaws.com",
    "ap-south-1": "991648021394.dkr.ecr.ap-south-1.amazonaws.com",
    "ap-southeast-1": "475088953585.dkr.ecr.ap-southeast-1.amazonaws.com",
    "ap-southeast-2": "544295431143.dkr.ecr.ap-southeast-2.amazonaws.com",
    "ca-central-1": "469771592824.dkr.ecr.ca-central-1.amazonaws.com",
    "cn-north-1": "390948362332.dkr.ecr.cn-north-1.amazonaws.com",
    "cn-northwest-1": "387376663083.dkr.ecr.cn-northwest-1.amazonaws.com",
    "eu-central-1": "813361260812.dkr.ecr.eu-central-1.amazonaws.com",
    "eu-north-1": "669576153137.dkr.ecr.eu-north-1.amazonaws.com",
    "eu-west-1": "685385470294.dkr.ecr.eu-west-1.amazonaws.com",
    "eu-west-2": "644912444149.dkr.ecr.eu-west-2.amazonaws.com",
    "eu-west-3": "749696950732.dkr.ecr.eu-west-3.amazonaws.com",
    "me-south-1": "249704162688.dkr.ecr.me-south-1.amazonaws.com",
    "sa-east-1": "855470959533.dkr.ecr.sa-east-1.amazonaws.com",
}


PYTORCH_TRAIN_IMAGE_URIS = {
    "us-east-1": "763104351884.dkr.ecr.us-east-1.amazonaws.com",
    "us-east-2": "763104351884.dkr.ecr.us-east-2.amazonaws.com",
    "us-west-1": "763104351884.dkr.ecr.us-west-1.amazonaws.com",
    "us-west-2": "763104351884.dkr.ecr.us-west-2.amazonaws.com",
    "af-south-1": "626614931356.dkr.ecr.af-south-1.amazonaws.com",
    "ap-east-1": "871362719292.dkr.ecr.ap-east-1.amazonaws.com",
    "ap-south-1": "763104351884.dkr.ecr.ap-south-1.amazonaws.com",
    "ap-northeast-2": "763104351884.dkr.ecr.ap-northeast-2.amazonaws.com",
    "ap-southeast-1": "763104351884.dkr.ecr.ap-southeast-1.amazonaws.com",
    "ap-southeast-2": "763104351884.dkr.ecr.ap-southeast-2.amazonaws.com",
    "ap-northeast-1": "763104351884.dkr.ecr.ap-northeast-1.amazonaws.com",
    "ca-central-1": "763104351884.dkr.ecr.ca-central-1.amazonaws.com",
    "eu-central-1": "763104351884.dkr.ecr.eu-central-1.amazonaws.com",
    "eu-west-1": "763104351884.dkr.ecr.eu-west-1.amazonaws.com",
    "eu-west-2": "763104351884.dkr.ecr.eu-west-2.amazonaws.com",
    "eu-south-1": "692866216735.dkr.ecr.eu-south-1.amazonaws.com",
    "eu-west-3": "763104351884.dkr.ecr.eu-west-3.amazonaws.com",
    "eu-north-1": "763104351884.dkr.ecr.eu-north-1.amazonaws.com",
    "me-south-1": "217643126080.dkr.ecr.me-south-1.amazonaws.com",
    "sa-east-1": "763104351884.dkr.ecr.sa-east-1.amazonaws.com",
    "cn-north-1": "727897471807.dkr.ecr.cn-north-1.amazonaws.com.cn",
    "cn-northwest-1": "727897471807.dkr.ecr.cn-northwest-1.amazonaws.com.cn",
}

# https://docs.aws.amazon.com/sagemaker/latest/dg/model-monitor-pre-built-container.html
MODEL_MONITOR_IMAGE_URIS = {
    "us-east-1": "156813124566.dkr.ecr.us-east-1.amazonaws.com",
    "us-east-2": "777275614652.dkr.ecr.us-east-2.amazonaws.com",
    "us-west-1": "890145073186.dkr.ecr.us-west-1.amazonaws.com",
    "us-west-2": "159807026194.dkr.ecr.us-west-2.amazonaws.com",
    "af-south-1": "875698925577.dkr.ecr.af-south-1.amazonaws.com",
    "ap-east-1": "001633400207.dkr.ecr.ap-east-1.amazonaws.com",
    "ap-northeast-1": "574779866223.dkr.ecr.ap-northeast-1.amazonaws.com",
    "ap-northeast-2": "709848358524.dkr.ecr.ap-northeast-2.amazonaws.com",
    "ap-south-1": "126357580389.dkr.ecr.ap-south-1.amazonaws.com",
    "ap-southeast-1": "245545462676.dkr.ecr.ap-southeast-1.amazonaws.com",
    "ap-southeast-2": "563025443158.dkr.ecr.ap-southeast-2.amazonaws.com",
    "ca-central-1": "536280801234.dkr.ecr.ca-central-1.amazonaws.com",
    "cn-north-1": "453000072557.dkr.ecr.cn-north-1.amazonaws.com.cn",
    "cn-northwest-1": "453252182341.dkr.ecr.cn-northwest-1.amazonaws.com.cn",
    "eu-central-1": "048819808253.dkr.ecr.eu-central-1.amazonaws.com",
    "eu-north-1": "895015795356.dkr.ecr.eu-north-1.amazonaws.com",
    "eu-south-1": "933208885752.dkr.ecr.eu-south-1.amazonaws.com",
    "eu-west-1": "468650794304.dkr.ecr.eu-west-1.amazonaws.com",
    "eu-west-2": "749857270468.dkr.ecr.eu-west-2.amazonaws.com",
    "eu-west-3": "680080141114.dkr.ecr.eu-west-3.amazonaws.com",
    "me-south-1": "607024016150.dkr.ecr.me-south-1.amazonaws.com",
    "sa-east-1": "539772159869.dkr.ecr.sa-east-1.amazonaws.com",
    "us-gov-west-1": "362178532790.dkr.ecr.us-gov-west-1.amazonaws.com",
}

# https://docs.aws.amazon.com/sagemaker/latest/dg/clarify-configure-processing-jobs.html#clarify-processing-job-configure-container
CLARIFY_IMAGE_URIS = {
    "us-east-1": "205585389593.dkr.ecr.us-east-1.amazonaws.com",
    "us-east-2": "211330385671.dkr.ecr.us-east-2.amazonaws.com",
    "us-west-1": "740489534195.dkr.ecr.us-west-1.amazonaws.com",
    "us-west-2": "306415355426.dkr.ecr.us-west-2.amazonaws.com",
    "ap-east-1": "098760798382.dkr.ecr.ap-east-1.amazonaws.com",
    "ap-south-1": "452307495513.dkr.ecr.ap-south-1.amazonaws.com",
    "ap-northeast-2": "263625296855.dkr.ecr.ap-northeast-2.amazonaws.com",
    "ap-southeast-1": "834264404009.dkr.ecr.ap-southeast-1.amazonaws.com",
    "ap-southeast-2": "007051062584.dkr.ecr.ap-southeast-2.amazonaws.com",
    "ap-northeast-1": "377024640650.dkr.ecr.ap-northeast-1.amazonaws.com",
    "ca-central-1": "675030665977.dkr.ecr.ca-central-1.amazonaws.com",
    "eu-central-1": "017069133835.dkr.ecr.eu-central-1.amazonaws.com",
    "eu-west-1": "131013547314.dkr.ecr.eu-west-1.amazonaws.com",
    "eu-west-2": "440796970383.dkr.ecr.eu-west-2.amazonaws.com",
    "eu-west-3": "341593696636.dkr.ecr.eu-west-3.amazonaws.com",
    "eu-north-1": "763603941244.dkr.ecr.eu-north-1.amazonaws.com",
    "me-south-1": "835444307964.dkr.ecr.me-south-1.amazonaws.com",
    "sa-east-1": "520018980103.dkr.ecr.sa-east-1.amazonaws.com",
    "af-south-1": "811711786498.dkr.ecr.af-south-1.amazonaws.com",
    "eu-south-1": "638885417683.dkr.ecr.eu-south-1.amazonaws.com",
}

ENDPOINT_INSTANCE_TYPES = {
    "us-east-1": "ml.c5.large",
    "us-east-2": "ml.c5.large",
    "us-west-1": "ml.c5.large",
    "us-west-2": "ml.c5.large",
    "ap-east-1": "ml.c5.large",
    "ap-south-1": "ml.c5.large",
    "ap-northeast-2": "ml.c5.large",
    "ap-southeast-1": "ml.c5.large",
    "ap-southeast-2": "ml.c5.large",
    "ap-northeast-1": "ml.c5.large",
    "ca-central-1": "ml.c5.large",
    "eu-central-1": "ml.c5.large",
    "eu-west-1": "ml.c5.large",
    "eu-west-2": "ml.c5.large",
    "eu-west-3": "ml.m5.large",
    "eu-north-1": "ml.m5.large",
    "me-south-1": "ml.c5.large",
    "sa-east-1": "ml.c5.large",
    "af-south-1": "ml.c5.large",
    "eu-south-1": "ml.c5.large",
}

TRAINING_JOB_INSTANCE_TYPES = {
    "us-east-1": "ml.m4.xlarge",
    "us-east-2": "ml.m4.xlarge",
    "us-west-1": "ml.m4.xlarge",
    "us-west-2": "ml.m4.xlarge",
    "ap-east-1": "ml.m4.xlarge",
    "ap-south-1": "ml.m4.xlarge",
    "ap-northeast-2": "ml.m4.xlarge",
    "ap-southeast-1": "ml.m4.xlarge",
    "ap-southeast-2": "ml.m4.xlarge",
    "ap-northeast-1": "ml.m4.xlarge",
    "ca-central-1": "ml.m4.xlarge",
    "eu-central-1": "ml.m4.xlarge",
    "eu-west-1": "ml.m4.xlarge",
    "eu-west-2": "ml.m4.xlarge",
    "eu-west-3": "ml.m5.xlarge",
    "eu-north-1": "ml.m5.xlarge",
    "me-south-1": "ml.m4.xlarge",
    "sa-east-1": "ml.m4.xlarge",
    "af-south-1": "ml.m4.xlarge",
    "eu-south-1": "ml.m4.xlarge",
}

REPLACEMENT_VALUES = {
    "SAGEMAKER_DATA_BUCKET": get_bootstrap_resources().DataBucketName,
    "XGBOOST_IMAGE_URI": f"{XGBOOST_IMAGE_URIS[get_region()]}/sagemaker-xgboost:1.0-1-cpu-py3",
    "DEBUGGER_IMAGE_URI": f"{DEBUGGER_IMAGE_URIS[get_region()]}/sagemaker-debugger-rules:latest",
    "XGBOOST_V1_IMAGE_URI": f"{XGBOOST_V1_IMAGE_URIS[get_region()]}/xgboost:latest",
    "PYTORCH_TRAIN_IMAGE_URI": f"{PYTORCH_TRAIN_IMAGE_URIS[get_region()]}/pytorch-training:1.5.0-cpu-py36-ubuntu16.04",
    "SAGEMAKER_EXECUTION_ROLE_ARN": get_bootstrap_resources().ExecutionRoleARN,
    "MODEL_MONITOR_ANALYZER_IMAGE_URI": f"{MODEL_MONITOR_IMAGE_URIS[get_region()]}/sagemaker-model-monitor-analyzer",
    "CLARIFY_IMAGE_URI": f"{CLARIFY_IMAGE_URIS[get_region()]}/sagemaker-clarify-processing:1.0",
    "ENDPOINT_INSTANCE_TYPE": ENDPOINT_INSTANCE_TYPES[get_region()],
    "TRAINING_JOB_INSTANCE_TYPE": TRAINING_JOB_INSTANCE_TYPES[get_region()],
}
