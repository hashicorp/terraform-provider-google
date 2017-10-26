package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleProjectServices_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_services.acceptance"
	projectId := "terraform-" + acctest.RandString(10)
	services := []string{"iam.googleapis.com", "cloudresourcemanager.googleapis.com", "servicemanagement.googleapis.com"}

	conf := testAccGoogleProjectAssociateServicesBasic(services, projectId, pname, org)

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
