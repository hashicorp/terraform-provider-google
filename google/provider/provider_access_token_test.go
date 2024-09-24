// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccSdkProvider_access_token is a series of acc tests asserting how the SDK provider handles access_token arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_access_token
func TestAccSdkProvider_access_token(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                        testAccSdkProvider_access_token_configPrecedenceOverEnvironmentVariables,
		"when access_token is unset in the config, environment variables are used in a given order": testAccSdkProvider_access_token_precedenceOrderEnvironmentVariables, // GOOGLE_OAUTH_ACCESS_TOKEN

		// Schema-level validation
		"when access_token is set to an empty string in the config the value isn't ignored and results in an error": testAccSdkProvider_access_token_emptyStringValidation,
		"access_token conflicts with credentials": testAccSdkProvider_access_token_conflictsWithCredentials,

		// Usage
		"access_token can be used to authenticate the provider": testAccSdkProvider_access_token_authInUse,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccSdkProvider_access_token_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	accessToken := "my-access-token"

	// ensure all possible access_token env vars set; show they aren't used instead
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", accessToken)

	// ensure credentials ENVs not set; they're used over access_token
	for _, env := range envvar.CredsEnvVars {
		t.Setenv(env, "")
	}

	providerAccessToken := "foobar"

	context := map[string]interface{}{
		"access_token": providerAccessToken,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Apply-time error; bad value in config is used over of good values in ENVs
				Config: testAccSdkProvider_access_tokenInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "access_token", providerAccessToken),
				),
			},
		},
	})
}

func testAccSdkProvider_access_token_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for access_token
		GOOGLE_OAUTH_ACCESS_TOKEN
	*/

	GOOGLE_OAUTH_ACCESS_TOKEN := "GOOGLE_OAUTH_ACCESS_TOKEN"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// GOOGLE_OAUTH_ACCESS_TOKEN is used if config doesn't provide a value
				PreConfig: func() {
					t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", GOOGLE_OAUTH_ACCESS_TOKEN) //used
				},
				Config: testAccSdkProvider_access_tokenInEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "access_token", GOOGLE_OAUTH_ACCESS_TOKEN),
				),
			},
		},
	})
}

func testAccSdkProvider_access_token_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	accessToken := "my-access-token"

	// ensure all access_token env vars set
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", accessToken)

	context := map[string]interface{}{
		"access_token": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_access_tokenInProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccSdkProvider_access_token_conflictsWithCredentials(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	goodCredentials := envvar.GetTestCredsFromEnv()

	// unset ENVs for both access_token and credentials
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "")
	for _, env := range envvar.CredsEnvVars {
		t.Setenv(env, "")
	}

	accessToken := "my-access-token"
	credentials := "./path/to/fake/credentials.json"

	providerBlockBoth := fmt.Sprintf(`
access_token = "%s"
credentials = "%s"
`, accessToken, credentials)

	providerAccessTokenOnly := fmt.Sprintf(`
	access_token = "%s"
	`, accessToken)

	providerCredentialsOnly := fmt.Sprintf(`
	credentials = "%s"
	`, credentials)

	contextBoth := map[string]interface{}{
		"fields": providerBlockBoth,
	}

	contextAccessTokenOnly := map[string]interface{}{
		"fields": providerAccessTokenOnly,
	}

	contextCredentialsOnly := map[string]interface{}{
		"fields": providerCredentialsOnly,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Conflicts when both set in the provider block in the configuration
				Config:      testAccSdkProvider_access_token_conflictingFields(contextBoth),
				ExpectError: regexp.MustCompile("Attribute \"access_token\" cannot be specified when \"credentials\" is specified"),
			},
			{
				// No conflict when access_token in the provider block, credentials in ENVs.
				PreConfig: func() {
					t.Setenv("GOOGLE_CREDENTIALS", credentials)
					t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "")
				},
				Config: testAccSdkProvider_access_token_conflictingFields(contextAccessTokenOnly),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "access_token", accessToken),
					// not set, as ENV not used
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "credentials", ""),
				),
			},
			{
				// No conflict when credentials in the provider block, access_token in ENVs.
				PreConfig: func() {
					t.Setenv("GOOGLE_CREDENTIALS", "")
					t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", accessToken)
				},
				Config: testAccSdkProvider_access_token_conflictingFields(contextCredentialsOnly),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "credentials", credentials),
					// not set, as ENV not used
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "access_token", ""),
				),
				ExpectError: regexp.MustCompile("JSON credentials are not valid"),
			},
			{
				PreConfig: func() {
					t.Setenv("GOOGLE_CREDENTIALS", goodCredentials)
				},
				Destroy: true,
				Config:  "// Empty config and good credentials in this step to avoid post-test destroy error",
			},
		},
	})
}

func testAccSdkProvider_access_token_authInUse(t *testing.T) {
	// Skip if not running in a acc test environment,
	// as acc test environment variables needed to get accessToken
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	// Access token to pass in via config
	accessToken := acctest.GetAccessTokenFromTestCredsFromEnv(t)

	// Inputs ready
	context := map[string]interface{}{
		"access_token":  accessToken,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// unset all relevant ENVs, so value passed in via
					// context is the only usable input
					t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "")
					for _, env := range envvar.CredsEnvVars {
						t.Setenv(env, "")
					}
				},
				Config: testAccSdkProvider_access_token_useAccessToken(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Assert provider is using access_token argument for auth
					resource.TestCheckResourceAttrSet("data.google_provider_config_sdk.default", "access_token"),
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "credentials", ""),
				),
			},
		},
	})
}

// testAccSdkProvider_access_tokenInProviderBlock allows setting the access_token argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_access_tokenInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	access_token = "%{access_token}"
}

data "google_provider_config_sdk" "default" {}

output "access_token" {
  value = data.google_provider_config_sdk.default.access_token
  sensitive = true
}
`, context)
}

// testAccSdkProvider_access_tokenInEnvsOnly allows testing when the access_token argument
// is only supplied via ENVs
func testAccSdkProvider_access_tokenInEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}

output "access_token" {
  value = data.google_provider_config_sdk.default.access_token
  sensitive = true
}
`, context)
}

// testAccSdkProvider_access_token_conflictingFields allows setting multiple fields in the provider
// block to test conflict validation in the provider schema
func testAccSdkProvider_access_token_conflictingFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
%{fields}
}

data "google_provider_config_sdk" "default" {}

output "access_token" {
  value = data.google_provider_config_sdk.default.access_token
  sensitive = true
}
`, context)
}

func testAccSdkProvider_access_token_useAccessToken(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  access_token = "%{access_token}"
}

data "google_provider_config_sdk" "default" {
}

resource "google_service_account" "default" {
  account_id   = "tf-test-%{random_suffix}"
  display_name = "AccTest Service Account"
}
`, context)
}
