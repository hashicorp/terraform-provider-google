package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleActiveFolder(t *testing.T) {
	org := getTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	suffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceGoogleActiveFolderConfig(parent, "terraform-test-"+suffix, "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleActiveFolderCheck("default"),
				),
			},
			resource.TestStep{
				Config: testAccDataSourceGoogleActiveFolderConfig(parent, "terraform test "+suffix, "space"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleActiveFolderCheck("space"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleActiveFolderCheck(resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds_name := "data.google_active_folder." + resource_name
		ds, ok := s.RootModule().Resources[ds_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", ds_name)
		}

		rs_name := "google_folder." + resource_name
		rs, ok := s.RootModule().Resources[rs_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", rs_name)
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

func testAccDataSourceGoogleActiveFolderConfig(parent string, displayName string, resourceName string) string {
	return fmt.Sprintf(`
resource "google_folder" "%s" {
  parent = "%s"
  display_name = "%s"
}

data "google_active_folder" "%s" {
  parent = "${google_folder.%s.parent}"
  display_name = "${google_folder.%s.display_name}"
}
`, resourceName, parent, displayName, resourceName, resourceName, resourceName)
}
