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
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	mocksvcsdkapi "github.com/aws-controllers-k8s/sagemaker-controller/test/mocks/aws-sdk-go/sagemaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var RecoverPanicString = "\t--- PANIC ON ERROR:"

// TestSuiteRunner runs given test suite config with the help of delegate supplied to it
type TestSuiteRunner struct {
	TestSuite *TestSuite
	Delegate  TestRunnerDelegate
}

// fixtureContext is runtime context for test scenario given fixture.
type fixtureContext struct {
	desired         acktypes.AWSResource
	latest          acktypes.AWSResource
	mocksdkapi      *mocksvcsdkapi.SageMakerAPI
	resourceManager acktypes.AWSResourceManager
}

//TODO: remove if no longer used
// expectContext is runtime context for test scenario expectation fixture.
type expectContext struct {
	latest acktypes.AWSResource
	err    error
}

// TestRunnerDelegate provides interface for custom resource tests to implement.
// TestSuiteRunner depends on it to run tests for custom resource.
type TestRunnerDelegate interface {
	ResourceDescriptor() acktypes.AWSResourceDescriptor
	Equal(desired acktypes.AWSResource, latest acktypes.AWSResource) bool // remove it when ResourceDescriptor.Delta() is available
	YamlEqual(expected string, actual acktypes.AWSResource) bool          // new
	ResourceManager(*mocksvcsdkapi.SageMakerAPI) acktypes.AWSResourceManager
	EmptyServiceAPIOutput(apiName string) (interface{}, error)
	GoTestRunner() *testing.T
}

// RunTests runs the tests from the test suite
func (runner *TestSuiteRunner) RunTests() {
	if runner.TestSuite == nil || runner.Delegate == nil {
		panic(errors.New("failed to run test suite"))
	}

	for _, test := range runner.TestSuite.Tests {
		fmt.Printf("Starting test: %s\n", test.Name)
		for _, scenario := range test.Scenarios {
			runner.startScenario(scenario)
		}
		fmt.Printf("Test: %s completed.\n", test.Name)
	}
}

// Wrapper for running test scenarios that catches any panics thrown.
func (runner *TestSuiteRunner) startScenario(scenario TestScenario) {
	t := runner.Delegate.GoTestRunner()
	t.Run(scenario.Name, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(RecoverPanicString, r)
				t.Fail()
			}
		}()
		fmt.Printf("Running test scenario: %s\n", scenario.Name)
		fixtureCxt := runner.setupFixtureContext(&scenario.Fixture)
		runner.runTestScenario(t, scenario.Name, fixtureCxt, scenario.UnitUnderTest, &scenario.Expect)
	})
}

// runTestScenario runs given test scenario which is expressed as: given fixture context, unit to test, expected fixture context.
func (runner *TestSuiteRunner) runTestScenario(t *testing.T, scenarioName string, fixtureCxt *fixtureContext, unitUnderTest string, expectation *Expect) {
	rm := fixtureCxt.resourceManager
	assert := assert.New(t)

	var actual acktypes.AWSResource = nil
	var err error = nil
	switch unitUnderTest {
	case "ReadOne":
		actual, err = rm.ReadOne(context.Background(), fixtureCxt.desired)
	case "Create":
		actual, err = rm.Create(context.Background(), fixtureCxt.desired)
	case "Update":
		delta := runner.Delegate.ResourceDescriptor().Delta(fixtureCxt.desired, fixtureCxt.latest)
		actual, err = rm.Update(context.Background(), fixtureCxt.desired, fixtureCxt.latest, delta)
	case "Delete":
		actual, err = rm.Delete(context.Background(), fixtureCxt.desired)
	default:
		panic(errors.New(fmt.Sprintf("unit under test: %s not supported", unitUnderTest)))
	}
	runner.assertExpectations(assert, expectation, actual, err)
}

