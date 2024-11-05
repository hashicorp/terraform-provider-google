// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/transport"
)

// TestAccSdkProvider_scopes is a series of acc tests asserting how the SDK provider handles scopes arguments
// It is SDK specific because the HCL used provisions SDK-implemented resources
// It is a counterpart to TestAccFwProvider_scopes
func TestAccSdkProvider_scopes(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"default scopes are used when there are no user inputs": testAccSdkProvider_scopes_providerDefault,
		"scopes can be set in config":                           testAccSdkProvider_scopes_setInConfig,
		//no ENVs to test

		// Schema-level validation
		"when scopes is set to an empty array in the config the value is ignored and default scopes are used": testAccSdkProvider_scopes_emptyArray,

		// Usage
		"the scopes argument impacts provisioning resources": testAccSdkProvider_scopes_usage,
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

func testAccSdkProvider_scopes_providerDefault(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_scopes_unset(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.#", fmt.Sprintf("%d", len(transport.DefaultClientScopes))),
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.0", transport.DefaultClientScopes[0]),
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.1", transport.DefaultClientScopes[1]),
				),
			},
		},
	})
}

func testAccSdkProvider_scopes_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	scopes := []string{"https://www.googleapis.com/auth/cloud-platform"} // first of the two default scopes
	context := map[string]interface{}{
		"scopes": fmt.Sprintf("[\"%s\"]", scopes[0]),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_scopes_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.#", fmt.Sprintf("%d", len(scopes))),
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.0", scopes[0]),
				),
			},
		},
	})
}

func testAccSdkProvider_scopes_emptyArray(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"scopes": "[]",
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSdkProvider_scopes_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.#", fmt.Sprintf("%d", len(transport.DefaultClientScopes))),
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.0", transport.DefaultClientScopes[0]),
					resource.TestCheckResourceAttr("data.google_provider_config_sdk.default", "scopes.1", transport.DefaultClientScopes[1]),
				),
			},
		},
	})
}

func testAccSdkProvider_scopes_usage(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	// We include scopes that aren't sufficient to enable provisioning the resources in the config below
	context := map[string]interface{}{
		"scopes":        "[\"https://www.googleapis.com/auth/pubsub\"]",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSdkProvider_scopes_affectsProvisioning(context),
				ExpectError: regexp.MustCompile("Request had insufficient authentication scopes"),
			},
		},
	})
}

// testAccSdkProvider_scopes_inProviderBlock allows setting the scopes argument in a provider block.
// This function uses data.google_provider_config_sdk because it is implemented with the SDK
func testAccSdkProvider_scopes_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	scopes = %{scopes}
}

data "google_provider_config_sdk" "default" {}
`, context)
}

// testAccSdkProvider_scopes_inEnvsOnly allows testing when the scopes argument is not set
func testAccSdkProvider_scopes_unset() string {
	return `
data "google_provider_config_sdk" "default" {}
`
}

// testAccSdkProvider_scopes_affectsProvisioning allows testing the impact of the scopes argument on provisioning
func testAccSdkProvider_scopes_affectsProvisioning(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	scopes = %{scopes}
}

data "google_provider_config_sdk" "default" {}

resource "google_service_account" "default" {
  account_id   = "tf-test-%{random_suffix}"
  display_name = "AccTest Service Account"
}
`, context)
}
