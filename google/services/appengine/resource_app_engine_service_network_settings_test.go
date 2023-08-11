// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package appengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAppEngineServiceNetworkSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineServiceNetworkSettings_basic(context),
			},
			{
				ResourceName:      "google_app_engine_service_network_settings.main",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppEngineServiceNetworkSettings_update(context),
			},
			{
				ResourceName:      "google_app_engine_service_network_settings.main",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppEngineServiceNetworkSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-%{random_suffix}-ae-networksettings"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name = "hello-world.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world.zip"
}

resource "google_app_engine_standard_app_version" "app" {
  version_id = "v1"
  service = "app-%{random_suffix}"
  delete_service_on_destroy = true

  runtime = "nodejs10"
  entrypoint {
    shell = "node ./app.js"
  }
  deployment {
    zip {
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
    }
  }
  env_variables = {
    port = "8080"
  }
}

resource "google_app_engine_service_network_settings" "main" {
  service = google_app_engine_standard_app_version.app.service
  network_settings {
    ingress_traffic_allowed = "INGRESS_TRAFFIC_ALLOWED_ALL"
  }
}`, context)
}

func testAccAppEngineServiceNetworkSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-%{random_suffix}-ae-networksettings"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name = "hello-world.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/hello-world.zip"
}

resource "google_app_engine_standard_app_version" "app" {
  version_id = "v1"
  service = "app-%{random_suffix}"
  delete_service_on_destroy = true

  runtime = "nodejs10"
  entrypoint {
    shell = "node ./app.js"
  }
  deployment {
    zip {
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
    }
  }
  env_variables = {
    port = "8080"
  }
}

resource "google_app_engine_service_network_settings" "main" {
  service = google_app_engine_standard_app_version.app.service
  network_settings {
    ingress_traffic_allowed = "INGRESS_TRAFFIC_ALLOWED_INTERNAL_ONLY"
  }
}`, context)
}
