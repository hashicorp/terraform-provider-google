// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsManagedZone_basic(acctest.RandString(t, 10)),
				Check: acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
					"data.google_dns_managed_zone.qa",
					"google_dns_managed_zone.foo",
					map[string]struct{}{
						"dnssec_config.#":                       {},
						"private_visibility_config.#":           {},
						"peering_config.#":                      {},
						"forwarding_config.#":                   {},
						"force_destroy":                         {},
						"labels.#":                              {},
						"terraform_labels.%":                    {},
						"effective_labels.%":                    {},
						"creation_time":                         {},
						"cloud_logging_config.#":                {},
						"cloud_logging_config.0.%":              {},
						"cloud_logging_config.0.enable_logging": {},
					},
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZone_basic(managedZoneName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
  name        = "tf-test-qa-zone-%s"
  dns_name    = "qa.gcp.tfacc.hashicorptest.com."
  description = "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.foo.name
}
`, managedZoneName)
}
