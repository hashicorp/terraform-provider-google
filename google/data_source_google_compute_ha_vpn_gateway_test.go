package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeHaVpnGateway(t *testing.T) {
	t.Parallel()

	gwName := fmt.Sprintf("tf-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeHaVpnGatewayConfig(gwName),
				Check:  acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_ha_vpn_gateway.ha_gateway", "google_compute_ha_vpn_gateway.ha_gateway"),
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
