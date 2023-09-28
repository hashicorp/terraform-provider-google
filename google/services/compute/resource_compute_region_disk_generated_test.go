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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeRegionDisk_regionDiskBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_regionDiskBasicExample(context),
			},
			{
				ResourceName:            "google_compute_region_disk.regiondisk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type", "region", "snapshot", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccComputeRegionDisk_regionDiskBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_disk" "regiondisk" {
  name                      = "tf-test-my-region-disk%{random_suffix}"
  snapshot                  = google_compute_snapshot.snapdisk.id
  type                      = "pd-ssd"
  region                    = "us-central1"
  physical_block_size_bytes = 4096

  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_disk" "disk" {
  name  = "tf-test-my-disk%{random_suffix}"
  image = "debian-cloud/debian-11"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name        = "tf-test-my-snapshot%{random_suffix}"
  source_disk = google_compute_disk.disk.name
  zone        = "us-central1-a"
}
`, context)
}

func TestAccComputeRegionDisk_regionDiskAsyncExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_regionDiskAsyncExample(context),
			},
			{
				ResourceName:            "google_compute_region_disk.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type", "region", "snapshot", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccComputeRegionDisk_regionDiskAsyncExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_disk" "primary" {
  name                      = "tf-test-primary-region-disk%{random_suffix}"
  type                      = "pd-ssd"
  region                    = "us-central1"
  physical_block_size_bytes = 4096

  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_region_disk" "secondary" {
  name                      = "tf-test-secondary-region-disk%{random_suffix}"
  type                      = "pd-ssd"
  region                    = "us-east1"
  physical_block_size_bytes = 4096

  async_primary_disk {
    disk = google_compute_region_disk.primary.id
  }

  replica_zones = ["us-east1-b", "us-east1-c"]
}
`, context)
}

func TestAccComputeRegionDisk_regionDiskFeaturesExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_regionDiskFeaturesExample(context),
			},
			{
				ResourceName:            "google_compute_region_disk.regiondisk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"type", "region", "snapshot", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccComputeRegionDisk_regionDiskFeaturesExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_disk" "regiondisk" {
  name                      = "tf-test-my-region-features-disk%{random_suffix}"
  type                      = "pd-ssd"
  region                    = "us-central1"
  physical_block_size_bytes = 4096

  guest_os_features {
    type = "SECURE_BOOT"
  }

  guest_os_features {
    type = "MULTI_IP_SUBNET"
  }

  guest_os_features {
    type = "WINDOWS"
  }

  licenses = ["https://www.googleapis.com/compute/v1/projects/windows-cloud/global/licenses/windows-server-core"]

  replica_zones = ["us-central1-a", "us-central1-f"]
}
`, context)
}

func testAccCheckComputeRegionDiskDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_region_disk" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/disks/{{name}}")
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
				return fmt.Errorf("ComputeRegionDisk still exists at %s", url)
			}
		}

		return nil
	}
}
