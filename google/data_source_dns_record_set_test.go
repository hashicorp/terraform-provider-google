package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcccDataSourceDnsRecordSet_basic(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsRecordSet_basic(RandString(t, 10), RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_dns_record_set.rs", "google_dns_record_set.rs"),
				),
			},
		},
	})
}

func testAccDataSourceDnsRecordSet_basic(zoneName, recordSetName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "zone" {
  name     = "test-zone"
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
`, zoneName, recordSetName)
}
