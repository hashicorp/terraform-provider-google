package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeRegionAutoscaler_basic(t *testing.T) {
	var ascaler compute.Autoscaler

	var it_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var tp_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var igm_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var autoscaler_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionAutoscalerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeRegionAutoscaler_basic(it_name, tp_name, igm_name, autoscaler_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionAutoscalerExists(
						"google_compute_region_autoscaler.foobar", &ascaler),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionAutoscaler_update(t *testing.T) {
	var ascaler compute.Autoscaler

	var it_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var tp_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var igm_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var autoscaler_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionAutoscalerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeRegionAutoscaler_basic(it_name, tp_name, igm_name, autoscaler_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionAutoscalerExists(
						"google_compute_region_autoscaler.foobar", &ascaler),
				),
			},
			resource.TestStep{
				Config: testAccComputeRegionAutoscaler_update(it_name, tp_name, igm_name, autoscaler_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionAutoscalerExists(
						"google_compute_region_autoscaler.foobar", &ascaler),
					testAccCheckComputeRegionAutoscalerUpdated(
						"google_compute_region_autoscaler.foobar", 10),
				),
			},
		},
	})
}

func testAccCheckComputeRegionAutoscalerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_region_autoscaler" {
			continue
		}

		idParts := strings.Split(rs.Primary.ID, "/")
		region, name := idParts[0], idParts[1]
		_, err := config.clientCompute.RegionAutoscalers.Get(config.Project, region, name).Do()
		if err == nil {
			return fmt.Errorf("Autoscaler still exists")
		}
	}

	return nil
}

func testAccCheckComputeRegionAutoscalerExists(n string, ascaler *compute.Autoscaler) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		idParts := strings.Split(rs.Primary.ID, "/")
		region, name := idParts[0], idParts[1]
		found, err := config.clientCompute.RegionAutoscalers.Get(config.Project, region, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("Autoscaler not found")
		}

		*ascaler = *found

		return nil
	}
}

func testAccCheckComputeRegionAutoscalerUpdated(n string, max int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		idParts := strings.Split(rs.Primary.ID, "/")
		region, name := idParts[0], idParts[1]
		ascaler, err := config.clientCompute.RegionAutoscalers.Get(config.Project, region, name).Do()
		if err != nil {
			return err
		}

		if ascaler.AutoscalingPolicy.MaxNumReplicas != max {
			return fmt.Errorf("maximum replicas incorrect")
		}

		return nil
	}
}

func testAccComputeRegionAutoscaler_basic(it_name, tp_name, igm_name, autoscaler_name string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
	name = "%s"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "${data.google_compute_image.my_image.self_link}"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_target_pool" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "%s"
	instance_template = "${google_compute_instance_template.foobar.self_link}"
	target_pools = ["${google_compute_target_pool.foobar.self_link}"]
	base_instance_name = "foobar"
	region = "us-central1"
}

resource "google_compute_region_autoscaler" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	region = "us-central1"
	target = "${google_compute_region_instance_group_manager.foobar.self_link}"
	autoscaling_policy = {
		max_replicas = 5
		min_replicas = 1
		cooldown_period = 60
		cpu_utilization = {
			target = 0.5
		}
	}

}
`, it_name, tp_name, igm_name, autoscaler_name)
}

func testAccComputeRegionAutoscaler_update(it_name, tp_name, igm_name, autoscaler_name string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
	name = "%s"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "${data.google_compute_image.my_image.self_link}"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_target_pool" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "foobar" {
	description = "Terraform test instance group manager"
	name = "%s"
	instance_template = "${google_compute_instance_template.foobar.self_link}"
	target_pools = ["${google_compute_target_pool.foobar.self_link}"]
	base_instance_name = "foobar"
	region = "us-central1"
}

resource "google_compute_region_autoscaler" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	region = "us-central1"
	target = "${google_compute_region_instance_group_manager.foobar.self_link}"
	autoscaling_policy = {
		max_replicas = 10
		min_replicas = 1
		cooldown_period = 60
		cpu_utilization = {
			target = 0.5
		}
	}

}
`, it_name, tp_name, igm_name, autoscaler_name)
}
