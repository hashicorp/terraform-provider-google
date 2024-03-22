// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestFunctionRun_location_from_id(t *testing.T) {
	t.Parallel()

	location := "us-central1"

	// Happy path inputs
	validId := fmt.Sprintf("projects/my-project/locations/%s/services/my-service", location)
	validSelfLink := fmt.Sprintf("https://run.googleapis.com/v2/%s", validId)
	validOpStyleResourceName := fmt.Sprintf("//run.googleapis.com/v2/%s", validId)

	// Unhappy path inputs
	repetitiveInput := fmt.Sprintf("https://run.googleapis.com/v2/projects/my-project/locations/%s/locations/not-this-one/services/my-service", location) // Multiple /locations/{{location}}/
	invalidInput := "zones/us-central1-c/instances/my-instance"

	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"it returns the expected output value when given a valid resource id input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validId)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(location)),
			},
		},
		"it returns the expected output value when given a valid resource self_link input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validSelfLink)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(location)),
			},
		},
		"it returns the expected output value when given a valid OP style resource name input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validOpStyleResourceName)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(location)),
			},
		},
		"it returns the first submatch (with no error) when given repetitive input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(repetitiveInput)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(location)),
			},
		},
		"it returns an error when given input with no submatches": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(invalidInput)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringNull()),
				Error:  function.NewArgumentFuncError(0, fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"locations/{location}/\".", invalidInput)),
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
			NewLocationFromIdFunction().Run(context.Background(), tc.request, &got)

			// Assert
			if diff := cmp.Diff(got.Result, tc.expected.Result); diff != "" {
				t.Errorf("unexpected diff between expected and received result: %s", diff)
			}
			if diff := cmp.Diff(got.Error, tc.expected.Error); diff != "" {
				t.Errorf("unexpected diff between expected and received errors: %s", diff)
			}
		})
	}
}
