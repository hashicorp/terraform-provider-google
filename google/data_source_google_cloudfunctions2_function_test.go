package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleCloudFunctions2Function_basic(t *testing.T) {
	t.Parallel()

	funcDataNameHttp := "data.google_cloudfunctions2_function.function_http_v2"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	zipFilePath := "./test-fixtures/cloudfunctions2/function-source.zip"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudfunctions2functionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFunctions2FunctionConfig(functionName,
					bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(funcDataNameHttp,
						"google_cloudfunctions2_function.function_http_v2", map[string]struct{}{"build_config.0.source.0.storage_source.0.bucket": {}, "build_config.0.source.0.storage_source.0.object": {}}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudFunctions2FunctionConfig(functionName, bucketName, zipFilePath string) string {
	return fmt.Sprintf(`
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
