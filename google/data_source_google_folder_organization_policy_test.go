package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
				Check: testAccDataSourceGoogleOrganizationPolicyCheck(
					"data.google_folder_organization_policy.data",
					"google_folder_organization_policy.resource"),
			},
		},
	})
}

func testAccDataSourceGoogleOrganizationPolicyCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		cloudFuncAttrToCheck := []string{
			"name",
			"folder",
			"constraint",
			"version",
			"list_policy",
			"restore_policy",
			"boolean_policy",
		}

		for _, attr := range cloudFuncAttrToCheck {
			if dsAttr[attr] != rsAttr[attr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					dsAttr[attr],
					rsAttr[attr],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceGoogleFolderOrganizationPolicy_basic(org, folder string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "resource" {
    folder     = "${google_folder.orgpolicy.name}"
    constraint = "serviceuser.services"

    restore_policy {
        default = true
    }
}

data "google_folder_organization_policy" "data" {
  folder     = "${google_folder.orgpolicy.name}"
  constraint = "serviceuser.services"
}
	`, folder, "organizations/"+org)
}
