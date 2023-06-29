// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestProvider_validateCredentials(t *testing.T) {
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
				contents, err := ioutil.ReadFile(transport_tpg.TestFakeCredentialsPath)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}
				return string(contents)
			},
		},
		"configuring credentials as an empty string is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return ""
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
			ws, es := validateCredentials(configValue, "")

			// Assert
			if len(ws) != len(tc.ExpectedWarnings) {
				t.Errorf("Expected %d warnings, got %d: %v", len(tc.ExpectedWarnings), len(ws), ws)
			}
			if len(es) != len(tc.ExpectedErrors) {
				t.Errorf("Expected %d errors, got %d: %v", len(tc.ExpectedErrors), len(es), es)
			}

			if len(tc.ExpectedErrors) > 0 {
				if es[0].Error() != tc.ExpectedErrors[0].Error() {
					t.Errorf("Expected first error to be \"%s\", got \"%s\"", tc.ExpectedErrors[0], es[0])
				}
			}
		})
	}
}

// ProviderConfigEnvNames returns a list of all the environment variables that could be set by a user to configure the provider
func ProviderConfigEnvNames() []string {

	envs := []string{}

	// Use existing collections of ENV names
	envVarsSets := [][]string{
		envvar.CredsEnvVars,   // credentials field
		envvar.ProjectEnvVars, // project field
		envvar.RegionEnvVars,  //region field
		envvar.ZoneEnvVars,    // zone field
	}
	for _, set := range envVarsSets {
		envs = append(envs, set...)
	}

	// Add remaining ENVs
	envs = append(envs, "GOOGLE_OAUTH_ACCESS_TOKEN")          // access_token field
	envs = append(envs, "GOOGLE_BILLING_PROJECT")             // billing_project field
	envs = append(envs, "GOOGLE_IMPERSONATE_SERVICE_ACCOUNT") // impersonate_service_account field
	envs = append(envs, "USER_PROJECT_OVERRIDE")              // user_project_override field
	envs = append(envs, "CLOUDSDK_CORE_REQUEST_REASON")       // request_reason field

	return envs
}

// unsetProviderConfigEnvs unsets any ENVs in the test environment that
// configure the provider.
// The testing package will restore the original values after the test
func unsetTestProviderConfigEnvs(t *testing.T) {
	envs := ProviderConfigEnvNames()
	if len(envs) > 0 {
		for _, k := range envs {
			t.Setenv(k, "")
		}
	}
}

func setupTestEnvs(t *testing.T, envValues map[string]string) {
	// Set ENVs
	if len(envValues) > 0 {
		for k, v := range envValues {
			t.Setenv(k, v)
		}
	}
}

// Returns a fake credentials JSON string with the client_email set to a test-specific value
func generateFakeCredentialsJson(testId string) string {
	json := fmt.Sprintf(`{"private_key_id": "foo","private_key": "bar","client_email": "%s@example.com","client_id": "id@foo.com","type": "service_account"}`, testId)
	return json
}

func TestProvider_providerConfigure_credentials(t *testing.T) {

	const pathToMissingFile string = "./this/path/doesnt/exist.json" // Doesn't exist

	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectError         bool
		ExpectFieldUnset    bool
		ExpectedSchemaValue string
		ExpectedConfigValue string
	}{
		"credentials can be configured as a path to a credentials JSON file": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables:        map[string]string{},
			ExpectedSchemaValue: transport_tpg.TestFakeCredentialsPath,
			ExpectedConfigValue: transport_tpg.TestFakeCredentialsPath,
		},
		"configuring credentials as a path to a non-existant file results in an error": {
			ConfigValues: map[string]interface{}{
				"credentials": pathToMissingFile,
			},
			ExpectError:         true,
			ExpectedSchemaValue: pathToMissingFile,
			ExpectedConfigValue: pathToMissingFile,
		},
		"credentials set in the config are not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"credentials": generateFakeCredentialsJson("test"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_CREDENTIALS":             generateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
				"GOOGLE_CLOUD_KEYFILE_JSON":      generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: generateFakeCredentialsJson("test"),
			ExpectedConfigValue: generateFakeCredentialsJson("test"),
		},
		"when credentials is unset in the config, environment variables are used: GOOGLE_CREDENTIALS used first": {
			EnvVariables: map[string]string{
				"GOOGLE_CREDENTIALS":             generateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
				"GOOGLE_CLOUD_KEYFILE_JSON":      generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: generateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
		},
		"when credentials is unset in the config, environment variables are used: GOOGLE_CLOUD_KEYFILE_JSON used second": {
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				"GOOGLE_CLOUD_KEYFILE_JSON":      generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
		},
		"when credentials is unset in the config, environment variables are used: GCLOUD_KEYFILE_JSON used third": {
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				// GOOGLE_CLOUD_KEYFILE_JSON not set
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
		},
		"when credentials is unset in the config (and access_token unset), GOOGLE_APPLICATION_CREDENTIALS is used for auth but not to set values in the config": {
			EnvVariables: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
		},
		"when credentials is set to an empty string in the config (and access_token unset), GOOGLE_APPLICATION_CREDENTIALS is used": {
			ConfigValues: map[string]interface{}{
				"credentials": "",
			},
			EnvVariables: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
		},
		// NOTE: these tests can't run in Cloud Build due to ADC locating credentials despite `GOOGLE_APPLICATION_CREDENTIALS` being unset
		// See https://cloud.google.com/docs/authentication/application-default-credentials#search_order
		// "error returned if credentials is set as an empty string and GOOGLE_APPLICATION_CREDENTIALS is unset": {
		// 	ConfigValues: map[string]interface{}{
		// 		"credentials": "",
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": "", // setting to empty string to help test run in CI
		// 	},
		// 	ExpectError: true,
		// },
		// "error returned if neither credentials nor access_token set in the provider config, and GOOGLE_APPLICATION_CREDENTIALS is unset": {
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": "", // setting to empty string to help test run in CI
		// 	},
		// 	ExpectError: true,
		// },
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("credentials")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSchemaValue {
						t.Fatalf("expected credentials value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected credentials value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("credentials")
			val := v.(string)
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if v != tc.ExpectedSchemaValue {
				t.Fatalf("expected credentials value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.Credentials != tc.ExpectedConfigValue {
				t.Fatalf("expected credentials value in provider struct to be %s, got %s", tc.ExpectedConfigValue, config.Credentials)
			}
		})
	}
}

