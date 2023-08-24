// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflow_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowEntityType_update(t *testing.T) {
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
				Config: testAccDialogflowEntityType_full1(context),
			},
			{
				ResourceName:      "google_dialogflow_entity_type.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowEntityType_full2(context),
			},
			{
				ResourceName:      "google_dialogflow_entity_type.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowEntityType_full1(context map[string]interface{}) string {
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
	
	resource "google_dialogflow_entity_type" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-entity-%{random_suffix}"
		kind = "KIND_MAP"
		enable_fuzzy_extraction = true
		entities {
		  value = "value1"
		  synonyms = ["synonym1","synonym2"]
		}
		entities {
		  value = "value2"
		  synonyms = ["synonym3","synonym4"]
		}
	}
	`, context)
}

func testAccDialogflowEntityType_full2(context map[string]interface{}) string {
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
	
	resource "google_dialogflow_entity_type" "foobar" {
		depends_on = [google_dialogflow_agent.agent]
		project = google_project.agent_project.project_id
		display_name = "tf-test-entity2-%{random_suffix}"
		kind = "KIND_LIST"
		enable_fuzzy_extraction = false
		entities {
		  value = "value1"
		  synonyms = ["value1"]
		}
	}
	`, context)
}
