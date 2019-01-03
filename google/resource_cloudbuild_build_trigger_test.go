package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCloudBuildTrigger_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudbuildTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleCloudBuildTrigger_basic(),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleCloudBuildTrigger_updated(),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_filename(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudbuildTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleCloudBuildTrigger_filename(),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.filename_build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testGoogleCloudBuildTrigger_basic() string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags = ["team-a", "service-b"]
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
    }
  }
}
  `)
}

func testGoogleCloudBuildTrigger_updated() string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  description = "acceptance test build trigger updated"
  trigger_template {
    branch_name = "master-updated"
    repo_name   = "some-repo-updated"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"]
    tags = ["team-a", "service-b", "updated"]
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile-updated.zip"]
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package_updated"]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA", "-f", "Dockerfile", "."]
    }
    step {
      name = "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"
      args = ["test"]
    }
  }
}
  `)
}

func testGoogleCloudBuildTrigger_filename() string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "filename_build_trigger" {
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  substitutions = {
    _FOO = "bar"
    _BAZ = "qux"
  }
  filename = "cloudbuild.yaml"
}
  `)
}
