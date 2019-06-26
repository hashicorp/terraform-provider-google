package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	funcDataNameHttp := "data.google_cloudfunctions_function.function_http"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIPArchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFunctionsFunctionConfig(functionName,
					bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState(funcDataNameHttp,
						"google_cloudfunctions_function.function_http"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFunctionsFunctionConfig(functionName, bucketName, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function_http" {
  name                  = "%s-http"
  description           = "test function"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
}

data "google_cloudfunctions_function" "function_http" {
  name = "${google_cloudfunctions_function.function_http.name}"
}
`, bucketName, zipFilePath, functionName)
}
