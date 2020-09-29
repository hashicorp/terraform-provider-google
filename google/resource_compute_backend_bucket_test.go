package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeBackendBucket_basicModified(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	storageName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	secondStorageName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucket_basic(backendName, storageName),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_basicModified(
					backendName, storageName, secondStorageName),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendBucket_withCdnPolicy(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	storageName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucket_withCdnPolicy(backendName, storageName),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeBackendBucket_basic(backendName, storageName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  bucket_name = google_storage_bucket.bucket_one.name
}

resource "google_storage_bucket" "bucket_one" {
  name     = "%s"
  location = "EU"
}
`, backendName, storageName)
}

func testAccComputeBackendBucket_basicModified(backendName, bucketOne, bucketTwo string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  bucket_name = google_storage_bucket.bucket_two.name
}

resource "google_storage_bucket" "bucket_one" {
  name     = "%s"
  location = "EU"
}

resource "google_storage_bucket" "bucket_two" {
  name     = "%s"
  location = "EU"
}
`, backendName, bucketOne, bucketTwo)
}

func testAccComputeBackendBucket_withCdnPolicy(backendName, storageName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  bucket_name = google_storage_bucket.bucket.name
  enable_cdn  = true
  cdn_policy {
    signed_url_cache_max_age_sec = 1000
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "EU"
}
`, backendName, storageName)
}
