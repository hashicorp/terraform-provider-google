package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeSnapshot_encryption(t *testing.T) {
	t.Parallel()

	snapshotName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	diskName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSnapshotDestroyProducer(t),
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

func testAccComputeSnapshot_encryption(snapshotName string, diskName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "%s"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
  disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}

resource "google_compute_snapshot" "foobar" {
  name        = "%s"
  source_disk = google_compute_disk.foobar.name
  zone        = "us-central1-a"
  snapshot_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }

  source_disk_encryption_key {
    raw_key = "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
  }
}
`, diskName, snapshotName)
}
