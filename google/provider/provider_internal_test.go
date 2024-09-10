// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestProvider_ValidateCredentials(t *testing.T) {
	cases := map[string]struct {
		ConfigValue      func(t *testing.T) interface{}
		ValueNotProvided bool
		ExpectedWarnings []string
		ExpectedErrors   []error
	}{
		"configuring credentials as a path to a credentials JSON file is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return transport_tpg.TestFakeCredentialsPath // Path to a test fixture
			},
		},
		"configuring credentials as a path to a non-existant file is NOT valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return "./this/path/doesnt/exist.json" // Doesn't exist
			},
			ExpectedErrors: []error{
				// As the file doesn't exist, so the function attempts to parse it as a JSON
				errors.New("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
		"configuring credentials as a credentials JSON string is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				contents, err := os.ReadFile(transport_tpg.TestFakeCredentialsPath)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}
				return string(contents)
			},
		},
		"configuring credentials as an empty string is not valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return ""
			},
			ExpectedErrors: []error{
				errors.New("expected a non-empty string"),
			},
		},
		"leaving credentials unconfigured is valid": {
			ValueNotProvided: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			var configValue interface{}
			if !tc.ValueNotProvided {
				configValue = tc.ConfigValue(t)
			}

			// Act
			// Note: second argument is currently unused by the function but is necessary to fulfill the SchemaValidateFunc type's function signature
			ws, es := provider.ValidateCredentials(configValue, "")

			// Assert
			if len(ws) != len(tc.ExpectedWarnings) {
				t.Fatalf("Expected %d warnings, got %d: %v", len(tc.ExpectedWarnings), len(ws), ws)
			}
			if len(es) != len(tc.ExpectedErrors) {
				t.Fatalf("Expected %d errors, got %d: %v", len(tc.ExpectedErrors), len(es), es)
			}

			if len(tc.ExpectedErrors) > 0 && len(es) > 0 {
				if es[0].Error() != tc.ExpectedErrors[0].Error() {
					t.Fatalf("Expected first error to be \"%s\", got \"%s\"", tc.ExpectedErrors[0], es[0])
				}
			}
		})
	}
}

func TestProvider_ValidateEmptyStrings(t *testing.T) {
	cases := map[string]struct {
		ConfigValue      interface{}
		ValueNotProvided bool
		ExpectedWarnings []string
		ExpectedErrors   []error
	}{
		"non-empty strings are valid": {
			ConfigValue: "foobar",
		},
		"unconfigured values are valid": {
			ValueNotProvided: true,
		},
		"empty strings are not valid": {
			ConfigValue: "",
			ExpectedErrors: []error{
				errors.New("expected a non-empty string"),
			},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			var configValue interface{}
			if !tc.ValueNotProvided {
				configValue = tc.ConfigValue
			}

			// Act
			// Note: second argument is currently unused by the function but is necessary to fulfill the SchemaValidateFunc type's function signature
			ws, es := provider.ValidateEmptyStrings(configValue, "")

			// Assert
			if len(ws) != len(tc.ExpectedWarnings) {
				t.Fatalf("Expected %d warnings, got %d: %v", len(tc.ExpectedWarnings), len(ws), ws)
			}
			if len(es) != len(tc.ExpectedErrors) {
				t.Fatalf("Expected %d errors, got %d: %v", len(tc.ExpectedErrors), len(es), es)
			}

			if len(tc.ExpectedErrors) > 0 && len(es) > 0 {
				if es[0].Error() != tc.ExpectedErrors[0].Error() {
					t.Fatalf("Expected first error to be \"%s\", got \"%s\"", tc.ExpectedErrors[0], es[0])
				}
			}
		})
	}
}

