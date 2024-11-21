// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccFwProvider_billing_project is a series of acc tests asserting how the PF provider handles billing_project arguments
// It is PF specific because the HCL used provisions PF-implemented resources
// It is a counterpart to TestAccSdkProvider_billing_project
func TestAccFwProvider_billing_project(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                           testAccFwProvider_billing_project_configPrecedenceOverEnvironmentVariables,
		"when billing_project is unset in the config, environment variables are used in a given order": testAccFwProvider_billing_project_precedenceOrderEnvironmentVariables, // GOOGLE_BILLING_PROJECT

		// Schema-level validation
		"when billing_project is set to an empty string in the config the value isn't ignored and results in an error": testAccFwProvider_billing_project_emptyStringValidation,

		// Usage
		// Usage test cases are implemented in a Beta-only way
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

func testAccFwProvider_billing_project_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	billingProject := "my-billing-project-id"

	// ensure all possible billing_project env vars set; show they aren't used instead
	t.Setenv("GOOGLE_BILLING_PROJECT", billingProject)

	providerBillingProject := "foobar"

	context := map[string]interface{}{
		"billing_project": providerBillingProject,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_billing_project_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "billing_project", providerBillingProject),
				)},
		},
	})
}

func testAccFwProvider_billing_project_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for billing_project
		GOOGLE_BILLING_PROJECT

		GOOGLE_CLOUD_QUOTA_PROJECT - NOT used by provider, but is in client libraries we use
	*/

	GOOGLE_BILLING_PROJECT := "GOOGLE_BILLING_PROJECT"
	GOOGLE_CLOUD_QUOTA_PROJECT := "GOOGLE_CLOUD_QUOTA_PROJECT"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// GOOGLE_BILLING_PROJECT is used if set
				PreConfig: func() {
					t.Setenv("GOOGLE_BILLING_PROJECT", GOOGLE_BILLING_PROJECT) //used
					t.Setenv("GOOGLE_CLOUD_QUOTA_PROJECT", GOOGLE_CLOUD_QUOTA_PROJECT)
				},
				Config: testAccFwProvider_billing_project_inEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "billing_project", GOOGLE_BILLING_PROJECT),
				),
			},
			{
				// GOOGLE_CLOUD_QUOTA_PROJECT is NOT used here
				PreConfig: func() {
					t.Setenv("GOOGLE_BILLING_PROJECT", "")
					t.Setenv("GOOGLE_CLOUD_QUOTA_PROJECT", GOOGLE_CLOUD_QUOTA_PROJECT) // NOT used
				},
				Config: testAccFwProvider_billing_project_inEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "billing_project", ""),
				),
			},
		},
	})
}

func testAccFwProvider_billing_project_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	billingProject := "my-billing-project-id"

	// ensure all billing_project env vars set
	t.Setenv("GOOGLE_BILLING_PROJECT", billingProject)

	context := map[string]interface{}{
		"billing_project": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFwProvider_billing_project_inProviderBlock(context),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected a non-empty string"),
			},
		},
	})
}

// testAccFwProvider_billing_project_inProviderBlock allows setting the billing_project argument in a provider block.
// This function uses data.google_provider_config_plugin_framework because it is implemented with the plugin-framework
func testAccFwProvider_billing_project_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	billing_project = "%{billing_project}"
}

data "google_provider_config_plugin_framework" "default" {}
`, context)
}

// testAccFwProvider_billing_project_inEnvsOnly allows testing when the billing_project argument
// is only supplied via ENVs
func testAccFwProvider_billing_project_inEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_plugin_framework" "default" {}
`, context)
}
