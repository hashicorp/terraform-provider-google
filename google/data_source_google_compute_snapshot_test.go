package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleSnapshot(t *testing.T) {
	snapshotName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	snapshotLabelsName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDataSnapshot_basic(snapshotName, snapshotLabelsName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSnapshotCheck(
						"data.google_compute_snapshot.foobar", "google_compute_snapshot.foobar"),
					testAccDataSourceGoogleSnapshotCheck(
						"data.google_compute_snapshot.foobarLabels", "google_compute_snapshot.foobarLabels"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSnapshotCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no datasource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes
		snapshotAttrsToTest := []string{
			"id",
			"self_link",
			"name",
			"snapshot_encryption_key_sha256",
			"source_disk_link",
			"source_disk_encryption_key_sha256",
		}

		for _, attrToCheck := range snapshotAttrsToTest {
			if dsAttr[attrToCheck] != rsAttr[attrToCheck] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attrToCheck,
					dsAttr[attrToCheck],
					rsAttr[attrToCheck],
				)
			}
		}
		return nil
	}
}

func testAccComputeDataSnapshot_basic(snapshotName string, snapshotLabelsName string, diskName string) string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "debian-8-jessie-v20160921"
	size = 10
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_snapshot" "foobar" {
	name = "%s"
	source_disk = "${google_compute_disk.foobar.name}"
	zone = "us-central1-a"
}

resource "google_compute_snapshot" "foobarLabels" {
	name = "%s"
	source_disk = "${google_compute_disk.foobar.name}"
	zone = "us-central1-a"
	labels {
		my_key       = "my_value"
		my_other_key = "my_other_value"
	}
}

data "google_compute_snapshot" "foobar" {
	name = "${google_compute_snapshot.foobar.name}"
}

data "google_compute_snapshot" "foobarLabels" {
	labels {
		my_key       = "my_value"
		my_other_key = "my_other_value"
	}
}

`, diskName, snapshotName, snapshotLabelsName)
}
