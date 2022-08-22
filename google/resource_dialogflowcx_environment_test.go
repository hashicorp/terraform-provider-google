package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowCXEnvironment_update(t *testing.T) {
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
				Config: testAccDialogflowCXEnvironment_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_environment.development",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXEnvironment_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_environment.development",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXEnvironment_basic(context map[string]interface{}) string {
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

	resource "google_dialogflow_cx_agent" "agent_version" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
		depends_on = [google_project_iam_member.agent_create]
	}
    
	resource "google_dialogflow_cx_version" "version1" {
		parent       = google_dialogflow_cx_agent.agent_version.start_flow
		display_name = "1.0.0"
		description  = "version 1.0.0"
	}	

	resource "google_dialogflow_cx_environment" "development" {
        parent       = google_dialogflow_cx_agent.agent_version.id
        display_name = "Development"
        version_configs {
            version = google_dialogflow_cx_version.version1.id
        }
    }
    `, context)
}

func testAccDialogflowCXEnvironment_full(context map[string]interface{}) string {
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

	resource "google_dialogflow_cx_agent" "agent_version" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
		depends_on = [google_project_iam_member.agent_create]
	}

	resource "google_dialogflow_cx_version" "version1" {
		parent       = google_dialogflow_cx_agent.agent_version.start_flow
		display_name = "1.0.0"
		description  = "version 1.0.0"
	}

	resource "google_dialogflow_cx_version" "version2" {
		parent       = google_dialogflow_cx_agent.agent_version.start_flow
		display_name = "2.0.0"
		description  = "version 2.0.0"
	}

	resource "google_dialogflow_cx_environment" "development" {
        parent       = google_dialogflow_cx_agent.agent_version.id
        display_name = "Development"
        version_configs {
            version = google_dialogflow_cx_version.version2.id
        }
    }
	  `, context)
}
