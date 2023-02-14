package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
						"force_destroy":               {},
						"labels.#":                    {},
						"creation_time":               {},
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

// testAccCheckDNSManagedZoneDestroyProducerFramework is the framework version of the generated testAccCheckDNSManagedZoneDestroyProducer
// when we automate this, we'll use the automated version and can get rid of this
func testAccCheckDNSManagedZoneDestroyProducerFramework(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dns_managed_zone" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			p := getTestFwProvider(t)

			url, err := replaceVarsForFrameworkTest(&p.ProdProvider, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if !p.ProdProvider.billingProject.IsNull() && p.ProdProvider.billingProject.String() != "" {
				billingProject = p.ProdProvider.billingProject.String()
			}

			_, diags := sendFrameworkRequest(&p.ProdProvider, "GET", billingProject, url, p.ProdProvider.userAgent, nil)
			if !diags.HasError() {
				return fmt.Errorf("DNSManagedZone still exists at %s", url)
			}
		}

		return nil
	}
}
