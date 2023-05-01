package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeRoute_defaultInternetGateway(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouteDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_defaultInternetGateway(RandString(t, 10)),
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
	instanceName := "tf-test-" + RandString(t, 10)
	zone := "us-central1-b"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouteDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_hopInstance(instanceName, zone, RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_route.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRoute_defaultInternetGateway(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_route" "foobar" {
  name             = "route-test-%s"
  dest_range       = "0.0.0.0/0"
  network          = "default"
  next_hop_gateway = "default-internet-gateway"
  priority         = 100
}
`, suffix)
}

func testAccComputeRoute_hopInstance(instanceName, zone, suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foo" {
  name         = "%s"
  machine_type = "e2-medium"
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
`, instanceName, zone, suffix)
}
