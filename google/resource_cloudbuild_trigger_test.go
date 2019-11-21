package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCloudBuildTrigger_basic(t *testing.T) {
	t.Parallel()
	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_basic(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_updated(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_disable(t *testing.T) {
	t.Parallel()
	name := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_basic(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudBuildTrigger_basicDisabled(name),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudBuildTrigger_fullStep(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildTrigger_fullStep(),
			},
			{
				ResourceName:      "google_cloudbuild_trigger.build_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudBuildTrigger_basic(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name        = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags   = ["team-a", "service-b"]
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env  = ["env1=two"]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
    }
  }
}
`, name)
}

func testAccCloudBuildTrigger_basicDisabled(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  disabled    = true
  name        = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags   = ["team-a", "service-b"]
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env  = ["env1=two"]
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
    }
  }
}
`, name)
}

func testAccCloudBuildTrigger_fullStep() string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags   = ["team-a", "service-b"]
    step {
      name       = "gcr.io/cloud-builders/go"
      args       = ["build", "my_package"]
      env        = ["env1=two"]
      dir        = "directory"
      id         = "12345"
      secret_env = ["fooo"]
      timeout    = "100s"
      wait_for   = ["something"]
    }
  }
}
`)
}

func testAccCloudBuildTrigger_updated(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  description = "acceptance test build trigger updated"
  name        = "%s"
  trigger_template {
    branch_name = "master-updated"
    repo_name   = "some-repo-updated"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"]
    tags   = ["team-a", "service-b", "updated"]
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
`, name)
}
