// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package testutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"path"
	"strings"
)

// LoadFromFixture fills an empty pointer variable with the
// data from a fixture JSON/YAML file.
func LoadFromFixture(
	fixturePath string,
	output interface{}, // output should be an addressable type (i.e. a pointer)
) {
	contents, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(fixturePath, ".json") {
		err = json.Unmarshal(contents, output)
	} else if strings.HasSuffix(fixturePath, ".yaml") ||
		strings.HasSuffix(fixturePath, ".yml") {
		err = yaml.Unmarshal(contents, output)
	} else {
		panic(errors.New(
			fmt.Sprintf("fixture file format not supported: %s", path.Base(fixturePath))))
	}
	if err != nil {
		panic(err)
	}
}

// CreateAWSError is used for mocking the types of errors received from aws-sdk-go
// so that the expected code path executes. Support for specifying the HTTP status code and request ID
// can be added in the future if needed
func CreateAWSError(awsError ServiceAPIError) awserr.RequestFailure {
	error := awserr.New(awsError.Code, awsError.Message, nil)
	return awserr.NewRequestFailure(error, 0, "")
}
