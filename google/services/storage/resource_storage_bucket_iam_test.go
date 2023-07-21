// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageBucketIamPolicy(t *testing.T) {
	t.Parallel()

	serviceAcct := envvar.GetTestServiceAccountFromEnv(t)
	bucket := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Policy creation
				Config: testAccStorageBucketIamPolicy_basic(bucket, account, serviceAcct),
			},
			{
				ResourceName:      "google_storage_bucket_iam_policy.bucket-binding",
				ImportStateId:     bucket,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Policy update
				Config: testAccStorageBucketIamPolicy_update(bucket, account, serviceAcct),
			},
			{
				ResourceName:      "google_storage_bucket_iam_policy.bucket-binding",
				ImportStateId:     bucket,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Policy with member 'allUsers'
				Config: testAccStorageBucketIamPolicy_allUsers(bucket),
			},
		},
	})
}

func testAccStorageBucketIamPolicy_update(bucket, account, serviceAcct string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Storage Bucket Iam Testing Account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Storage Bucket Iam Testing Account"
}

data "google_iam_policy" "foo-policy" {
  binding {
    role = "roles/storage.objectViewer"

    members = [
      "serviceAccount:${google_service_account.test-account-1.email}",
      "serviceAccount:${google_service_account.test-account-2.email}",
    ]
  }

  binding {
    role = "roles/storage.admin"
    members = [
      "serviceAccount:%s",
    ]
  }
}

resource "google_storage_bucket_iam_policy" "bucket-binding" {
  bucket      = google_storage_bucket.bucket.name
  policy_data = data.google_iam_policy.foo-policy.policy_data
}
`, bucket, account, account, serviceAcct)
}

func testAccStorageBucketIamPolicy_basic(bucket, account, serviceAcct string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_service_account" "test-account-1" {
  account_id   = "%s-1"
  display_name = "Storage Bucket Iam Testing Account"
}

data "google_iam_policy" "foo-policy" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "serviceAccount:${google_service_account.test-account-1.email}",
    ]
  }

  binding {
    role = "roles/storage.admin"
    members = [
      "serviceAccount:%s",
    ]
  }
}

resource "google_storage_bucket_iam_policy" "bucket-binding" {
  bucket      = google_storage_bucket.bucket.name
  policy_data = data.google_iam_policy.foo-policy.policy_data
}
`, bucket, account, serviceAcct)
}

func testAccStorageBucketIamPolicy_allUsers(bucket string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

data "google_iam_policy" "foo-policy" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
      "allAuthenticatedUsers",
    ]
  }
}

resource "google_storage_bucket_iam_policy" "bucket-binding" {
  bucket      = google_storage_bucket.bucket.name
  policy_data = data.google_iam_policy.foo-policy.policy_data
}
`, bucket)
}
