package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleProjectOrganizationPolicy_basic(t *testing.T) {
	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleProjectOrganizationPolicy_basic(project),
				Check: checkDataSourceStateMatchesResourceState(
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
