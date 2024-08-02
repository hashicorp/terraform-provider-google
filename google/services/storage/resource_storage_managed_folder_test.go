// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageManagedFolder_storageManagedFolderUpdate(t *testing.T) {
	t.Parallel()
	bucketName := fmt.Sprintf("tf-test-managed-folder-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(false),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config:      testAccStorageManagedFolder_bucket(bucketName),
				ExpectError: regexp.MustCompile(`Error 409: The managed folder you tried to delete is not empty.`),
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName) + testAccStorageManagedFolder_managedFolder(true),
			},
			{
				ResourceName:            "google_storage_managed_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "force_destroy"},
			},
			{
				Config: testAccStorageManagedFolder_bucket(bucketName),
			},
		},
	})
}

func testAccStorageManagedFolder_bucket(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name       = "managed/folder/name/file.txt"
  content    = "This file will affect the folder being deleted if allowNonEmpty=false"
  bucket     = google_storage_bucket.bucket.name
}
`, bucketName)
}

func testAccStorageManagedFolder_managedFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_managed_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "managed/folder/name/"
  force_destroy = %t
}
`, forceDestroy)
}
