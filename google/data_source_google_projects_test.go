package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleProjects_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectsConfig(project),
				Check: resource.ComposeTestCheckFunc(
					// We can't guarantee no project won't have our project ID as a prefix, so we'll check set-ness rather than correctness
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.project_id"),
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.name"),
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.number"),
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.parent.id"),
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.parent.type"),
					resource.TestCheckResourceAttrSet("data.google_projects.my-project", "projects.0.create_time"),
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
