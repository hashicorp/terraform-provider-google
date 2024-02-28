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

func TestFunctionRun_project_from_id(t *testing.T) {
	t.Parallel()

	projectId := "my-project"

	// Happy path inputs
	validId := fmt.Sprintf("projects/%s/zones/us-central1-c/instances/my-instance", projectId)
	validSelfLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/%s", validId)
	validOpStyleResourceName := fmt.Sprintf("//gkehub.googleapis.com/projects/%s/locations/us-central1/memberships/my-membership", projectId)

	// Unhappy path inputs
	repetitiveInput := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/projects/not-this-1/projects/not-this-2/instances/my-instance", projectId) // Multiple /projects/{{project}}/
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
				Result: function.NewResultData(types.StringValue(projectId)),
			},
		},
		"it returns the expected output value when given a valid resource self_link input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validSelfLink)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(projectId)),
			},
		},
		"it returns the expected output value when given a valid OP style resource name input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(validOpStyleResourceName)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(projectId)),
			},
		},
		"it returns a warning and the first submatch when given repetitive input": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(repetitiveInput)}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue(projectId)),
				Diagnostics: diag.Diagnostics{
					diag.NewArgumentWarningDiagnostic(
						0,
						ambiguousMatchesWarningSummary,
						fmt.Sprintf("The input string \"%s\" contains more than one match for the pattern \"projects/{project}/\". Terraform will use the first found match.", repetitiveInput),
					),
				},
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
						fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"projects/{project}/\".", invalidInput),
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
			NewProjectFromIdFunction().Run(context.Background(), tc.request, &got)

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
