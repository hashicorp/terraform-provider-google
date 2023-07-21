// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccComputeAttachedDisk_basic(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-disk-%d", acctest.RandInt(t))
	instanceName := fmt.Sprintf("tf-test-inst-%d", acctest.RandInt(t))
	importID := fmt.Sprintf("%s/us-central1-a/%s/%s", envvar.GetTestProjectFromEnv(), instanceName, diskName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// Check destroy isn't a good test here, see comment on testCheckAttachedDiskIsNowDetached
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAttachedDiskResource(diskName, instanceName) + testAttachedDiskResourceAttachment(),
			},
			{
				ResourceName:      "google_compute_attached_disk.test",
				ImportStateId:     importID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAttachedDiskResource(diskName, instanceName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAttachedDiskIsNowDetached(t, instanceName, diskName),
				),
			},
		},
	})
}

func TestAccComputeAttachedDisk_full(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	instanceName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	importID := fmt.Sprintf("%s/us-central1-a/%s/%s", envvar.GetTestProjectFromEnv(), instanceName, diskName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// Check destroy isn't a good test here, see comment on testCheckAttachedDiskIsNowDetached
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAttachedDiskResource(diskName, instanceName) + testAttachedDiskResourceAttachmentFull(),
			},
			{
				ResourceName:      "google_compute_attached_disk.test",
				ImportStateId:     importID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccComputeAttachedDisk_region(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	instanceName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	importID := fmt.Sprintf("%s/us-central1-a/%s/%s", envvar.GetTestProjectFromEnv(), instanceName, diskName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// Check destroy isn't a good test here, see comment on testCheckAttachedDiskIsNowDetached
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAttachedDiskResource_region(diskName, instanceName),
			},
			{
				ResourceName:      "google_compute_attached_disk.test",
				ImportStateId:     importID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccComputeAttachedDisk_count(t *testing.T) {
	t.Parallel()

	diskPrefix := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	instanceName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	count := 2

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAttachedDiskResourceCount(diskPrefix, instanceName, count),
				Check: resource.ComposeTestCheckFunc(
					testCheckAttachedDiskContainsManyDisks(t, instanceName, count),
				),
			},
		},
	})

}

// testCheckAttachedDiskIsNowDetached queries a compute instance and iterates through the attached
// disks to confirm that a specific disk is no longer attached to the instance
//
// This is being used instead of a CheckDestroy method because destroy will delete both the compute
// instance and the disk, whereas destroying just the attached disk should only detach the disk but
// leave the instance and disk around. So just using a normal check destroy could end up with a
// situation where the detach fails but since the instance/disk get destroyed we wouldn't notice.
func testCheckAttachedDiskIsNowDetached(t *testing.T, instanceName, diskName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		instance, err := config.NewComputeClient(config.UserAgent).Instances.Get(envvar.GetTestProjectFromEnv(), "us-central1-a", instanceName).Do()
		if err != nil {
			return err
		}

		ad := compute.FindDiskByName(instance.Disks, diskName)
		if ad != nil {
			return fmt.Errorf("compute disk is still attached to compute instance")
		}

		return nil
	}
}

func testCheckAttachedDiskContainsManyDisks(t *testing.T, instanceName string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		instance, err := config.NewComputeClient(config.UserAgent).Instances.Get(envvar.GetTestProjectFromEnv(), "us-central1-a", instanceName).Do()
		if err != nil {
			return err
		}

		// There will be 1 extra disk because of the compute instance's boot disk
		if (count + 1) != len(instance.Disks) {
			return fmt.Errorf("expected %d disks to be attached, found %d", count+1, len(instance.Disks))
		}

		return nil
	}
}

func testAttachedDiskResourceAttachment() string {
	return fmt.Sprintf(`
resource "google_compute_attached_disk" "test" {
  disk     = google_compute_disk.test1.self_link
  instance = google_compute_instance.test.self_link
}
`)
}

func testAttachedDiskResourceAttachmentFull() string {
	return fmt.Sprintf(`
resource "google_compute_attached_disk" "test" {
  disk        = google_compute_disk.test1.self_link
  instance    = google_compute_instance.test.self_link
  mode        = "READ_ONLY"
  device_name = "test-device-name"
}
`)
}

func testAttachedDiskResource_region(diskName, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_attached_disk" "test" {
  disk     = google_compute_region_disk.region.self_link
  instance = google_compute_instance.test.self_link
}

resource "google_compute_region_disk" "region" {
  name          = "%s"
  region        = "us-central1"
  replica_zones = ["us-central1-b", "us-central1-a"]
}

resource "google_compute_instance" "test" {
  name         = "%s"
  machine_type = "f1-micro"
  zone         = "us-central1-a"

  lifecycle {
    ignore_changes = [attached_disk]
  }

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}
`, diskName, instanceName)
}

func testAttachedDiskResource(diskName, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "test1" {
  name = "%s"
  zone = "us-central1-a"
  size = 10
}

resource "google_compute_instance" "test" {
  name         = "%s"
  machine_type = "f1-micro"
  zone         = "us-central1-a"

  lifecycle {
    ignore_changes = [attached_disk]
  }

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}
`, diskName, instanceName)
}

func testAttachedDiskResourceCount(diskPrefix, instanceName string, count int) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "many" {
  name  = "%s-${count.index}"
  zone  = "us-central1-a"
  size  = 10
  count = %d
}

resource "google_compute_instance" "test" {
  name         = "%s"
  machine_type = "f1-micro"
  zone         = "us-central1-a"

  lifecycle {
    ignore_changes = [attached_disk]
  }

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_attached_disk" "test" {
  count    = length(google_compute_disk.many)
  disk     = google_compute_disk.many[count.index].self_link
  instance = google_compute_instance.test.self_link
}
`, diskPrefix, count, instanceName)
}
