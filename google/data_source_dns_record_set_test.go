package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDnsRecordSet_basic(t *testing.T) {
	t.Parallel()

	var ttl1, ttl2 string // ttl is a computed string-type attribute that is easy to compare in the test

	managedZoneName := fmt.Sprintf("tf-test-zone-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion450(),
				Config:            testAccDataSourceDnsRecordSet_basic(managedZoneName, randString(t, 10), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_dns_record_set.rs", "google_dns_record_set.rs"),
					testExtractResourceAttr("data.google_dns_record_set.rs", "ttl", &ttl1),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(t),
				Config:                   testAccDataSourceDnsRecordSet_basic(managedZoneName, randString(t, 10), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_dns_record_set.rs", "google_dns_record_set.rs"),
					testExtractResourceAttr("data.google_dns_record_set.rs", "ttl", &ttl2),
					testCheckAttributeValuesEqual(&ttl1, &ttl2),
				),
			},
		},
	})
}

func testAccDataSourceDnsRecordSet_basic(managedZoneName, zoneName, recordSetName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "zone" {
  name     = "%s"
  dns_name = "%s.hashicorptest.com."
}

resource "google_dns_record_set" "rs" {
  managed_zone = google_dns_managed_zone.zone.name
  name         = "%s.${google_dns_managed_zone.zone.dns_name}"
  type         = "A"
  ttl          = 300
  rrdatas      = [
	"192.168.1.0",
  ]
}

data "google_dns_record_set" "rs" {
  managed_zone = google_dns_record_set.rs.managed_zone
  name         = google_dns_record_set.rs.name
  type         = google_dns_record_set.rs.type
}
`, managedZoneName, zoneName, recordSetName)
}

// testAccCheckDnsRecordSetDestroyProducerFramework is the framework version of the generated testAccCheckDnsRecordSetDestroyProducer
func testAccCheckDnsRecordSetDestroyProducerFramework(t *testing.T) func(s *terraform.State) error {

	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dns_record_set" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			p := getTestFwProvider(t)

			url, err := replaceVarsForFrameworkTest(&p.ProdProvider, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if !p.ProdProvider.billingProject.IsNull() && p.ProdProvider.billingProject.String() != "" {
				billingProject = p.ProdProvider.billingProject.String()
			}

			_, diags := sendFrameworkRequest(&p.ProdProvider, "GET", billingProject, url, p.ProdProvider.userAgent, nil)
			if !diags.HasError() {
				return fmt.Errorf("DNSResourceDnsRecordSet still exists at %s", url)
			}
		}

		return nil
	}
}
