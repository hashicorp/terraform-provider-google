// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccStorageDefaultObjectAcl_basic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageDefaultObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectsAclBasic(bucketName, roleEntityBasic1, roleEntityBasic2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageDefaultObjectAcl_noRoleEntity(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageDefaultObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectsAclNoRoleEntity(bucketName),
			},
		},
	})
}

func TestAccStorageDefaultObjectAcl_upgrade(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageDefaultObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectsAclBasic(bucketName, roleEntityBasic1, roleEntityBasic2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageDefaultObjectsAclBasic(bucketName, roleEntityBasic2, roleEntityBasic3_owner),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageDefaultObjectsAclBasicDelete(bucketName, roleEntityBasic1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAclDelete(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAclDelete(t, bucketName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func TestAccStorageDefaultObjectAcl_downgrade(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageDefaultObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectsAclBasic(bucketName, roleEntityBasic2, roleEntityBasic3_owner),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageDefaultObjectsAclBasic(bucketName, roleEntityBasic2, roleEntityBasic3_reader),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic3_reader),
				),
			},

			{
				Config: testGoogleStorageDefaultObjectsAclBasicDelete(bucketName, roleEntityBasic1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageDefaultObjectAcl(t, bucketName, roleEntityBasic1),
					testAccCheckGoogleStorageDefaultObjectAclDelete(t, bucketName, roleEntityBasic2),
					testAccCheckGoogleStorageDefaultObjectAclDelete(t, bucketName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

// Test that we allow the API to reorder our role entities without perma-diffing.
func TestAccStorageDefaultObjectAcl_unordered(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageDefaultObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageDefaultObjectAclUnordered(bucketName),
			},
		},
	})
}

func testAccCheckGoogleStorageDefaultObjectAcl(t *testing.T, bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := storage.GetRoleEntityPair(roleEntityS)
		config := acctest.GoogleProviderConfig(t)

		res, err := config.NewStorageClient(config.UserAgent).DefaultObjectAccessControls.Get(bucket,
			roleEntity.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of storage default Acl for bucket %s: %s", bucket, err)
		}

		if res.Role != roleEntity.Role {
			return fmt.Errorf("Error, Role mismatch %s != %s", res.Role, roleEntity.Role)
		}

		return nil
	}
}

func testAccStorageDefaultObjectAclDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {

			if rs.Type != "google_storage_default_object_acl" {
				continue
			}

			bucket := rs.Primary.Attributes["bucket"]

			_, err := config.NewStorageClient(config.UserAgent).DefaultObjectAccessControls.List(bucket).Do()
			if err == nil {
				return fmt.Errorf("Default Storage Object Acl for bucket %s still exists", bucket)
			}
		}
		return nil
	}
}

func testAccCheckGoogleStorageDefaultObjectAclDelete(t *testing.T, bucket, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := storage.GetRoleEntityPair(roleEntityS)
		config := acctest.GoogleProviderConfig(t)

		_, err := config.NewStorageClient(config.UserAgent).DefaultObjectAccessControls.Get(bucket, roleEntity.Entity).Do()

		if err != nil {
			return nil
		}

		return fmt.Errorf("Error, Object Default Acl Entity still exists %s for bucket %s",
			roleEntity.Entity, bucket)
	}
}

func testGoogleStorageDefaultObjectsAclBasic(bucketName, roleEntity1, roleEntity2 string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_default_object_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, roleEntity1, roleEntity2)
}

func testGoogleStorageDefaultObjectsAclBasicDelete(bucketName, roleEntity string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_default_object_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s"]
}
`, bucketName, roleEntity)
}

func testGoogleStorageDefaultObjectsAclNoRoleEntity(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_default_object_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = []
}
`, bucketName)
}

func testGoogleStorageDefaultObjectAclUnordered(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_default_object_acl" "acl" {
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, roleEntityBasic1, roleEntityViewers, roleEntityOwners, roleEntityBasic2, roleEntityEditors)
}
