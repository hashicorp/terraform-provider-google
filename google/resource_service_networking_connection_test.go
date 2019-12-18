package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

// Standard checkDestroy cannot be used here because destroying the network will delete
// all the networking connections so this would return false positives.
func TestAccServiceNetworkingConnectionDestroy(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	addressRange := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(
					network,
					addressRange,
					"servicenetworking.googleapis.com",
				),
			},
			{
				Config: testAccServiceNetworkingConnectionDestroy(network, addressRange),
				Check: resource.ComposeTestCheckFunc(
					testServiceNetworkingConnectionDestroy("servicenetworking.googleapis.com", network, getTestProjectFromEnv()),
				),
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

func testServiceNetworkingConnectionDestroy(parent, network, project string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		parentService := "services/" + parent
		networkName := fmt.Sprintf("projects/%s/global/networks/%s", project, network)

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
resource "google_compute_network" "foobar" {
  name = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.foobar.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = google_compute_network.foobar.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}
`, networkName, addressRangeName, serviceName)
}

func testAccServiceNetworkingConnectionDestroy(networkName, addressRangeName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s"
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.foobar.self_link
}
`, networkName, addressRangeName)
}
