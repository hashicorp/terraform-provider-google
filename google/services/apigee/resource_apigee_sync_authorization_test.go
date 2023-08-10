// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeSyncAuthorization_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSyncAuthorization_basic(context),
			},
			{
				ResourceName:            "google_apigee_sync_authorization.apigee_sync_authorization",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccApigeeSyncAuthorization_multipleIdentities(context),
			},
			{
				ResourceName:            "google_apigee_sync_authorization.apigee_sync_authorization",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccApigeeSyncAuthorization_emptyIdentities(context),
			},
			{
				ResourceName:            "google_apigee_sync_authorization.apigee_sync_authorization",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccApigeeSyncAuthorization_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-my-project%{random_suffix}"
  name            = "tf-test-my-project%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project.project_id

  runtime_type       = "HYBRID"
  depends_on         = [google_project_service.apigee]
}

resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Service Account"
}

resource "google_project_iam_binding" "synchronizer-iam" {
  project = google_project.project.project_id
  role    = "roles/apigee.synchronizerManager"
  members = [
    "serviceAccount:${google_service_account.service_account.email}",
  ]
}

resource "google_apigee_sync_authorization" "apigee_sync_authorization" {
  name       = google_apigee_organization.apigee_org.name
  identities = [
    "serviceAccount:${google_service_account.service_account.email}",
  ]
  depends_on = [google_project_iam_binding.synchronizer-iam]
}
`, context)
}

func testAccApigeeSyncAuthorization_multipleIdentities(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-my-project%{random_suffix}"
  name            = "tf-test-my-project%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project.project_id

  runtime_type       = "HYBRID"
  depends_on         = [google_project_service.apigee]
}

resource "google_service_account" "service_account1" {
  account_id   = "tf-test-my-account1%{random_suffix}"
  display_name = "Service Account"
}

resource "google_service_account" "service_account2" {
  account_id   = "tf-test-my-account2%{random_suffix}"
  display_name = "Service Account"
}

resource "google_service_account" "service_account3" {
  account_id   = "tf-test-my-account3%{random_suffix}"
  display_name = "Service Account"
}

resource "google_project_iam_binding" "synchronizer-iam" {
  project = google_project.project.project_id
  role    = "roles/apigee.synchronizerManager"
  members = [
    "serviceAccount:${google_service_account.service_account1.email}",
    "serviceAccount:${google_service_account.service_account2.email}",
    "serviceAccount:${google_service_account.service_account3.email}",
  ]
}

resource "google_apigee_sync_authorization" "apigee_sync_authorization" {
  name       = google_apigee_organization.apigee_org.name
  identities = [
    "serviceAccount:${google_service_account.service_account1.email}",
    "serviceAccount:${google_service_account.service_account2.email}",
    "serviceAccount:${google_service_account.service_account3.email}"
  ]
  depends_on = [google_project_iam_binding.synchronizer-iam]
}
`, context)
}

func testAccApigeeSyncAuthorization_emptyIdentities(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test-my-project%{random_suffix}"
  name            = "tf-test-my-project%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region   = "us-central1"
  project_id         = google_project.project.project_id

  runtime_type       = "HYBRID"
  depends_on         = [google_project_service.apigee]
}

resource "google_apigee_sync_authorization" "apigee_sync_authorization" {
  name       = google_apigee_organization.apigee_org.name
  identities = []
}
`, context)
}
