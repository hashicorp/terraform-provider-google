package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceComputeHaVpnGateway(t *testing.T) {
	t.Parallel()

	gwName := fmt.Sprintf("tf-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeHaVpnGatewayConfig(gwName),
				Check:  checkDataSourceStateMatchesResourceState("data.google_compute_ha_vpn_gateway.ha_gateway", "google_compute_ha_vpn_gateway.ha_gateway"),
			},
		},
	})
}

func testAccDataSourceComputeHaVpnGatewayConfig(gwName string) string {
	return fmt.Sprintf(`
resource "google_compute_ha_vpn_gateway" "ha_gateway" {
  name     = "%s"
  network  = google_compute_network.network1.id
}

resource "google_compute_network" "network1" {
  name                    = "%s"
  auto_create_subnetworks = false
}

data "google_compute_ha_vpn_gateway" "ha_gateway" {
  name = google_compute_ha_vpn_gateway.ha_gateway.name
}
`, gwName, gwName)
}
