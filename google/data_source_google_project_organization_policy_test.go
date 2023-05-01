package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleProjectOrganizationPolicy_basic(t *testing.T) {
	project := acctest.GetTestProjectFromEnv()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleProjectOrganizationPolicy_basic(project),
				Check: acctest.CheckDataSourceStateMatchesResourceState(
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
