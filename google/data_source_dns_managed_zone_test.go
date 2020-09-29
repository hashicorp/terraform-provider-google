package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsManagedZone_basic(randString(t, 10)),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.google_dns_managed_zone.qa",
					"google_dns_managed_zone.foo",
					map[string]struct{}{
						"dnssec_config.#":             {},
						"private_visibility_config.#": {},
						"peering_config.#":            {},
						"forwarding_config.#":         {},
					},
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZone_basic(managedZoneName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
  name        = "qa-zone-%s"
  dns_name    = "qa.tf-test.club."
  description = "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.foo.name
}
`, managedZoneName)
}
