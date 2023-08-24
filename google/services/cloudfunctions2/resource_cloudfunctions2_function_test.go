// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudfunctions2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCloudFunctions2Function_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"zip_path":      "./test-fixtures/function-source.zip",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudfunctions2functionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudfunctions2function_basic(context),
			},
			{
				ResourceName:            "google_cloudfunctions2_function.terraform-test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "build_config.0.source.0.storage_source.0.object", "build_config.0.source.0.storage_source.0.bucket"},
			},
			{
				Config: testAccCloudFunctions2Function_test_update(context),
			},
			{
				ResourceName:            "google_cloudfunctions2_function.terraform-test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "build_config.0.source.0.storage_source.0.object", "build_config.0.source.0.storage_source.0.bucket"},
			},
			{
				Config: testAccCloudFunctions2Function_test_redeploy(context),
			},
			{
				ResourceName:            "google_cloudfunctions2_function.terraform-test2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "build_config.0.source.0.storage_source.0.object", "build_config.0.source.0.storage_source.0.bucket"},
			},
		},
	})
}

func testAccCloudfunctions2function_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-cloudfunctions2-function-bucket%{random_suffix}"
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%{zip_path}"
}

resource "google_cloudfunctions2_function" "terraform-test2" {
  name = "tf-test-test-function%{random_suffix}"
  location = "us-central1"
  description = "a new function"

  build_config {
    runtime = "nodejs12"
    entry_point = "helloHttp"
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count  = 1
    available_memory    = "1536Mi"
    timeout_seconds     = 30
  }
}
`, context)
}

func testAccCloudFunctions2Function_test_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-cloudfunctions2-function-bucket%{random_suffix}"
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%{zip_path}"
}

resource "google_cloudfunctions2_function" "terraform-test2" {
  name = "tf-test-test-function%{random_suffix}"
  location = "us-central1"
  description = "an updated function"

  build_config {
    runtime = "nodejs12"
    entry_point = "helloHttp"
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    min_instance_count = 1
  }
}
`, context)
}

func testAccCloudFunctions2Function_test_redeploy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "tf-test-cloudfunctions2-function-bucket%{random_suffix}"
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%{zip_path}"
}

resource "google_cloudfunctions2_function" "terraform-test2" {
  name = "tf-test-test-function%{random_suffix}"
  location = "us-west1"
  description = "function test"

  build_config {
    runtime = "nodejs16"
    entry_point = "helloHttp"
    environment_variables = {
        BUILD_CONFIG_TEST = "build_test"
    }
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count  = 5
    min_instance_count = 1
    available_memory    = "256M"
    timeout_seconds     = 60
    environment_variables = {
        SERVICE_CONFIG_TEST = "build_test"
    }
  }
}
`, context)
}

func TestAccCloudFunctions2Function_fullUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"zip_path":      "./test-fixtures/function-source-eventarc-gcs.zip",
		"random_suffix": acctest.RandString(t, 10),
	}

	if acctest.BootstrapPSARole(t, "service-", "gcp-sa-pubsub", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a binding was added.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudfunctions2functionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Re-use config from the generated tests
				Config: testAccCloudfunctions2function_cloudfunctions2BasicAuditlogsExample(context),
			},
			{
				Config: testAccCloudfunctions2function_cloudfunctions2BasicAuditlogsExample_update(context),
			},
		},
	})
}

func testAccCloudfunctions2function_cloudfunctions2BasicAuditlogsExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
# This example follows the examples shown in this Google Cloud Community blog post
# https://medium.com/google-cloud/applying-a-path-pattern-when-filtering-in-eventarc-f06b937b4c34
# and the docs:
# https://cloud.google.com/eventarc/docs/path-patterns

resource "google_storage_bucket" "source-bucket" {
  name     = "tf-test-gcf-source-bucket%{random_suffix}"
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.source-bucket.name
  source = "%{zip_path}"  # Add path to the zipped function source code
}

resource "google_service_account" "account" {
  account_id   = "tf-test-gcf-sa%{random_suffix}"
  display_name = "Test Service Account - used for both the cloud function and eventarc trigger in the test"
}

# Note: The right way of listening for Cloud Storage events is to use a Cloud Storage trigger.
# Here we use Audit Logs to monitor the bucket so path patterns can be used in the example of
# google_cloudfunctions2_function below (Audit Log events have path pattern support)
resource "google_storage_bucket" "audit-log-bucket" {
  name     = "tf-test-gcf-auditlog-bucket%{random_suffix}"
  location = "us-central1"  # The trigger must be in the same location as the bucket
  uniform_bucket_level_access = true
}

# Permissions on the service account used by the function and Eventarc trigger
resource "google_project_iam_member" "invoking" {
  project = "%{project}"
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.account.email}"
}

resource "google_project_iam_member" "event-receiving" {
  project = "%{project}"
  role    = "roles/eventarc.eventReceiver"
  member  = "serviceAccount:${google_service_account.account.email}"
  depends_on = [google_project_iam_member.invoking]
}

resource "google_project_iam_member" "artifactregistry-reader" {
  project = "%{project}"
  role     = "roles/artifactregistry.reader"
  member   = "serviceAccount:${google_service_account.account.email}"
  depends_on = [google_project_iam_member.event-receiving]
}

resource "google_cloudfunctions2_function" "function" {
  depends_on = [
    google_project_iam_member.event-receiving,
    google_project_iam_member.artifactregistry-reader,
  ]
  name = "tf-test-gcf-function%{random_suffix}"
  location = "us-central1"
  description = "a new function"

  build_config {
    runtime     = "nodejs12"
    entry_point = "entryPoint" # Set the entry point in the code
    environment_variables = {
      BUILD_CONFIG_TEST = "build_test"
    }
    source {
      storage_source {
        bucket = google_storage_bucket.source-bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }

  service_config {
    max_instance_count  = 3
    min_instance_count = 1
    available_memory    = "256M"
    timeout_seconds     = 60
    environment_variables = {
        SERVICE_CONFIG_TEST = "config_test"
    }
    ingress_settings = "ALLOW_INTERNAL_ONLY"
    all_traffic_on_latest_revision = true
    service_account_email = google_service_account.account.email
  }

  event_trigger {
    trigger_region = "us-central1" # The trigger must be in the same location as the bucket
    event_type = "google.cloud.audit.log.v1.written"
    retry_policy = "RETRY_POLICY_RETRY"
    service_account_email = google_service_account.account.email
    event_filters {
      attribute = "serviceName"
      value = "storage.googleapis.com"
    }
    event_filters {
      attribute = "methodName"
      value = "storage.objects.get" # Update: change value
    }
    event_filters {
      attribute = "resourceName"
      value = google_storage_bucket.audit-log-bucket.name # Update: stops using path pattern operator
    }
  }
}`, context)
}
