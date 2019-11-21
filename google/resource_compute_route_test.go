package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeRoute_defaultInternetGateway(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_defaultInternetGateway(),
			},
			{
				ResourceName:      "google_compute_route.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRoute_hopInstance(t *testing.T) {
	instanceName := "tf" + acctest.RandString(10)
	zone := "us-central1-b"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_hopInstance(instanceName, zone),
			},
			{
				ResourceName:      "google_compute_route.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRoute_defaultInternetGateway() string {
	return fmt.Sprintf(`
resource "google_compute_route" "foobar" {
  name             = "route-test-%s"
  dest_range       = "0.0.0.0/0"
  network          = "default"
  next_hop_gateway = "default-internet-gateway"
  priority         = 100
}
`, acctest.RandString(10))
}

func testAccComputeRoute_hopInstance(instanceName, zone string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "foo" {
  name         = "%s"
  machine_type = "n1-standard-1"
  zone         = "%s"

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_route" "foobar" {
  name                   = "route-test-%s"
  dest_range             = "0.0.0.0/0"
  network                = "default"
  next_hop_instance      = google_compute_instance.foo.name
  next_hop_instance_zone = google_compute_instance.foo.zone
  priority               = 100
}
`, instanceName, zone, acctest.RandString(10))
}
