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
					t.Logf("unexpected error #%d : %s : %s", num, err.Summary(), err.Detail())
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
		// "when billing_project is set as an empty string the field is treated as if it's unset, without error": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		BillingProject: types.StringValue(""),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
		// 	ExpectedConfigStructValue: types.StringNull(),
		// },
		// "when billing_project is set as an empty string an environment variable will be used": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		BillingProject: types.StringValue(""),
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_BILLING_PROJECT": "billing-project-from-env",
		// 	},
		// 	ExpectedDataModelValue:    types.StringValue("billing-project-from-env"),
		// 	ExpectedConfigStructValue: types.StringValue("billing-project-from-env"),
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
		// This test currently fails - PF code doesn't behave like SDK code
		// TODO(SarahFrench) - address https://github.com/hashicorp/terraform-provider-google/issues/15714
		// "region values can be supplied as a self link": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Region: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
		// 	},
		// 	ExpectedDataModelValue:    types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
		// 	ExpectedConfigStructValue: types.StringValue("us-central1"),
		// },
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
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14255
		// "when region is set as an empty string the field is treated as if it's unset, without error": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Region: types.StringValue(""),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
		// 	ExpectedConfigStructValue: types.StringNull(),
		// },
		// "when region is set as an empty string an environment variable will be used": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Region: types.StringValue(""),
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_REGION": "region-from-env",
		// 	},
		// 	ExpectedDataModelValue:    types.StringValue("region-from-env"),
		// 	ExpectedConfigStructValue: types.StringValue("region-from-env"),
		// },
		// Handling unknown values
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when region is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Region: types.StringUnknown(),
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
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14255
		// "when zone is set as an empty string the field is treated as if it's unset, without error": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Zone: types.StringValue(""),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
		// 	ExpectedConfigStructValue: types.StringNull(),
		// },
		// "when zone is set as an empty string an environment variable will be used": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Zone: types.StringValue(""),
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_ZONE": "zone-from-env",
		// 	},
		// 	ExpectedDataModelValue:    types.StringValue("zone-from-env"),
		// 	ExpectedConfigStructValue: types.StringValue("zone-from-env"),
		// },
		// Handling unknown values
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when zone is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		Zone: types.StringUnknown(),
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
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14255
		// "when access_token is set as an empty string the field is treated as if it's unset, without error (as long as credentials supplied in its absence)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		AccessToken: types.StringValue(""),
		// 		Credentials: types.StringValue(transport_tpg.TestFakeCredentialsPath),
		// 	},
		// 	ExpectedDataModelValue: types.StringNull(),
		// },
		// "when access_token is set as an empty string in the config, an environment variable is used": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		AccessToken: types.StringValue(""),
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
		// 	},
		// 	ExpectedDataModelValue: types.StringValue("value-from-GOOGLE_OAUTH_ACCESS_TOKEN"),
		// },
		// Handling unknown values
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when access_token is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		AccessToken: types.StringUnknown(),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
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
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when user_project_override is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		UserProjectOverride: types.BoolUnknown(),
		// 	},
		// 	ExpectedDataModelValue:    types.BoolNull(),
		// 	ExpectedConfigStructValue: types.BoolNull(),
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
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14255
		// "when impersonate_service_account is set as an empty array the field is treated as if it's unset, without error": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		ImpersonateServiceAccount: types.StringValue(""),
		// 	},
		// 	ExpectedDataModelValue: types.StringNull(),
		// },
		// Handling unknown values
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when impersonate_service_account is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	ConfigValues: fwmodels.ProviderModel{
		// 		ImpersonateServiceAccount: types.StringUnknown(),
		// 	},
		// 	ExpectedDataModelValue:    types.StringNull(),
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
		ExpectedDataModelValue                  []string
		// ExpectedConfigStructValue not used here, as impersonate_service_account_delegates info isn't stored in the config struct
		ExpectError bool
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
			SetAsNull:              true, // not setting impersonate_service_account_delegates
			ExpectedDataModelValue: nil,
		},
		// Handling empty values in config
		"when impersonate_service_account_delegates is set as an empty array the field is treated as if it's unset, without error": {
			ImpersonateServiceAccountDelegatesValue: []string{},
			ExpectedDataModelValue:                  []string{},
		},
		// Handling unknown values
		// TODO(SarahFrench) make these tests pass to address: https://github.com/hashicorp/terraform-provider-google/issues/14444
		// "when impersonate_service_account_delegates is an unknown value, the provider treats it as if it's unset (align to SDK behaviour)": {
		// 	SetAsUnknown: true,
		// 	// Currently this causes an error at google/fwtransport/framework_config.go:1518
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
					t.Logf("unexpected error #%d : %s", num, err.Summary())
				}
				t.Fatalf("did not expect error, but [%d] error(s) occurred", diags.ErrorsCount())
			}
			// Checking mutation of the data model
			expected, _ := types.ListValueFrom(ctx, types.StringType, tc.ExpectedDataModelValue)
			if !data.ImpersonateServiceAccountDelegates.Equal(expected) {
				t.Fatalf("want impersonate_service_account in the `fwmodels.ProviderModel` struct to be `%s`, but got the value `%s`", expected, data.ImpersonateServiceAccountDelegates.String())
			}
			// fwtransport.FrameworkProviderConfig does not store impersonate_service_account info, so test does not make assertions on config struct
		})
	}
}
