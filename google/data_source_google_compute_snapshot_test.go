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
	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeDataSnapshot_basic(snapshotName, diskName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSnapshotCheck(
						"data.google_compute_snapshot.foobar", "google_compute_snapshot.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSnapshotCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no datasource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes
		snapshot_attrs_to_test := []string{
			"id",
			"self_link",
			"name",
			"snapshot_encryption_key_sha256",
			"source_disk_link",
			"source_disk_encryption_key_sha256",
		}

		for _, attr_to_check := range snapshot_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}
		return nil
	}
}

func testAccComputeDataSnapshot_basic(snapshotName string, diskName string) string {
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

data "google_compute_snapshot" "foobar" {
	name = "${google_compute_snapshot.foobar.name}"
}
`, diskName, snapshotName)
}
