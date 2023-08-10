// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXIntent_update(t *testing.T) {
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
				Config: testAccDialogflowCXIntent_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_intent.my_intent",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXIntent_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_intent.my_intent",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXIntent_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	data "google_project" "project" {}

	resource "google_service_account" "dialogflowcx_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}

	resource "google_project_iam_member" "agent_create" {
		project = data.google_project.project.project_id
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflowcx_service_account.email}"
	}

	resource "google_dialogflow_cx_agent" "agent_intent" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
		depends_on = [google_project_iam_member.agent_create]
	}
    
	resource "google_dialogflow_cx_intent" "my_intent" {
        parent       = google_dialogflow_cx_agent.agent_intent.id
        display_name = "Example"
        priority     = 1
        description  = "Intent example"
        training_phrases {
            parts {
                text = "training"
            }

            parts {
                text = "phrase"
            }

            parts {
                text = "example"
            }

			repeat_count = 1
        }

        parameters {
            id          = "param1"
            entity_type = "projects/-/locations/-/agents/-/entityTypes/sys.date"
        }

        labels  = {
            label1 = "value1",
            label2 = "value2"
        } 
    } 
    `, context)
}

func testAccDialogflowCXIntent_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
	data "google_project" "project" {}

	resource "google_service_account" "dialogflowcx_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}

	resource "google_project_iam_member" "agent_create" {
		project = data.google_project.project.project_id
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflowcx_service_account.email}"
	}

	resource "google_dialogflow_cx_agent" "agent_intent" {
		display_name = "tf-test-%{random_suffix}update"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["no"]
		time_zone = "Europe/London"
		description = "Description 2!"
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo-2.png"
		enable_stackdriver_logging = true
        enable_spell_correction    = true
		speech_to_text_settings {
			enable_speech_adaptation = true
		}
		depends_on = [google_project_iam_member.agent_create]
	}
    
	resource "google_dialogflow_cx_intent" "my_intent" {
        parent       = google_dialogflow_cx_agent.agent_intent.id
        display_name = "Example"
        priority     = 1
        description  = "Intent example"
        training_phrases {
            parts {
                text = "training"
            }

            parts {
                text = "phrase"
            }

            parts {
                text = "example"
            }

			repeat_count = 1
        }

        parameters {
            id          = "param1"
            entity_type = "projects/-/locations/-/agents/-/entityTypes/sys.date"
        }

        labels  = {
            label1 = "value1",
            label2 = "value2"
        } 
    } 
	  `, context)
}
