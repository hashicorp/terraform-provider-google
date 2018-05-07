package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudBuildTrigger_basic(t *testing.T) {
	t.Parallel()

	projectID := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleCloudBuildTriggerVersionsDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleCloudBuildTrigger_basic(projectID, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleCloudBuildTriggerExists("google_cloudbuild_trigger.build_trigger"),
				),
			},
			resource.TestStep{
				Config: testGoogleCloudBuildTrigger_removed(projectID, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleCloudBuildTriggerWasRemovedFromState("google_cloudbuild_trigger.build_trigger"),
				),
			},
		},
	})
}

func testAccCheckGoogleCloudBuildTriggerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		project := rs.Primary.Attributes["project"]

		_, err := config.clientBuild.Projects.Triggers.Get(project, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Trigger does not exist")
		}
		return nil
	}
}

func testAccCheckGoogleCloudBuildTriggerWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

func testAccCheckGoogleCloudBuildTriggerVersionsDestroyed(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudbuild_trigger" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		project := rs.Primary.Attributes["project"]

		_, err := config.clientBuild.Projects.Triggers.Get(project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Trigger still exists")
		}

	}

	return nil
}

/*
  This test runs in its own project, otherwise the test project would start to get filled
  with undeletable resources
*/
func testGoogleCloudBuildTrigger_basic(projectID, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudbuild.googleapis.com",
    "containerregistry.googleapis.com",
    "logging.googleapis.com",
    "pubsub.googleapis.com",
    "storage-api.googleapis.com",
  ]
}

resource "google_cloudbuild_trigger" "build_trigger" {
  project  = "${google_project_services.acceptance.project}"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    project     = "${google_project_services.acceptance.project}"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags = ["team-a", "service-b"]
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = "cp gs://mybucket/remotefile.zip localfile.zip "
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = "build my_package"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = "build -t gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA -f Dockerfile ."
    }
  }
}
  `, projectID, projectID, projectOrg, projectBillingAccount)
}

func testGoogleCloudBuildTrigger_removed(projectID, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudbuild.googleapis.com",
    "containerregistry.googleapis.com",
    "logging.googleapis.com",
    "pubsub.googleapis.com",
    "storage-api.googleapis.com",
  ]
}
  `, projectID, projectID, projectOrg, projectBillingAccount)
}
