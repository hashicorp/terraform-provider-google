// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudfunctions2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleCloudFunctions2Function_basic(t *testing.T) {
	t.Parallel()

	funcDataNameHttp := "data.google_cloudfunctions2_function.function_http_v2"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt(t))
	zipFilePath := "./test-fixtures/function-source.zip"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudfunctions2functionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFunctions2FunctionConfig(functionName,
					bucketName, zipFilePath),
				// As the value of "labels" and "terraform_labels" in the state of the data source are all labels,
				// but the "labels" field in resource are user defined labels, which is the reason for the mismatch.
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(funcDataNameHttp,
						"google_cloudfunctions2_function.function_http_v2", map[string]struct{}{"build_config.0.source.0.storage_source.0.bucket": {}, "build_config.0.source.0.storage_source.0.object": {}, "labels.%": {}, "terraform_labels.%": {}}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFunctions2FunctionConfig(functionName, bucketName, zipFilePath string) string {
	return fmt.Sprintf(`
provider "google" {
  default_labels = {
    default_key1 = "default_value1"
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}
	
resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions2_function" "function_http_v2" {
  name = "%s"
  location = "us-central1"
  description = "a new function"
  labels = {
    env = "test"
  }
  build_config {
    runtime = "nodejs18"
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
    available_memory    = "256Mi"
    timeout_seconds     = 60
  }
}
data "google_cloudfunctions2_function" "function_http_v2" {
  name = google_cloudfunctions2_function.function_http_v2.name
  location = "us-central1"
}
`, bucketName, zipFilePath, functionName)
}
