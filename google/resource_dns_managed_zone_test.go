package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/dns/v1"
)

func TestAccDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	var zone dns.ManagedZone

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsManagedZone_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsManagedZoneExists(
						"google_dns_managed_zone.foobar", &zone),
				),
			},
		},
	})
}

func testAccCheckDnsManagedZoneDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_dns_zone" {
			continue
		}

		_, err := config.clientDns.ManagedZones.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("DNS ManagedZone still exists")
		}
	}

	return nil
}

func testAccCheckDnsManagedZoneExists(n string, zone *dns.ManagedZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientDns.ManagedZones.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("DNS Zone not found")
		}

		*zone = *found

		return nil
	}
}

func testAccDnsManagedZone_basic() string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
	name = "mzone-test-%s"
	dns_name = "tf-acctest-%s.hashicorptest.com."
}`, acctest.RandString(10), acctest.RandString(10))
}
