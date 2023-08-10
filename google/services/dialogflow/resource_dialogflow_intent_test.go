// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflow_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowIntent_basic(t *testing.T) {
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
				Config: testAccDialogflowIntent_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_intent.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDialogflowIntent_update(t *testing.T) {
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
				Config: testAccDialogflowIntent_full1(context),
			},
			{
				ResourceName:      "google_dialogflow_intent.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowIntent_full2(context),
			},
			{
				ResourceName:      "google_dialogflow_intent.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowIntent_basic(context map[string]interface{}) string {
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

	resource "google_dialogflow_agent" "agent" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-agent-%{random_suffix}"
		default_language_code = "en"
		time_zone = "America/New_York"
		depends_on = [google_project_iam_member.agent_create]
	}
	
	resource "google_dialogflow_intent" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-%{random_suffix}"
	}
	`, context)
}

func testAccDialogflowIntent_full1(context map[string]interface{}) string {
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

	resource "google_dialogflow_agent" "agent" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-agent-%{random_suffix}"
		default_language_code = "en"
		time_zone = "America/New_York"
		depends_on = [google_project_iam_member.agent_create]
	}

	resource "google_dialogflow_intent" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-intent-%{random_suffix}"
		webhook_state = "WEBHOOK_STATE_ENABLED"
		priority = 1
		is_fallback = false
		ml_disabled = true
		action = "some_action"
		reset_contexts = true
		input_context_names = ["projects/${google_project.agent_project.project_id}/agent/sessions/-/contexts/some_id"]
		events = ["some_event"]
		default_response_platforms = ["FACEBOOK","SLACK"]
	}
	`, context)
}

func testAccDialogflowIntent_full2(context map[string]interface{}) string {
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

	resource "google_dialogflow_agent" "agent" {
		project = google_project.agent_project.project_id
		display_name = "tf-test-agent-%{random_suffix}"
		default_language_code = "en"
		time_zone = "America/New_York"
		depends_on = [google_project_iam_member.agent_create]
	}

	resource "google_dialogflow_intent" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-intent-%{random_suffix}2"
		webhook_state = "WEBHOOK_STATE_ENABLED_FOR_SLOT_FILLING"
		priority = 2
		is_fallback = false
		ml_disabled = false
		action = "some_other_action"
		reset_contexts = false
		input_context_names = ["projects/${google_project.agent_project.project_id}/agent/sessions/-/contexts/some_other_id"]
		events = ["some_other_event"]
		default_response_platforms = ["SKYPE"]
	}
	`, context)
}
