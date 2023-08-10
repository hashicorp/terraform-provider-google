// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDialogflowCXEnvironment_update(t *testing.T) {
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

func TestAccDialogflowCXEnvironment_dialogflowcxEnvironmentFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDialogflowCXEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXEnvironment_dialogflowcxEnvironmentFullExample(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_environment.development",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccDialogflowCXEnvironment_dialogflowcxEnvironmentFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name = "tf-test-dialogflowcx-agent%{random_suffix}"
  location = "global"
  default_language_code = "en"
  supported_language_codes = ["fr","de","es"]
  time_zone = "America/New_York"
  description = "Example description."
  avatar_uri = "https://cloud.google.com/_static/images/cloud/icons/favicons/onecloud/super_cloud.png"
  enable_stackdriver_logging = true
  enable_spell_correction    = true
	speech_to_text_settings {
		enable_speech_adaptation = true
	}
}

resource "google_dialogflow_cx_version" "version_1" {
  parent       = google_dialogflow_cx_agent.agent.start_flow
  display_name = "1.0.0"
  description  = "version 1.0.0"
}

resource "google_dialogflow_cx_environment" "development" {
  parent       = google_dialogflow_cx_agent.agent.id
  display_name = "Development"
  description  = "Development Environment"
  version_configs {
    version = google_dialogflow_cx_version.version_1.id
  }
}
`, context)
}

func TestAccDialogflowCXEnvironment_dialogflowcxEnvironmentRegional(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDialogflowCXEnvironmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXEnvironment_dialogflowcxEnvironmentFRegional(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_environment.development",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccDialogflowCXEnvironment_dialogflowcxEnvironmentFRegional(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
	display_name = "issue12880"
	location = "europe-west2"
	default_language_code = "en"
	supported_language_codes = ["fr","de","es"]
	time_zone = "Europe/London"
	description = "CX BOT Agent"
	enable_stackdriver_logging = true
		speech_to_text_settings {
			enable_speech_adaptation = true
		}
	}

resource "google_dialogflow_cx_version" "version_1" {
	parent       = google_dialogflow_cx_agent.agent.start_flow
	display_name = "1.0.0"
	description  = "version 1.0.0"
}

resource "google_dialogflow_cx_environment" "development" {
	parent       = google_dialogflow_cx_agent.agent.id
	display_name = "Development"
	description  = "Development Environment"
	version_configs {
		version = google_dialogflow_cx_version.version_1.id
	}
}
`, context)
}

func testAccCheckDialogflowCXEnvironmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dialogflow_cx_environment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DialogflowCXBasePath}}{{parent}}/environments/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DialogflowCXEnvironment still exists at %s", url)
			}
		}

		return nil
	}
}
