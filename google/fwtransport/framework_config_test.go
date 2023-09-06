// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwtransport_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFrameworkProvider_LoadAndValidateFramework_project(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues              fwmodels.ProviderModel
		EnvVariables              map[string]string
		ExpectedDataModelValue    basetypes.StringValue // Sometimes the value is mutated, and no longer matches the original value we supply
		ExpectedConfigStructValue basetypes.StringValue // Sometimes the value in config struct differs from what is in the data model
		ExpectError               bool
	}{
		"project value set in the provider schema is not overridden by environment variables": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringValue("my-project-from-config"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("my-project-from-config"),
			ExpectedConfigStructValue: types.StringValue("my-project-from-config"),
		},
		"project value can be set by environment variable: GOOGLE_PROJECT is used first": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("project-from-GOOGLE_PROJECT"),
			ExpectedConfigStructValue: types.StringValue("project-from-GOOGLE_PROJECT"),
		},
		"project value can be set by environment variable: GOOGLE_CLOUD_PROJECT is used second": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("project-from-GOOGLE_CLOUD_PROJECT"),
			ExpectedConfigStructValue: types.StringValue("project-from-GOOGLE_CLOUD_PROJECT"),
		},
		"project value can be set by environment variable: GCLOUD_PROJECT is used third": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				// GOOGLE_CLOUD_PROJECT unset
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("project-from-GCLOUD_PROJECT"),
			ExpectedConfigStructValue: types.StringValue("project-from-GCLOUD_PROJECT"),
		},
		"project value can be set by environment variable: CLOUDSDK_CORE_PROJECT is used fourth": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				// GOOGLE_CLOUD_PROJECT unset
				// GCLOUD_PROJECT unset
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("project-from-CLOUDSDK_CORE_PROJECT"),
			ExpectedConfigStructValue: types.StringValue("project-from-CLOUDSDK_CORE_PROJECT"),
		},
		"when no project values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringNull(), // unset
			},
			ExpectedDataModelValue:    types.StringNull(),
			ExpectedConfigStructValue: types.StringNull(),
		},
		// Handling empty strings in config
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14255
		// "when project is set as an empty string the field is treated as if it's unset, without error": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Project: types.StringValue(""),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
		// 	ExpectedConfigStructValue: types.StringNull(),
		// },
		// "when project is set as an empty string an environment variable will be used": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Project: types.StringValue(""),
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_PROJECT": "project-from-GOOGLE_PROJECT",
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
		// 	ExpectedConfigStructValue: types.StringValue("project-from-GOOGLE_PROJECT"),
		// },
		// Handling unknown values
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when project is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Project: types.StringUnknown(),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
		// 	ExpectedConfigStructValue: types.StringNull(),
		// },
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)

			ctx := context.Background()
			tfVersion := "foobar"
			providerversion := "999"
			diags := diag.Diagnostics{}

			data := tc.ConfigValues
			data.Credentials = types.StringValue(transport_tpg.TestFakeCredentialsPath)
			impersonateServiceAccountDelegates, _ := types.ListValue(types.StringType, []attr.Value{}) // empty list
			data.ImpersonateServiceAccountDelegates = impersonateServiceAccountDelegates

			p := fwtransport.FrameworkProviderConfig{}

			// Act
			p.LoadAndValidateFramework(ctx, &data, tfVersion, &diags, providerversion)

			// Assert
			if diags.HasError() && tc.ExpectError {
				return
			}
			if diags.HasError() && !tc.ExpectError {
				for i, err := range diags.Errors() {
					num := i + 1
					t.Logf("unexpected error #%d : %s", num, err.Summary())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.Project.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want project in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.Project.String())
			}
			// Checking the value passed to the config structs
			if !p.Project.Equal(tc.ExpectedConfigStructValue) {
				t.Fatalf("want project in the `FrameworkProviderConfig` struct to be `%s`, but got the value `%s`", tc.ExpectedConfigStructValue, p.Project.String())
			}
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_credentials(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	const pathToMissingFile string = "./this/path/doesnt/exist.json" // Doesn't exist

	cases := map[string]struct {
		ConfigValues           fwmodels.ProviderModel
		EnvVariables           map[string]string
		ExpectedDataModelValue basetypes.StringValue
		// ExpectedConfigStructValue not used here, as credentials info isn't stored in the config struct
		ExpectError bool
	}{
		"credentials can be configured as a path to a credentials JSON file": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringValue(transport_tpg.TestFakeCredentialsPath),
			},
			ExpectedDataModelValue: types.StringValue(transport_tpg.TestFakeCredentialsPath),
		},
		"configuring credentials as a path to a non-existent file results in an error": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringValue(pathToMissingFile),
			},
			ExpectError: true,
		},
		"credentials set in the config are not overridden by environment variables": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringValue(acctest.GenerateFakeCredentialsJson("test")),
			},
			EnvVariables: map[string]string{
				"GOOGLE_CREDENTIALS":             acctest.GenerateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
				"GOOGLE_CLOUD_KEYFILE_JSON":      acctest.GenerateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            acctest.GenerateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": acctest.GenerateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedDataModelValue: types.StringValue(acctest.GenerateFakeCredentialsJson("test")),
		},
		"when credentials is unset in the config, environment variables are used: GOOGLE_CREDENTIALS used first": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				"GOOGLE_CREDENTIALS":             acctest.GenerateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
				"GOOGLE_CLOUD_KEYFILE_JSON":      acctest.GenerateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            acctest.GenerateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": acctest.GenerateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedDataModelValue: types.StringValue(acctest.GenerateFakeCredentialsJson("GOOGLE_CREDENTIALS")),
		},
		"when credentials is unset in the config, environment variables are used: GOOGLE_CLOUD_KEYFILE_JSON used second": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				"GOOGLE_CLOUD_KEYFILE_JSON":      acctest.GenerateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            acctest.GenerateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": acctest.GenerateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedDataModelValue: types.StringValue(acctest.GenerateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON")),
		},
		"when credentials is unset in the config, environment variables are used: GCLOUD_KEYFILE_JSON used third": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				// GOOGLE_CLOUD_KEYFILE_JSON not set
				"GCLOUD_KEYFILE_JSON":            acctest.GenerateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": acctest.GenerateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedDataModelValue: types.StringValue(acctest.GenerateFakeCredentialsJson("GCLOUD_KEYFILE_JSON")),
		},
		"when credentials is unset in the config (and access_token unset), GOOGLE_APPLICATION_CREDENTIALS is used for auth but not to set values in the config": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringNull(), // unset
			},
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				// GOOGLE_CLOUD_KEYFILE_JSON not set
				// GCLOUD_KEYFILE_JSON not set
				"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath, // needs to be a path to a file when used by code
			},
			ExpectedDataModelValue: types.StringNull(),
		},
		// Handling empty strings in config
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14255
		// "when credentials is set to an empty string in the config (and access_token unset), GOOGLE_APPLICATION_CREDENTIALS is used": {
		// ConfigValues: fwmodels.ProviderModel{
		// 	Credentials: types.StringValue(""),
		// },
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath, // needs to be a path to a file when used by code
		// 	},
		// 	ExpectedDataModelValue: types.StringNull(),
		// },
		// NOTE: these tests can't run in Cloud Build due to ADC locating credentials despite `GOOGLE_APPLICATION_CREDENTIALS` being unset
		// See https://cloud.google.com/docs/authentication/application-default-credentials#search_order
		// Also, when running these tests locally you need to run `gcloud auth application-default revoke` to ensure your machine isn't supplying ADCs
		// "error returned if credentials is set as an empty string and GOOGLE_APPLICATION_CREDENTIALS is unset": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Credentials: types.StringValue(""),
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": "",
		// 	},
		// 	ExpectError: true,
		// },
		// "error returned if neither credentials nor access_token set in the provider config, and GOOGLE_APPLICATION_CREDENTIALS is unset": {
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": "",
		// 	},
		// 	ExpectError: true,
		// },
		// Handling unknown values - see separate `TestFrameworkProvider_LoadAndValidateFramework_credentials_unknown` test
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			acctest.UnsetTestProviderConfigEnvs(t)
			acctest.SetupTestEnvs(t, tc.EnvVariables)

			ctx := context.Background()
			tfVersion := "foobar"
			providerversion := "999"
			diags := diag.Diagnostics{}

			data := tc.ConfigValues
			impersonateServiceAccountDelegates, _ := types.ListValue(types.StringType, []attr.Value{}) // empty list
			data.ImpersonateServiceAccountDelegates = impersonateServiceAccountDelegates

			p := fwtransport.FrameworkProviderConfig{}

			// Act
			p.LoadAndValidateFramework(ctx, &data, tfVersion, &diags, providerversion)

			// Assert
			if diags.HasError() && tc.ExpectError {
				return
			}
			if diags.HasError() && !tc.ExpectError {
				for i, err := range diags.Errors() {
					num := i + 1
					t.Logf("unexpected error #%d : %s", num, err.Summary())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			if !data.Credentials.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want credentials to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.Credentials.String())
			}
			// fwtransport.FrameworkProviderConfig does not store the credentials info, so test does not make assertions on config struct
		})
	}
}

