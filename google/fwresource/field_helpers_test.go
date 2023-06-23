// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGetProjectFramework(t *testing.T) {
	cases := map[string]struct {
		ResourceProject types.String
		ProviderProject types.String
		ExpectedProject types.String
		ExpectedError   bool
	}{
		"project is pulled from the resource config value instead of the provider config value, even if both set": {
			ResourceProject: types.StringValue("foo"),
			ProviderProject: types.StringValue("bar"),
			ExpectedProject: types.StringValue("foo"),
		},
		"project is pulled from the provider config value when unset on the resource": {
			ResourceProject: types.StringNull(),
			ProviderProject: types.StringValue("bar"),
			ExpectedProject: types.StringValue("bar"),
		},
		"error when project is not set on the provider or the resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var diags diag.Diagnostics

			// Act
			project := GetProjectFramework(tc.ResourceProject, tc.ProviderProject, &diags)

			// Assert
			if diags.HasError() {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Got %d unexpected error(s) during test: %s", diags.ErrorsCount(), diags.Errors())
			}

			if project != tc.ExpectedProject {
				t.Fatalf("Incorrect project: got %s, want %s", project, tc.ExpectedProject)
			}
		})
	}
}
