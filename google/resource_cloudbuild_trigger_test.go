package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudBuildTrigger_basic(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroyProducer(t),
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

func TestAccCloudBuildTrigger_customizeDiffTimeoutSum(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudBuildTrigger_customizeDiffTimeoutSum(name),
				ExpectError: regexp.MustCompile("cannot be greater than build timeout"),
			},
		},
	})
}

func TestAccCloudBuildTrigger_customizeDiffTimeoutFormat(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudBuildTrigger_customizeDiffTimeoutFormat(name),
				ExpectError: regexp.MustCompile("Error parsing build timeout"),
			},
		},
	})
}

func TestAccCloudBuildTrigger_disable(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroyProducer(t),
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

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudBuildTriggerDestroyProducer(t),
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
    tags   = ["team-a", "service-b"]
    timeout = "1800s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env  = ["env1=two"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
      timeout = "300s"
    }
    artifacts {
      images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
      objects {
        location = "gs://bucket/path/to/somewhere/"
        paths = ["path"]
      }
    }
    options {
      source_provenance_hash = ["MD5"]
      requested_verify_option = "VERIFIED"
      machine_type = "N1_HIGHCPU_8"
      disk_size_gb = 100
      substitution_option = "ALLOW_LOOSE"
      dynamic_substitutions = false
      log_streaming_option = "STREAM_OFF"
      worker_pool = "pool"
      logging = "LEGACY"
      env = ["ekey = evalue"]
      secret_env = ["secretenv = svalue"]
      volumes {
        name = "v1"
        path = "v1"
      }
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
	invert_regex = false
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
	invert_regex = true
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"]
    tags   = ["team-a", "service-b", "updated"]
    timeout = "2100s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile-updated.zip"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package_updated"]
      timeout = "300s"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA", "-f", "Dockerfile", "."]
      timeout = "300s"
    }
    step {
      name = "gcr.io/$PROJECT_ID/$REPO_NAME:$SHORT_SHA"
      args = ["test"]
      timeout = "300s"
    }
    logs_bucket = "gs://mybucket/logs"
    options {
      # this field is always enabled for triggered build and cannot be overridden in the build configuration file.
      dynamic_substitutions = true
    }
  }
}
  `, name)
}

func testAccCloudBuildTrigger_customizeDiffTimeoutSum(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags = ["team-a", "service-b"]
    timeout = "900s"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "500s"
    }
    step {
      name = "gcr.io/cloud-builders/go"
      args = ["build", "my_package"]
      env = ["env1=two"]
      timeout = "500s"
    }
    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", "gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA", "-f", "Dockerfile", "."]
      timeout = "500s"
    }
  }
}
  `, name)
}

func testAccCloudBuildTrigger_customizeDiffTimeoutFormat(name string) string {
	return fmt.Sprintf(`
resource "google_cloudbuild_trigger" "build_trigger" {
  name = "%s"
  description = "acceptance test build trigger"
  trigger_template {
    branch_name = "master"
    repo_name   = "some-repo"
  }
  build {
    images = ["gcr.io/$PROJECT_ID/$REPO_NAME:$COMMIT_SHA"]
    tags = ["team-a", "service-b"]
    timeout = "1200"
    step {
      name = "gcr.io/cloud-builders/gsutil"
      args = ["cp", "gs://mybucket/remotefile.zip", "localfile.zip"]
      timeout = "500s"
    }
  }
}
`, name)
}