// TODO(SarahFrench) make this test pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
// func TestFrameworkProvider_LoadAndValidateFramework_credentials_unknown(t *testing.T) {
// 	// This test case is kept separate from other credentials tests, as it requires comparing
// 	// error messages returned by two different error states:
// 	// - When credentials = Null
// 	// - When credentials = Unknown

// 	t.Run("when project is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)", func(t *testing.T) {

// 		// Arrange
// 		acctest.UnsetTestProviderConfigEnvs(t)

// 		ctx := context.Background()
// 		tfVersion := "foobar"
// 		providerversion := "999"

// 		impersonateServiceAccountDelegates, _ := types.ListValue(types.StringType, []attr.Value{}) // empty list

// 		// Null data and error collection
// 		diagsNull := diag.Diagnostics{}
// 		dataNull := fwmodels.ProviderModel{
// 			Credentials: types.StringNull(),
// 		}
// 		dataNull.ImpersonateServiceAccountDelegates = impersonateServiceAccountDelegates

// 		// Unknown data and error collection
// 		diagsUnknown := diag.Diagnostics{}
// 		dataUnknown := fwmodels.ProviderModel{
// 			Credentials: types.StringUnknown(),
// 		}
// 		dataUnknown.ImpersonateServiceAccountDelegates = impersonateServiceAccountDelegates

// 		pNull := fwtransport.FrameworkProviderConfig{}
// 		pUnknown := fwtransport.FrameworkProviderConfig{}

// 		// Act
// 		pNull.LoadAndValidateFramework(ctx, &dataNull, tfVersion, &diagsNull, providerversion)
// 		pUnknown.LoadAndValidateFramework(ctx, &dataUnknown, tfVersion, &diagsUnknown, providerversion)

// 		// Assert
// 		if !diagsNull.HasError() {
// 			t.Fatalf("expect errors when credentials is null, but [%d] errors occurred", diagsNull.ErrorsCount())
// 		}
// 		if !diagsUnknown.HasError() {
// 			t.Fatalf("expect errors when credentials is unknown, but [%d] errors occurred", diagsUnknown.ErrorsCount())
// 		}

// 		errNull := diagsNull.Errors()
// 		errUnknown := diagsUnknown.Errors()
// 		for i := 0; i < len(errNull); i++ {
// 			if errNull[i] != errUnknown[i] {
// 				t.Fatalf("expect errors to be the same for null and unknown credentials values, instead got \nnull=`%s` \nunknown=%s", errNull[i], errUnknown[i])
// 			}
// 		}
// 	})
// }
