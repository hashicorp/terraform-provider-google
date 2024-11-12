// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwvalidators_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFrameworkProvider_CredentialsValidator(t *testing.T) {
	cases := map[string]struct {
		ConfigValue          types.String
		ExpectedWarningCount int
		ExpectedErrorCount   int
	}{
		"configuring credentials as a path to a credentials JSON file is valid": {
			ConfigValue: types.StringValue(transport_tpg.TestFakeCredentialsPath), // Path to a test fixture
		},
		"configuring credentials as a path to a non-existant file is NOT valid": {
			ConfigValue:        types.StringValue("./this/path/doesnt/exist.json"), // Doesn't exist
			ExpectedErrorCount: 1,
		},
		"configuring credentials as a credentials JSON string is valid": {
			ConfigValue: types.StringValue(acctest.GenerateFakeCredentialsJson("CredentialsValidator")),
		},
		"configuring credentials as an empty string is not valid": {
			ConfigValue:        types.StringValue(""),
			ExpectedErrorCount: 1,
		},
		"leaving credentials unconfigured is valid": {
			ConfigValue: types.StringNull(),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			req := validator.StringRequest{
				ConfigValue: tc.ConfigValue,
			}

			resp := validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			cv := fwvalidators.CredentialsValidator()

			// Act
			cv.ValidateString(context.Background(), req, &resp)

			// Assert
			if resp.Diagnostics.WarningsCount() > tc.ExpectedWarningCount {
				t.Errorf("Expected %d warnings, got %d", tc.ExpectedWarningCount, resp.Diagnostics.WarningsCount())
			}
			if resp.Diagnostics.ErrorsCount() > tc.ExpectedErrorCount {
				t.Errorf("Expected %d errors, got %d", tc.ExpectedErrorCount, resp.Diagnostics.ErrorsCount())
			}
		})
	}
}

func TestServiceAccountEmailValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		value         types.String
		expectError   bool
		errorContains string
	}

	tests := map[string]testCase{
		"correct service account name": {
			value:       types.StringValue("test@test.iam.gserviceaccount.com"),
			expectError: false,
		},
		"developer service account": {
			value:       types.StringValue("test@developer.gserviceaccount.com"),
			expectError: false,
		},
		"app engine service account": {
			value:       types.StringValue("test@appspot.gserviceaccount.com"),
			expectError: false,
		},
		"cloud services service account": {
			value:       types.StringValue("test@cloudservices.gserviceaccount.com"),
			expectError: false,
		},
		"cloud build service account": {
			value:       types.StringValue("test@cloudbuild.gserviceaccount.com"),
			expectError: false,
		},
		"compute engine service account": {
			value:       types.StringValue("service-123456@compute-system.iam.gserviceaccount.com"),
			expectError: false,
		},
		"incorrect service account name": {
			value:         types.StringValue("test"),
			expectError:   true,
			errorContains: "Service account name must match one of the expected patterns for Google service accounts",
		},
		"empty string": {
			value:         types.StringValue(""),
			expectError:   true,
			errorContains: "Service account name must not be empty",
		},
		"null value": {
			value:       types.StringNull(),
			expectError: false,
		},
		"unknown value": {
			value:       types.StringUnknown(),
			expectError: false,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.value,
			}
			response := validator.StringResponse{}
			validator := fwvalidators.ServiceAccountEmailValidator{}

			validator.ValidateString(context.Background(), request, &response)

			if test.expectError && !response.Diagnostics.HasError() {
				t.Errorf("expected error, got none")
			}

			if !test.expectError && response.Diagnostics.HasError() {
				t.Errorf("got unexpected error: %s", response.Diagnostics.Errors())
			}

			if test.errorContains != "" {
				foundError := false
				for _, err := range response.Diagnostics.Errors() {
					if err.Detail() == test.errorContains {
						foundError = true
						break
					}
				}
				if !foundError {
					t.Errorf("expected error with summary %q, got none", test.errorContains)
				}
			}
		})
	}
}

