// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceDNSKeys_basic(t *testing.T) {
	t.Parallel()

	dnsZoneName := fmt.Sprintf("tf-test-dnskey-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSKeysConfig(dnsZoneName, "on"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDNSKeysDSRecordCheck("data.google_dns_keys.foo_dns_key"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "key_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "zone_signing_keys.#", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceDNSKeys_noDnsSec(t *testing.T) {
	t.Parallel()

	dnsZoneName := fmt.Sprintf("tf-test-dnskey-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSKeysConfig(dnsZoneName, "off"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "0"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceDNSKeysDSRecordCheck(datasourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		}

		if ds.Primary.Attributes["key_signing_keys.0.ds_record"] == "" {
			return fmt.Errorf("DS record not found in data source")
		}

		return nil
	}
}

func testAccDataSourceDNSKeysConfig(dnsZoneName, dnssecStatus string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
  name     = "%s"
  dns_name = "dnssec.gcp.tfacc.hashicorptest.com."

  dnssec_config {
    state         = "%s"
    non_existence = "nsec3"
  }
}

data "google_dns_keys" "foo_dns_key" {
  managed_zone = google_dns_managed_zone.foo.name
}

data "google_dns_keys" "foo_dns_key_id" {
  managed_zone = google_dns_managed_zone.foo.id
}
`, dnsZoneName, dnssecStatus)
}

// TestAccDataSourceDNSKeys_basic_AdcAuth is the same as TestAccDataSourceDNSKeys_basic but the test enforces that a developer runs this using
// ADCs, supplied via GOOGLE_APPLICATION_CREDENTIALS. If any other credentials ENVs are set the PreCheck will fail.
// Commented out until this test can run in TeamCity/CI.
// func TestAccDataSourceDNSKeys_basic_AdcAuth(t *testing.T) {
// 	acctest.SkipIfVcr(t) // Uses external providers
// 	t.Parallel()

// 	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") // PreCheck assertion handles checking this is set

// 	dnsZoneName := fmt.Sprintf("tf-test-dnskey-test-%s", acctest.RandString(t, 10))

// 	context := map[string]interface{}{
// 		"credentials_path": creds,
// 		"dns_zone_name":    dnsZoneName,
// 		"dnssec_status":    "on",
// 	}

// 	acctest.VcrTest(t, resource.TestCase{
// 		PreCheck:     func() { acctest.AccTestPreCheck_AdcCredentialsOnly(t) }, // Note different than default
// 		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducer(t),
// 		Steps: []resource.TestStep{
// 			// Check test fails with version of provider where data source is implemented with PF
// 			{
// 				ExternalProviders: map[string]resource.ExternalProvider{
// 					"google": {
// 						VersionConstraint: "4.60.0", // Muxed provider with dns data sources migrated to PF
// 						Source:            "hashicorp/google",
// 					},
// 				},
// 				ExpectError: regexp.MustCompile("Post \"https://oauth2.googleapis.com/token\": context canceled"),
// 				Config:      testAccDataSourceDNSKeysConfig_AdcCredentials(context),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccDataSourceDNSKeysDSRecordCheck("data.google_dns_keys.foo_dns_key"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "1"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "1"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "key_signing_keys.#", "1"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "zone_signing_keys.#", "1"),
// 				),
// 			},
// 			// Test should pass with more recent code
// 			{
// 				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
// 				Config:                   testAccDataSourceDNSKeysConfig_AdcCredentials(context),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccDataSourceDNSKeysDSRecordCheck("data.google_dns_keys.foo_dns_key"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "1"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "1"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "key_signing_keys.#", "1"),
// 					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "zone_signing_keys.#", "1"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccDataSourceDNSKeysConfig_AdcCredentials(context map[string]interface{}) string {
	return acctest.Nprintf(`

// The auth problem isn't triggered unless provider block is
// present in the test config.

provider "google" {
  credentials = "%{credentials_path}"
}

resource "google_dns_managed_zone" "foo" {
  name     = "%{dns_zone_name}"
  dns_name = "dnssec.gcp.tfacc.hashicorptest.com."

  dnssec_config {
    state         = "%{dnssec_status}"
    non_existence = "nsec3"
  }
}

data "google_dns_keys" "foo_dns_key" {
  managed_zone = google_dns_managed_zone.foo.name
}

data "google_dns_keys" "foo_dns_key_id" {
  managed_zone = google_dns_managed_zone.foo.id
}
`, context)
}
