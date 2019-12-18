package google

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

func TestAccComputeRegionDisk_basic(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	var disk computeBeta.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_basic(diskName, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
			{
				ResourceName:      "google_compute_region_disk.regiondisk",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionDisk_basic(diskName, "name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
			{
				ResourceName:      "google_compute_region_disk.regiondisk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionDisk_basicUpdate(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	var disk computeBeta.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_basic(diskName, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
			{
				ResourceName:      "google_compute_region_disk.regiondisk",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionDisk_basicUpdated(diskName, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
					resource.TestCheckResourceAttr("google_compute_region_disk.regiondisk", "size", "100"),
					testAccCheckComputeRegionDiskHasLabel(&disk, "my-label", "my-updated-label-value"),
					testAccCheckComputeRegionDiskHasLabel(&disk, "a-new-label", "a-new-label-value"),
					testAccCheckComputeRegionDiskHasLabelFingerprint(&disk, "google_compute_region_disk.regiondisk"),
				),
			},
			{
				ResourceName:      "google_compute_region_disk.regiondisk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionDisk_encryption(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk computeBeta.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_encryption(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
					testAccCheckRegionDiskEncryptionKey(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
		},
	})
}

func TestAccComputeRegionDisk_deleteDetach(t *testing.T) {
	t.Parallel()

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	regionDiskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	regionDiskName2 := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk computeBeta.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDisk_deleteDetach(instanceName, diskName, regionDiskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
			// this needs to be an additional step so we refresh and see the instance
			// listed as attached to the disk; the instance is created after the
			// disk. and the disk's properties aren't refreshed unless there's
			// another step
			{
				Config: testAccComputeRegionDisk_deleteDetach(instanceName, diskName, regionDiskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
					testAccCheckComputeRegionDiskInstances(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
			// Change the disk name to destroy it, which detaches it from the instance
			{
				Config: testAccComputeRegionDisk_deleteDetach(instanceName, diskName, regionDiskName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
			// Add the extra step like before
			{
				Config: testAccComputeRegionDisk_deleteDetach(instanceName, diskName, regionDiskName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionDiskExists(
						"google_compute_region_disk.regiondisk", &disk),
					testAccCheckComputeRegionDiskInstances(
						"google_compute_region_disk.regiondisk", &disk),
				),
			},
		},
	})
}

func testAccCheckComputeRegionDiskExists(n string, disk *computeBeta.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		p := getTestProjectFromEnv()
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientComputeBeta.RegionDisks.Get(
			p, rs.Primary.Attributes["region"], rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("RegionDisk not found")
		}

		*disk = *found

		return nil
	}
}

func testAccCheckComputeRegionDiskHasLabel(disk *computeBeta.Disk, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := disk.Labels[key]
		if !ok {
			return fmt.Errorf("Label with key %s not found", key)
		}

		if val != value {
			return fmt.Errorf("Label value did not match for key %s: expected %s but found %s", key, value, val)
		}
		return nil
	}
}

func testAccCheckComputeRegionDiskHasLabelFingerprint(disk *computeBeta.Disk, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		state := s.RootModule().Resources[resourceName]
		if state == nil {
			return fmt.Errorf("Unable to find resource named %s", resourceName)
		}

		labelFingerprint := state.Primary.Attributes["label_fingerprint"]
		if labelFingerprint != disk.LabelFingerprint {
			return fmt.Errorf("Label fingerprints do not match: api returned %s but state has %s",
				disk.LabelFingerprint, labelFingerprint)
		}

		return nil
	}
}

func testAccCheckRegionDiskEncryptionKey(n string, disk *computeBeta.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["disk_encryption_key.0.sha256"]
		if disk.DiskEncryptionKey == nil {
			return fmt.Errorf("RegionDisk %s has mismatched encryption key.\nTF State: %+v\nGCP State: <empty>", n, attr)
		} else if attr != disk.DiskEncryptionKey.Sha256 {
			return fmt.Errorf("RegionDisk %s has mismatched encryption key.\nTF State: %+v.\nGCP State: %+v",
				n, attr, disk.DiskEncryptionKey.Sha256)
		}
		return nil
	}
}

func testAccCheckComputeRegionDiskInstances(n string, disk *computeBeta.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["users.#"]
		if strconv.Itoa(len(disk.Users)) != attr {
			return fmt.Errorf("RegionDisk %s has mismatched users.\nTF State: %+v\nGCP State: %+v", n, rs.Primary.Attributes["users"], disk.Users)
		}

		for pos, user := range disk.Users {
			if ConvertSelfLinkToV1(rs.Primary.Attributes["users."+strconv.Itoa(pos)]) != ConvertSelfLinkToV1(user) {
				return fmt.Errorf("RegionDisk %s has mismatched users.\nTF State: %+v.\nGCP State: %+v",
					n, rs.Primary.Attributes["users"], disk.Users)
			}
		}
		return nil
	}
}

func testAccComputeRegionDisk_basic(diskName, refSelector string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "disk" {
  name  = "%s"
  image = "debian-cloud/debian-9"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name        = "%s"
  source_disk = google_compute_disk.disk.name
  zone        = "us-central1-a"
}

resource "google_compute_region_disk" "regiondisk" {
  name     = "%s"
  snapshot = google_compute_snapshot.snapdisk.%s
  type     = "pd-ssd"
  region   = "us-central1"

  replica_zones = ["us-central1-a", "us-central1-f"]
}
`, diskName, diskName, diskName, refSelector)
}

func testAccComputeRegionDisk_basicUpdated(diskName, refSelector string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "disk" {
  name  = "%s"
  image = "debian-cloud/debian-9"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name        = "%s"
  source_disk = google_compute_disk.disk.name
  zone        = "us-central1-a"
}

resource "google_compute_region_disk" "regiondisk" {
  name     = "%s"
  snapshot = google_compute_snapshot.snapdisk.%s
  type     = "pd-ssd"
  region   = "us-central1"

  replica_zones = ["us-central1-a", "us-central1-f"]

  size = 100
  labels = {
    my-label    = "my-updated-label-value"
    a-new-label = "a-new-label-value"
  }
}
`, diskName, diskName, diskName, refSelector)
}

func testAccComputeRegionDisk_encryption(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "disk" {
  name  = "%s"
  image = "debian-cloud/debian-9"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name = "%s"
  zone = "us-central1-a"

  source_disk = google_compute_disk.disk.name
}

resource "google_compute_region_disk" "regiondisk" {
  name     = "%s"
  snapshot = google_compute_snapshot.snapdisk.self_link
  type     = "pd-ssd"
  region   = "us-central1"

  replica_zones = ["us-central1-a", "us-central1-f"]

  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
`, diskName, diskName, diskName)
}

func testAccComputeRegionDisk_deleteDetach(instanceName, diskName, regionDiskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "disk" {
  name  = "%s"
  image = "debian-cloud/debian-9"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name        = "%s"
  source_disk = google_compute_disk.disk.name
  zone        = "us-central1-a"
}

resource "google_compute_region_disk" "regiondisk" {
  name     = "%s"
  snapshot = google_compute_snapshot.snapdisk.self_link
  type     = "pd-ssd"
  region   = "us-central1"

  replica_zones = ["us-central1-a", "us-central1-f"]
}

resource "google_compute_instance" "inst" {
  name         = "%s"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  attached_disk {
    source = google_compute_region_disk.regiondisk.self_link
  }

  network_interface {
    network = "default"
  }
}
`, diskName, diskName, regionDiskName, instanceName)
}
