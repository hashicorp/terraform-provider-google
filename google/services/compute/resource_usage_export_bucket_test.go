// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeResourceUsageExportBucket(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)

	baseProject := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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
