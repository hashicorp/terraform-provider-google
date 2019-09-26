package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleProjects_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectsConfig(project),
				Check: resource.ComposeTestCheckFunc(
					// We can't guarantee no project won't have our project ID as a prefix, so we'll check set-ness rather than correctness
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.project_id"),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectsConfig(project string) string {
	return fmt.Sprintf(`
data "google_projects" "my-project" {
  filter = "projectId:%s"
}
`, project)
}
