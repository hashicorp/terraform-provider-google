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

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccDataSourceDnsManagedZone_basic(RandString(t, 10)),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.google_dns_managed_zone.qa",
					"google_dns_managed_zone.foo",
					map[string]struct{}{
						"dnssec_config.#":                       {},
						"private_visibility_config.#":           {},
						"peering_config.#":                      {},
						"forwarding_config.#":                   {},
						"force_destroy":                         {},
						"labels.#":                              {},
						"creation_time":                         {},
						"cloud_logging_config.#":                {},
						"cloud_logging_config.0.%":              {},
						"cloud_logging_config.0.enable_logging": {},
					},
				),
			},
			{
				ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
				Config:                   testAccDataSourceDnsManagedZone_basic(RandString(t, 10)),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.google_dns_managed_zone.qa",
					"google_dns_managed_zone.foo",
					map[string]struct{}{
						"dnssec_config.#":                       {},
						"private_visibility_config.#":           {},
						"peering_config.#":                      {},
						"forwarding_config.#":                   {},
						"force_destroy":                         {},
						"labels.#":                              {},
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
  name        = "tf-test-zone-%s"
  dns_name    = "tf-test-zone-%s.hashicorptest.com."
  description = "tf test DNS zone"
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.foo.name
}
`, managedZoneName, managedZoneName)
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

			p := GetFwTestProvider(t)

			url, err := replaceVarsForFrameworkTest(&p.frameworkProvider, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if !p.billingProject.IsNull() && p.billingProject.String() != "" {
				billingProject = p.billingProject.String()
			}

			_, diags := sendFrameworkRequest(&p.frameworkProvider, "GET", billingProject, url, p.userAgent, nil)
			if !diags.HasError() {
				return fmt.Errorf("DNSManagedZone still exists at %s", url)
			}
		}

		return nil
	}
}
