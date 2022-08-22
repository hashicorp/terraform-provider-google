package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowCXEntityType_update(t *testing.T) {
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
				Config: testAccDialogflowCXEntityType_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_entity_type.my_entity",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXEntityType_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_entity_type.my_entity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXEntityType_basic(context map[string]interface{}) string {
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
    
	resource "google_dialogflow_cx_entity_type" "my_entity" {
        parent       = google_dialogflow_cx_agent.agent_entity.id
        display_name  = "MyEntity"
        kind         = "KIND_MAP"
        entities {
            value = "value1"
            synonyms = ["synonym1","synonym2"]
        }
        entities {
            value = "value2"
            synonyms = ["synonym3","synonym4"]
        }
        enable_fuzzy_extraction = false
    } 
    `, context)
}

func testAccDialogflowCXEntityType_full(context map[string]interface{}) string {
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
    
	resource "google_dialogflow_cx_entity_type" "my_entity" {
        parent           = google_dialogflow_cx_agent.agent_entity.id
        display_name     = "MyEntity"
        kind             = "KIND_MAP"
        entities {
            value = "value1"
            synonyms = ["synonym1","synonym2","synonym11","synonym22"]
        }
        entities {
            value = "value2"
            synonyms = ["synonym3","synonym4"]
        }
        enable_fuzzy_extraction = false
        redact                  = true
        auto_expansion_mode     = "AUTO_EXPANSION_MODE_DEFAULT"
        excluded_phrases {
			value = "excluded1"
        }
		
    } 
	  `, context)
}
