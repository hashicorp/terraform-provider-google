package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowCXWebhook_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXWebhook_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_webhook.my_webhook",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXWebhook_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_webhook.my_webhook",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXWebhook_basic(context map[string]interface{}) string {
	return Nprintf(`
	data "google_project" "project" {}

	resource "google_dialogflow_cx_agent" "agent_entity" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["it","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
	}

	resource "google_dialogflow_cx_webhook" "my_webhook" {
		parent       = google_dialogflow_cx_agent.agent_entity.id
		display_name = "MyWebhook"
		generic_web_service {
			uri = "https://example.com"
		}
	}
	`, context)
}

func testAccDialogflowCXWebhook_full(context map[string]interface{}) string {
	return Nprintf(`
	data "google_project" "project" {}

	resource "google_dialogflow_cx_agent" "agent_entity" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["it","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
	}

	resource "google_dialogflow_cx_webhook" "my_webhook" {
		parent       = google_dialogflow_cx_agent.agent_entity.id
		display_name = "MyWebhook"
		timeout      = "20s"
		disabled     = false
		generic_web_service {
			uri = "https://example.com"
			request_headers = {
				"Authorization": "Bearer {{access_token}}"
			}
		}
	}
	`, context)
}
