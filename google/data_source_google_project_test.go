package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleProject_basic(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	project := "terraform-" + acctest.RandString(10)

	projectAttrToCheck := []string{
		"project_id",
		"name",
		"billing_account",
		"org_id",
		"folder_id",
		"number",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectConfig(project, org),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceMatchesResourceCheck("data.google_project.project", "google_project.project", projectAttrToCheck),
				),
			},
		},
	})
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
