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

// TestSuite represents instructions to run unit tests using test fixtures and mock service apis
type TestSuite struct {
	Tests []TestConfig `json:"tests"`
}

// TestConfig represents declarative unit test
type TestConfig struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Scenarios   []TestScenario `json:"scenarios"`
}

// TestScenario represents declarative test scenario details
type TestScenario struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Fixture lets you specify test scenario given input fixtures
	Fixture Fixture `json:"given"`
	// UnitUnderTest lets you specify the unit to test
	// For example resource manager API: ReadOne, Create, Update, Delete
	UnitUnderTest string `json:"invoke"`
	// Expect lets you specify test scenario expected outcome fixtures
	Expect Expect `json:"expect"`
}

// Fixture represents test scenario fixture to load from file paths
type Fixture struct {
	// DesiredState lets you specify fixture path to load the desired state fixture
	DesiredState string `json:"desired_state"`
	// LatestState lets you specify fixture path to load the current state fixture
	LatestState string `json:"latest_state"`
	// ServiceAPIs lets you specify fixture path to mock service sdk api response
	ServiceAPIs []ServiceAPI `json:"svc_api"`
}

// ServiceAPI represents details about the the service sdk api and fixture path to mock its response
type ServiceAPI struct {
	Operation       string           `json:"operation"`
	Output          string           `json:"output_fixture"`
	ServiceAPIError *ServiceAPIError `json:"error,omitempty"`
}

// ServiceAPIError contains the specification for the error of the mock API response
type ServiceAPIError struct {
	// Code here is usually the type of fault/error, not the HTTP status code
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Expect represents test scenario expected outcome fixture to load from file path
type Expect struct {
	LatestState string `json:"latest_state"`
	// Error is a string matching the message of the expected error returned from the ResourceManager operation.
	// Possible errors can be found in runtime/pkg/errors/error.go
	Error string `json:"error"`
}
