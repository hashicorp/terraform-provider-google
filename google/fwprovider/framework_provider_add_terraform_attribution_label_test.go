// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccFwProvider_add_terraform_attribution_label is a series of acc tests asserting how the plugin-framework provider handles add_terraform_attribution_label arguments
// It is plugin-framework specific because the HCL used provisions plugin-framework-implemented resources
// It is a counterpart to TestAccSdkProvider_add_terraform_attribution_label
func TestAccFwProvider_add_terraform_attribution_label(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		"when add_terraform_attribution_label is set in the config, the value is set in the provider meta data":                      testAccFwProvider_add_terraform_attribution_label_configUsed,
		"when add_terraform_attribution_label is unset in the config, the default value 'true' is NOT set on the provider meta data": testAccFwProvider_add_terraform_attribution_label_defaultValue,
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

func testAccFwProvider_add_terraform_attribution_label_configUsed(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	contextFalse := map[string]interface{}{
		"add_terraform_attribution_label": "false",
	}
	contextTrue := map[string]interface{}{
		"add_terraform_attribution_label": "true",
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_add_terraform_attribution_label_inProviderBlock(contextFalse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "add_terraform_attribution_label", "false"),
				),
			},
			{
				Config: testAccFwProvider_add_terraform_attribution_label_inProviderBlock(contextTrue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "add_terraform_attribution_label", "true"),
				),
			},
		},
	})
}

func testAccFwProvider_add_terraform_attribution_label_defaultValue(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	context := map[string]interface{}{}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_add_terraform_attribution_label_inEnvsOnly(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.google_provider_config_plugin_framework.default", "add_terraform_attribution_label"),
				),
			},
		},
	})
}

// testAccFwProvider_add_terraform_attribution_label_inProviderBlock allows setting the add_terraform_attribution_label argument in a provider block.
func testAccFwProvider_add_terraform_attribution_label_inProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	add_terraform_attribution_label = "%{add_terraform_attribution_label}"
}

data "google_provider_config_plugin_framework" "default" {}
`, context)
}

// testAccFwProvider_add_terraform_attribution_label_inEnvsOnly allows testing when the add_terraform_attribution_label argument
// is only supplied via ENVs
func testAccFwProvider_add_terraform_attribution_label_inEnvsOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_provider_config_plugin_framework" "default" {}
`, context)
}
