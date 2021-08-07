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
	"os"
	"os/exec"
)

var (
	TestDataDirectory = "testdata"
	DefaultTimestamp = "0001-01-01T00:00:00Z"
	ReplaceTimestampRegExp = "s/\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z\"/\"" + DefaultTimestamp + "\"/"
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

// Checks to see if the contents of a given yaml file, with name stored
// in expectation, matches the given actualYamlByteArray.
func IsYamlEqual(expectation *string, actualYamlByteArray *[]byte) bool {
     	// Get the file name of the expected yaml.
	expectedYamlFileName := TestDataDirectory + "/" + *expectation

	// Build a tmp file for the actual yaml.
	actualYamlFileName := buildTmpFile("actualYaml", *actualYamlByteArray)
	defer os.Remove(actualYamlFileName)
	if "" == actualYamlFileName {
	  fmt.Printf("Could not create temporary actual file.\n")
	  return false
	}

	// Replace Timestamps that would show up as different.
	_, err := exec.Command("sed", "-r", "-i", ReplaceTimestampRegExp, actualYamlFileName).Output()
	if isExecCommandError(err) {
	    return false
	}
	
	output,err := exec.Command("diff", "-c", expectedYamlFileName, actualYamlFileName).Output()
	if isExecCommandError(err) {
	   return false
	}

	if len(output) > 0 {
	   actualOutput,err := exec.Command("cat", actualYamlFileName).Output()
	   if isExecCommandError(err) {
	      return false
	   }
	   fmt.Printf("\nExpected Yaml File Name: " + expectedYamlFileName + "\n")
	   fmt.Printf("\nActual Output Yaml:\n" + string(actualOutput) + "\n")
	   fmt.Printf("Diff From Expected:\n" + string(output) + "\n")
	   return false
	}
	return true
}

func buildTmpFile(fileNameBase string, contents []byte) string {
     newTmpFile, err := ioutil.TempFile(TestDataDirectory, fileNameBase)
     if err != nil {
        fmt.Println(err)
	return ""
     }
     if _, err := newTmpFile.Write(contents); err != nil {
        fmt.Println(err)
	return ""
     }
     if err := newTmpFile.Close(); err != nil {
        fmt.Println(err)
	return ""
     }		
     return newTmpFile.Name()
}

// isExecCommandError returns true if an error
// that is not an ExitError is found.
func isExecCommandError(err error) bool {
     if err == nil {
     	return false
     }
     switch err.(type) {
	case *exec.ExitError:
	       // ExitError is expected.
	       return false
	default:
	       // Couldn't run diff.
	       fmt.Printf("Exec Command Error: ")
	       fmt.Println(err)
	       return true
     }
}
