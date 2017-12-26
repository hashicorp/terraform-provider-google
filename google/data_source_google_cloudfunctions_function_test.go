package google

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
)

func TestAccDataSourceGoogleCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	topicName := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))
	zipFilePath, err := createZIParchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFunctionsFunctionConfig(functionName,
					bucketName, zipFilePath, topicName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"name", fmt.Sprintf("%s-http", functionName)),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_bucket",
						"name", fmt.Sprintf("%s-bucket", functionName)),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_pubsub",
						"name", fmt.Sprintf("%s-pubsub", functionName)),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"description", "test function"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"memory", "128"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"region", "us-central1"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"timeout", "61"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"storage_bucket", bucketName),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"storage_object", "index.zip"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"trigger_http", "true"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_http",
						"entry_point", "helloGET"),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_bucket",
						"trigger_bucket", bucketName),
					resource.TestCheckResourceAttr("data.google_cloudfunctions_function.function_pubsub",
						"trigger_topic", topicName),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFunctionsFunctionConfig(functionName string,
	bucketName string, zipFilePath string, topicName string) string {
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
  name           = "%s-http"
  description    = "test function"
  memory		 = 128
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_http   = true
  timeout		 = 61
  entry_point    = "helloGET"
}

resource "google_cloudfunctions_function" "function_bucket" {
  name           = "%s-bucket"
  memory		 = 128
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_bucket  = "${google_storage_bucket.bucket.name}"
  timeout		 = 61
  entry_point    = "helloGET"
}

resource "google_pubsub_topic" "sub" {
	name = "%s"
}

resource "google_cloudfunctions_function" "function_pubsub" {
  name           = "%s-pubsub"
  memory		 = 128
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_topic  = "${google_pubsub_topic.sub.name}"
  timeout		 = 61
  entry_point    = "helloGET"
}

data "google_cloudfunctions_function" "function_http" {
	name = "${google_cloudfunctions_function.function_http.name}"
}

data "google_cloudfunctions_function" "function_bucket" {
	name = "${google_cloudfunctions_function.function_bucket.name}"
}

data "google_cloudfunctions_function" "function_pubsub" {
	name = "${google_cloudfunctions_function.function_pubsub.name}"
}
`, bucketName, zipFilePath, functionName, functionName,
		topicName, functionName)
}
