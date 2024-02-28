// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	tpg_functions "github.com/hashicorp/terraform-provider-google/google/functions"
)

func TestFunctionInternals_ValidateElementFromIdArguments(t *testing.T) {

	// Values here are matched to test case values below
	regex := regexp.MustCompile("two/(?P<Element>[^/]+)/")
	pattern := "two/{two}/"

	cases := map[string]struct {
		Input           string
		ExpectedElement string
		ExpectError     bool
		ExpectWarning   bool
	}{
		"it sets an error in diags if no match is found": {
			Input:       "one/element-1/three/element-3",
			ExpectError: true,
		},
		"it sets a warning in diags if more than one match is found": {
			Input:           "two/element-2/two/element-2/two/element-2",
			ExpectedElement: "element-2",
			ExpectWarning:   true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			resp := function.RunResponse{
				Result: function.NewResultData(basetypes.StringValue{}),
			}

			// Act
			tpg_functions.ValidateElementFromIdArguments(tc.Input, regex, pattern, &resp)

			// Assert
			if resp.Diagnostics.HasError() && !tc.ExpectError {
				t.Fatalf("Unexpected error(s) were set in response diags: %s", resp.Diagnostics.Errors())
			}
			if !resp.Diagnostics.HasError() && tc.ExpectError {
				t.Fatal("Expected error(s) to be set in response diags, but there were none.")
			}
			if (resp.Diagnostics.WarningsCount() > 0) && !tc.ExpectWarning {
				t.Fatalf("Unexpected warning(s) were set in response diags: %s", resp.Diagnostics.Warnings())
			}
			if (resp.Diagnostics.WarningsCount() == 0) && tc.ExpectWarning {
				t.Fatal("Expected warning(s) to be set in response diags, but there were none.")
			}
		})
	}
}

func TestFunctionInternals_GetElementFromId(t *testing.T) {

	// Values here are matched to test case values below
	regex := regexp.MustCompile("two/(?P<Element>[^/]+)/")
	template := "$Element"

	cases := map[string]struct {
		Input           string
		ExpectedElement string
	}{
		"it can pull out a value from a string using a regex with a submatch": {
			Input:           "one/element-1/two/element-2/three/element-3",
			ExpectedElement: "element-2",
		},
		"it will pull out the first value from a string with more than one submatch": {
			Input:           "one/element-1/two/element-2/two/not-this-one/three/element-3",
			ExpectedElement: "element-2",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Act
			result := tpg_functions.GetElementFromId(tc.Input, regex, template)

			// Assert
			if result != tc.ExpectedElement {
				t.Fatalf("Expected function logic to retrieve %s from input %s, got %s", tc.ExpectedElement, tc.Input, result)
			}
		})
	}
}
