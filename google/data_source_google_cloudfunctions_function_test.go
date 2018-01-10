package google

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	funcDataNameHttp := "data.google_cloudfunctions_function.function_http"
	funcDataNamePubSub := "data.google_cloudfunctions_function.function_pubsub"
	funcDataNameBucket := "data.google_cloudfunctions_function.function_bucket"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	topicName := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))
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
					bucketName, zipFilePath, topicName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleCloudFunctionsFunctionCheck(funcDataNameHttp,
						"google_cloudfunctions_function.function_http"),
					testAccDataSourceGoogleCloudFunctionsFunctionCheck(funcDataNamePubSub,
						"google_cloudfunctions_function.function_pubsub"),
					testAccDataSourceGoogleCloudFunctionsFunctionCheck(funcDataNameBucket,
						"google_cloudfunctions_function.function_bucket"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFunctionsFunctionCheck(dataSourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		cloudFuncAttrToCheck := []string{
			"name",
			"region",
			"description",
			"available_memory_mb",
			"timeout",
			"storage_bucket",
			"storage_object",
			"entry_point",
			"trigger_http",
			"trigger_bucket",
			"trigger_topic",
		}

		for _, attr := range cloudFuncAttrToCheck {
			if dsAttr[attr] != rsAttr[attr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr,
					dsAttr[attr],
					rsAttr[attr],
				)
			}
		}

		return nil
	}
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
  name                  = "%s-http"
  description           = "test function"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
}

resource "google_cloudfunctions_function" "function_bucket" {
  name                  = "%s-bucket"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_bucket        = "${google_storage_bucket.bucket.name}"
  timeout               = 61
  entry_point           = "helloGET"
}

resource "google_pubsub_topic" "sub" {
  name = "%s"
}

resource "google_cloudfunctions_function" "function_pubsub" {
  name                  = "%s-pubsub"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_topic         = "${google_pubsub_topic.sub.name}"
  timeout               = 61
  entry_point           = "helloGET"
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