func TestProvider_providerConfigure_accessToken(t *testing.T) {

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
		"when access_token is unset in the config, an environment variable is used but doesn't update the schema data": {
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
		},
		"when access_token is set as an empty string in the config, an environment variable is used but doesn't update the schema data": {
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
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

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
						t.Fatalf("expected access_token value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected access_token value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("access_token")
			val := v.(string)
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if val != tc.ExpectedSchemaValue {
				t.Fatalf("expected access_token value set in provider data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.AccessToken != tc.ExpectedConfigValue {
				t.Fatalf("expected access_token value in provider struct to be %s, got %s", tc.ExpectedConfigValue, config.AccessToken)
			}
		})
	}
}

func TestProvider_providerConfigure_impersonateServiceAccount(t *testing.T) {

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
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedValue: "value-from-env@example.com",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectFieldUnset: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

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
						t.Fatalf("expected impersonate_service_account value set in provider data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected impersonate_service_account value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("impersonate_service_account")
			val := v.(string)
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if val != tc.ExpectedValue {
				t.Fatalf("expected impersonate_service_account value set in provider data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.ImpersonateServiceAccount != tc.ExpectedValue {
				t.Fatalf("expected impersonate_service_account value in provider struct to be %s, got %s", tc.ExpectedValue, config.ImpersonateServiceAccount)
			}
		})
	}
}

func TestProvider_providerConfigure_impersonateServiceAccountDelegates(t *testing.T) {

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
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

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
						t.Fatalf("expected impersonate_service_account_delegates value to not be set in provider data, got %#v", val)
					}
					if len(val) != len(tc.ExpectedValue) {
						t.Fatalf("expected impersonate_service_account_delegates value set in provider data to be %#v, got %#v", tc.ExpectedValue, val)
					}
					for i := 0; i < len(val); i++ {
						if val[i].(string) != tc.ExpectedValue[i] {
							t.Fatalf("expected impersonate_service_account_delegates value set in provider data to be %#v, got %#v", tc.ExpectedValue, val)
						}
					}

				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("impersonate_service_account_delegates")
			val := v.([]interface{})
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if len(val) != len(tc.ExpectedValue) {
				t.Fatalf("expected impersonate_service_account_delegates value set in provider data to be %#v, got %#v", tc.ExpectedValue, val)
			}
			for i := 0; i < len(val); i++ {
				if val[i].(string) != tc.ExpectedValue[i] {
					t.Fatalf("expected impersonate_service_account_delegates value set in provider data to be %#v, got %#v", tc.ExpectedValue, val)
				}
				if config.ImpersonateServiceAccountDelegates[i] != tc.ExpectedValue[i] {
					t.Fatalf("expected impersonate_service_account_delegates value in provider struct to be %#v, got %#v", tc.ExpectedValue, config.ImpersonateServiceAccountDelegates)
				}
			}
		})
	}
}

func TestProvider_providerConfigure_project(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"project value set in the provider schema is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"project":     "my-project-from-config",
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "my-project-from-config",
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
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

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
						t.Fatalf("expected project value set in provider data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected project value to not be set in provider data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			v := d.Get("project")
			val := v.(string)
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			if val != tc.ExpectedValue {
				t.Fatalf("expected project value set in provider data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.Project != tc.ExpectedValue {
				t.Fatalf("expected project value in provider struct to be %s, got %s", tc.ExpectedValue, config.Project)
			}
		})
	}
}
