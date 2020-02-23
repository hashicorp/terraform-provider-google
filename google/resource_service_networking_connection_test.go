package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccServiceNetworkingConnection_create(t *testing.T) {
	t.Parallel()

	network := BootstrapSharedTestNetwork(t, "service-networking-connection-create")
	addr := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	service := "servicenetworking.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testServiceNetworkingConnectionDestroy(service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(network, addr, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceNetworkingConnection_update(t *testing.T) {
	t.Parallel()

	network := BootstrapSharedTestNetwork(t, "service-networking-connection-update")
	addr1 := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	addr2 := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	service := "servicenetworking.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testServiceNetworkingConnectionDestroy(service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(network, addr1, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceNetworkingConnection(network, addr2, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testServiceNetworkingConnectionDestroy(parent, network string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		parentService := "services/" + parent
		networkName := fmt.Sprintf("projects/%s/global/networks/%s", getTestProjectFromEnv(), network)

		response, err := config.clientServiceNetworking.Services.Connections.List(parentService).
			Network(networkName).Do()
		if err != nil {
			return err
		}

		for _, c := range response.Connections {
			if c.Network == networkName {
				return fmt.Errorf("Found %s which should have been destroyed.", networkName)
			}
		}

		return nil
	}
}

func testAccServiceNetworkingConnection(networkName, addressRangeName, serviceName string) string {
	return fmt.Sprintf(`
data "google_compute_network" "servicenet" {
  name = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = data.google_compute_network.servicenet.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}
`, networkName, addressRangeName, serviceName)
}
