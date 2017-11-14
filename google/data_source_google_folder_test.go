package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleFolder(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
		}...,
	)

	parent := fmt.Sprintf("organizations/%s", os.Getenv("GOOGLE_ORG"))
	displayName := "terraform-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceGoogleFolderConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleFolderCheck("data.google_folder.my_folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleFolderCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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
		folder_attrs_to_test := []string{"parent", "display_name", "name"}

		for _, attr_to_check := range folder_attrs_to_test {
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

func testAccDataSourceGoogleFolderConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
  parent = "%s"
  display_name = "%s"
}

data "google_folder" "my_folder" {
  parent = "${google_folder.foobar.parent}"
  display_name = "${google_folder.foobar.display_name}"
}
`, parent, displayName)
}
