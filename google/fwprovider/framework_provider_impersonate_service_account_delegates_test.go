// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwprovider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// TestAccFwProvider_impersonate_service_account_delegates is a series of acc tests asserting how the plugin-framework provider handles impersonate_service_account_delegates arguments
// It is plugin-framework specific because the HCL used provisions plugin-framework-implemented resources
// It is a counterpart to TestAccSdkProvider_impersonate_service_account_delegates
func TestAccFwProvider_impersonate_service_account_delegates(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		// Configuring the provider using inputs
		//     There are no environment variables for this field
		"impersonate_service_account_delegates can be set in config": testAccFwProvider_impersonate_service_account_delegates_setInConfig,

		// Schema-level validation
		"when impersonate_service_account_delegates is set to an empty list in the config the value IS ignored": testAccFwProvider_impersonate_service_account_delegates_emptyListUsage,

		// Usage
		// We need to wait for a non-Firebase resource to be migrated to the plugin-framework to enable writing this test
		// "impersonate_service_account_delegates controls which service account is used for actions"
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

func testAccFwProvider_impersonate_service_account_delegates_setInConfig(t *testing.T) {
	acctest.SkipIfVcr(t) // Test doesn't interact with API

	delegates := []string{
		"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
		"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
	}
	delegatesString := fmt.Sprintf(`["%s","%s"]`, delegates[0], delegates[1])

	// There are no ENVs for this provider argument

	context := map[string]interface{}{
		"impersonate_service_account_delegates": delegatesString,
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_impersonate_service_account_delegatesInProviderBlock(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "impersonate_service_account_delegates.#", fmt.Sprintf("%d", len(delegates))),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "impersonate_service_account_delegates.0", delegates[0]),
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "impersonate_service_account_delegates.1", delegates[1]),
				),
			},
		},
	})
}

func testAccFwProvider_impersonate_service_account_delegates_emptyListUsage(t *testing.T) {

	context := map[string]interface{}{
		"impersonate_service_account_delegates": "[]", // empty array
		"random_suffix":                         acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		// No PreCheck for checking ENVs
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFwProvider_impersonate_service_account_delegates_testProvisioning(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_provider_config_plugin_framework.default", "impersonate_service_account_delegates.#", "0"),
				),
				// No error expected as empty array is ignored
			},
		},
	})
}

// testAccFwProvider_impersonate_service_account_delegatesInProviderBlock allows setting the impersonate_service_account_delegates argument in a provider block.
func testAccFwProvider_impersonate_service_account_delegatesInProviderBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	impersonate_service_account_delegates = %{impersonate_service_account_delegates}
}

data "google_provider_config_plugin_framework" "default" {}

output "impersonate_service_account_delegates" {
  value = data.google_provider_config_plugin_framework.default.impersonate_service_account_delegates
  sensitive = true
}
`, context)
}

// testAccFwProvider_impersonate_service_account_delegates_testProvisioning allows setting the impersonate_service_account_delegates argument in a provider block
// and testing its impact on provisioning a resource
func testAccFwProvider_impersonate_service_account_delegates_testProvisioning(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
	impersonate_service_account_delegates = %{impersonate_service_account_delegates}
}

data "google_provider_config_plugin_framework" "default" {}

resource "google_pubsub_topic" "example" {
  name = "tf-test-%{random_suffix}"
}
`, context)
}
