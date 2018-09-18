package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description1"),
			},
			resource.TestStep{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDnsManagedZone_update(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description1"),
			},
			resource.TestStep{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description2"),
			},
			resource.TestStep{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccDnsManagedZone_basic(suffix, description string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
	name = "mzone-test-%s"
	dns_name = "tf-acctest-%s.hashicorptest.com."
	description = "%s"
}`, suffix, suffix, description)
}
