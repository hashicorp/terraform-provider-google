package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccServiceNetworkingConnectionCreate(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(
					fmt.Sprintf("tf-test-%s", acctest.RandString(10)),
					fmt.Sprintf("tf-test-%s", acctest.RandString(10)),
					"servicenetworking.googleapis.com",
				),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccServiceNetworkingConnectionUpdate(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(
					network,
					fmt.Sprintf("tf-test-%s", acctest.RandString(10)),
					"servicenetworking.googleapis.com",
				),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceNetworkingConnection(
					network,
					fmt.Sprintf("tf-test-%s", acctest.RandString(10)),
					"servicenetworking.googleapis.com",
				),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccServiceNetworkingConnection(networkName, addressRangeName, serviceName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
	name       = "%s"
}

resource "google_compute_global_address" "foobar" {
	name          = "%s"
	purpose       = "VPC_PEERING"
	address_type = "INTERNAL"
	prefix_length = 16
	network       = "${google_compute_network.foobar.self_link}"
}

resource "google_service_networking_connection" "foobar" {
	network       = "${google_compute_network.foobar.self_link}"
	service       = "%s"
	reserved_peering_ranges = ["${google_compute_global_address.foobar.name}"]
}
`, networkName, addressRangeName, serviceName)
}
