package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDialogflowAgent_update(t *testing.T) {
	t.Parallel()

	agentName := acctest.RandomWithPrefix("tf-test")
	agentNameUpdate := acctest.RandomWithPrefix("tf-test")
	projectID := acctest.RandomWithPrefix("tf-test")
	orgID := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowAgent_full1(projectID, orgID, agentName),
			},
			{
				ResourceName:            "google_dialogflow_agent.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"avatar_uri"},
			},
			{
				Config: testAccDialogflowAgent_full2(projectID, orgID, agentNameUpdate),
			},
			{
				ResourceName:            "google_dialogflow_agent.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"avatar_uri"},
			},
		},
	})
}

func testAccDialogflowAgent_full1(projectID string, orgID string, agentName string) string {
	return fmt.Sprintf(`
	resource "google_project" "agent_project" {
		project_id = "%s"
		name       = "%s"
		org_id     = "%s"
	  }

	resource "google_project_service" "agent_project" {
		project = google_project.agent_project.project_id
		service = "dialogflow.googleapis.com"
	  }
	  
	resource "google_project_iam_member" "agent_create" {
		project = google_project_service.agent_project.project
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:service-${google_project.agent_project.number}@gcp-sa-dialogflow.iam.gserviceaccount.com"
		depends_on = [google_project_service.agent_project]
	  }

	resource "google_dialogflow_agent" "foobar" {
		project = "%s"
		display_name = "%s"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://cloud.google.com/_static/images/cloud/icons/favicons/onecloud/super_cloud.png"
		enable_logging = true
		match_mode = "MATCH_MODE_ML_ONLY"
		classification_threshold = 0.3
		api_version = "API_VERSION_V2_BETA_1"
		tier = "TIER_STANDARD"
		depends_on = [google_project_iam_member.agent_create]
	  }
	`, projectID, projectID, orgID, projectID, agentName)
}

func testAccDialogflowAgent_full2(projectID string, orgID string, agentName string) string {
	return fmt.Sprintf(`
	resource "google_project" "agent_project" {
		project_id = "%s"
		name       = "%s"
		org_id     = "%s"
	  }

	  resource "google_project_service" "agent_project" {
		project = google_project.agent_project.project_id
		service = "dialogflow.googleapis.com"
	  }
	  
	resource "google_project_iam_member" "agent_create" {
		project = google_project_service.agent_project.project
		role    = "roles/dialogflow.admin"
		member  = "serviceAccount:service-${google_project.agent_project.number}@gcp-sa-dialogflow.iam.gserviceaccount.com"
		depends_on = [google_project_service.agent_project]
	  }

	resource "google_dialogflow_agent" "foobar" {
		project = "%s"
		display_name = "%s"
		default_language_code = "en"
		supported_language_codes = ["no"]
		time_zone = "Europe/London"
		description = "Description 2!"
		avatar_uri = "https://storage.googleapis.com/gweb-cloudblog-publish/images/f4xvje.max-200x200.PNG"
		enable_logging = false
		match_mode = "MATCH_MODE_HYBRID"
		classification_threshold = 0.7
		api_version = "API_VERSION_V2"
		tier = "TIER_ENTERPRISE"
		depends_on = [google_project_iam_member.agent_create]
	  }
	  `, projectID, projectID, orgID, projectID, agentName)
}
