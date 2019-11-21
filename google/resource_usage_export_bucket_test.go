package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeResourceUsageExportBucket(t *testing.T) {
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)

	baseProject := "ub-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUsageExportBucket(baseProject, org, billingId),
			},
			// Test import.
			{
				ResourceName:      "google_project_usage_export_bucket.ueb",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceUsageExportBucket(baseProject, org, billingId string) string {
	return fmt.Sprintf(`
resource "google_project" "base" {
  project_id      = "%s"
  name            = "Export Bucket Base"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "service" {
  project = google_project.base.project_id
  service = "compute.googleapis.com"
}

resource "google_storage_bucket" "bucket" {
  name    = "b-${google_project.base.project_id}"
  project = google_project_service.service.project
}

resource "google_project_usage_export_bucket" "ueb" {
  project     = google_project.base.project_id
  bucket_name = google_storage_bucket.bucket.name
  prefix      = "foobar"
}
`, baseProject, org, billingId)
}