func TestBoundedDuration(t *testing.T) {
	t.Parallel()

	type testCase struct {
		value         types.String
		minDuration   time.Duration
		maxDuration   time.Duration
		expectError   bool
		errorContains string
	}

	tests := map[string]testCase{
		"valid duration between min and max": {
			value:       types.StringValue("1800s"),
			minDuration: time.Hour / 2,
			maxDuration: time.Hour,
			expectError: false,
		},
		"valid duration at min": {
			value:       types.StringValue("1800s"),
			minDuration: 30 * time.Minute,
			maxDuration: time.Hour,
			expectError: false,
		},
		"valid duration at max": {
			value:       types.StringValue("3600s"),
			minDuration: time.Hour / 2,
			maxDuration: time.Hour,
			expectError: false,
		},
		"valid duration with different unit": {
			value:       types.StringValue("1h"),
			minDuration: 30 * time.Minute,
			maxDuration: 2 * time.Hour,
			expectError: false,
		},
		"duration below min": {
			value:         types.StringValue("900s"),
			minDuration:   30 * time.Minute,
			maxDuration:   time.Hour,
			expectError:   true,
			errorContains: "Invalid Duration",
		},
		"duration exceeds max - seconds": {
			value:         types.StringValue("7200s"),
			minDuration:   30 * time.Minute,
			maxDuration:   time.Hour,
			expectError:   true,
			errorContains: "Invalid Duration",
		},
		"duration exceeds max - minutes": {
			value:         types.StringValue("120m"),
			minDuration:   30 * time.Minute,
			maxDuration:   time.Hour,
			expectError:   true,
			errorContains: "Invalid Duration",
		},
		"duration exceeds max - hours": {
			value:         types.StringValue("2h"),
			minDuration:   30 * time.Minute,
			maxDuration:   time.Hour,
			expectError:   true,
			errorContains: "Invalid Duration",
		},
		"invalid duration format": {
			value:         types.StringValue("invalid"),
			minDuration:   30 * time.Minute,
			maxDuration:   time.Hour,
			expectError:   true,
			errorContains: "Invalid Duration Format",
		},
		"setting min to 0": {
			value:       types.StringValue("10s"),
			minDuration: 0,
			maxDuration: time.Hour,
			expectError: false,
		},
		"setting max to be less than min": {
			value:         types.StringValue("10s"),
			minDuration:   30 * time.Minute,
			maxDuration:   10 * time.Second,
			expectError:   true,
			errorContains: "Invalid Duration",
		},
		"empty string": {
			value:         types.StringValue(""),
			minDuration:   30 * time.Minute,
			maxDuration:   time.Hour,
			expectError:   true,
			errorContains: "Invalid Duration Format",
		},
		"null value": {
			value:       types.StringNull(),
			minDuration: 30 * time.Minute,
			maxDuration: time.Hour,
			expectError: false,
		},
		"unknown value": {
			value:       types.StringUnknown(),
			minDuration: 30 * time.Minute,
			maxDuration: time.Hour,
			expectError: false,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.value,
			}
			response := validator.StringResponse{}
			validator := fwvalidators.BoundedDuration{
				MinDuration: test.minDuration,
				MaxDuration: test.maxDuration,
			}

			validator.ValidateString(context.Background(), request, &response)

			if test.expectError && !response.Diagnostics.HasError() {
				t.Errorf("expected error, got none")
			}

			if !test.expectError && response.Diagnostics.HasError() {
				t.Errorf("got unexpected error: %s", response.Diagnostics.Errors())
			}

			if test.errorContains != "" {
				foundError := false
				for _, err := range response.Diagnostics.Errors() {
					if err.Summary() == test.errorContains {
						foundError = true
						break
					}
				}
				if !foundError {
					t.Errorf("expected error with summary %q, got none", test.errorContains)
				}
			}
		})
	}
}
