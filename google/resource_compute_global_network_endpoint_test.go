package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeGlobalNetworkEndpoint_networkEndpointsBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"default_port":  90,
		"modified_port": 100,
	}
	negId := fmt.Sprintf("projects/%s/global/networkEndpointGroups/neg-%s",
		getTestProjectFromEnv(), context["random_suffix"])

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputeGlobalNetworkEndpoint_networkEndpointsBasic(context),
			},
			{
				ResourceName:      "google_compute_global_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Force-recreate old endpoint
				Config: testAccComputeGlobalNetworkEndpoint_networkEndpointsModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortsDestroyed(t, negId, "90"),
				),
			},
			{
				ResourceName:      "google_compute_global_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// delete all endpoints
				Config: testAccComputeGlobalNetworkEndpoint_noNetworkEndpoints(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortsDestroyed(t, negId, "100"),
				),
			},
		},
	})
}

func testAccComputeGlobalNetworkEndpoint_networkEndpointsBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint" "default" {
  global_network_endpoint_group = google_compute_global_network_endpoint_group.neg.id

  ip_address = "8.8.8.8"
  port       = google_compute_global_network_endpoint_group.neg.default_port
}
`, context) + testAccComputeGlobalNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeGlobalNetworkEndpoint_networkEndpointsModified(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint" "default" {
  global_network_endpoint_group = google_compute_global_network_endpoint_group.neg.name

  ip_address = "8.8.8.8"
  port = "%{modified_port}"
}
`, context) + testAccComputeGlobalNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeGlobalNetworkEndpoint_noNetworkEndpoints(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint_group" "neg" {
  name                  = "neg-%{random_suffix}"
  default_port          = "%{default_port}"
  network_endpoint_type = "INTERNET_IP_PORT"
}
`, context)
}