/* assertExpectations validates the actual outcome against the expected outcome.
There are two components to the expected outcome, corresponding to the return values of the resource manager's CRUD operation:
	1) the actual return value of type AWSResource ("expect.latest_state" in test_suite.yaml)
	2) the error ("expect.error" in test_suite.yaml)
With each of these components, there are three possibilities in test_suite.yaml, which are interpreted as follows:
	1) the key does not exist, or was provided with no value: no explicit expectations, don't assert anything
	2) the key was provided with value "nil": explicit expectation; assert that the error or return value is nil
	3) the key was provided with value other than "nil": explicit expectation; assert that the value matches the
		expected value
However, if neither expect.latest_state nor error are provided, assertExpectations will fail the test case.
*/
func (runner *TestSuiteRunner) assertExpectations(assert *assert.Assertions, expectation *Expect, actual acktypes.AWSResource, err error) {
	if expectation.LatestState == "" && expectation.Error == "" {
		fmt.Println("Invalid test case: no expectation given for either latest_state or error")
		assert.True(false)
		return
	}

	// expectation exists for at least one of LatestState and Error; assert results independently
	if expectation.LatestState == "nil" {
		assert.Nil(actual)
	} else if expectation.LatestState != "" {
		expectedLatest := runner.loadAWSResource(expectation.LatestState)
		assert.NotNil(actual)

		delta := runner.Delegate.ResourceDescriptor().Delta(expectedLatest, actual)
		assert.Equal(0, len(delta.Differences))
		if len(delta.Differences) > 0 {
			fmt.Println("Unexpected differences:")
			for _, difference := range delta.Differences {
				fmt.Printf("Path: %v, expected: %v, actual: %v\n", difference.Path, difference.A, difference.B)
				fmt.Printf("See expected differences below:\n")
			}
		}

		// Check that the yaml files are equivalent.
		// This makes it easier to make changes to unit test cases.
		assert.True(runner.Delegate.YamlEqual(expectation.LatestState, actual))
		// Delta only contains `Spec` differences. Thus, we need Delegate.Equal to compare `Status`.
		assert.True(runner.Delegate.Equal(expectedLatest, actual))
	}

	if expectation.Error == "nil" {
		assert.Nil(err)
	} else if expectation.Error != "" {
		expectedError := errors.New(expectation.Error)
		assert.NotNil(err)

		assert.Equal(expectedError.Error(), err.Error())
	}
}

// setupFixtureContext provides runtime context for test scenario given fixture.
func (runner *TestSuiteRunner) setupFixtureContext(fixture *Fixture) *fixtureContext {
	if fixture == nil {
		return nil
	}
	var cxt = fixtureContext{}
	if fixture.DesiredState != "" {
		cxt.desired = runner.loadAWSResource(fixture.DesiredState)
	}
	if fixture.LatestState != "" {
		cxt.latest = runner.loadAWSResource(fixture.LatestState)
	}
	mocksdkapi := &mocksvcsdkapi.SageMakerAPI{}
	for _, serviceApi := range fixture.ServiceAPIs {
		if serviceApi.Operation != "" {

			if serviceApi.ServiceAPIError != nil {
				mockError := CreateAWSError(*serviceApi.ServiceAPIError)
				mocksdkapi.On(serviceApi.Operation, mock.Anything, mock.Anything).Return(nil, mockError)
			} else if serviceApi.Operation != "" && serviceApi.Output != "" {
				var outputObj, err = runner.Delegate.EmptyServiceAPIOutput(serviceApi.Operation)
				apiOutputFixturePath := append([]string{"testdata"}, strings.Split(serviceApi.Output, "/")...)
				LoadFromFixture(filepath.Join(apiOutputFixturePath...), outputObj)
				mocksdkapi.On(serviceApi.Operation, mock.Anything, mock.Anything).Return(outputObj, nil)
				if err != nil {
					panic(err)
				}
			} else if serviceApi.ServiceAPIError == nil && serviceApi.Output == "" {
				// Default case for no defined output fixture or error.
				mocksdkapi.On(serviceApi.Operation, mock.Anything, mock.Anything).Return(nil, nil)
			}
		}
	}
	cxt.mocksdkapi = mocksdkapi
	cxt.resourceManager = runner.Delegate.ResourceManager(mocksdkapi)
	return &cxt
}

// loadAWSResource loads AWSResource from the supplied fixture file.
func (runner *TestSuiteRunner) loadAWSResource(resourceFixtureFilePath string) acktypes.AWSResource {
	if resourceFixtureFilePath == "" {
		panic(errors.New(fmt.Sprintf("resourceFixtureFilePath not specified")))
	}
	var rd = runner.Delegate.ResourceDescriptor()
	ro := rd.EmptyRuntimeObject()
	path := append([]string{"testdata"}, strings.Split(resourceFixtureFilePath, "/")...)
	LoadFromFixture(filepath.Join(path...), ro)
	return rd.ResourceFromRuntimeObject(ro)
}
