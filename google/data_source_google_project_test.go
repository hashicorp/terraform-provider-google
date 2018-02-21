package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleProject_basic(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	project := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectConfig(project, org),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleProjectCheck("data.google_project.project", "google_project.project"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleProjectCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
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

		projectAttrToCheck := []string{
			"project_id",
			"name",
			"billing_account",
			"org_id",
			"folder_id",
			"number",
		}

		for _, attr := range projectAttrToCheck {
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

func testAccCheckGoogleProjectConfig(project, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id = "%s"
	name = "%s"
	org_id = "%s"
}
	
data "google_project" "project" {
	project_id = "${google_project.project.project_id}"
}`, project, project, org)
}
