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

func TestAccDataSourceGoogleCloudRunV2Service_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	name := fmt.Sprintf("tf-test-cloud-run-v2-service-%d", acctest.RandInt(t))
	location := "us-central1"
	id := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, location, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunV2Service_basic(name, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_service.hello", "id", id),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_service.hello", "name", name),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_service.hello", "location", location),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudRunV2Service_basic(name, location string) string {
	return fmt.Sprintf(`
resource "google_cloud_run_v2_service" "hello" {
  name     = "%s"
  location = "%s"
  ingress  = "INGRESS_TRAFFIC_ALL"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }
}

data "google_cloud_run_v2_service" "hello" {
  name     = google_cloud_run_v2_service.hello.name
  location = google_cloud_run_v2_service.hello.location
}
`, name, location)
}

func TestAccDataSourceGoogleCloudRunV2Service_bindIAMPermission(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()

	name := fmt.Sprintf("tf-test-cloud-run-v2-service-%d", acctest.RandInt(t))
	location := "us-central1"
	id := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, location, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunV2Service_bindIAMPermission(name, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_service.hello", "id", id),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_service.hello", "name", name),
					resource.TestCheckResourceAttr("data.google_cloud_run_v2_service.hello", "location", location),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudRunV2Service_bindIAMPermission(name, location string) string {
	return fmt.Sprintf(`
resource "google_cloud_run_v2_service" "hello" {
  name     = "%s"
  location = "%s"
  ingress  = "INGRESS_TRAFFIC_ALL"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }
}

data "google_cloud_run_v2_service" "hello" {
  name     = google_cloud_run_v2_service.hello.name
  location = google_cloud_run_v2_service.hello.location
}

resource "google_service_account" "foo" {
  account_id   = "foo-service-account"
  display_name = "foo-service-account"
}

resource "google_cloud_run_v2_service_iam_binding" "foo_run_invoker" {
  name     = data.google_cloud_run_v2_service.hello.name
  location = data.google_cloud_run_v2_service.hello.location

  role = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.foo.email}",
  ]
}
`, name, location)
}
