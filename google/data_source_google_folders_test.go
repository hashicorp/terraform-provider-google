package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleFolders_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFoldersConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.name"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.state"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.create_time"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.update_time"),
					resource.TestCheckResourceAttrSet("data.google_folders.root-test", "folders.0.etag"),
				),
			},
		},
	})
}

func testAccCheckGoogleFoldersConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
		parent       = "%s"
		display_name = "%s"
}

data "google_folders" "root-test" {
  parent_id = "%s"
}
`, parent, displayName, parent)
}
