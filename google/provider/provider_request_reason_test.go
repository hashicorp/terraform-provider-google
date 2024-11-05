// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccSdkProvider_request_reason is a series of acc tests asserting how the SDK provider handles request_reason arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_request_reason
func TestAccSdkProvider_request_reason(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"config takes precedence over environment variables":                                          testAccSdkProvider_request_reason_configPrecedenceOverEnvironmentVariables,
		"when request_reason is unset in the config, environment variables are used in a given order": testAccSdkProvider_request_reason_precedenceOrderEnvironmentVariables, // CLOUDSDK_CORE_REQUEST_REASON

		// Schema-level validation
		// TODO: https://github.com/hashicorp/terraform-provider-google/issues/19643
		"when request_reason is set to an empty string in the config the value IS ignored, allowing environment values to be used": testAccSdkProvider_request_reason_emptyStringValidation,

		// Usage
		// We cannot test the impact of this field in an acc test, as it sets the X-Goog-Request-Reason value for audit logging purposes in GCP
		// See: https://cloud.google.com/apis/docs/system-parameters#definitions
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

func testAccSdkProvider_request_reason_configPrecedenceOverEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	envReason := "environment-variables"

	// ensure all possible request_reason env vars set; show they aren't used instead
	t.Setenv("CLOUDSDK_CORE_REQUEST_REASON", envReason)

	providerReason := "provider-config"

	context := map[string]interface{}{
		"request_reason": providerReason,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_request_reason_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_reason", providerReason),
				),
			},
		},
	})
}

func testAccSdkProvider_request_reason_precedenceOrderEnvironmentVariables(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API
	/*
		These are all the ENVs for request_reason
		CLOUDSDK_CORE_REQUEST_REASON
	*/

	CLOUDSDK_CORE_REQUEST_REASON := "CLOUDSDK_CORE_REQUEST_REASON"

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// CLOUDSDK_CORE_REQUEST_REASON is used if config doesn't provide a value
				PreConfig: func() {
					t.Setenv("CLOUDSDK_CORE_REQUEST_REASON", CLOUDSDK_CORE_REQUEST_REASON) //used
				},
				Config: testAccSdkProvider_request_reason_inEnvsOnly(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_reason", CLOUDSDK_CORE_REQUEST_REASON),
				),
			},
		},
	})
}

func testAccSdkProvider_request_reason_emptyStringValidation(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	envReason := "environment-variables"

	// ensure all request_reason env vars set
	t.Setenv("CLOUDSDK_CORE_REQUEST_REASON", envReason)

	context := map[string]interface{}{
		"request_reason": "", // empty string used
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_request_reason_inProviderBlock(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Currently the SDK provider silently ignores empty strings, treating them similar to being unset
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "request_reason", envReason),
				),
			},
		},
	})
}

// testAccSdkProvider_request_reason_inProviderBlock allows setting the request_reason argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDKv2
func testAccSdkProvider_request_reason_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	request_reason = "%{request_reason}"
}

data "google_provider_config_sdk" "default" {}
`, context)
}

// testAccSdkProvider_request_reason_inEnvsOnly allows testing when the request_reason argument
// is only supplied via ENVs
func testAccSdkProvider_request_reason_inEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_sdk" "default" {}
`, context)
}
