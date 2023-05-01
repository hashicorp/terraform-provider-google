package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeResourceUsageExportBucket(t *testing.T) {
	org := acctest.GetTestOrgFromEnv(t)
	billingId := acctest.GetTestBillingAccountFromEnv(t)

	baseProject := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
  name     = "b-${google_project.base.project_id}"
  project  = google_project_service.service.project
  location = "US"
}

resource "google_project_usage_export_bucket" "ueb" {
  project     = google_project.base.project_id
  bucket_name = google_storage_bucket.bucket.name
  prefix      = "foobar"
}
`, baseProject, org, billingId)
}
