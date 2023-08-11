// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeBackendBucket_basicModified(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	storageName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	secondStorageName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendBucketDestroyProducer(t),
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

	backendName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	storageName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucket_withCdnPolicy(backendName, storageName),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withCdnPolicy2(backendName, storageName, 1000, 301, 2, 1),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withCdnPolicy2(backendName, storageName, 0, 404, 86400, 0),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withCdnPolicy(backendName, storageName),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withCdnPolicy3(backendName, storageName, 0, 404, 0, 0),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withCdnPolicy4(backendName, storageName, 0, 404, 0),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendBucket_withSecurityPolicy(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	polName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucket_withSecurityPolicy(bucketName, polName, "google_compute_security_policy.policy.self_link"),
			},
			{
				ResourceName:      "google_compute_backend_bucket.image_backend",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withSecurityPolicy(bucketName, polName, "\"\""),
			},
			{
				ResourceName:      "google_compute_backend_bucket.image_backend",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeBackendBucket_withCompressionMode(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	storageName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendBucket_withCompressionMode(backendName, storageName, "DISABLED"),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_withCompressionMode(backendName, storageName, "AUTOMATIC"),
			},
			{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeBackendBucket_basic(backendName, storageName),
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
	negative_caching = false
  }
}
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "EU"
}
`, backendName, storageName)
}

func testAccComputeBackendBucket_withCdnPolicy2(backendName, storageName string, age, code, max_ttl, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  bucket_name = google_storage_bucket.bucket.name
  enable_cdn  = true
  cdn_policy {
	cache_mode                   = "CACHE_ALL_STATIC"
	signed_url_cache_max_age_sec = %d
	max_ttl                      = %d
	default_ttl                  = %d
	client_ttl                   = %d
	serve_while_stale            = %d
	negative_caching_policy {
		code = %d
		ttl = %d
	}
	negative_caching = true
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "EU"
}
`, backendName, age, max_ttl, ttl, ttl, ttl, code, ttl, storageName)
}

func testAccComputeBackendBucket_withCdnPolicy3(backendName, storageName string, age, code, max_ttl, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  bucket_name = google_storage_bucket.bucket.name
  enable_cdn  = true
  cdn_policy {
	cache_mode                   = "FORCE_CACHE_ALL"
	signed_url_cache_max_age_sec = %d
	max_ttl                      = %d
	default_ttl                  = %d
	client_ttl                   = %d
	serve_while_stale            = %d
	negative_caching_policy {
		code = %d
		ttl = %d
	}
	negative_caching = true
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "EU"
}
`, backendName, age, max_ttl, ttl, ttl, ttl, code, ttl, storageName)
}

func testAccComputeBackendBucket_withSecurityPolicy(bucketName, polName, polLink string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "image_backend" {
  name        = "%s"
  description = "Contains beautiful images"
  bucket_name = google_storage_bucket.image_bucket.name
  enable_cdn  = true
  edge_security_policy = %s
}

resource "google_storage_bucket" "image_bucket" {
  name     = "%s"
  location = "EU"
}


resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type = "CLOUD_ARMOR_EDGE"
}
`, bucketName, polLink, bucketName, polName)
}

func testAccComputeBackendBucket_withCompressionMode(backendName, storageName, compressionMode string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name             = "%s"
  bucket_name      = google_storage_bucket.bucket_one.name
  enable_cdn       = true
  compression_mode = "%s"
}

resource "google_storage_bucket" "bucket_one" {
  name     = "%s"
  location = "EU"
}
`, backendName, compressionMode, storageName)
}

func testAccComputeBackendBucket_withCdnPolicy4(backendName, storageName string, age, code, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_backend_bucket" "foobar" {
  name        = "%s"
  bucket_name = google_storage_bucket.bucket.name
  enable_cdn  = true
  cdn_policy {
	cache_mode                   = "USE_ORIGIN_HEADERS"
	signed_url_cache_max_age_sec = %d
	serve_while_stale            = %d
	negative_caching_policy {
		code = %d
		ttl = %d
	}
	negative_caching = true
  }
}
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "EU"
}
`, backendName, age, ttl, code, ttl, storageName)
}
