// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXVersion_update(t *testing.T) {
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
				Config: testAccDialogflowCXVersion_basic(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_version.version1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXVersion_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_version.version1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXVersion_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dialogflow_cx_agent" "agent_version" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
	}
    
	resource "google_dialogflow_cx_version" "version1" {
		parent       = google_dialogflow_cx_agent.agent_version.start_flow
		display_name = "1.0.0"
		description  = "version 1.0.0"
	}
    `, context)
}

func testAccDialogflowCXVersion_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dialogflow_cx_agent" "agent_version" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
	}

	resource "google_dialogflow_cx_version" "version1" {
		parent       = google_dialogflow_cx_agent.agent_version.start_flow
		display_name = "1.0.0 updated"
		description  = "version 1.0.0 updated"
	}
	  `, context)
}
