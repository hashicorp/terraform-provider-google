// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXWebhook_update(t *testing.T) {
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
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
