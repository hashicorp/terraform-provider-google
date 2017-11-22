package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeGlobalAddress_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_basic,
			},
			{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeGlobalAddress_ipv6(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_ipv6,
			},
			{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeGlobalAddressDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_global_address" {
			continue
		}

		_, err := config.clientCompute.GlobalAddresses.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Address still exists")
		}
	}

	return nil
}

var testAccComputeGlobalAddress_basic = fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
	name = "address-test-%s"
}`, acctest.RandString(10))

var testAccComputeGlobalAddress_ipv6 = fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
	name = "address-test-%s"
	ip_version = "IPV6"
}`, acctest.RandString(10))
