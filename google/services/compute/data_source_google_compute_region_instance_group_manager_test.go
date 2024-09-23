// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleComputeRegionInstanceGroupManager(t *testing.T) {
	t.Parallel()

	regionName := "us-central1"
	igmName := "tf-test-igm" + acctest.RandString(t, 6)

	context := map[string]interface{}{
		"regionName":   regionName,
		"igmName":      igmName,
		"baseName":     "tf-test-igm-base" + acctest.RandString(t, 6),
		"poolName":     "tf-test-pool" + acctest.RandString(t, 6),
		"templateName": "tf-test-templt" + acctest.RandString(t, 6),
		"autoHealName": "tf-test-ah-name" + acctest.RandString(t, 6),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeRegionInstanceGroupManager_usingSelfLink(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group_manager.data_source", "project", envvar.GetTestProjectFromEnv()),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group_manager.data_source", "region", regionName),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group_manager.data_source", "name", igmName)),
			},
			{
				Config: testAccDataSourceGoogleComputeRegionInstanceGroupManager_usingNameAndRegion(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group_manager.data_source", "project", envvar.GetTestProjectFromEnv()),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group_manager.data_source", "region", regionName),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group_manager.data_source", "name", igmName)),
			},
		},
	})
}

func testAccDataSourceGoogleComputeRegionInstanceGroupManager_usingSelfLink(context map[string]interface{}) string {
	return acctest.Nprintf(`
    resource "google_compute_health_check" "autohealing" {
        name                = "%{autoHealName}"
        check_interval_sec  = 5
        timeout_sec         = 5
        healthy_threshold   = 2
        unhealthy_threshold = 10 # 50 seconds

        http_health_check {
          request_path = "/healthz"
          port         = "8080"
        }
    }

    resource "google_compute_region_instance_group_manager" "appserver" {
        name = "%{igmName}"
        base_instance_name = "%{baseName}"
        region             = "us-central1"

        version {
          instance_template  = google_compute_instance_template.igm-basic.id
          name = "primary"
        }

        target_pools = [google_compute_target_pool.igm-basic.id]
        target_size  = 2

        named_port {
          name = "customhttp"
          port = 8888
        }

        auto_healing_policies {
          health_check      = google_compute_health_check.autohealing.id
          initial_delay_sec = 300
        }
    }

    data "google_compute_region_instance_group_manager" "data_source" {
        self_link = google_compute_region_instance_group_manager.appserver.instance_group
    }

    resource "google_compute_target_pool" "igm-basic" {
        description      = "Resource created for Terraform acceptance testing"
        name             = "%{poolName}"
        session_affinity = "CLIENT_IP_PROTO"
    }

    data "google_compute_image" "my_image" {
        family  = "debian-11"
        project = "debian-cloud"
    }

    resource "google_compute_instance_template" "igm-basic" {
        name           = "%{templateName}"
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
    }`, context)
}

func testAccDataSourceGoogleComputeRegionInstanceGroupManager_usingNameAndRegion(context map[string]interface{}) string {
	return acctest.Nprintf(`
    resource "google_compute_health_check" "autohealing" {
        name                = "%{autoHealName}"
        check_interval_sec  = 5
        timeout_sec         = 5
        healthy_threshold   = 2
        unhealthy_threshold = 10 # 50 seconds

        http_health_check {
          request_path = "/healthz"
          port         = "8080"
        }
    }

    resource "google_compute_region_instance_group_manager" "appserver" {
        name = "%{igmName}"
        base_instance_name = "%{baseName}"
        region             = "us-central1"

        version {
          instance_template  = google_compute_instance_template.igm-basic.id
          name = "primary"
        }

        target_pools = [google_compute_target_pool.igm-basic.id]
        target_size  = 2

        named_port {
          name = "customhttp"
          port = 8888
        }

        auto_healing_policies {
          health_check      = google_compute_health_check.autohealing.id
          initial_delay_sec = 300
        }
    }

    data "google_compute_region_instance_group_manager" "data_source" {
        name   = "%{igmName}"
        region = "us-central1"
    }

    resource "google_compute_target_pool" "igm-basic" {
        description      = "Resource created for Terraform acceptance testing"
        name             = "%{poolName}"
        session_affinity = "CLIENT_IP_PROTO"
    }

    data "google_compute_image" "my_image" {
        family  = "debian-11"
        project = "debian-cloud"
    }

    resource "google_compute_instance_template" "igm-basic" {
        name           = "%{templateName}"
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
    }`, context)
}
