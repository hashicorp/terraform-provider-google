// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestFunctionRun_name_from_id(t *testing.T) {
	t.Parallel()

	name := "foobar"

	// Happy path inputs
	validId := fmt.Sprintf("projects/my-project/zones/us-central1-c/instances/%s", name)
	validSelfLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", validId)
	validOpStyleResourceName := fmt.Sprintf("//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/%s", name)

	// Unhappy path inputs
	invalidInput := "this isn't a URI or id"

	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"it returns the expected output value when given a valid resource id input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validId)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(name)),
			},
		},
		"it returns the expected output value when given a valid resource self_link input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validSelfLink)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(name)),
			},
		},
		"it returns the expected output value when given a valid OP style resource name input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validOpStyleResourceName)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(name)),
			},
		},
		"it returns an error when given input with no submatches": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(invalidInput)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringNull()),
				Diagnostics: diag.Diagnostics{
					diag.NewArgumentErrorDiagnostic(
						0,
						noMatchesErrorSummary,
						fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"resourceType/{name}$\".", invalidInput),
					),
				},
			},
		},
	}

	for name, testCase := range testCases {
		tn, tc := name, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			// Arrange
			got := function.RunResponse{
				Result: function.NewResultData(basetypes.StringValue{}),
			}

			// Act
			NewNameFromIdFunction().Run(context.Background(), tc.request, &got)

			// Assert
			if diff := cmp.Diff(got.Result, tc.expected.Result); diff != "" {
				t.Errorf("unexpected diff between expected and received result: %s", diff)
			}
			if diff := cmp.Diff(got.Diagnostics, tc.expected.Diagnostics); diff != "" {
				t.Errorf("unexpected diff between expected and received diagnostics: %s", diff)
			}
		})
	}
}
