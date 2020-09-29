package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceComputeBackendBucket_basic(t *testing.T) {
	t.Parallel()

	backendBucketName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeBackendBucket_basic(backendBucketName, bucketName),
				Check:  checkDataSourceStateMatchesResourceState("data.google_compute_backend_bucket.baz", "google_compute_backend_bucket.foobar"),
			},
		},
	})
}

func testAccDataSourceComputeBackendBucket_basic(backendBucketName, bucketName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  description = "Contains beautiful images"
  bucket_name = google_storage_bucket.image_bucket.name
  enable_cdn  = true
}
resource "google_storage_bucket" "image_bucket" {
  name     = "%s"
  location = "EU"
}
data "google_compute_backend_bucket" "baz" {
  name = google_compute_backend_bucket.foobar.name
}
`, backendBucketName, bucketName)
}
