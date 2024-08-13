// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeProjectCloudArmorTier_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org":       envvar.GetTestOrgFromEnv(t),
		"billingId": envvar.GetTestBillingAccountFromEnv(t),
		"projectID": fmt.Sprintf("tf-test-%d", acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_cloudArmorTier_standard(context),
			},
			{
				ResourceName:      "google_compute_project_cloud_armor_tier.cloud_armor_tier_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectCloudArmorTier_modify(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org":       envvar.GetTestOrgFromEnv(t),
		"billingId": envvar.GetTestBillingAccountFromEnv(t),
		"projectID": fmt.Sprintf("tf-test-%d", acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_cloudArmorTier_standard(context),
			},
			{
				ResourceName:      "google_compute_project_cloud_armor_tier.cloud_armor_tier_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeProject_cloudArmorTier_enterprise_paygo(context),
			},
			{
				ResourceName:      "google_compute_project_cloud_armor_tier.cloud_armor_tier_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeProject_cloudArmorTier_standard(context),
			},
			{
				ResourceName:      "google_compute_project_cloud_armor_tier.cloud_armor_tier_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeProject_cloudArmorTier_enterprise_paygo(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "%{projectID}"
  name            = "%{projectID}"
  org_id          = "%{org}"
  billing_account = "%{billingId}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_project_cloud_armor_tier" "cloud_armor_tier_config" {
  project      = google_project.project.project_id
  cloud_armor_tier = "CA_ENTERPRISE_PAYGO"
  depends_on   = [google_project_service.compute]
}
`, context)
}

func testAccComputeProject_cloudArmorTier_standard(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "%{projectID}"
  name            = "%{projectID}"
  org_id          = "%{org}"
  billing_account = "%{billingId}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_project_cloud_armor_tier" "cloud_armor_tier_config" {
  project      = google_project.project.project_id
  cloud_armor_tier = "CA_STANDARD"
  depends_on   = [google_project_service.compute]
}
`, context)
}
