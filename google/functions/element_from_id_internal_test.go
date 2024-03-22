// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions_test

import (
	"context"
	"regexp"
	"testing"

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
	}{
		"it sets an error if no match is found": {
			Input:       "one/element-1/three/element-3",
			ExpectError: true,
		},
		"it doesn't set an error if more than one match is found": {
			Input:           "two/element-2/two/element-2/two/element-2",
			ExpectedElement: "element-2",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()

			// Act
			err := tpg_functions.ValidateElementFromIdArguments(ctx, tc.Input, regex, pattern, "function-name-here") // last arg value is inconsequential for this test

			// Assert
			if err != nil && !tc.ExpectError {
				t.Fatalf("Unexpected error(s) were set in response diags: %s", err.Text)
			}
			if err == nil && tc.ExpectError {
				t.Fatal("Expected error(s) to be set in response diags, but there were none.")
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
