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

package endpoint_config

import (
	"errors"
	"fmt"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	"github.com/aws-controllers-k8s/sagemaker-controller/pkg/testutil"
	mocksvcsdkapi "github.com/aws-controllers-k8s/sagemaker-controller/test/mocks/aws-sdk-go/sagemaker"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	ctrlrtzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
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

// TestEndpointConfigTestSuite runs the test suite for endpoint config.
func TestEndpointConfigTestSuite(t *testing.T) {
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
	case "CreateEndpointConfigWithContext":
		var output svcsdk.CreateEndpointConfigOutput
		return &output, nil
	case "DescribeEndpointConfigWithContext":
		var output svcsdk.DescribeEndpointConfigOutput
		return &output, nil
	case "DeleteEndpointConfigWithContext":
		var output svcsdk.DeleteEndpointConfigOutput
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
	// Build a tmp file for the actual yaml.
	actualResource := actual.(*resource)
	actualYamlByteArray, _ := yaml.Marshal(actualResource.ko)
	return testutil.IsYamlEqual(&expectation, &actualYamlByteArray)
}
