// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCompileUserAgentString(t *testing.T) {
	cases := map[string]struct {
		Name              string // Name of the provider
		TerraformVersion  string
		ProviderVersion   string
		EnvValue          string
		ExpectedUserAgent string
	}{
		"the expected user agent is returned for given inputs": {
			Name:              "terraform-provider-foobar",
			TerraformVersion:  "1.2.3",
			ProviderVersion:   "9.9.9",
			ExpectedUserAgent: "Terraform/1.2.3 (+https://www.terraform.io) Terraform-Plugin-SDK/terraform-plugin-framework terraform-provider-foobar/9.9.9",
		},
		"the user agent can have values appended via an environment variable": {
			Name:              "terraform-provider-foobar",
			TerraformVersion:  "1.2.3",
			ProviderVersion:   "9.9.9",
			EnvValue:          "I'm appended at the end!",
			ExpectedUserAgent: "Terraform/1.2.3 (+https://www.terraform.io) Terraform-Plugin-SDK/terraform-plugin-framework terraform-provider-foobar/9.9.9 I'm appended at the end!",
		},
		"values appended via an environment variable have whitespace trimmed": {
			Name:              "terraform-provider-foobar",
			TerraformVersion:  "1.2.3",
			ProviderVersion:   "9.9.9",
			EnvValue:          "              my surrounding white space is removed              ",
			ExpectedUserAgent: "Terraform/1.2.3 (+https://www.terraform.io) Terraform-Plugin-SDK/terraform-plugin-framework terraform-provider-foobar/9.9.9 my surrounding white space is removed",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			ctx := context.Background()

			t.Setenv(uaEnvVar, tc.EnvValue) // Use same global const as the CompileUserAgentString function

			// Act

			ua := CompileUserAgentString(ctx, tc.Name, tc.TerraformVersion, tc.ProviderVersion)

			// Assert
			if ua != tc.ExpectedUserAgent {
				t.Fatalf("Incorrect user agent output: got %s, want %s", ua, tc.ExpectedUserAgent)
			}
		})
	}
}

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
			project := getProjectFramework(tc.ResourceProject, tc.ProviderProject, &diags)

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
