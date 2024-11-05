// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/transport"
)

// TestAccFwProvider_scopes is a series of acc tests asserting how the PF provider handles scopes arguments
// It is PF specific because the HCL used provisions PF-implemented resources
// It is a counterpart to TestAccSdkProvider_scopes
func TestAccFwProvider_scopes(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"default scopes are used when there are no user inputs": testAccFwProvider_scopes_providerDefault,
		"scopes can be set in config":                           testAccFwProvider_scopes_setInConfig,
		//no ENVs to test

		// Schema-level validation
		"when scopes is set to an empty array in the config the value is ignored and default scopes are used": testAccFwProvider_scopes_emptyArray,

		// Usage
		// No usage test cases are implemented in the GA provider because the only PF-implemented data sources are Beta-only
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

func testAccFwProvider_scopes_providerDefault(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_scopes_unset(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.#", fmt.Sprintf("%d", len(transport.DefaultClientScopes))),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.0", transport.DefaultClientScopes[0]),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.1", transport.DefaultClientScopes[1]),
				),
			},
		},
	})
}

func testAccFwProvider_scopes_setInConfig(t *testing.T) {
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
				Config: testAccFwProvider_scopes_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.#", fmt.Sprintf("%d", len(scopes))),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.0", scopes[0]),
				),
			},
		},
	})
}

func testAccFwProvider_scopes_emptyArray(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{
		"scopes": "[]",
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_scopes_inProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.#", fmt.Sprintf("%d", len(transport.DefaultClientScopes))),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.0", transport.DefaultClientScopes[0]),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "scopes.1", transport.DefaultClientScopes[1]),
				),
			},
		},
	})
}

// testAccFwProvider_scopes_inProviderBlock allows setting the scopes argument in a provider block.
// This function uses data.google_provider_config_plugin_framework because it is implemented with the PF
func testAccFwProvider_scopes_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	scopes = %{scopes}
}

data "google_provider_config_plugin_framework" "default" {}
`, context)
}

// testAccFwProvider_scopes_inEnvsOnly allows testing when the scopes argument is not set
func testAccFwProvider_scopes_unset() string {
	return `
data "google_provider_config_plugin_framework" "default" {}
`
}
