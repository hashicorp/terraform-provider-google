package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDialogflowCXVersion_update(t *testing.T) {
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
    `, context)
}

func testAccDialogflowCXVersion_full(context map[string]interface{}) string {
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
	resource "google_dialogflow_cx_version" "version1" {
		parent       = google_dialogflow_cx_agent.agent_version.start_flow
		display_name = "1.0.0"
		description  = "version 1.0.0"
	}
	  `, context)
}

func TestAccDialogflowCXVersion_dialogflowcxVersionFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckDialogflowCXVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXVersion_dialogflowcxVersionFullExample(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_version.version_1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccDialogflowCXVersion_dialogflowcxVersionFullExample(context map[string]interface{}) string {
	return Nprintf(`
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
`, context)
}

func testAccCheckDialogflowCXVersionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dialogflow_cx_version" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{DialogflowCXBasePath}}{{parent}}/versions/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("DialogflowCXVersion still exists at %s", url)
			}
		}

		return nil
	}
}
