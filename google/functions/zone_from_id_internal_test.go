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

func TestFunctionRun_zone_from_id(t *testing.T) {
	t.Parallel()

	zone := "us-central1-a"

	// Happy path inputs
	validId := fmt.Sprintf("projects/my-project/zones/%s/networkEndpointGroups/my-neg", zone)
	validSelfLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/my-project/zones/%s/networkEndpointGroups/my-neg", zone)
	validOpStyleResourceName := fmt.Sprintf("//compute.googleapis.com/projects/my-project/zones/%s/instances/my-instance", zone)

	// Unhappy path inputs
	repetitiveInput := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/my-project/zones/%s/zones/not-this-one/networkEndpointGroups/my-neg", zone)
	invalidInput := "projects/my-project/regions/us-central1/subnetworks/my-subnetwork"

	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"it returns the expected output value when given a valid resource id input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validId)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(zone)),
			},
		},
		"it returns the expected output value when given a valid resource self_link input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validSelfLink)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(zone)),
			},
		},
		"it returns the expected output value when given a valid OP style resource name input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validOpStyleResourceName)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(zone)),
			},
		},
		"it returns the first submatch (with no error) when given repetitive input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(repetitiveInput)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(zone)),
			},
		},
		"it returns an error when given input with no submatches": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(invalidInput)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringNull()),
				Error: function.NewArgumentFuncError(
					0,
					fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"zones/{zone}/\".", invalidInput),
				),
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
			NewZoneFromIdFunction().Run(context.Background(), tc.request, &got)

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