func TestProvider_ProviderConfigure_accessToken(t *testing.T) {

	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectedSchemaValue string
		ExpectedConfigValue string
		ExpectError         bool
		ExpectFieldUnset    bool
	}{
		"access_token configured in the provider can be invalid without resulting in errors": {
			ConfigValues: map[string]interface{}{
				"access_token": "This is not a valid token string",
			},
			EnvVariables:        map[string]string{},
			ExpectedSchemaValue: "This is not a valid token string",
			ExpectedConfigValue: "This is not a valid token string",
		},
		"access_token set in the provider config is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"access_token": "value-from-config",
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-env",
			},
			ExpectedSchemaValue: "value-from-config",
			ExpectedConfigValue: "value-from-config",
		},
		"when access_token is unset in the config, the GOOGLE_OAUTH_ACCESS_TOKEN environment variable is used": {
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
		},
		"when no access_token values are provided via config or environment variables there's no error": {
			ConfigValues: map[string]interface{}{
				// access_token unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:         false,
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "",
		},
		// Handle empty strings in config
		"when access_token is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"access_token": "",
				"credentials":  transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:         false,
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "",
		},
		"when access_token is set as an empty string in the config, an environment variable is used": {
			ConfigValues: map[string]interface{}{
				"access_token": "",
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("access_token")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSchemaValue {
						t.Fatalf("expected access_token value set in provider config data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected access_token value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("access_token")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected access_token value to be unset in provider config data")
			}
			if val != tc.ExpectedSchemaValue {
				t.Fatalf("expected access_token value set in provider config data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.AccessToken != tc.ExpectedConfigValue {
				t.Fatalf("expected access_token value set in Config struct to be to be %s, got %s", tc.ExpectedConfigValue, config.AccessToken)
			}
		})
	}
}

func TestProvider_ProviderConfigure_impersonateServiceAccount(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"impersonate_service_account value set in the provider schema is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account": "value-from-config@example.com",
				"credentials":                 transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedValue: "value-from-config@example.com",
		},
		"impersonate_service_account value can be set by environment variable": {
			ConfigValues: map[string]interface{}{
				// impersonate_service_account unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedValue: "value-from-env@example.com",
		},
		"when no impersonate_service_account values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// impersonate_service_account unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		// Handling empty strings in config
		"when impersonate_service_account is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account": "",
				"credentials":                 transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when impersonate_service_account is set as an empty string in the config, an environment variable is used": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account": "",
				"credentials":                 transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedValue: "value-from-env@example.com",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("impersonate_service_account")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected impersonate_service_account value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected impersonate_service_account value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("impersonate_service_account")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected impersonate_service_account value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected impersonate_service_account value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.ImpersonateServiceAccount != tc.ExpectedValue {
				t.Fatalf("expected impersonate_service_account value in Config struct to be %s, got %s", tc.ExpectedValue, config.ImpersonateServiceAccount)
			}
		})
	}
}

func TestProvider_ProviderConfigure_impersonateServiceAccountDelegates(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    []string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"impersonate_service_account_delegates value can be set in the provider schema": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account_delegates": []string{
					"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
					"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
				},
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue: []string{
				"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
				"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
			},
		},
		// No environment variables can be used for impersonate_service_account_delegates
		"when no impersonate_service_account_delegates value is provided via config, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// impersonate_service_account_delegates unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    nil,
		},
		// Handling empty values in config
		"when impersonate_service_account_delegates is set as an empty array the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account_delegates": []string{},
				"credentials":                           transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("impersonate_service_account_delegates")
				if ok {
					val := v.([]interface{})
					if tc.ExpectFieldUnset {
						t.Fatalf("expected impersonate_service_account_delegates value to not be set in provider config data, got %#v", val)
					}
					if len(val) != len(tc.ExpectedValue) {
						t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
					}
					for i := 0; i < len(val); i++ {
						if val[i].(string) != tc.ExpectedValue[i] {
							t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
						}
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors
			v, ok := d.GetOk("impersonate_service_account_delegates")
			val := v.([]interface{})
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected impersonate_service_account_delegates value to be unset in provider config data")
			}
			if len(val) != len(tc.ExpectedValue) {
				t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
			}
			for i := 0; i < len(val); i++ {
				if val[i].(string) != tc.ExpectedValue[i] {
					t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
				}
				if config.ImpersonateServiceAccountDelegates[i] != tc.ExpectedValue[i] {
					t.Fatalf("expected impersonate_service_account_delegates value set in Config struct to be to be %#v, got %#v", tc.ExpectedValue, config.ImpersonateServiceAccountDelegates)
				}
			}
		})
	}
}

