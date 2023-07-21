// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflow_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowAgent_update(t *testing.T) {
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
				Config: testAccDialogflowAgent_full1(context),
			},
			{
				ResourceName:            "google_dialogflow_agent.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"avatar_uri", "tier"},
			},
			{
				Config: testAccDialogflowAgent_full2(context),
			},
			{
				ResourceName:            "google_dialogflow_agent.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"avatar_uri", "tier"},
			},
		},
	})
}

func testAccDialogflowAgent_full1(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "agent_project" {
		name = "tf-test-dialogflow-%{random_suffix}"
		project_id = "tf-test-dialogflow-%{random_suffix}"
		org_id     = "%{org_id}"
		billing_account = "%{billing_account}"
	}

	resource "google_project_service" "agent_project" {
		project = google_project.agent_project.project_id
		service = "dialogflow.googleapis.com"
		disable_dependent_services = false
	}

	resource "google_service_account" "dialogflow_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}
	  
	resource "google_project_iam_member" "agent_create" {
		project = google_project_service.agent_project.project
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflow_service_account.email}"
	}

	resource "google_dialogflow_agent" "foobar" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-%{random_suffix}"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
		enable_logging = true
		match_mode = "MATCH_MODE_ML_ONLY"
		classification_threshold = 0.3
		api_version = "API_VERSION_V2_BETA_1"
		tier = "TIER_STANDARD"
		depends_on = [google_project_iam_member.agent_create]
	}
	`, context)
}

func testAccDialogflowAgent_full2(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "agent_project" {
		name = "tf-test-dialogflow-%{random_suffix}"
		project_id = "tf-test-dialogflow-%{random_suffix}"
		org_id     = "%{org_id}"
		billing_account = "%{billing_account}"
	}

	resource "google_project_service" "agent_project" {
		project = google_project.agent_project.project_id
		service = "dialogflow.googleapis.com"
		disable_dependent_services = false
	}

	resource "google_service_account" "dialogflow_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}
	  
	resource "google_project_iam_member" "agent_create" {
		project = google_project_service.agent_project.project
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflow_service_account.email}"
	}

	resource "google_dialogflow_agent" "foobar" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-%{random_suffix}update"
		default_language_code = "en"
		supported_language_codes = ["no"]
		time_zone = "Europe/London"
		description = "Description 2!"
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo-2.png"
		enable_logging = false
		match_mode = "MATCH_MODE_HYBRID"
		classification_threshold = 0.7
		api_version = "API_VERSION_V2"
		tier = "TIER_ENTERPRISE"
		depends_on = [google_project_iam_member.agent_create]
	}
	  `, context)
}
