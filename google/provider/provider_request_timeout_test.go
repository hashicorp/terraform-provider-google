// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccSdkProvider_request_timeout is a series of acc tests asserting how the SDK provider handles request_timeout arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_request_timeout
func TestAccSdkProvider_request_timeout(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"a default value of 0s is used when there are no user inputs (it is overridden downstream)":       testAccSdkProvider_request_timeout_providerDefault,
		"request_timeout can be set in config in different formats, are normalized to full-length format": testAccSdkProvider_request_timeout_setInConfig,
		//no ENVs to test

		// Schema-level validation
		"when request_timeout is set to an empty string in the config the value fails validation, as it is not a duration": testAccSdkProvider_request_timeout_emptyStringValidation,

		// Usage
		"short timeouts impact provisioning resources": testAccSdkProvider_request_timeout_usage,
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

// In the SDK version of the provider config code request_timeout has a zero value of "0s" and no default.
// The final 'effective' value is "120s", matching the default used in the plugin-framework version of the provider config code.
// See : https://github.com/hashicorp/terraform-provider-google/blob/09cb850ee64bcd78e4457df70905530c1ed75f19/google/transport/config.go#L1228-L1233
func testAccSdkProvider_request_timeout_providerDefault(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_request_timeout_unset(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_timeout", "0s"),
				),
			},
		},
	})
}

func testAccSdkProvider_request_timeout_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	providerTimeout1 := "3m0s"
	providerTimeout2 := "3m"

	// All inputs are normalised to this
	expectedValue := "3m0s"

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
				Config: testAccSdkProvider_request_timeout_inProviderBlock(context1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_timeout", expectedValue),
				),
			},
			{
				Config: testAccSdkProvider_request_timeout_inProviderBlock(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_timeout", expectedValue),
				),
			},
		},
	})
}

func testAccSdkProvider_request_timeout_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"request_timeout": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_request_timeout_inProviderBlock(context),
				ExpectError: regexp.MustCompile("invalid duration"),
			},
		},
	})
}

func testAccSdkProvider_request_timeout_usage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	shortTimeout := "10ms" // short time that will result in an error
	longTimeout := "120s"

	randomString := acctest.RandString(t, 10)
	context1 := map[string]interface{}{
		"request_timeout": shortTimeout,
		"random_suffix":   randomString,
	}
	context2 := map[string]interface{}{
		"request_timeout": longTimeout,
		"random_suffix":   randomString,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_request_timeout_provisionWithTimeout(context1),
				// Effect of request_timeout value
				ExpectError: regexp.MustCompile("context deadline exceeded"),
			},
			{
				Config: testAccSdkProvider_request_timeout_provisionWithTimeout(context2),
				// No error; everything is fine with an appropriate timeout value
			},
		},
	})
}

// testAccSdkProvider_request_timeout_inProviderBlock allows setting the request_timeout argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_request_timeout_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	request_timeout = "%{request_timeout}"
}

data "google_provider_config_sdk" "default" {}
`, context)
}

// testAccSdkProvider_request_timeout_provisionWithTimeout allows testing the effects of request_timeout on
// provisioning a resource.
func testAccSdkProvider_request_timeout_provisionWithTimeout(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	request_timeout = "%{request_timeout}"
}

data "google_provider_config_sdk" "default" {}

resource "google_service_account" "default" {
  account_id   = "tf-test-%{random_suffix}"
  display_name = "AccTest Service Account"
}
`, context)
}

// testAccSdkProvider_request_timeout_inEnvsOnly allows testing when the request_timeout argument is not set
func testAccSdkProvider_request_timeout_unset() string {
	return `
data "google_provider_config_sdk" "default" {}
`
}
