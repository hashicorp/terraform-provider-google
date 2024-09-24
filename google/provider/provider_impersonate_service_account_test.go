// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccSdkProvider_impersonate_service_account is a series of acc tests asserting how the SDK provider handles impersonate_service_account arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_impersonate_service_account
func TestAccSdkProvider_impersonate_service_account(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                                       testAccSdkProvider_impersonate_service_account_configPrecedenceOverEnvironmentVariables,
		"when impersonate_service_account is unset in the config, environment variables are used in a given order": testAccSdkProvider_impersonate_service_account_precedenceOrderEnvironmentVariables, // GOOGLE_IMPERSONATE_SERVICE_ACCOUNT

		// Schema-level validation
		"when impersonate_service_account is set to an empty string in the config the value isn't ignored and results in an error": testAccSdkProvider_impersonate_service_account_emptyStringValidation,

		// Usage
		"impersonate_service_account controls which service account is used for actions": testAccSdkProvider_impersonate_service_account_usage,
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

func testAccSdkProvider_impersonate_service_account_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	impersonateServiceAccountEnvironment := "value-from-envs@example.com"
	impersonateServiceAccountProviderBlock := "value-from-provider-block@example.com"

	// ensure all possible impersonate_service_account env vars set; show they aren't used
	t.Setenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT", impersonateServiceAccountEnvironment)

	context := map[string]interface{}{
		"impersonate_service_account": impersonateServiceAccountProviderBlock,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_impersonate_service_account_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "impersonate_service_account", impersonateServiceAccountProviderBlock),
				)},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for impersonate_service_account, and they are in order of precedence.
		GOOGLE_IMPERSONATE_SERVICE_ACCOUNT
	*/

	impersonateServiceAccount := "foobar@example.com"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT", impersonateServiceAccount)
				},
				Config: testAccSdkProvider_impersonate_service_account_inEnvsOnly(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "impersonate_service_account", impersonateServiceAccount),
				)},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	impersonateServiceAccountEnvironment := "value-from-envs@example.com"

	// ensure all possible impersonate_service_account env vars set; show they aren't used
	t.Setenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT", impersonateServiceAccountEnvironment)

	context := map[string]interface{}{
		"impersonate_service_account": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_impersonate_service_account_inProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

func testAccSdkProvider_impersonate_service_account_usage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	// ensure env vars unset
	t.Setenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT", "")

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_impersonate_service_account_testViaFailure(context),
				ExpectError: regexp.MustCompile("Error creating Topic: googleapi: Error 403: User not authorized"),
			},
		},
	})
}

// testAccSdkProvider_impersonate_service_account_inProviderBlock allows setting the impersonate_service_account argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_impersonate_service_account_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	impersonate_service_account = "%{impersonate_service_account}"
}

data "google_provider_config_sdk" "default" {}

`, context)
}

// testAccSdkProvider_impersonate_service_account_inEnvsOnly allows testing when the impersonate_service_account argument
// is only supplied via ENVs
func testAccSdkProvider_impersonate_service_account_inEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}

`, context)
}

func testAccSdkProvider_impersonate_service_account_testViaFailure(context map[string]interface{}) string {
	return acctest.Nprintf(`
// This will succeed due to the Terraform identity having necessary permissions
resource "google_pubsub_topic" "ok" {
  name = "tf-test-%{random_suffix}-ok"
}

//  Create a service account and ensure the Terraform identity can make tokens for it
resource "google_service_account" "default" {
  account_id   = "tf-test-%{random_suffix}"
  display_name = "Acceptance test impersonated service account"
}

data "google_client_openid_userinfo" "me" {
}

resource "google_service_account_iam_member" "token" {
  service_account_id = google_service_account.default.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${data.google_client_openid_userinfo.me.email}"
}

// Impersonate the created service account
provider "google" {
  alias = "impersonation"
  impersonate_service_account = google_service_account.default.email
}

// This will fail due to the impersonated service account not having any permissions
resource "google_pubsub_topic" "fail" {
  provider = google.impersonation
  name = "tf-test-%{random_suffix}-fail"
}
`, context)
}
