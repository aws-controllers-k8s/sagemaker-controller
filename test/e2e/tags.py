# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Utilities for working with tags"""

from typing import List

ACK_SYSTEM_TAG_PREFIX = "services.k8s.aws/"

def clean(tags: List[object]) -> List[object]:
    """Returns supplied tags collection, stripped of ACK system tags.
    """
    return [
        t for t in tags if not t['key'].startswith(ACK_SYSTEM_TAG_PREFIX)
    ]