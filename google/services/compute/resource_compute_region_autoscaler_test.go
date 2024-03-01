// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeRegionAutoscaler_update(t *testing.T) {
	var itName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerName = fmt.Sprintf("tf-test-region-autoscaler-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionAutoscalerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionAutoscaler_basic(itName, tpName, igmName, autoscalerName),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionAutoscaler_update(itName, tpName, igmName, autoscalerName),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionAutoscaler_scaleDownControl(t *testing.T) {
	t.Parallel()

	var itName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerName = fmt.Sprintf("tf-test-region-autoscaler-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionAutoscalerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionAutoscaler_scaleDownControl(itName, tpName, igmName, autoscalerName),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionAutoscaler_scalingSchedule(t *testing.T) {
	t.Parallel()

	var itName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerName = fmt.Sprintf("tf-test-region-autoscaler-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionAutoscalerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionAutoscaler_scalingSchedule(itName, tpName, igmName, autoscalerName),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionAutoscaler_scaleInControl(t *testing.T) {
	t.Parallel()

	var itName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerName = fmt.Sprintf("tf-test-region-autoscaler-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionAutoscalerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionAutoscaler_scaleInControl(itName, tpName, igmName, autoscalerName),
			},
			{
				ResourceName:      "google_compute_region_autoscaler.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar" {
  name           = "%s"
  machine_type   = "e2-medium"
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
  base_instance_name = "tf-test-foobar"
  region             = "us-central1"
}

`, itName, tpName, igmName)
}

func testAccComputeRegionAutoscaler_basic(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
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
`, autoscalerName)
}

func testAccComputeRegionAutoscaler_update(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
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
`, autoscalerName)
}

func testAccComputeRegionAutoscaler_scaleDownControl(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
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
      predictive_method = "OPTIMIZE_AVAILABILITY"
    }
  }
}
`, autoscalerName)
}

func testAccComputeRegionAutoscaler_scaleInControl(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
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
    scale_in_control {
      max_scaled_in_replicas {
        percent = 80
      }
      time_window_sec = 300
    }
  }
}
`, autoscalerName)
}

func testAccComputeRegionAutoscaler_scalingSchedule(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
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
    scaling_schedules {
      name = "every-weekday-morning"
      description = "Increase to 2 every weekday at 7AM for 6 hours."
      min_required_replicas = 0
      schedule = "0 7 * * MON-FRI"
      time_zone = "America/New_York"
      duration_sec = 21600
    }
    scaling_schedules {
      name = "every-weekday-afternoon"
      description = "Increase to 2 every weekday at 7PM for 6 hours."
      min_required_replicas = 2
      schedule = "0 19 * * MON-FRI"
      time_zone = "America/New_York"
      duration_sec = 21600
    }
  }
}
`, autoscalerName)
}
