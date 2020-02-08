package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeRegionAutoscaler_update(t *testing.T) {
	var it_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var tp_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var igm_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))
	var autoscaler_name = fmt.Sprintf("region-autoscaler-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionAutoscalerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionAutoscaler_basic(it_name, tp_name, igm_name, autoscaler_name),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionAutoscaler_update(it_name, tp_name, igm_name, autoscaler_name),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionAutoscaler_basic(it_name, tp_name, igm_name, autoscaler_name string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "foobar" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.foobar.self_link
    name              = "primary"
  }
  target_pools       = [google_compute_target_pool.foobar.self_link]
  base_instance_name = "foobar"
  region             = "us-central1"
}

resource "google_compute_region_autoscaler" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  region      = "us-central1"
  target      = google_compute_region_instance_group_manager.foobar.self_link
  autoscaling_policy {
    max_replicas    = 5
    min_replicas    = 0
    cooldown_period = 60
    cpu_utilization {
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
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "foobar" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.foobar.self_link
    name              = "primary"
  }
  target_pools       = [google_compute_target_pool.foobar.self_link]
  base_instance_name = "foobar"
  region             = "us-central1"
}

resource "google_compute_region_autoscaler" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  region      = "us-central1"
  target      = google_compute_region_instance_group_manager.foobar.self_link
  autoscaling_policy {
    max_replicas    = 10
    min_replicas    = 1
    cooldown_period = 60
    cpu_utilization {
      target = 0.5
    }
  }
}
`, it_name, tp_name, igm_name, autoscaler_name)
}
