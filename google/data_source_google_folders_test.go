package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleFolders_basic(t *testing.T) {
	t.Parallel()

	org_id := getTestOrgFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFoldersConfig(org_id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.name"),
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.state"),
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.parent.id"),
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.parent.type"),
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.create_time"),
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.update_time"),
					// resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.delete_time"),
					// deleteTime will only be set on a deleted folder
					resource.TestCheckResourceAttrSet("data.google_folders.my-folder", "folders.0.etag"),
				),
			},
		},
	})
}

func testAccCheckGoogleFoldersConfig(org_id string) string {
	return fmt.Sprintf(`
data "google_folders" "my-folder" {
  parent_id = "organizations/%s"
}
`, org_id)
}
