package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeGlobalAddress_ipv6(t *testing.T) {
	t.Parallel()

	var addr compute.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_ipv6(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeGlobalAddressExists(
						"google_compute_global_address.foobar", &addr),
					testAccCheckComputeGlobalAddressIpVersion("google_compute_global_address.foobar", "IPV6"),
				),
			},
			{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeGlobalAddress_internal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_internal(),
			},
			{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeGlobalAddressExists(n string, addr *compute.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.GlobalAddresses.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Addr not found")
		}

		*addr = *found

		return nil
	}
}

func testAccCheckComputeGlobalAddressIpVersion(n, version string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		addr, err := config.clientCompute.GlobalAddresses.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if addr.IpVersion != version {
			return fmt.Errorf("Expected IP version to be %s, got %s", version, addr.IpVersion)
		}

		return nil
	}
}

func testAccComputeGlobalAddress_ipv6() string {
	return fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
	name = "address-test-%s"
	description = "Created for Terraform acceptance testing"
	ip_version = "IPV6"
}`, acctest.RandString(10))
}

func testAccComputeGlobalAddress_internal() string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "address-test-%s"
}


resource "google_compute_global_address" "foobar" {
  name = "address-test-%s"
  address_type = "INTERNAL"
  purpose = "VPC_PEERING"
  prefix_length = 24
	address = "172.20.181.0"
  network = "${google_compute_network.foobar.self_link}"
}`, acctest.RandString(10), acctest.RandString(10))
}
