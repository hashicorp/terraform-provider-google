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
	// this is to stop the code under test experiencing errors, and could be addressed in future refactoring.
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
				Project: types.StringValue("project-from-config"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("project-from-config"),
			ExpectedConfigStructValue: types.StringValue("project-from-config"),
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
		"when project is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringValue(""),
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		"when project is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT": "project-from-GOOGLE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		// Handling unknown values
		"when project is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				Project: types.StringUnknown(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT": "project-from-GOOGLE_PROJECT",
			},
			ExpectedDataModelValue:    types.StringValue("project-from-GOOGLE_PROJECT"),
			ExpectedConfigStructValue: types.StringValue("project-from-GOOGLE_PROJECT"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
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
	// this is to stop the code under test experiencing errors, and could be addressed in future refactoring.
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
		// Error states
		"when credentials is set to an empty string in the config the value isn't ignored and results in an error": {
			ConfigValues: fwmodels.ProviderModel{
				Credentials: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath, // needs to be a path to a file when used by code
			},
			ExpectError: true,
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			if tc.ExpectError && !diags.HasError() {
				t.Fatalf("expected error, but no errors occurred")
			}
			if !data.Credentials.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want credentials to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.Credentials.String())
			}
			// fwtransport.FrameworkProviderConfig does not store the credentials info, so test does not make assertions on config struct
		})
	}
}

// NOTE: these tests can't run in Cloud Build due to ADC locating credentials despite `GOOGLE_APPLICATION_CREDENTIALS` being unset
// See https://cloud.google.com/docs/authentication/application-default-credentials#search_order
// Also, when running these tests locally you need to run `gcloud auth application-default revoke` to ensure your machine isn't supplying ADCs
// func TestFrameworkProvider_LoadAndValidateFramework_credentials_unknown(t *testing.T) {
// 	// This test case is kept separate from other credentials tests, as it requires comparing
// 	// error messages returned by two different error states:
// 	// - When credentials = Null
// 	// - When credentials = Unknown

// 	t.Run("the same error is returned whether credentials is set as a null or unknown value (and access_token isn't set)", func(t *testing.T) {
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

