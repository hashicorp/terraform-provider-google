// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrunv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleCloudRunV2Job_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	name := fmt.Sprintf("tf-test-cloud-run-v2-job-%d", acctest.RandInt(t))
	location := "us-central1"
	id := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, location, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunV2Job_basic(name, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_job.hello", "id", id),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_job.hello", "name", name),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_job.hello", "location", location),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudRunV2Job_basic(name, location string) string {
	return fmt.Sprintf(`
resource "google_cloud_run_v2_job" "hello" {
  name     = "%s"
  location = "%s"

  template {
    template {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  labels = {
    "key" = "value"
  }

  annotations = {
    "key" = "value"
  }

  lifecycle {
    ignore_changes = [
      launch_stage,
    ]
  }
}

data "google_cloud_run_v2_job" "hello" {
  name     = google_cloud_run_v2_job.hello.name
  location = google_cloud_run_v2_job.hello.location
}
`, name, location)
}

func TestAccDataSourceGoogleCloudRunV2Job_bindIAMPermission(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	name := fmt.Sprintf("tf-test-cloud-run-v2-job-%d", acctest.RandInt(t))
	location := "us-central1"
	id := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, location, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunV2Job_bindIAMPermission(name, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_job.hello", "id", id),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_job.hello", "name", name),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_job.hello", "location", location),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudRunV2Job_bindIAMPermission(name, location string) string {
	return fmt.Sprintf(`
resource "google_cloud_run_v2_job" "hello" {
  name     = "%s"
  location = "%s"

  template {
    template {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  lifecycle {
    ignore_changes = [
      launch_stage,
    ]
  }
}

data "google_cloud_run_v2_job" "hello" {
  name     = google_cloud_run_v2_job.hello.name
  location = google_cloud_run_v2_job.hello.location
}

resource "google_service_account" "foo" {
  account_id   = "foo-service-account"
  display_name = "foo-service-account"
}

resource "google_cloud_run_v2_job_iam_binding" "foo_run_invoker" {
  name     = data.google_cloud_run_v2_job.hello.name
  location = data.google_cloud_run_v2_job.hello.location

  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.foo.email}",
  ]
}
`, name, location)
}
