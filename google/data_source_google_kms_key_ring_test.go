package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleKmsKeyRing_basic(t *testing.T) {
	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	billingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	folderId := os.Getenv("GOOGLE_FOLDER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRingConfig(projectId, projectOrg, billingAccount, folderId, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceMatchesResourceCheck("data.google_kms_key_ring.test", "google_kms_key_ring.key_ring", []string{"project", "location", "self_link"}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleKmsKeyRingConfig(projectId, projectOrg, billingAccount string, folderId string, keyRingName string) string {
	var parent string
	if len(folderId) != 0 {
		parent = fmt.Sprintf(`folder_id = "%s"`, folderId)
	} else if len(projectOrg) != 0 {
		parent = fmt.Sprintf(`org_id = "%s"`, projectOrg)
	}

	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  billing_account = "%s"
  %s
}

resource "google_project_service" "acceptance" {
  project  = "${google_project.project.project_id}"
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  project  = "${google_project_service.acceptance.project}"
  name     = "%s"
  location = "us-central1"
}

data "google_kms_key_ring" "test" {
  name     = "${google_kms_key_ring.key_ring.name}"
  project  = "${google_project.project.project_id}"
  location = "us-central1"
}`, projectId, projectId, billingAccount, parent, keyRingName)
}
