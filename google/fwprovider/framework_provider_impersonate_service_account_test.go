// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccFwProvider_impersonate_service_account is a series of acc tests asserting how the plugin-framework provider handles impersonate_service_account arguments
// It is plugin-framework specific because the HCL used provisions plugin-framework-implemented resources
// It is a counterpart to TestAccSdkProvider_impersonate_service_account
func TestAccFwProvider_impersonate_service_account(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                                       testAccFwProvider_impersonate_service_account_configPrecedenceOverEnvironmentVariables,
		"when impersonate_service_account is unset in the config, environment variables are used in a given order": testAccFwProvider_impersonate_service_account_precedenceOrderEnvironmentVariables, // GOOGLE_IMPERSONATE_SERVICE_ACCOUNT

		// Schema-level validation
		"when impersonate_service_account is set to an empty string in the config the value isn't ignored and results in an error": testAccFwProvider_impersonate_service_account_emptyStringValidation,

		// Usage
		// We need to wait for a non-Firebase resource to be migrated to the plugin-framework to enable writing this test
		// "impersonate_service_account controls which service account is used for actions"
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

func testAccFwProvider_impersonate_service_account_configPrecedenceOverEnvironmentVariables(t *testing.T) {
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
				Config: testAccFwProvider_impersonate_service_account_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "impersonate_service_account", impersonateServiceAccountProviderBlock),
				),
			},
		},
	})
}

func testAccFwProvider_impersonate_service_account_precedenceOrderEnvironmentVariables(t *testing.T) {
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
				Config: testAccFwProvider_impersonate_service_account_inEnvsOnly(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "impersonate_service_account", impersonateServiceAccount),
				),
			},
		},
	})
}

func testAccFwProvider_impersonate_service_account_emptyStringValidation(t *testing.T) {
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
				Config:      testAccFwProvider_impersonate_service_account_inProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

// testAccFwProvider_impersonate_service_account_inProviderBlock allows setting the impersonate_service_account argument in a provider block.
func testAccFwProvider_impersonate_service_account_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	impersonate_service_account = "%{impersonate_service_account}"
}

data "google_provider_config_plugin_framework" "default" {}

`, context)
}

// testAccFwProvider_impersonate_service_account_inEnvsOnly allows testing when the impersonate_service_account argument
// is only supplied via ENVs
func testAccFwProvider_impersonate_service_account_inEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_plugin_framework" "default" {}

`, context)
}