func TestFrameworkProvider_LoadAndValidateFramework_billingProject(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under test experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues              fwmodels.ProviderModel
		EnvVariables              map[string]string
		ExpectedDataModelValue    basetypes.StringValue
		ExpectedConfigStructValue basetypes.StringValue
		ExpectError               bool
	}{
		"billing_project value set in the provider schema is not overridden by environment variables": {
			ConfigValues: fwmodels.ProviderModel{
				BillingProject: types.StringValue("billing-project-from-config"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("billing-project-from-config"),
			ExpectedConfigStructValue: types.StringValue("billing-project-from-config"),
		},
		"billing_project can be set by environment variable, when no value supplied via the config": {
			ConfigValues: fwmodels.ProviderModel{
				BillingProject: types.StringNull(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("billing-project-from-env"),
			ExpectedConfigStructValue: types.StringValue("billing-project-from-env"),
		},
		"when no billing_project values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				BillingProject: types.StringNull(),
			},
			ExpectedDataModelValue:    types.StringNull(),
			ExpectedConfigStructValue: types.StringNull(),
		},
		// Handling empty strings in config
		"when billing_project is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				BillingProject: types.StringValue(""),
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		"when billing_project is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				BillingProject: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.BillingProject.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want billing_project in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.BillingProject.String())
			}
			// Checking the value passed to the config structs
			if !p.BillingProject.Equal(tc.ExpectedConfigStructValue) {
				t.Fatalf("want billing_project in the `FrameworkProviderConfig` struct to be `%s`, but got the value `%s`", tc.ExpectedConfigStructValue, p.BillingProject.String())
			}
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_region(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under test experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues              fwmodels.ProviderModel
		EnvVariables              map[string]string
		ExpectedDataModelValue    basetypes.StringValue
		ExpectedConfigStructValue basetypes.StringValue
		ExpectError               bool
	}{
		"region value set in the provider config is not overridden by ENVs": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringValue("region-from-config"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("region-from-config"),
			ExpectedConfigStructValue: types.StringValue("region-from-config"),
		},
		"region values can be supplied as a self link": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
			},
			ExpectedDataModelValue:    types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
			ExpectedConfigStructValue: types.StringValue("us-central1"),
		},
		"region value can be set by environment variable: GOOGLE_REGION is used": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringNull(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("region-from-env"),
			ExpectedConfigStructValue: types.StringValue("region-from-env"),
		},
		"when no region values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringNull(),
			},
			ExpectedDataModelValue:    types.StringNull(),
			ExpectedConfigStructValue: types.StringNull(),
		},
		// Handling empty strings in config
		"when region is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringValue(""),
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		"when region is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		// Handling unknown values
		"when region is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				Region: types.StringUnknown(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("region-from-env"),
			ExpectedConfigStructValue: types.StringValue("region-from-env"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.Region.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want region in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.Region.String())
			}
			// Checking the value passed to the config structs
			if !p.Region.Equal(tc.ExpectedConfigStructValue) {
				t.Fatalf("want region in the `FrameworkProviderConfig` struct to be `%s`, but got the value `%s`", tc.ExpectedConfigStructValue, p.Region.String())
			}
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_zone(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under test experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues              fwmodels.ProviderModel
		EnvVariables              map[string]string
		ExpectedDataModelValue    basetypes.StringValue
		ExpectedConfigStructValue basetypes.StringValue
		ExpectError               bool
	}{
		"zone value set in the provider config is not overridden by ENVs": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringValue("zone-from-config"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE": "zone-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("zone-from-config"),
			ExpectedConfigStructValue: types.StringValue("zone-from-config"),
		},
		"does not shorten zone values when provided as a self link": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1"),
			},
			ExpectedDataModelValue:    types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1"),
			ExpectedConfigStructValue: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1"), // Value is not shortened from URI to name
		},
		"when multiple zone environment variables are provided, `GOOGLE_ZONE` is used first": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringNull(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE":           "zone-from-GOOGLE_ZONE",
				"GCLOUD_ZONE":           "zone-from-GCLOUD_ZONE",
				"CLOUDSDK_COMPUTE_ZONE": "zone-from-CLOUDSDK_COMPUTE_ZONE",
			},
			ExpectedDataModelValue:    types.StringValue("zone-from-GOOGLE_ZONE"),
			ExpectedConfigStructValue: types.StringValue("zone-from-GOOGLE_ZONE"),
		},
		"when multiple zone environment variables are provided, `GCLOUD_ZONE` is used second": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringNull(),
			},
			EnvVariables: map[string]string{
				// GOOGLE_ZONE unset
				"GCLOUD_ZONE":           "zone-from-GCLOUD_ZONE",
				"CLOUDSDK_COMPUTE_ZONE": "zone-from-CLOUDSDK_COMPUTE_ZONE",
			},
			ExpectedDataModelValue:    types.StringValue("zone-from-GCLOUD_ZONE"),
			ExpectedConfigStructValue: types.StringValue("zone-from-GCLOUD_ZONE"),
		},
		"when multiple zone environment variables are provided, `CLOUDSDK_COMPUTE_ZONE` is used third": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringNull(),
			},
			EnvVariables: map[string]string{
				// GOOGLE_ZONE unset
				// GCLOUD_ZONE unset
				"CLOUDSDK_COMPUTE_ZONE": "zone-from-CLOUDSDK_COMPUTE_ZONE",
			},
			ExpectedDataModelValue:    types.StringValue("zone-from-CLOUDSDK_COMPUTE_ZONE"),
			ExpectedConfigStructValue: types.StringValue("zone-from-CLOUDSDK_COMPUTE_ZONE"),
		},
		"when no zone values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringNull(),
			},
			ExpectedDataModelValue:    types.StringNull(),
			ExpectedConfigStructValue: types.StringNull(),
		},
		// Handling empty strings in config
		"when zone is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringValue(""),
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		"when zone is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE": "zone-from-env",
			},
			ExpectedDataModelValue:    types.StringValue(""),
			ExpectedConfigStructValue: types.StringValue(""),
		},
		// Handling unknown values
		"when zone is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				Zone: types.StringUnknown(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE": "zone-from-env",
			},
			ExpectedDataModelValue:    types.StringValue("zone-from-env"),
			ExpectedConfigStructValue: types.StringValue("zone-from-env"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.Zone.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want zone in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.Zone.String())
			}
			// Checking the value passed to the config structs
			if !p.Zone.Equal(tc.ExpectedConfigStructValue) {
				t.Fatalf("want zone in the `FrameworkProviderConfig` struct to be `%s`, but got the value `%s`", tc.ExpectedConfigStructValue, p.Zone.String())
			}
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_accessToken(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues           fwmodels.ProviderModel
		EnvVariables           map[string]string
		ExpectedDataModelValue basetypes.StringValue // Sometimes the value is mutated, and no longer matches the original value we supply
		// ExpectedConfigStructValue not used here, as credentials info isn't stored in the config struct
		ExpectError bool
	}{
		"access_token configured in the provider can be invalid without resulting in errors": {
			ConfigValues: fwmodels.ProviderModel{
				AccessToken: types.StringValue("This is not a valid token string"),
			},
			ExpectedDataModelValue: types.StringValue("This is not a valid token string"),
		},
		"access_token set in the provider config is not overridden by environment variables": {
			ConfigValues: fwmodels.ProviderModel{
				AccessToken: types.StringValue("value-from-config"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-env",
			},
			ExpectedDataModelValue: types.StringValue("value-from-config"),
		},
		"when access_token is unset in the config, the GOOGLE_OAUTH_ACCESS_TOKEN environment variable is used": {
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedDataModelValue: types.StringValue("value-from-GOOGLE_OAUTH_ACCESS_TOKEN"),
		},
		"when no access_token values are provided via config or environment variables there's no error (as long as credentials supplied in its absence)": {
			ConfigValues: fwmodels.ProviderModel{
				AccessToken: types.StringNull(),
				Credentials: types.StringValue(transport_tpg.TestFakeCredentialsPath),
			},
			ExpectedDataModelValue: types.StringNull(),
		},
		// Handling empty strings in config
		"when access_token is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				AccessToken: types.StringValue(""),
			},
			ExpectedDataModelValue: types.StringValue(""),
		},
		"when access_token is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				AccessToken: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedDataModelValue: types.StringValue(""),
		},
		// Handling unknown values
		"when access_token is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				AccessToken: types.StringUnknown(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedDataModelValue: types.StringValue("value-from-GOOGLE_OAUTH_ACCESS_TOKEN"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.AccessToken.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want project in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.AccessToken.String())
			}
			// fwtransport.FrameworkProviderConfig does not store the credentials info, so test does not make assertions on config struct
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_userProjectOverride(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues              fwmodels.ProviderModel
		EnvVariables              map[string]string
		ExpectedDataModelValue    basetypes.BoolValue
		ExpectedConfigStructValue basetypes.BoolValue
		ExpectError               bool
	}{
		"user_project_override value set in the provider schema is not overridden by ENVs": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolValue(false),
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedDataModelValue:    types.BoolValue(false),
			ExpectedConfigStructValue: types.BoolValue(false),
		},
		"user_project_override can be set by environment variable: value = true": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolNull(), // not set
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedDataModelValue:    types.BoolValue(true),
			ExpectedConfigStructValue: types.BoolValue(true),
		},
		"user_project_override can be set by environment variable: value = false": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolNull(), // not set
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "false",
			},
			ExpectedDataModelValue:    types.BoolValue(false),
			ExpectedConfigStructValue: types.BoolValue(false),
		},
		"user_project_override can be set by environment variable: value = 1": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolNull(), // not set
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "1",
			},
			ExpectedDataModelValue:    types.BoolValue(true),
			ExpectedConfigStructValue: types.BoolValue(true),
		},
		"user_project_override can be set by environment variable: value = 0": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolNull(), // not set
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "0",
			},
			ExpectedDataModelValue:    types.BoolValue(false),
			ExpectedConfigStructValue: types.BoolValue(false),
		},
		"setting user_project_override using a non-boolean environment variables results in an error": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "I'm not a boolean",
			},
			ExpectError: true,
		},
		"when no user_project_override values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolNull(), // not set
			},
			ExpectedDataModelValue:    types.BoolNull(),
			ExpectedConfigStructValue: types.BoolNull(),
		},
		// Handling unknown values
		"when user_project_override is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				UserProjectOverride: types.BoolUnknown(),
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedDataModelValue:    types.BoolValue(true),
			ExpectedConfigStructValue: types.BoolValue(true),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.UserProjectOverride.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want user_project_override in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.UserProjectOverride.String())
			}
			// Checking the value passed to the config structs
			if !p.UserProjectOverride.Equal(tc.ExpectedConfigStructValue) {
				t.Fatalf("want user_project_override in the `FrameworkProviderConfig` struct to be `%s`, but got the value `%s`", tc.ExpectedConfigStructValue, p.UserProjectOverride.String())
			}
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_impersonateServiceAccount(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues              fwmodels.ProviderModel
		EnvVariables              map[string]string
		ExpectedDataModelValue    basetypes.StringValue
		ExpectedConfigStructValue basetypes.StringValue
		ExpectError               bool
	}{
		"impersonate_service_account value set in the provider schema is not overridden by environment variables": {
			ConfigValues: fwmodels.ProviderModel{
				ImpersonateServiceAccount: types.StringValue("value-from-config@example.com"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedDataModelValue: types.StringValue("value-from-config@example.com"),
		},
		"impersonate_service_account value can be set by environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				ImpersonateServiceAccount: types.StringNull(), // not set
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedDataModelValue: types.StringValue("value-from-env@example.com"),
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				ImpersonateServiceAccount: types.StringNull(), // not set
			},
			ExpectedDataModelValue: types.StringNull(),
		},
		// Handling empty strings in config
		"when impersonate_service_account is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				ImpersonateServiceAccount: types.StringValue(""),
			},
			ExpectedDataModelValue: types.StringValue(""),
		},
		"when impersonate_service_account is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				ImpersonateServiceAccount: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedDataModelValue: types.StringValue(""),
		},
		// Handling unknown values
		"when impersonate_service_account is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				ImpersonateServiceAccount: types.StringUnknown(),
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedDataModelValue: types.StringValue("value-from-env@example.com"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.ImpersonateServiceAccount.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want impersonate_service_account in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.ImpersonateServiceAccount.String())
			}
			// fwtransport.FrameworkProviderConfig does not store impersonate_service_account info, so test does not make assertions on config struct
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_impersonateServiceAccountDelegates(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test

	cases := map[string]struct {
		// It's not easy to define basetypes.ListValue values directly in test case, so instead
		// pass values into test function to control construction of basetypes.ListValue there.
		SetAsNull                               bool
		SetAsUnknown                            bool
		ImpersonateServiceAccountDelegatesValue []string
		EnvVariables                            map[string]string

		ExpectedNull           bool
		ExpectedUnknown        bool
		ExpectedDataModelValue []string
		ExpectError            bool
	}{
		"impersonate_service_account_delegates value can be set in the provider schema": {
			ImpersonateServiceAccountDelegatesValue: []string{
				"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
				"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
			},
			ExpectedDataModelValue: []string{
				"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
				"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
			},
		},
		// Note: no environment variables can be used for impersonate_service_account_delegates
		"when no impersonate_service_account_delegates value is provided via config, the field remains unset without error": {
			SetAsNull:    true, // not setting impersonate_service_account_delegates
			ExpectedNull: true,
		},
		// Handling empty values in config
		"when impersonate_service_account_delegates is set as an empty array, that value isn't ignored": {
			ImpersonateServiceAccountDelegatesValue: []string{},
			ExpectedDataModelValue:                  []string{},
		},
		// Handling unknown values
		"when impersonate_service_account_delegates is an unknown value, the provider treats it as if it's unset, without error": {
			SetAsUnknown:    true,
			ExpectedUnknown: true,
		},
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

			data := fwmodels.ProviderModel{}
			data.Credentials = types.StringValue(transport_tpg.TestFakeCredentialsPath)
			// Set ImpersonateServiceAccountDelegates depending on test case
			if !tc.SetAsNull && !tc.SetAsUnknown {
				isad, _ := types.ListValueFrom(ctx, types.StringType, tc.ImpersonateServiceAccountDelegatesValue)
				data.ImpersonateServiceAccountDelegates = isad
			}
			if tc.SetAsNull {
				data.ImpersonateServiceAccountDelegates = types.ListNull(types.StringType)
			}
			if tc.SetAsUnknown {
				data.ImpersonateServiceAccountDelegates = types.ListUnknown(types.StringType)
			}

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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			var expected attr.Value
			if !tc.ExpectedNull && !tc.ExpectedUnknown {
				expected, _ = types.ListValueFrom(ctx, types.StringType, tc.ExpectedDataModelValue)
			}
			if tc.ExpectedNull {
				expected = types.ListNull(types.StringType)
			}
			if tc.ExpectedUnknown {
				expected = types.ListUnknown(types.StringType)
			}
			if !data.ImpersonateServiceAccountDelegates.Equal(expected) {
				t.Fatalf("want impersonate_service_account in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", expected, data.ImpersonateServiceAccountDelegates.String())
			}
			// fwtransport.FrameworkProviderConfig does not store impersonate_service_account info, so test does not make assertions on config struct
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_scopes(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ScopesValue               []string
		EnvVariables              map[string]string
		ExpectedDataModelValue    []string
		ExpectedConfigStructValue []string
		SetAsNull                 bool
		SetAsUnknown              bool
		ExpectError               bool
	}{
		"scopes are set in the provider config as a list": {
			ScopesValue:               []string{"fizz", "buzz", "baz"},
			ExpectedDataModelValue:    []string{"fizz", "buzz", "baz"},
			ExpectedConfigStructValue: []string{"fizz", "buzz", "baz"},
		},
		"scopes can be left unset in the provider config without any issues, and a default value is used": {
			SetAsNull:                 true,
			ExpectedDataModelValue:    transport_tpg.DefaultClientScopes,
			ExpectedConfigStructValue: transport_tpg.DefaultClientScopes,
		},
		// Handling empty values in config
		"scopes set as an empty list the field is treated as if it's unset and a default value is used without errors": {
			ScopesValue:               []string{},
			ExpectedDataModelValue:    transport_tpg.DefaultClientScopes,
			ExpectedConfigStructValue: transport_tpg.DefaultClientScopes,
		},
		// Handling unknown values
		"when scopes is an unknown value, the provider treats it as if it's unset and a default value is used without errors": {
			SetAsUnknown:              true,
			ExpectedDataModelValue:    transport_tpg.DefaultClientScopes,
			ExpectedConfigStructValue: transport_tpg.DefaultClientScopes,
		},
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

			data := fwmodels.ProviderModel{}
			data.Credentials = types.StringValue(transport_tpg.TestFakeCredentialsPath)
			impersonateServiceAccountDelegates, _ := types.ListValue(types.StringType, []attr.Value{}) // empty list
			data.ImpersonateServiceAccountDelegates = impersonateServiceAccountDelegates
			// Set ImpersonateServiceAccountDelegates depending on test case
			if !tc.SetAsNull && !tc.SetAsUnknown {
				s, _ := types.ListValueFrom(ctx, types.StringType, tc.ScopesValue)
				data.Scopes = s
			}
			if tc.SetAsNull {
				data.Scopes = types.ListNull(types.StringType)
			}
			if tc.SetAsUnknown {
				data.Scopes = types.ListUnknown(types.StringType)
			}

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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			expectedDm, _ := types.ListValueFrom(ctx, types.StringType, tc.ExpectedDataModelValue)
			if !data.Scopes.Equal(expectedDm) {
				t.Fatalf("want project in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.Scopes.String())
			}
			// Checking the value passed to the config structs
			expectedFpc, _ := types.ListValueFrom(ctx, types.StringType, tc.ExpectedConfigStructValue)
			if !p.Scopes.Equal(expectedFpc) {
				t.Fatalf("want project in the `FrameworkProviderConfig` struct to be `%s`, but got the value `%s`", tc.ExpectedConfigStructValue, p.Scopes.String())
			}
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_requestReason(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues           fwmodels.ProviderModel
		EnvVariables           map[string]string
		ExpectedDataModelValue basetypes.StringValue
		// ExpectedConfigStructValue not used here, as credentials info isn't stored in the config struct
		ExpectError bool
	}{
		"when request_reason is unset in the config, environment variable CLOUDSDK_CORE_REQUEST_REASON is used": {
			ConfigValues: fwmodels.ProviderModel{
				RequestReason: types.StringNull(),
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "foo",
			},
			ExpectedDataModelValue: types.StringValue("foo"),
		},
		"request_reason set in the config is not overridden by environment variables": {
			ConfigValues: fwmodels.ProviderModel{
				RequestReason: types.StringValue("value-from-config"),
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "value-from-env",
			},
			ExpectedDataModelValue: types.StringValue("value-from-config"),
		},
		"when no request_reason is provided via config or environment variables, the field remains unset without error": {
			ConfigValues: fwmodels.ProviderModel{
				RequestReason: types.StringNull(),
			},
			ExpectedDataModelValue: types.StringNull(),
		},
		// Handling empty strings in config
		"when request_reason is set as an empty string, the empty string is not ignored in favor of an environment variable": {
			ConfigValues: fwmodels.ProviderModel{
				RequestReason: types.StringValue(""),
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "foo",
			},
			ExpectedDataModelValue: types.StringValue(""),
		},
		"when request_reason is set as an empty string the empty string is used and not ignored": {
			ConfigValues: fwmodels.ProviderModel{
				RequestReason: types.StringValue(""),
			},
			ExpectedDataModelValue: types.StringValue(""),
		},
		// Handling unknown values
		"when request_reason is an unknown value, the provider treats it as if it's unset and uses an environment variable instead": {
			ConfigValues: fwmodels.ProviderModel{
				RequestReason: types.StringUnknown(),
			},
			EnvVariables: map[string]string{
				"CLOUDSDK_CORE_REQUEST_REASON": "foo",
			},
			ExpectedDataModelValue: types.StringValue("foo"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.RequestReason.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want request_reason in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.RequestReason.String())
			}
			// fwtransport.FrameworkProviderConfig does not store the request reason info, so test does not make assertions on config struct
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_requestTimeout(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		ConfigValues           fwmodels.ProviderModel
		EnvVariables           map[string]string
		ExpectedDataModelValue basetypes.StringValue
		// ExpectedConfigStructValue not used here, as credentials info isn't stored in the config struct
		ExpectError bool
	}{
		"if a valid request_timeout is configured in the provider, no error will occur": {
			ConfigValues: fwmodels.ProviderModel{
				RequestTimeout: types.StringValue("10s"),
			},
			ExpectedDataModelValue: types.StringValue("10s"),
		},
		"if an invalid request_timeout is configured in the provider, an error will occur": {
			ConfigValues: fwmodels.ProviderModel{
				RequestTimeout: types.StringValue("timeout"),
			},
			ExpectError: true,
		},
		"when request_timeout is set as an empty string, the empty string isn't ignored and an error will occur": {
			ConfigValues: fwmodels.ProviderModel{
				RequestTimeout: types.StringValue(""),
			},
			ExpectError: true,
		},
		// In the SDK version of the provider config code, this scenario results in a value of "0s"
		// instead of "120s", but the final 'effective' value is also "120s"
		// See : https://github.com/hashicorp/terraform-provider-google/blob/09cb850ee64bcd78e4457df70905530c1ed75f19/google/transport/config.go#L1228-L1233
		"when request_timeout is unset in the config, the default value is 120s.": {
			ConfigValues: fwmodels.ProviderModel{
				RequestTimeout: types.StringNull(),
			},
			ExpectedDataModelValue: types.StringValue("120s"),
		},
		// Handling unknown values
		"when request_timeout is an unknown value, the provider treats it as if it's unset and uses the default value 120s": {
			ConfigValues: fwmodels.ProviderModel{
				RequestTimeout: types.StringUnknown(),
			},
			ExpectedDataModelValue: types.StringValue("120s"),
		},
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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			if !data.RequestTimeout.Equal(tc.ExpectedDataModelValue) {
				t.Fatalf("want request_timeout in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectedDataModelValue, data.RequestTimeout.String())
			}
			// fwtransport.FrameworkProviderConfig does not store the request timeout info, so test does not make assertions on config struct
		})
	}
}

func TestFrameworkProvider_LoadAndValidateFramework_batching(t *testing.T) {

	// Note: In the test function we need to set the below fields in test case's fwmodels.ProviderModel value
	// this is to stop the code under tests experiencing errors, and could be addressed in future refactoring.
	// - Credentials: If we don't set this then the test looks for application default credentials and can fail depending on the machine running the test
	// - ImpersonateServiceAccountDelegates: If we don't set this, we get a nil pointer exception ¯\_(ツ)_/¯

	cases := map[string]struct {
		// It's not easy to create the value of Batching in the test case, so these inputs are used in the test function
		SetBatchingAsNull    bool
		SetBatchingAsUnknown bool
		EnableBatchingValue  basetypes.BoolValue
		SendAfterValue       basetypes.StringValue

		EnvVariables map[string]string

		ExpectBatchingNull        bool
		ExpectBatchingUnknown     bool
		ExpectEnableBatchingValue basetypes.BoolValue
		ExpectSendAfterValue      basetypes.StringValue
		ExpectError               bool
	}{
		"batching can be configured with values for enable_batching and send_after": {
			EnableBatchingValue:       types.BoolValue(true),
			SendAfterValue:            types.StringValue("45s"),
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("45s"),
		},
		"if batching is an empty block, it will set the default values for enable_batching and send_after": {
			// In this test, we try to create a list containing only null values
			EnableBatchingValue:       types.BoolNull(),
			SendAfterValue:            types.StringNull(),
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("10s"),
		},
		"when batching is configured with only enable_batching, send_after will be set to a default value": {
			EnableBatchingValue:       types.BoolValue(true),
			SendAfterValue:            types.StringNull(),
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("10s"),
		},
		"when batching is configured with only send_after, enable_batching will be set to a default value": {
			EnableBatchingValue:       types.BoolNull(),
			SendAfterValue:            types.StringValue("45s"),
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("45s"),
		},
		"when the whole batching block is a null value, the provider provides default values for send_after and enable_batching": {
			SetBatchingAsNull:         true,
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("3s"),
		},
		// Handling unknown values
		"when batching is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
			SetBatchingAsUnknown:      true,
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("3s"),
		},
		"when batching is configured with send_after as an unknown value, send_after will be set to a default value": {
			EnableBatchingValue:       types.BoolValue(true),
			SendAfterValue:            types.StringUnknown(),
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("10s"),
		},
		"when batching is configured with enable_batching as an unknown value, enable_batching will be set to a default value": {
			EnableBatchingValue:       types.BoolUnknown(),
			SendAfterValue:            types.StringValue("45s"),
			ExpectEnableBatchingValue: types.BoolValue(true),
			ExpectSendAfterValue:      types.StringValue("45s"),
		},
		// Error states
		"when batching is configured with send_after as an empty string, the empty string is not ignored and results in an error": {
			EnableBatchingValue: types.BoolValue(true),
			SendAfterValue:      types.StringValue(""),
			ExpectError:         true,
		},
		"if batching is configured with send_after as an invalid value, there's an error": {
			SendAfterValue: types.StringValue("invalid value"),
			ExpectError:    true,
		},
		"if batching is configured with send_after as number value without seconds (s), there's an error": {
			SendAfterValue: types.StringValue("123"),
			ExpectError:    true,
		},
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

			data := fwmodels.ProviderModel{}
			data.Credentials = types.StringValue(transport_tpg.TestFakeCredentialsPath)
			impersonateServiceAccountDelegates, _ := types.ListValue(types.StringType, []attr.Value{}) // empty list
			data.ImpersonateServiceAccountDelegates = impersonateServiceAccountDelegates

			// TODO(SarahFrench) - this code will change when batching is reworked
			// See https://github.com/GoogleCloudPlatform/magic-modules/pull/7668
			if !tc.SetBatchingAsNull && !tc.SetBatchingAsUnknown {
				b, _ := types.ObjectValue(
					map[string]attr.Type{
						"enable_batching": types.BoolType,
						"send_after":      types.StringType,
					},
					map[string]attr.Value{
						"enable_batching": tc.EnableBatchingValue,
						"send_after":      tc.SendAfterValue,
					},
				)
				batching, _ := types.ListValue(types.ObjectType{}.WithAttributeTypes(fwmodels.ProviderBatchingAttributes), []attr.Value{b})
				data.Batching = batching
			}
			if tc.SetBatchingAsNull {
				data.Batching = types.ListNull(types.ObjectType{}.WithAttributeTypes(fwmodels.ProviderBatchingAttributes))
			}
			if tc.SetBatchingAsUnknown {
				data.Batching = types.ListUnknown(types.ObjectType{}.WithAttributeTypes(fwmodels.ProviderBatchingAttributes))
			}

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
			if !data.Batching.IsNull() && tc.ExpectBatchingNull {
				t.Fatalf("want batching in the `fwmodels.ProviderModel` struct to be null, but got the value `%s`", data.Batching.String())
			}
			if !data.Batching.IsUnknown() && tc.ExpectBatchingUnknown {
				t.Fatalf("want batching in the `fwmodels.ProviderModel` struct to be unknown, but got the value `%s`", data.Batching.String())
			}

			// The code doesn't mutate values in the fwmodels.ProviderModel struct if the whole batching block is null/unknown,
			// so run these checks below only if we're not setting the whole batching block is null/unknown
			if !tc.SetBatchingAsNull && !tc.SetBatchingAsUnknown {
				var pbConfigs []fwmodels.ProviderBatching
				_ = data.Batching.ElementsAs(ctx, &pbConfigs, true)
				if !pbConfigs[0].EnableBatching.Equal(tc.ExpectEnableBatchingValue) {
					t.Fatalf("want batching.enable_batching in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectEnableBatchingValue.String(), pbConfigs[0].EnableBatching.String())
				}
				if !pbConfigs[0].SendAfter.Equal(tc.ExpectSendAfterValue) {
					t.Fatalf("want batching.send_after in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", tc.ExpectSendAfterValue.String(), pbConfigs[0].SendAfter.String())
				}
			}

			// Check how the batching block's values are used to configure other parts of the `FrameworkProviderConfig` struct
			// - RequestBatcherServiceUsage
			// - RequestBatcherIam
			if p.RequestBatcherServiceUsage.BatchingConfig.EnableBatching != tc.ExpectEnableBatchingValue.ValueBool() {
				t.Fatalf("want batching.enable_batching to be `%s`, but got the value `%v`", tc.ExpectEnableBatchingValue.String(), p.RequestBatcherServiceUsage.BatchingConfig.EnableBatching)
			}
			if !types.StringValue(p.RequestBatcherServiceUsage.BatchingConfig.SendAfter.String()).Equal(tc.ExpectSendAfterValue) {
				t.Fatalf("want batching.send_after to be `%s`, but got the value `%s`", tc.ExpectSendAfterValue.String(), p.RequestBatcherServiceUsage.BatchingConfig.SendAfter.String())
			}
			if p.RequestBatcherIam.BatchingConfig.EnableBatching != tc.ExpectEnableBatchingValue.ValueBool() {
				t.Fatalf("want batching.enable_batching to be `%s`, but got the value `%v`", tc.ExpectEnableBatchingValue.String(), p.RequestBatcherIam.BatchingConfig.EnableBatching)
			}
			if !types.StringValue(p.RequestBatcherIam.BatchingConfig.SendAfter.String()).Equal(tc.ExpectSendAfterValue) {
				t.Fatalf("want batching.send_after to be `%s`, but got the value `%s`", tc.ExpectSendAfterValue.String(), p.RequestBatcherIam.BatchingConfig.SendAfter.String())
			}
		})
	}
}

func TestGetRegionFromRegionSelfLink(t *testing.T) {
	cases := map[string]struct {
		Input          basetypes.StringValue
		ExpectedOutput basetypes.StringValue
	}{
		"A short region name is returned unchanged": {
			Input:          types.StringValue("us-central1"),
			ExpectedOutput: types.StringValue("us-central1"),
		},
		"A selflink is shortened to a region name": {
			Input:          types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
			ExpectedOutput: types.StringValue("us-central1"),
		},
		"Logic is specific to region selflinks; zone selflinks are not shortened": {
			Input:          types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/asia-east1-a"),
			ExpectedOutput: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/asia-east1-a"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			region := fwtransport.GetRegionFromRegionSelfLink(tc.Input)

			if region != tc.ExpectedOutput {
				t.Fatalf("want %s,  got %s", region, tc.ExpectedOutput)
			}
		})
	}
}
