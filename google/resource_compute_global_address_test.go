package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeGlobalAddress_ipv6(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalAddress_ipv6(),
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

func testAccComputeGlobalAddress_ipv6() string {
	return fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
  name        = "address-test-%s"
  description = "Created for Terraform acceptance testing"
  ip_version  = "IPV6"
}
`, acctest.RandString(10))
}

func testAccComputeGlobalAddress_internal() string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "address-test-%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "address-test-%s"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 24
  address       = "172.20.181.0"
  network       = google_compute_network.foobar.self_link
}
`, acctest.RandString(10), acctest.RandString(10))
}
