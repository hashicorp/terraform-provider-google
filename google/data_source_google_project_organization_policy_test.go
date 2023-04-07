package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleProjectOrganizationPolicy_basic(t *testing.T) {
	project := GetTestProjectFromEnv()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleProjectOrganizationPolicy_basic(project),
				Check: CheckDataSourceStateMatchesResourceState(
					"data.google_project_organization_policy.data",
					"google_project_organization_policy.resource"),
			},
		},
	})
}

func testAccDataSourceGoogleProjectOrganizationPolicy_basic(project string) string {
	return fmt.Sprintf(`
resource "google_project_organization_policy" "resource" {
  project    = "%s"
  constraint = "constraints/compute.trustedImageProjects"

  list_policy {
    allow {
      all = true
    }
  }
}

data "google_project_organization_policy" "data" {
  project    = google_project_organization_policy.resource.project
  constraint = "constraints/compute.trustedImageProjects"
}
`, project)
}
