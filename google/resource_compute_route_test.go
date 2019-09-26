package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeRoute_defaultInternetGateway(t *testing.T) {
	t.Parallel()

	var route compute.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_defaultInternetGateway(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouteExists(
						"google_compute_route.foobar", &route),
				),
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
	var route compute.Route

	instanceName := "tf" + acctest.RandString(10)
	zone := "us-central1-b"
	instanceNameRegexp := regexp.MustCompile(fmt.Sprintf("projects/(.+)/zones/%s/instances/%s$", zone, instanceName))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRoute_hopInstance(instanceName, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRouteExists(
						"google_compute_route.foobar", &route),
					resource.TestMatchResourceAttr("google_compute_route.foobar", "next_hop_instance", instanceNameRegexp),
					resource.TestMatchResourceAttr("google_compute_route.foobar", "next_hop_instance", instanceNameRegexp),
				),
			},
			{
				ResourceName:      "google_compute_route.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeRouteExists(n string, route *compute.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Routes.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Route not found")
		}

		*route = *found

		return nil
	}
}

func testAccComputeRoute_defaultInternetGateway() string {
	return fmt.Sprintf(`
resource "google_compute_route" "foobar" {
	name = "route-test-%s"
	dest_range = "0.0.0.0/0"
	network = "default"
	next_hop_gateway = "default-internet-gateway"
	priority = 100
}`, acctest.RandString(10))
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
    initialize_params{
      image = "${data.google_compute_image.my_image.self_link}"
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_route" "foobar" {
	name = "route-test-%s"
	dest_range = "0.0.0.0/0"
	network = "default"
  	next_hop_instance = "${google_compute_instance.foo.name}"
  	next_hop_instance_zone = "${google_compute_instance.foo.zone}"
	priority = 100
}`, instanceName, zone, acctest.RandString(10))
}
