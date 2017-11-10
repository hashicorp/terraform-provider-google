package google

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeDisk_basic(t *testing.T) {
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_basic(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					testAccCheckComputeDiskHasLabel(&disk, "my-label", "my-label-value"),
					testAccCheckComputeDiskHasLabelFingerprint(&disk, "google_compute_disk.foobar"),
				),
			},
		},
	})
}

func TestAccComputeDisk_update(t *testing.T) {
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeDisk_basic(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "size", "50"),
					testAccCheckComputeDiskHasLabel(&disk, "my-label", "my-label-value"),
					testAccCheckComputeDiskHasLabelFingerprint(&disk, "google_compute_disk.foobar"),
				),
			},
			{
				Config: testAccComputeDisk_updated(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					resource.TestCheckResourceAttr("google_compute_disk.foobar", "size", "100"),
					testAccCheckComputeDiskHasLabel(&disk, "my-label", "my-updated-label-value"),
					testAccCheckComputeDiskHasLabel(&disk, "a-new-label", "a-new-label-value"),
					testAccCheckComputeDiskHasLabelFingerprint(&disk, "google_compute_disk.foobar"),
				),
			},
		},
	})
}

func TestAccComputeDisk_fromSnapshot(t *testing.T) {
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	firstDiskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var xpn_host = os.Getenv("GOOGLE_XPN_HOST_PROJECT")

	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_fromSnapshot(firstDiskName, snapshotName, diskName, xpn_host, "self_link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.seconddisk", &disk),
				),
			},
			resource.TestStep{
				Config: testAccComputeDisk_fromSnapshot(firstDiskName, snapshotName, diskName, xpn_host, "name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.seconddisk", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_encryption(t *testing.T) {
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_encryption(diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foobar", &disk),
					testAccCheckEncryptionKey(
						"google_compute_disk.foobar", &disk),
				),
			},
		},
	})
}

func TestAccComputeDisk_deleteDetach(t *testing.T) {
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var disk compute.Disk

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetach(instanceName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
				),
			},
			// this needs to be a second step so we refresh and see the instance
			// listed as attached to the disk; the instance is created after the
			// disk. and the disk's properties aren't refreshed unless there's
			// another step
			resource.TestStep{
				Config: testAccComputeDisk_deleteDetach(instanceName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeDiskExists(
						"google_compute_disk.foo", &disk),
					testAccCheckComputeDiskInstances(
						"google_compute_disk.foo", &disk),
				),
			},
		},
	})
}

func testAccCheckComputeDiskDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_disk" {
			continue
		}

		_, err := config.clientCompute.Disks.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Disk still exists")
		}
	}

	return nil
}

func testAccCheckComputeDiskExists(n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Disks.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = *found

		return nil
	}
}

func testAccCheckComputeDiskHasLabel(disk *compute.Disk, key, value string) resource.TestCheckFunc {
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

func testAccCheckComputeDiskHasLabelFingerprint(disk *compute.Disk, resourceName string) resource.TestCheckFunc {
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

func testAccCheckEncryptionKey(n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["disk_encryption_key_sha256"]
		if disk.DiskEncryptionKey == nil && attr != "" {
			return fmt.Errorf("Disk %s has mismatched encryption key.\nTF State: %+v\nGCP State: <empty>", n, attr)
		}

		if attr != disk.DiskEncryptionKey.Sha256 {
			return fmt.Errorf("Disk %s has mismatched encryption key.\nTF State: %+v.\nGCP State: %+v",
				n, attr, disk.DiskEncryptionKey.Sha256)
		}
		return nil
	}
}

func testAccCheckComputeDiskInstances(n string, disk *compute.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		attr := rs.Primary.Attributes["users.#"]
		if strconv.Itoa(len(disk.Users)) != attr {
			return fmt.Errorf("Disk %s has mismatched users.\nTF State: %+v\nGCP State: %+v", n, rs.Primary.Attributes["users"], disk.Users)
		}

		for pos, user := range disk.Users {
			if rs.Primary.Attributes["users."+strconv.Itoa(pos)] != user {
				return fmt.Errorf("Disk %s has mismatched users.\nTF State: %+v.\nGCP State: %+v",
					n, rs.Primary.Attributes["users"], disk.Users)
			}
		}
		return nil
	}
}

func testAccComputeDisk_basic(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "debian-8-jessie-v20160803"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	labels {
		my-label = "my-label-value"
	}
}`, diskName)
}

func testAccComputeDisk_updated(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "debian-8-jessie-v20160803"
	size = 100
	type = "pd-ssd"
	zone = "us-central1-a"
	labels {
		my-label = "my-updated-label-value"
		a-new-label = "a-new-label-value"
	}
}`, diskName)
}

func testAccComputeDisk_fromSnapshot(firstDiskName, snapshotName, diskName, xpn_host string, ref_selector string) string {
	return fmt.Sprintf(`
		resource "google_compute_disk" "foobar" {
			name = "%s"
			image = "debian-8-jessie-v20160803"
			size = 50
			type = "pd-ssd"
			zone = "us-central1-a"
			project = "%s"
		}

resource "google_compute_snapshot" "snapdisk" {
  name = "%s"
  source_disk = "${google_compute_disk.foobar.name}"
  zone = "us-central1-a"
	project = "%s"
}
resource "google_compute_disk" "seconddisk" {
	name = "%s"
	snapshot = "${google_compute_snapshot.snapdisk.%s}"
	type = "pd-ssd"
	zone = "us-central1-a"
}`, firstDiskName, xpn_host, snapshotName, xpn_host, diskName, ref_selector)
}

func testAccComputeDisk_encryption(diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "debian-8-jessie-v20160803"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
	disk_encryption_key_raw = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
}`, diskName)
}

func testAccComputeDisk_deleteDetach(instanceName, diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foo" {
	name = "%s"
	image = "debian-8-jessie-v20170523"
	size = 50
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_instance" "bar" {
	name = "%s"
	machine_type = "n1-standard-1"
	zone = "us-central1-a"

	disk {
		image = "debian-8-jessie-v20170523"
	}

	disk {
		disk = "${google_compute_disk.foo.name}"
		auto_delete = false
	}

	network_interface {
		network = "default"
	}
}`, diskName, instanceName)
}
