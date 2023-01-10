package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGoogleProjectService_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", randInt(t))
	services := []string{"iam.googleapis.com", "cloudresourcemanager.googleapis.com"}
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleProjectService_basic(services, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleProjectServiceCheck("data.google_project_service.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleProjectService_basic(services []string, pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_service" "foo" {
  project = google_project.acceptance.project_id
  service = "%s"
}

data "google_project_service" "foo" {
  project = google_project.acceptance.project_id
  service = google_project_service.foo.service
}
`, pid, name, org, services[0])
}

func testAccDataSourceGoogleProjectServiceCheck(datasourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		}

		if ds.Primary.Attributes["id"] == "" {
			return fmt.Errorf("specified API service is not enabled")
		}

		return nil
	}
}
