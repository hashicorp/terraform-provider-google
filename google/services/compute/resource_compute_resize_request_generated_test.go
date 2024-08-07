// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeResizeRequest_computeMigResizeRequestExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeResizeRequestDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeResizeRequest_computeMigResizeRequestExample(context),
			},
			{
				ResourceName:            "google_compute_resize_request.a3_resize_request",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_group_manager", "zone"},
			},
		},
	})
}

func testAccComputeResizeRequest_computeMigResizeRequestExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_instance_template" "a3_dws" {
  name                 = "a3-dws"
  region               = "us-central1"
  description          = "This template is used to create a mig instance that is compatible with DWS resize requests."
  instance_description = "A3 GPU"
  machine_type         = "a3-highgpu-8g"
  can_ip_forward       = false

  scheduling {
    automatic_restart   = false
    on_host_maintenance = "TERMINATE"
  }

  disk {
    source_image = "cos-cloud/cos-105-lts"
    auto_delete  = true
    boot         = true
    disk_type    = "pd-ssd"
    disk_size_gb = "960"
    mode         = "READ_WRITE"
  }

  guest_accelerator {
    type  = "nvidia-h100-80gb"
    count = 8
  }

  reservation_affinity {
    type = "NO_RESERVATION"
  }

  shielded_instance_config {
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_instance_group_manager" "a3_dws" {
  name               = "a3-dws"
  base_instance_name = "a3-dws"
  zone               = "us-central1-a"

  version {
    instance_template = google_compute_region_instance_template.a3_dws.self_link
  }

  instance_lifecycle_policy {
    default_action_on_failure = "DO_NOTHING"
  }

  wait_for_instances = false

}

resource "google_compute_resize_request" "a3_resize_request" {
  name                   = "tf-test-a3-dws%{random_suffix}"
  instance_group_manager = google_compute_instance_group_manager.a3_dws.name
  zone                   = "us-central1-a"
  description            = "Test resize request resource"
  resize_by              = 2
  requested_run_duration {
    seconds = 14400
    nanos   = 0
  }
}
`, context)
}

func testAccCheckComputeResizeRequestDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_resize_request" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{instance_group_manager}}/resizeRequests/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeResizeRequest still exists at %s", url)
			}
		}

		return nil
	}
}
