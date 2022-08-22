package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowCXFlow_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
		"random_suffix":   randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXFlow_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_flow.my_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXFlow_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_flow.my_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXFlow_basic(context map[string]interface{}) string {
	return Nprintf(`
	data "google_project" "project" {}

	resource "google_service_account" "dialogflowcx_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}

	resource "google_project_iam_member" "agent_create" {
		project = data.google_project.project.project_id
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflowcx_service_account.email}"
	}

	resource "google_dialogflow_cx_agent" "agent_entity" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
		depends_on = [google_project_iam_member.agent_create]
	}
    
	resource "google_dialogflow_cx_flow" "my_flow" {
        parent       = google_dialogflow_cx_agent.agent_entity.id
        display_name = "MyFlow"

        nlu_settings {
           classification_threshold = 0.3 
           model_type               = "MODEL_TYPE_STANDARD"
	    }
    } 
    `, context)
}

func testAccDialogflowCXFlow_full(context map[string]interface{}) string {
	return Nprintf(`
	data "google_project" "project" {}

	resource "google_service_account" "dialogflowcx_service_account" {
		account_id = "tf-test-dialogflow-%{random_suffix}"
	}

	resource "google_project_iam_member" "agent_create" {
		project = data.google_project.project.project_id
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:${google_service_account.dialogflowcx_service_account.email}"
	}

	resource "google_dialogflow_cx_agent" "agent_entity" {
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
    
	resource "google_dialogflow_cx_flow" "my_flow" {
        parent       = google_dialogflow_cx_agent.agent_entity.id
        display_name = "MyFlow"

        nlu_settings {
           classification_threshold = 0.3 
           model_type               = "MODEL_TYPE_STANDARD"
	    }

        event_handlers {
		   event                    = "custom-event"
		   trigger_fulfillment {
			    return_partial_responses = false
				messages {
					text {
						text  = ["I didn't get that. Can you say it again?"]
					}
				}
		    }
		}

		event_handlers {
			event                    = "sys.no-match-default"
			trigger_fulfillment {
				 return_partial_responses = false
				 messages {
					 text {
						 text  = ["Sorry, could you say that again?"]
					 }
				 }
			 }
		 }

		 event_handlers {
			event                    = "sys.no-input-default"
			trigger_fulfillment {
				 return_partial_responses = false
				 messages {
					 text {
						 text  = ["One more time?"]
					 }
				 }
			 }
		 }
    } 
	  `, context)
}
