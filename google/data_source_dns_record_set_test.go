package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceDnsRecordSet_basic(t *testing.T) {
	// TODO: https://github.com/hashicorp/terraform-provider-google/issues/14158
	acctest.SkipIfVcr(t)
	t.Parallel()

	var ttl1, ttl2 string // ttl is a computed string-type attribute that is easy to compare in the test

	managedZoneName := fmt.Sprintf("tf-test-zone-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.58.0",
						Source:            "hashicorp/google",
					},
				},
				Config: testAccDataSourceDnsRecordSet_basic(managedZoneName, RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_dns_record_set.rs", "google_dns_record_set.rs"),
					testExtractResourceAttr("data.google_dns_record_set.rs", "ttl", &ttl1),
				),
			},
			{
				ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
				Config:                   testAccDataSourceDnsRecordSet_basic(managedZoneName, RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_dns_record_set.rs", "google_dns_record_set.rs"),
					testExtractResourceAttr("data.google_dns_record_set.rs", "ttl", &ttl2),
					testCheckAttributeValuesEqual(&ttl1, &ttl2),
				),
			},
		},
	})
}

func testAccDataSourceDnsRecordSet_basic(managedZoneName, recordSetName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "zone" {
  name     = "%s-hashicorptest-com"
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
`, managedZoneName, managedZoneName, recordSetName)
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

			p := GetFwTestProvider(t)

			url, err := replaceVarsForFrameworkTest(&p.frameworkProvider, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if !p.billingProject.IsNull() && p.billingProject.String() != "" {
				billingProject = p.billingProject.String()
			}

			_, diags := sendFrameworkRequest(&p.frameworkProvider, "GET", billingProject, url, p.userAgent, nil)
			if !diags.HasError() {
				return fmt.Errorf("DNSResourceDnsRecordSet still exists at %s", url)
			}
		}

		return nil
	}
}