func TestProvider_ProviderConfigure_project(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"project value set in the provider schema is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"project":     "project-from-config",
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-config",
		},
		"project value can be set by environment variable: GOOGLE_PROJECT is used first": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-GOOGLE_PROJECT",
		},
		"project value can be set by environment variable: GOOGLE_CLOUD_PROJECT is used second": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-GOOGLE_CLOUD_PROJECT",
		},
		"project value can be set by environment variable: GCLOUD_PROJECT is used third": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				// GOOGLE_CLOUD_PROJECT unset
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-GCLOUD_PROJECT",
		},
		"project value can be set by environment variable: CLOUDSDK_CORE_PROJECT is used fourth": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				// GOOGLE_CLOUD_PROJECT unset
				// GCLOUD_PROJECT unset
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-CLOUDSDK_CORE_PROJECT",
		},
		"when no project values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue: "",
		},
		// Handling empty strings in config
		"when project is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"project":     "",
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when project is set as an empty string an environment variable will be used": {
			ConfigValues: map[string]interface{}{
				"project":     "",
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT": "project-from-GOOGLE_PROJECT",
			},
			ExpectedValue: "project-from-GOOGLE_PROJECT",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("project")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected project value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("project")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected project value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.Project != tc.ExpectedValue {
				t.Fatalf("expected project value set in Config struct to be to be %s, got %s", tc.ExpectedValue, config.Project)
			}
		})
	}
}

func TestProvider_ProviderConfigure_billingProject(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"billing_project value set in the provider config is not overridden by ENVs": {
			ConfigValues: map[string]interface{}{
				"billing_project": "billing-project-from-config",
				"credentials":     transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
			},
			ExpectedValue: "billing-project-from-config",
		},
		"billing_project can be set by environment variable, when no value supplied via the config": {
			ConfigValues: map[string]interface{}{
				// billing_project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
			},
			ExpectedValue: "billing-project-from-env",
		},
		"when no billing_project values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// billing_project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		// Handling empty strings in config
		"when billing_project is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"billing_project": "",
				"credentials":     transport_tpg.TestFakeCredentialsPath,
			},
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when billing_project is set as an empty string an environment variable will be used": {
			ConfigValues: map[string]interface{}{
				"billing_project": "",
				"credentials":     transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
			},
			ExpectedValue: "billing-project-from-env",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("billing_project")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected billing_project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected billing_project value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("billing_project")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected billing_project value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected billing_project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.BillingProject != tc.ExpectedValue {
				t.Fatalf("expected billing_project value set in Config struct to be to be %s, got %s", tc.ExpectedValue, config.BillingProject)
			}
		})
	}
}

