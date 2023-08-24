// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceDNSKeys_basic(t *testing.T) {
	// TODO: https://github.com/hashicorp/terraform-provider-google/issues/14158
	acctest.SkipIfVcr(t)
	t.Parallel()

	dnsZoneName := fmt.Sprintf("tf-test-dnskey-test-%s", acctest.RandString(t, 10))

	var kskDigest1, kskDigest2, zskPubKey1, zskPubKey2, kskAlg1, kskAlg2 string

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccDataSourceDNSKeysConfigWithOutputs(dnsZoneName, "on"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDNSKeysDSRecordCheck("data.google_dns_keys.foo_dns_key"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "key_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key_id", "zone_signing_keys.#", "1"),
					acctest.TestExtractResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.0.digests.0.digest", &kskDigest1),
					acctest.TestExtractResourceAttr("data.google_dns_keys.foo_dns_key_id", "zone_signing_keys.0.public_key", &zskPubKey1),
					acctest.TestExtractResourceAttr("data.google_dns_keys.foo_dns_key_id", "key_signing_keys.0.algorithm", &kskAlg1),
				),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccDataSourceDNSKeysConfigWithOutputs(dnsZoneName, "on"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDNSKeysDSRecordCheck("data.google_dns_keys.foo_dns_key"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "1"),
					acctest.TestExtractResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.0.digests.0.digest", &kskDigest2),
					acctest.TestExtractResourceAttr("data.google_dns_keys.foo_dns_key_id", "zone_signing_keys.0.public_key", &zskPubKey2),
					acctest.TestExtractResourceAttr("data.google_dns_keys.foo_dns_key_id", "key_signing_keys.0.algorithm", &kskAlg2),
					acctest.TestCheckAttributeValuesEqual(&kskDigest1, &kskDigest2),
					acctest.TestCheckAttributeValuesEqual(&zskPubKey1, &zskPubKey2),
					acctest.TestCheckAttributeValuesEqual(&kskAlg1, &kskAlg2),
				),
			},
		},
	})
}

func TestAccDataSourceDNSKeys_noDnsSec(t *testing.T) {
	// TODO: https://github.com/hashicorp/terraform-provider-google/issues/14158
	acctest.SkipIfVcr(t)
	t.Parallel()

	dnsZoneName := fmt.Sprintf("tf-test-dnskey-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccDataSourceDNSKeysConfig(dnsZoneName, "off"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "key_signing_keys.#", "0"),
					resource.TestCheckResourceAttr("data.google_dns_keys.foo_dns_key", "zone_signing_keys.#", "0"),
				),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccDataSourceDNSKeysConfig(dnsZoneName, "off"),
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
  dns_name = "%s.hashicorptest.com."

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
`, dnsZoneName, dnsZoneName, dnssecStatus)
}

// This function extends the config returned from the `testAccDataSourceDNSKeysConfig` function
// to include output blocks that access the `key_signing_keys` and `zone_signing_keys` attributes.
// These are null if DNSSEC is not enabled.
func testAccDataSourceDNSKeysConfigWithOutputs(dnsZoneName, dnssecStatus string) string {

	config := testAccDataSourceDNSKeysConfig(dnsZoneName, dnssecStatus)
	config = config + `
# These outputs will cause an error if google_dns_managed_zone.foo.dnssec_config.state == "off"

output "test_access_google_dns_keys_key_signing_keys" {
  description = "Testing that we can access a value in key_signing_keys ok as a computed block"
  value       = data.google_dns_keys.foo_dns_key_id.key_signing_keys[0].ds_record
}

output "test_access_google_dns_keys_zone_signing_keys" {
  description = "Testing that we can access a value in zone_signing_keys ok as a computed block"
  value       = data.google_dns_keys.foo_dns_key_id.zone_signing_keys[0].id
}
`
	return config
}
