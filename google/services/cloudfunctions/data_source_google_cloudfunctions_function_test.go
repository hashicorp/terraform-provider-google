// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudfunctions_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	funcDataNameHttp := "data.google_cloudfunctions_function.function_http"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt(t))
	zipFilePath := acctest.CreateZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFunctionsFunctionConfig(functionName,
					bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(funcDataNameHttp,
						"google_cloudfunctions_function.function_http"),
					resource.TestCheckResourceAttr(funcDataNameHttp,
						"status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFunctionsFunctionConfig(functionName, bucketName, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function_http" {
  name                  = "%s-http"
  runtime               = "nodejs14"
  description           = "test function"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
}

data "google_cloudfunctions_function" "function_http" {
  name = google_cloudfunctions_function.function_http.name
}
`, bucketName, zipFilePath, functionName)
}