func TestProvider_ProviderConfigure_userProjectOverride(t *testing.T) {
	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    bool
		ExpectFieldUnset bool
		ExpectError      bool
	}{
		"user_project_override value set in the provider schema is not overridden by ENVs": {
			ConfigValues: map[string]interface{}{
				"user_project_override": false,
				"credentials":           transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedValue: false,
		},
		"user_project_override can be set by environment variable: value = true": {
			ConfigValues: map[string]interface{}{
				// user_project_override not set
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedValue: true,
		},
		"user_project_override can be set by environment variable: value = false": {
			ConfigValues: map[string]interface{}{
				// user_project_override not set
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "false",
			},
			ExpectedValue: false,
		},
		"user_project_override can be set by environment variable: value = 1": {
			ConfigValues: map[string]interface{}{
				// user_project_override not set
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "1",
			},
			ExpectedValue: true,
		},
		"user_project_override can be set by environment variable: value = 0": {
			ConfigValues: map[string]interface{}{
				// user_project_override not set
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "0",
			},
			ExpectedValue: false,
		},
		"setting user_project_override using a non-boolean environment variables results in an error": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "I'm not a boolean",
			},
			ExpectError: true,
		},
		"when no user_project_override values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// user_project_override unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    false,
		},
		// There isn't an equivalent test case for 'user sets value as empty string' because user_project_override is a boolean; true/false both valid.
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("user_project_override")
				if ok {
					val := v.(bool)
					if tc.ExpectFieldUnset {
						t.Fatalf("expected user_project_override value to not be set in provider config data, got %v", val)
					}
					if val != tc.ExpectedValue {
						t.Fatalf("expected user_project_override value set in provider config data to be %v, got %v", tc.ExpectedValue, val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("user_project_override")
			val := v.(bool)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected user_project_override value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected user_project_override value set in provider config data to be %v, got %v", tc.ExpectedValue, val)
			}
			if config.UserProjectOverride != tc.ExpectedValue {
				t.Fatalf("expected user_project_override value set in Config struct to be to be %v, got %v", tc.ExpectedValue, config.UserProjectOverride)
			}
		})
	}
}

func TestProvider_ProviderConfigure_scopes(t *testing.T) {
	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectedSchemaValue []string
		ExpectedConfigValue []string
		ExpectFieldUnset    bool
		ExpectError         bool
	}{
		"scopes are set in the provider config as a list": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				"scopes":      []string{"fizz", "buzz", "fizzbuzz"},
			},
			ExpectedSchemaValue: []string{"fizz", "buzz", "fizzbuzz"},
			ExpectedConfigValue: []string{"fizz", "buzz", "fizzbuzz"},
		},
		"scopes can be left unset in the provider config without any issues, and a default value is used": {
			ConfigValues: map[string]interface{}{
				// scopes unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedSchemaValue: nil,
			ExpectedConfigValue: transport_tpg.DefaultClientScopes,
		},
		// Handling empty values in config
		"scopes set as an empty list the field is treated as if it's unset and a default value is used without errors": {
			ConfigValues: map[string]interface{}{
				"scopes":      []string{},
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:         false,
			ExpectFieldUnset:    true, //unset in provider config data, not the subsequent Config struct
			ExpectedSchemaValue: nil,
			ExpectedConfigValue: transport_tpg.DefaultClientScopes,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("scopes")
				if ok {
					val := v.([]interface{})
					if tc.ExpectFieldUnset {
						t.Fatalf("expected scopes value to not be set in provider config data, got %#v", val)
					}
					if len(val) != len(tc.ExpectedSchemaValue) {
						t.Fatalf("expected scopes value set in provider config data to be %#v, got %#v", tc.ExpectedSchemaValue, val)
					}
					for i := 0; i < len(val); i++ {
						if val[i].(string) != tc.ExpectedSchemaValue[i] {
							t.Fatalf("expected scopes value set in provider config data to be %#v, got %#v", tc.ExpectedSchemaValue, val)
						}
					}
				}
				// Return early in tests where errors expected
				return
			}
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors
			v, ok := d.GetOk("scopes")
			if ok {
				val := v.([]interface{})

				if len(val) != len(tc.ExpectedSchemaValue) {
					t.Fatalf("expected %v scopes set in provider config data, got %v", len(tc.ExpectedSchemaValue), len(val))
				}
				for i, el := range val {
					scope := el.(string)
					if scope != tc.ExpectedSchemaValue[i] {
						t.Fatalf("expected scopes value set in provider config data to be %v, got %v", tc.ExpectedSchemaValue, val)
					}
				}
			}
			if !ok && (len(tc.ExpectedSchemaValue) > 0) {
				t.Fatalf("expected %v scopes to be set in the provider data, but is unset", tc.ExpectedSchemaValue)
			}

			if len(config.Scopes) != len(tc.ExpectedConfigValue) {
				t.Fatalf("expected %v scopes set in the config struct, got %v", len(tc.ExpectedConfigValue), len(config.Scopes))
			}
			for i, el := range config.Scopes {
				if el != tc.ExpectedConfigValue[i] {
					t.Fatalf("expected scopes value set in provider config data to be %v, got %v", tc.ExpectedConfigValue, config.Scopes)
				}
			}
		})
	}
}

