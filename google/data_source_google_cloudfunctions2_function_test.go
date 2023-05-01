package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleCloudFunctions2Function_basic(t *testing.T) {
	t.Parallel()

	funcDataNameHttp := "data.google_cloudfunctions2_function.function_http_v2"
	functionName := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", RandInt(t))
	zipFilePath := "./test-fixtures/cloudfunctions2/function-source.zip"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudfunctions2functionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudFunctions2FunctionConfig(functionName,
					bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(funcDataNameHttp,
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
