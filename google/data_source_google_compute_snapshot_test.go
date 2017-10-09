package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDataSourceGoogleSnapshot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: TestAccDataSourceGoogleSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSnapshotCheck("data.google_compute_snapshot.my_snapshot", "google_compute_snapshot.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSnapshotCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes
		snapshot_attrs_to_test := []string{
			"id",
			"self_link",
			"name",
			"description",
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

var TestAccDataSourceGoogleSnapshotConfig = `
resource "google_compute_snapshot" "foobar" {
	name = "snapshot-test"
	description = "my-description"
}

data "google_compute_snapshot" "my_snapshot" {
	name = "${google_compute_snapshot.foobar.name}"
}`
