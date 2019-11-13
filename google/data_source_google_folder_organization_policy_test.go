package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleFolderOrganizationPolicy_basic(t *testing.T) {
	folder := acctest.RandomWithPrefix("tf-test")
	org := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleFolderOrganizationPolicy_basic(org, folder),
				Check: checkDataSourceStateMatchesResourceState(
					"data.google_folder_organization_policy.data",
					"google_folder_organization_policy.resource",
				),
			},
		},
	})
}

func testAccDataSourceGoogleFolderOrganizationPolicy_basic(org, folder string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "resource" {
  folder     = google_folder.orgpolicy.name
  constraint = "serviceuser.services"

  restore_policy {
    default = true
  }
}

data "google_folder_organization_policy" "data" {
  folder     = google_folder_organization_policy.resource.folder
  constraint = "serviceuser.services"
}
`, folder, "organizations/"+org)
}
