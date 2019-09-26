package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleProjectServices_basic(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	project := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectServicesConfig(project, org),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_project_services.project_services",
						"google_project_services.project_services",
						map[string]struct{}{
							// Virtual fields
							"disable_on_destroy": {},
						},
					),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectServicesConfig(project, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id = "%s"
	name = "%s"
	org_id = "%s"
}

resource "google_project_services" "project_services" {
	project = "${google_project.project.project_id}"
	services = ["admin.googleapis.com"]
}

data "google_project_services" "project_services" {
	project = "${google_project_services.project_services.project}"
}`, project, project, org)
}
