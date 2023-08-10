// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXPage_update(t *testing.T) {
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
				Config: testAccDialogflowCXPage_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_page.my_page",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXPage_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_page.my_page",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXPage_basic(context map[string]interface{}) string {
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

	resource "google_dialogflow_cx_agent" "agent_page" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
		depends_on = [google_project_iam_member.agent_create]
	}
    
	resource "google_dialogflow_cx_page" "my_page" {
        parent       = google_dialogflow_cx_agent.agent_page.start_flow
        display_name  = "MyPage"
    } 
    `, context)
}

func testAccDialogflowCXPage_full(context map[string]interface{}) string {
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

	resource "google_dialogflow_cx_agent" "agent_page" {
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
    
	resource "google_dialogflow_cx_page" "my_page" {
        parent       = google_dialogflow_cx_agent.agent_page.start_flow
        display_name  = "MyPage"

		entry_fulfillment {
			messages {
				text {
					text = ["Welcome to page"]
				}
			}
		}

		form {
			parameters {
				display_name = "param1"
				entity_type  = "projects/-/locations/-/agents/-/entityTypes/sys.date"
				fill_behavior {
					initial_prompt_fulfillment {
						messages {
							text {
								text = ["Please provide param1"]
							}
						}
					}
				}
				required = "true"
				redact   = "true"
			}
		}

		transition_routes {
			condition = "$page.params.status = 'FINAL'"
			trigger_fulfillment {
				messages {
					text {
						text = ["information completed, navigating to page 2"]
					}
				}
			}
			target_page = google_dialogflow_cx_page.my_page2.id
		}
    } 

	resource "google_dialogflow_cx_page" "my_page2" {
        parent       = google_dialogflow_cx_agent.agent_page.start_flow
        display_name  = "MyPage2"
    } 
	  `, context)
}
