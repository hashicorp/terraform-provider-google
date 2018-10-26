package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeVpnGateway_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeVpnGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeVpnGateway_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeVpnGatewayExists(
						"google_compute_vpn_gateway.foobar"),
					testAccCheckComputeVpnGatewayExists(
						"google_compute_vpn_gateway.baz"),
				),
			},
		},
	})
}

func testAccCheckComputeVpnGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]
		region := rs.Primary.Attributes["region"]
		project := config.Project

		vpnGatewaysService := compute.NewTargetVpnGatewaysService(config.clientCompute)
		_, err := vpnGatewaysService.Get(project, region, name).Do()

		if err != nil {
			return fmt.Errorf("Error Reading VPN Gateway %s: %s", name, err)
		}

		return nil
	}
}

func testAccComputeVpnGateway_basic() string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
	name = "gateway-test-%s"
	auto_create_subnetworks = false
	ipv4_range = "10.0.0.0/16"
}

resource "google_compute_vpn_gateway" "foobar" {
	name = "gateway-test-%s"
	network = "${google_compute_network.foobar.self_link}"
	region = "us-central1"
}
resource "google_compute_vpn_gateway" "baz" {
	name = "gateway-test-%s"
	network = "${google_compute_network.foobar.name}"
	region = "us-central1"
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}
