package google

import (
	"fmt"
	"testing"

	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func TestAccComputeSnapshot_basic(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var snapshot compute.Snapshot
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSnapshotDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSnapshot_basic(snapshotName, diskName, "my-value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotExists(
						"google_compute_snapshot.foobar", &snapshot),
				),
			},
		},
	})
}

func TestAccComputeSnapshot_update(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var snapshot compute.Snapshot
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSnapshotDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSnapshot_basic(snapshotName, diskName, "my-value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotExists(
						"google_compute_snapshot.foobar", &snapshot),
				),
			},
			resource.TestStep{
				Config: testAccComputeSnapshot_basic(snapshotName, diskName, "my-updated-value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotExists(
						"google_compute_snapshot.foobar", &snapshot),
				),
			},
		},
	})
}

func TestAccComputeSnapshot_encryption(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	var snapshot compute.Snapshot

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSnapshotDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSnapshot_encryption(snapshotName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSnapshotExists(
						"google_compute_snapshot.foobar", &snapshot),
				),
			},
		},
	})
}

func testAccCheckComputeSnapshotDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_snapshot" {
			continue
		}

		_, err := config.clientCompute.Snapshots.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
				return nil
			} else if ok {
				return fmt.Errorf("Error while requesting Google Cloud Plateform: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
			}
			return fmt.Errorf("Error while requesting Google Cloud Plateform")
		}
		return fmt.Errorf("Snapshot still exists")
	}

	return nil
}

func testAccCheckComputeSnapshotExists(n string, snapshot *compute.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Snapshots.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Snapshot %s not found", n)
		}

		attr := rs.Primary.Attributes["snapshot_encryption_key_sha256"]
		if found.SnapshotEncryptionKey != nil && found.SnapshotEncryptionKey.Sha256 != attr {
			return fmt.Errorf("Snapshot %s has mismatched encryption key (Sha256).\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.SnapshotEncryptionKey.Sha256)
		} else if found.SnapshotEncryptionKey == nil && attr != "" {
			return fmt.Errorf("Snapshot %s has mismatched encryption key.\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.SnapshotEncryptionKey)
		}

		attr = rs.Primary.Attributes["source_disk_encryption_key_sha256"]
		if found.SourceDiskEncryptionKey != nil && found.SourceDiskEncryptionKey.Sha256 != attr {
			return fmt.Errorf("Snapshot %s has mismatched source disk encryption key (Sha256).\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.SourceDiskEncryptionKey.Sha256)
		} else if found.SourceDiskEncryptionKey == nil && attr != "" {
			return fmt.Errorf("Snapshot %s has mismatched source disk encryption key.\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.SourceDiskEncryptionKey)
		}

		attr = rs.Primary.Attributes["source_disk_link"]
		if found.SourceDisk != attr {
			return fmt.Errorf("Snapshot %s has mismatched source disk link.\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.SourceDisk)
		}

		foundDisk, errDisk := config.clientCompute.Disks.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["source_disk"]).Do()
		if errDisk != nil {
			return errDisk
		}
		if foundDisk.SelfLink != attr {
			return fmt.Errorf("Snapshot %s has mismatched source disk\nTF State: %+v.\nGCP State: %+v",
				n, attr, foundDisk.SelfLink)
		}

		attr = rs.Primary.Attributes["self_link"]
		if found.SelfLink != attr {
			return fmt.Errorf("Snapshot %s has mismatched self link.\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.SelfLink)
		}

		// We should have a map
		attr, ok = rs.Primary.Attributes["labels.%"]
		if !ok {
			return fmt.Errorf("Snapshot %s has no labels map in attributes", n)
		}
		// Parse out our map
		attrMap := make(map[string]string)
		for k, v := range rs.Primary.Attributes {
			if !strings.HasPrefix(k, "labels.") || k == "labels.%" {
				continue
			}
			key := k[len("labels."):]
			attrMap[key] = v
		}
		if (len(attrMap) != 0 || len(found.Labels) != 0) && !reflect.DeepEqual(attrMap, found.Labels) {
			return fmt.Errorf("Snapshot %s has mismatched labels.\nTF State: %+v\nGCP State: %+v",
				n, attrMap, found.Labels)
		}

		attr = rs.Primary.Attributes["label_fingerprint"]
		if found.LabelFingerprint != attr {
			return fmt.Errorf("Snapshot %s has mismatched label fingerprint\nTF State: %+v.\nGCP State: %+v",
				n, attr, found.LabelFingerprint)
		}

		*snapshot = *found

		return nil
	}
}

func testAccComputeSnapshot_basic(snapshotName, diskName, labelValue string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 10
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
	name = "%s"
	source_disk = "${google_compute_disk.foobar.name}"
	zone = "us-central1-a"
	labels = {
		my_label = "%s"
	}
}`, diskName, snapshotName, labelValue)
}

func testAccComputeSnapshot_encryption(snapshotName string, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 10
	type = "pd-ssd"
	zone = "us-central1-a"
	disk_encryption_key {
		raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	}
}
resource "google_compute_snapshot" "foobar" {
	name = "%s"
	source_disk = "${google_compute_disk.foobar.name}"
	zone = "us-central1-a"
	source_disk_encryption_key_raw = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	snapshot_encryption_key_raw = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
}`, diskName, snapshotName)
}
