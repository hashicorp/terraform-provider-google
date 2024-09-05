// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-google/google/fwprovider"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFrameworkProvider_CredentialsValidator(t *testing.T) {
	cases := map[string]struct {
		ConfigValue          func(t *testing.T) types.String
		ExpectedWarningCount int
		ExpectedErrorCount   int
	}{
		"configuring credentials as a path to a credentials JSON file is valid": {
			ConfigValue: func(t *testing.T) types.String {
				return types.StringValue(transport_tpg.TestFakeCredentialsPath) // Path to a test fixture
			},
		},
		"configuring credentials as a path to a non-existant file is NOT valid": {
			ConfigValue: func(t *testing.T) types.String {
				return types.StringValue("./this/path/doesnt/exist.json") // Doesn't exist
			},
			ExpectedErrorCount: 1,
		},
		"configuring credentials as a credentials JSON string is valid": {
			ConfigValue: func(t *testing.T) types.String {
				contents, err := ioutil.ReadFile(transport_tpg.TestFakeCredentialsPath)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}
				stringContents := string(contents)
				return types.StringValue(stringContents)
			},
		},
		"configuring credentials as an empty string is not valid": {
			ConfigValue: func(t *testing.T) types.String {
				return types.StringValue("")
			},
			ExpectedErrorCount: 1,
		},
		"leaving credentials unconfigured is valid": {
			ConfigValue: func(t *testing.T) types.String {
				return types.StringNull()
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			req := validator.StringRequest{
				ConfigValue: tc.ConfigValue(t),
			}

			resp := validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			cv := fwprovider.CredentialsValidator()

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
