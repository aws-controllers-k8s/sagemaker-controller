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

package feature_group

import (
	"errors"
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	"github.com/ghodss/yaml"
	mocksvcsdkapi "github.com/aws-controllers-k8s/sagemaker-controller/test/mocks/aws-sdk-go/sagemaker"
	"github.com/aws-controllers-k8s/sagemaker-controller/pkg/testutil"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"path/filepath"
	"testing"
	ctrlrtzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultTimestamp = "0001-01-01T00:00:00Z"
	ReplaceTimestampRegExp = "s/\"[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z\"/\"" + DefaultTimestamp + "\"/"
	TestDataDirectory = "testdata"
)

// provideResourceManagerWithMockSDKAPI accepts MockSageMakerAPI and returns pointer to resourceManager
// the returned resourceManager is configured to use mockapi api.
func provideResourceManagerWithMockSDKAPI(mockSageMakerAPI *mocksvcsdkapi.SageMakerAPI) *resourceManager {
	zapOptions := ctrlrtzap.Options{
		Development: true,
		Level:       zapcore.InfoLevel,
	}
	fakeLogger := ctrlrtzap.New(ctrlrtzap.UseFlagOptions(&zapOptions))
	return &resourceManager{
		rr:           nil,
		awsAccountID: "",
		awsRegion:    "",
		sess:         nil,
		sdkapi:       mockSageMakerAPI,
		log:          fakeLogger,
		metrics:      ackmetrics.NewMetrics("sagemaker"),
	}
}

// TestFeatureGroupTestSuite runs the test suite for feature group
func TestFeatureGroupTestSuite(t *testing.T) {
     	defer func() {
     	   if r := recover(); r != nil {
     	      fmt.Println(testutil.RecoverPanicString, r)
	      t.Fail()
     	   }
        }()
	var ts = testutil.TestSuite{}
	testutil.LoadFromFixture(filepath.Join("testdata", "test_suite.yaml"), &ts)
	var delegate = testRunnerDelegate{t: t}
	var runner = testutil.TestSuiteRunner{TestSuite: &ts, Delegate: &delegate}
	runner.RunTests()
}

// testRunnerDelegate implements testutil.TestRunnerDelegate
type testRunnerDelegate struct {
	t *testing.T
}

func (d *testRunnerDelegate) ResourceDescriptor() acktypes.AWSResourceDescriptor {
	return &resourceDescriptor{}
}

func (d *testRunnerDelegate) ResourceManager(mocksdkapi *mocksvcsdkapi.SageMakerAPI) acktypes.AWSResourceManager {
	return provideResourceManagerWithMockSDKAPI(mocksdkapi)
}

func (d *testRunnerDelegate) GoTestRunner() *testing.T {
	return d.t
}

func (d *testRunnerDelegate) EmptyServiceAPIOutput(apiName string) (interface{}, error) {
	if apiName == "" {
		return nil, errors.New("no API name specified")
	}
	//TODO: use reflection, template to auto generate this block/method.
	switch apiName {
	case "CreateFeatureGroupWithContext":
		var output svcsdk.CreateFeatureGroupOutput
		return &output, nil
	case "DeleteFeatureGroupWithContext":
		var output svcsdk.DeleteFeatureGroupOutput
		return &output, nil
	case "DescribeFeatureGroupWithContext":
		var output svcsdk.DescribeFeatureGroupOutput
		return &output, nil
	}
	return nil, errors.New(fmt.Sprintf("no matching API name found for: %s", apiName))
}

func (d *testRunnerDelegate) Equal(a acktypes.AWSResource, b acktypes.AWSResource) bool {
	ac := a.(*resource)
	bc := b.(*resource)
	// Ignore LastTransitionTime since it gets updated each run.
	opts := []cmp.Option{cmpopts.EquateEmpty(), cmpopts.IgnoreFields(ackv1alpha1.Condition{}, "LastTransitionTime")}

	if cmp.Equal(ac.ko.Status, bc.ko.Status, opts...) {
		return true
	} else {
		fmt.Printf("Difference (-expected +actual):\n\n")
		fmt.Println(cmp.Diff(ac.ko.Status, bc.ko.Status, opts...))
		return false
	}
}

// Checks to see if the given yaml file, with name stored as expectation,
// matches the yaml marshal of the AWSResource stored as actual.
func (d *testRunnerDelegate) YamlEqual(expectation string, actual acktypes.AWSResource) bool {
     	// Get the file name of the expected yaml.
	expectedYamlFileName := TestDataDirectory + "/" + expectation

	// Build a tmp file for the actual yaml.
	actualResource := actual.(*resource)
	actualYamlByteArray, _ := yaml.Marshal(actualResource.ko)
	actualYamlFileName := buildTmpFile("actualYaml", actualYamlByteArray)
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
