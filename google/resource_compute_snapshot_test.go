package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeSnapshot_update(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_basic(snapshotName, diskName, "my-value"),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccComputeSnapshot_basic(snapshotName, diskName, "my-updated-value"),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func TestAccComputeSnapshot_encryption(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshot_encryption(snapshotName, diskName),
			},
			{
				ResourceName:            "google_compute_snapshot.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "snapshot_encryption_key", "source_disk_encryption_key"},
			},
		},
	})
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
	snapshot_encryption_key {
		raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	}

	source_disk_encryption_key {
		raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	}
}`, diskName, snapshotName)
}
