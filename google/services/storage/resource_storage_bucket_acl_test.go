// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/storage"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	roleEntityBasic1        = "OWNER:user-gterraformtest1@gmail.com"
	roleEntityBasic1_reader = "READER:user-gterraformtest1@gmail.com"
	roleEntityBasic2        = "READER:user-gterraformtest2@gmail.com"
	roleEntityBasic3_owner  = "OWNER:user-paddy@paddy.io"
	roleEntityBasic3_reader = "READER:user-foran.paddy@gmail.com"

	roleEntityOwners  = "OWNER:project-owners-" + os.Getenv("GOOGLE_PROJECT_NUMBER")
	roleEntityEditors = "OWNER:project-editors-" + os.Getenv("GOOGLE_PROJECT_NUMBER")
	roleEntityViewers = "READER:project-viewers-" + os.Getenv("GOOGLE_PROJECT_NUMBER")
)

func TestAccStorageBucketAcl_basic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_upgrade(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasic2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic3_owner),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_upgradeSingleUser(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic1_reader(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic1_reader),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasic1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic1_reader),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_downgrade(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclBasic2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasic3(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAcl(t, bucketName, roleEntityBasic3_reader),
				),
			},

			{
				Config: testGoogleStorageBucketsAclBasicDelete(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageBucketAclDelete(t, bucketName, roleEntityBasic3_owner),
				),
			},
		},
	})
}

func TestAccStorageBucketAcl_predefined(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclPredefined(bucketName),
			},
		},
	})
}

// Test that we allow the API to reorder our role entities without perma-diffing.
func TestAccStorageBucketAcl_unordered(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclUnordered(bucketName),
			},
		},
	})
}

// Test that project owner doesn't get removed or cause a diff
func TestAccStorageBucketAcl_RemoveOwner(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsAclRemoveOwner(bucketName),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketAclDelete(t *testing.T, bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := storage.GetRoleEntityPair(roleEntityS)
		config := acctest.GoogleProviderConfig(t)

		_, err := config.NewStorageClient(config.UserAgent).BucketAccessControls.Get(bucket, roleEntity.Entity).Do()

		if err != nil {
			return nil
		}

		return fmt.Errorf("Error, entity %s still exists", roleEntity.Entity)
	}
}

func testAccCheckGoogleStorageBucketAcl(t *testing.T, bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := storage.GetRoleEntityPair(roleEntityS)
		config := acctest.GoogleProviderConfig(t)

		res, err := config.NewStorageClient(config.UserAgent).BucketAccessControls.Get(bucket, roleEntity.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of acl for bucket %s: %s", bucket, err)
		}

		if res.Role != roleEntity.Role {
			return fmt.Errorf("Error, Role mismatch %s != %s", res.Role, roleEntity.Role)
		}

		return nil
	}
}

func testAccStorageBucketAclDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_bucket_acl" {
				continue
			}

			bucket := rs.Primary.Attributes["bucket"]

			_, err := config.NewStorageClient(config.UserAgent).BucketAccessControls.List(bucket).Do()

			if err == nil {
				return fmt.Errorf("Acl for bucket %s still exists", bucket)
			}
		}

		return nil
	}
}

func testGoogleStorageBucketsAclBasic1_reader(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic1_reader, roleEntityBasic2)
}

func testGoogleStorageBucketsAclBasic1(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic1, roleEntityBasic2)
}

func testGoogleStorageBucketsAclBasic2(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic2, roleEntityBasic3_owner)
}

func testGoogleStorageBucketsAclBasicDelete(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = []
}
`, bucketName)
}

func testGoogleStorageBucketsAclBasic3(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityOwners, roleEntityEditors, roleEntityViewers, roleEntityBasic2, roleEntityBasic3_reader)
}

func testGoogleStorageBucketsAclUnordered(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityBasic1, roleEntityViewers, roleEntityOwners, roleEntityBasic2, roleEntityEditors)
}

func testGoogleStorageBucketsAclPredefined(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket         = google_storage_bucket.bucket.name
  predefined_acl = "projectPrivate"
  default_acl    = "projectPrivate"
}
`, bucketName)
}

func testGoogleStorageBucketsAclRemoveOwner(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_acl" "acl" {
  bucket         = google_storage_bucket.bucket.name
  role_entity = [
	"READER:user-gterraformtest2@gmail.com"
  ]
}
`, bucketName)
}
