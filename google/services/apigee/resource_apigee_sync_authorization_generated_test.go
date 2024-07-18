// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeSyncAuthorization_apigeeSyncAuthorizationBasicTestExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSyncAuthorization_apigeeSyncAuthorizationBasicTestExample(context),
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

func testAccApigeeSyncAuthorization_apigeeSyncAuthorizationBasicTestExample(context map[string]interface{}) string {
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

resource "google_project_iam_member" "synchronizer-iam" {
  project = google_project.project.project_id
  role    = "roles/apigee.synchronizerManager"
  member = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_apigee_sync_authorization" "apigee_sync_authorization" {
  name       = google_apigee_organization.apigee_org.name
  identities = [
    "serviceAccount:${google_service_account.service_account.email}",
  ]
  depends_on = [google_project_iam_member.synchronizer-iam]
}
`, context)
}