func TestProvider_ProviderConfigure_requestTimeout(t *testing.T) {
	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		ExpectedValue       string
		ExpectedSchemaValue string
		ExpectError         bool
		ExpectFieldUnset    bool
	}{
		"if a valid request_timeout is configured in the provider, no error will occur": {
			ConfigValues: map[string]interface{}{
				"request_timeout": "10s",
				"credentials":     transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue:       "10s",
			ExpectedSchemaValue: "10s",
		},
		"if an invalid request_timeout is configured in the provider, an error will occur": {
			ConfigValues: map[string]interface{}{
				"request_timeout": "timeout",
				"credentials":     transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue:       "timeout",
			ExpectedSchemaValue: "timeout",
			ExpectError:         true,
			ExpectFieldUnset:    false,
		},
		// RequestTimeout is "0s" if unset by the user, and logic elsewhere will supply a different value.
		// This can be seen in this part of the config code where the default value is set to "120s"
		// https://github.com/hashicorp/terraform-provider-google/blob/09cb850ee64bcd78e4457df70905530c1ed75f19/google/transport/config.go#L1228-L1233
		"when request_timeout is unset in the config, the default value is 0s. (initially; this value is subsequently overwritten)": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue:    "0s",
			ExpectFieldUnset: true,
		},
		"when request_timeout is set as an empty string, the default value is 0s. (initially; this value is subsequently overwritten)": {
			ConfigValues: map[string]interface{}{
				"request_timeout": "",
				"credentials":     transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue: "0s",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("request_timeout")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSchemaValue {
						t.Fatalf("expected request_timeout value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected request_timeout value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("request_timeout") // checks for an empty or "0" string in order to set the default value
			val := v.(string)
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if val != tc.ExpectedSchemaValue {
				t.Fatalf("expected request_timeout value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.RequestTimeout.String() != tc.ExpectedValue {
				t.Fatalf("expected request_timeout value in provider struct to be %s, got %v", tc.ExpectedValue, config.RequestTimeout.String())
			}
		})
	}
}

func TestProvider_ProviderConfigure_requestReason(t *testing.T) {

	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectError         bool
		ExpectFieldUnset    bool
		ExpectedSchemaValue string
		ExpectedConfigValue string
	}{
		"when request_reason is unset in the config, environment variable CLOUDSDK_CORE_REQUEST_REASON is used": {
			ConfigValues: map[string]interface{}{
				// request_reason unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "test",
			},
			ExpectedSchemaValue: "test",
			ExpectedConfigValue: "test",
		},
		"request_reason set in the config is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"request_reason": "request test",
				"credentials":    transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "test",
			},
			ExpectedSchemaValue: "request test",
			ExpectedConfigValue: "request test",
		},
		"when no request_reason is provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// request_reason unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "",
		},
		// Handling empty values in config
		"when request_reason is set as an empty string in the config it is overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"request_reason": "",
				"credentials":    transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "test",
			},
			ExpectedSchemaValue: "test",
			ExpectedConfigValue: "test",
		},
		"when request_reason is set as an empty string in the config, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				"request_reason": "",
				"credentials":    transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("request_reason")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSchemaValue {
						t.Fatalf("expected request_reason value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected request_reason value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("request_reason")
			val := v.(string)
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if v != tc.ExpectedSchemaValue {
				t.Fatalf("expected request_reason value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.RequestReason != tc.ExpectedConfigValue {
				t.Fatalf("expected request_reason value in provider struct to be %s, got %s", tc.ExpectedConfigValue, config.Credentials)
			}
		})
	}
}

