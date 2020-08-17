package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleProject_basic(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectConfig(project, org),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_project.project",
						"google_project.project",
						map[string]struct{}{
							// Virtual fields
							"auto_create_network": {},
							"skip_delete":         {},
						}),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectConfig(project, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

data "google_project" "project" {
  project_id = google_project.project.project_id
}
`, project, project, org)
}
