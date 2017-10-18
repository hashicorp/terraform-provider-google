package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleProjectServices_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_services.acceptance"
	projectId := "terraform-" + acctest.RandString(10)
	conf := testAccGoogleProjectServices_import(projectId, org, pname)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     projectId,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGoogleProjectServices_import(projectId, orgId, projectName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    org_id = "%s"
    name = "%s"
}

resource "google_project_services" "acceptance" {
    project = "${google_project.acceptance.project_id}"
    services = [
	  	"servicemanagement.googleapis.com",
	  	"iam.googleapis.com",
	  	"cloudresourcemanager.googleapis.com",
	]
}`, projectId, orgId, projectName)
}
