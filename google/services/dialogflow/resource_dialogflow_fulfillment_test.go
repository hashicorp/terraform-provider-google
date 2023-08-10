// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflow_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowFulfillment_update(t *testing.T) {
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
				Config: testAccDialogflowFulfillment_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_fulfillment.agent_fulfillment",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowEntityType_full(context),
			},
			{
				ResourceName:      "google_dialogflow_fulfillment.agent_fulfillment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowFulfillment_basic(context map[string]interface{}) string {
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
	
	resource "google_dialogflow_fulfillment" "agent_fulfillment" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-fulfillment-%{random_suffix}"
		enabled  = true
	}
	`, context)
}

func testAccDialogflowEntityType_full(context map[string]interface{}) string {
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
	
	resource "google_dialogflow_fulfillment" "agent_fulfillment" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-entity2-%{random_suffix}"
		enabled = true
		generic_web_service {
			uri      = "https://google.com"
			username = "admin"
			password = "password"
			request_headers = { 
                 "name" = "wrench"
			}
		}
	}
	`, context)
}
