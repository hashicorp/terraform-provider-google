package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsManagedZone_basic(),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.google_dns_managed_zone.qa",
					"google_dns_managed_zone.foo",
					map[string]struct{}{
						"dnssec_config.#":             {},
						"private_visibility_config.#": {},
					},
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZone_basic() string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
  name        = "qa-zone-%s"
  dns_name    = "qa.tf-test.club."
  description = "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.foo.name
}
`, acctest.RandString(10))
}
