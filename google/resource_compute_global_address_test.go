package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeGlobalAddress_basic(t *testing.T) {
	t.Parallel()

	var addr compute.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeGlobalAddress_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeGlobalAddressExists(
						"google_compute_global_address.foobar", &addr),

					// implicitly IPV4 - if we don't send an ip_version, we don't get one back.
					testAccCheckComputeGlobalAddressIpVersion("google_compute_global_address.foobar", ""),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_global_address.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeGlobalAddress_ipv6(t *testing.T) {
	t.Parallel()

	var addr compute.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeGlobalAddress_ipv6(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeGlobalAddressExists(
						"google_compute_global_address.foobar", &addr),
					testAccCheckComputeGlobalAddressIpVersion("google_compute_global_address.foobar", "IPV6"),
				),
			},
			resource.TestStep{
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
			resource.TestStep{
				Config: testAccComputeGlobalAddress_internal(),
			},
			resource.TestStep{
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

func testAccComputeGlobalAddress_basic() string {
	return fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
	name = "address-test-%s"
	description = "Created for Terraform acceptance testing"
}`, acctest.RandString(10))
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
  network = "${google_compute_network.foobar.self_link}"
}`, acctest.RandString(10), acctest.RandString(10))
}
