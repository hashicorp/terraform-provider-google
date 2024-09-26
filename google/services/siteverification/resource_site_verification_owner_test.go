// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package siteverification_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSiteVerificationOwner_siteVerificationBucket(t *testing.T) {
	t.Parallel()

	account1 := "tf-test-" + acctest.RandString(t, 10)
	account2 := "tf-test-" + acctest.RandString(t, 10)

	bucket := "tf-siteverification-test-" + acctest.RandString(t, 10)
	context := map[string]interface{}{
		"bucket":   bucket,
		"account1": account1,
		"account2": account2,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Initial configuration with one owner resource.
				Config: testAccSiteVerificationOwner_siteVerificationBucket(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_site_verification_web_resource.example", "owners.#", "1"),
				),
			},
			{
				ResourceName:      "google_site_verification_owner.example1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add a second owner resource.
				Config: testAccSiteVerificationOwner_siteVerificationBucketSecondOwner(context),
			},
			{
				ResourceName:      "google_site_verification_owner.example1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_site_verification_owner.example2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove both owners.
				Config: testAccSiteVerificationOwner_siteVerificationBucketRemoveOwners(context),
			},
			{
				// Remove the object before the test deletes the web resource. The API will not
				// allow the web resource to be deleted if the object still exists, and this
				// ensures the proper order of deletion.
				Config: testAccSiteVerificationOwner_siteVerificationRemoveObject(context),
			},
		},
	})
}

func testAccSiteVerificationOwner_siteVerificationBucket(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  alias                 = "scoped"
  user_project_override = true
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_service_account" "test-account1" {
  account_id   = "%{account1}"
  display_name = "Site Verification Testing Account One"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account2}"
  display_name = "Site Verification Testing Account Two"
}

resource "google_storage_bucket" "bucket" {
  provider = google.scoped
  name     = "%{bucket}"
  location = "US"
}

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "https://${google_storage_bucket.bucket.name}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_storage_bucket_object" "object" {
  provider = google.scoped
  name     = "${data.google_site_verification_token.token.token}"
  content  = "google-site-verification: ${data.google_site_verification_token.token.token}"
  bucket   = google_storage_bucket.bucket.name
}

resource "google_storage_object_access_control" "public_rule" {
  provider = google.scoped
  bucket   = google_storage_bucket.bucket.name
  object   = google_storage_bucket_object.object.name
  role     = "READER"
  entity   = "allUsers"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}

resource "google_site_verification_owner" "example1" {
  provider        = google.scoped
  web_resource_id = google_site_verification_web_resource.example.id
  email           = "${google_service_account.test-account1.email}"
}
`, context)
}

func testAccSiteVerificationOwner_siteVerificationBucketSecondOwner(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  alias                 = "scoped"
  user_project_override = true
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_service_account" "test-account1" {
  account_id   = "%{account1}"
  display_name = "Site Verification Testing Account One"
}

resource "google_service_account" "test-account2" {
  account_id   = "%{account2}"
  display_name = "Site Verification Testing Account Two"
}

resource "google_storage_bucket" "bucket" {
  provider = google.scoped
  name     = "%{bucket}"
  location = "US"
}

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "https://${google_storage_bucket.bucket.name}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_storage_bucket_object" "object" {
  provider = google.scoped
  name     = "${data.google_site_verification_token.token.token}"
  content  = "google-site-verification: ${data.google_site_verification_token.token.token}"
  bucket   = google_storage_bucket.bucket.name
}

resource "google_storage_object_access_control" "public_rule" {
  provider = google.scoped
  bucket   = google_storage_bucket.bucket.name
  object   = google_storage_bucket_object.object.name
  role     = "READER"
  entity   = "allUsers"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}

resource "google_site_verification_owner" "example1" {
  provider        = google.scoped
  web_resource_id = google_site_verification_web_resource.example.id
  email           = "${google_service_account.test-account1.email}"
}

resource "google_site_verification_owner" "example2" {
  provider        = google.scoped
  web_resource_id = google_site_verification_web_resource.example.id
  email           = "${google_service_account.test-account2.email}"
}
`, context)
}

func testAccSiteVerificationOwner_siteVerificationBucketRemoveOwners(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  alias                 = "scoped"
  user_project_override = true
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_storage_bucket" "bucket" {
  provider = google.scoped
  name     = "%{bucket}"
  location = "US"
}

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "https://${google_storage_bucket.bucket.name}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_storage_bucket_object" "object" {
  provider = google.scoped
  name     = "${data.google_site_verification_token.token.token}"
  content  = "google-site-verification: ${data.google_site_verification_token.token.token}"
  bucket   = google_storage_bucket.bucket.name
}

resource "google_storage_object_access_control" "public_rule" {
  provider = google.scoped
  bucket   = google_storage_bucket.bucket.name
  object   = google_storage_bucket_object.object.name
  role     = "READER"
  entity   = "allUsers"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}
`, context)
}

func testAccSiteVerificationOwner_siteVerificationRemoveObject(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  alias                 = "scoped"
  user_project_override = true
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

resource "google_storage_bucket" "bucket" {
  provider = google.scoped
  name     = "%{bucket}"
  location = "US"
}

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "https://${google_storage_bucket.bucket.name}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}
`, context)
}
