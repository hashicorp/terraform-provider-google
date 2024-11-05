// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccFwProvider_request_timeout is a series of acc tests asserting how the PF provider handles request_timeout arguments
// It is PF specific because the HCL used provisions PF-implemented resources
// It is a counterpart to TestAccSdkProvider_request_timeout
func TestAccFwProvider_request_timeout(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"a default value of 120s is used when there are no user inputs":                                       testAccFwProvider_request_timeout_providerDefault,
		"request_timeout can be set in config in different formats, are NOT normalized to full-length format": testAccFwProvider_request_timeout_setInConfig,
		//no ENVs to test

		// Schema-level validation
		"when request_timeout is set to an empty string in the config the value fails validation, as it is not a duration": testAccFwProvider_request_timeout_emptyStringValidation,

		// Usage
		// We cannot test the impact of this field in an acc test until more resources/data sources are implemented with the plugin-framework
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

func testAccFwProvider_request_timeout_providerDefault(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	defaultValue := "120s"

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_request_timeout_unset(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "request_timeout", defaultValue),
				),
			},
		},
	})
}

func testAccFwProvider_request_timeout_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	providerTimeout1 := "3m0s"
	providerTimeout2 := "3m"

	// In the SDK version of the test expectedValue = "3m0s" only
	expectedValue1 := "3m0s"
	expectedValue2 := "3m"

	context1 := map[string]interface{}{
		"request_timeout": providerTimeout1,
	}
	context2 := map[string]interface{}{
		"request_timeout": providerTimeout2,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_request_timeout_inProviderBlock(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "request_timeout", expectedValue1),
				),
			},
			{
				Config: testAccFwProvider_request_timeout_inProviderBlock(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "request_timeout", expectedValue2),
				),
			},
		},
	})
}

func testAccFwProvider_request_timeout_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"request_timeout": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccFwProvider_request_timeout_inProviderBlock(context),
				ExpectError: regexp.MustCompile("invalid duration"),
			},
		},
	})
}

// testAccFwProvider_request_timeout_inProviderBlock allows setting the request_timeout argument in a provider block.
// This function uses data.google_provider_config_plugin_framework because it is implemented with the PF
func testAccFwProvider_request_timeout_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	request_timeout = "%{request_timeout}"
}

data "google_provider_config_plugin_framework" "default" {}
`, context)
}

// testAccFwProvider_request_timeout_inEnvsOnly allows testing when the request_timeout argument is not set
func testAccFwProvider_request_timeout_unset() string {
	return `
data "google_provider_config_plugin_framework" "default" {}
`
}