func TestProvider_ProviderConfigure_batching(t *testing.T) {
	//var batch []interface{}
	cases := map[string]struct {
		ConfigValues                map[string]interface{}
		EnvVariables                map[string]string
		ExpectError                 bool
		ExpectFieldUnset            bool
		ExpectedEnableBatchingValue bool
		ExpectedSendAfterValue      string
	}{
		"batching can be configured with values for enable_batching and send_after": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				"batching": []interface{}{
					map[string]interface{}{
						"enable_batching": true,
						"send_after":      "45s",
					},
				},
			},
			ExpectedEnableBatchingValue: true,
			ExpectedSendAfterValue:      "45s",
		},
		"if batching is an empty block, it will set the default values for enable_batching and send_after": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				// batching not set
			},
			// Although at the schema level it's shown that by default it's set to false, the actual default value
			// is true and can be seen in the `ExpanderProviderBatchingConfig` struct
			// https://github.com/GoogleCloudPlatform/magic-modules/blob/8cd4a506f0ac4db7b07a8cce914449d34df6f20b/mmv1/third_party/terraform/transport/config.go.erb#L504-L508
			ExpectedEnableBatchingValue: false,
			ExpectedSendAfterValue:      "", // uses "" value to be able to set the default value of 30s
			ExpectFieldUnset:            true,
		},
		"when batching is configured with only enable_batching, send_after will be set to a default value": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				"batching": []interface{}{
					map[string]interface{}{
						"enable_batching": true,
					},
				},
			},
			ExpectedEnableBatchingValue: true,
			ExpectedSendAfterValue:      "",
		},
		"when batching is configured with only send_after, enable_batching will be set to a default value": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				"batching": []interface{}{
					map[string]interface{}{
						"send_after": "45s",
					},
				},
			},
			ExpectedEnableBatchingValue: false,
			ExpectedSendAfterValue:      "45s",
		},
		// Error states
		"if batching is configured with send_after as an invalid value, there's an error": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				"batching": []interface{}{
					map[string]interface{}{
						"send_after": "invalid value",
					},
				},
			},
			ExpectedSendAfterValue: "invalid value",
			ExpectError:            true,
		},
		"if batching is configured with send_after as number value without seconds (s), there's an error": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
				"batching": []interface{}{
					map[string]interface{}{
						"send_after": "10",
					},
				},
			},
			ExpectedSendAfterValue: "10",
			ExpectError:            true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			acctest.UnsetTestProviderConfigEnvs(t)
			p := provider.Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			_, diags := provider.ProviderConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("batching.0.enable_batching")
				val := v.(bool)
				if ok {
					if val != tc.ExpectedEnableBatchingValue {
						t.Fatalf("expected request_timeout value set in provider data to be %v, got %v", tc.ExpectedEnableBatchingValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected request_timeout value to not be set in provider data, got %v", val)
					}
				}

				v, ok = d.GetOk("batching.0.send_after")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSendAfterValue {
						t.Fatalf("expected send_after value set in provider data to be %v, got %v", tc.ExpectedSendAfterValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected send_after value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("batching.0.enable_batching")
			enableBatching := v.(bool)
			if enableBatching != tc.ExpectedEnableBatchingValue {
				t.Fatalf("expected enable_batching value set in provider data to be %v, got %v", tc.ExpectedEnableBatchingValue, enableBatching)
			}

			v = d.Get("batching.0.send_after") // checks for an empty string in order to set the default value
			sendAfter := v.(string)
			if sendAfter != tc.ExpectedSendAfterValue {
				t.Fatalf("expected send_after value set in provider data to be %s, got %s", tc.ExpectedSendAfterValue, sendAfter)
			}
		})
	}
}
