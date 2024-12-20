// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"google.golang.org/api/storage/v1"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageFolder_storageFolderBasic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true, true) + testAccStorageFolder_storageFolder(true),
			},
			{
				ResourceName:            "google_storage_folder.folder",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket", "recursive", "force_destroy"},
			},
		},
	})
}

func TestAccStorageFolder_hnsDisabled(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageFolder_storageBucket(bucketName, false, true) + testAccStorageFolder_storageFolder(true),
				ExpectError: regexp.MustCompile("Error creating Folder: googleapi: Error 409: The bucket does not support hierarchical namespace., conflict"),
			},
		},
	})
}

func TestAccStorageFolder_FolderForceDestroy(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	data := []byte("data data data")

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucketObject(bucketName, true, true, testFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketUploadItem(t, bucketName),
				),
			},
		},
	})
}

func TestAccStorageFolder_DeleteEmptyFolderWithForceDestroyDefault(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucket(bucketName, true, true) + testAccStorageFolder_storageOneFolder(false),
			},
		},
	})
}

func TestAccStorageFolder_FailDeleteNonEmptyFolder(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolder_storageBucketAndFolders(bucketName, false),
			},
			{
				Config:      testAccStorageFolder_removeParentFolder(bucketName),
				ExpectError: regexp.MustCompile("use force_destroy to true to delete all subfolders"),
			},
			{
				Config: testAccStorageFolder_storageBucketAndFolders(bucketName, true),
			},
		},
	})
}

func testAccStorageFolder_storageBucket(bucketName string, hnsFlag bool, forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = %t
  }
  force_destroy = %t
}
`, bucketName, hnsFlag, forceDestroy)
}

func testAccStorageFolder_storageBucketObject(bucketName string, hnsFlag bool, forceDestroy bool, fileName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = %t
  }
  force_destroy = true
}
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
}
resource "google_storage_folder" "subfolder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "${google_storage_folder.folder.name}subfolder/"
  force_destroy = %t
}  
resource "google_storage_bucket_object" "object" {
  name   = "${google_storage_folder.subfolder.name}tffile"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}  
`, bucketName, hnsFlag, forceDestroy, forceDestroy, fileName)
}

func testAccStorageFolder_storageFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
}
resource "google_storage_folder" "subfolder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "${google_storage_folder.folder.name}name/"
  force_destroy = %t
}  
`, forceDestroy, forceDestroy)
}

func testAccStorageFolder_storageOneFolder(forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
} 
`, forceDestroy)
}

func testAccStorageFolder_storageBucketAndFolders(bucketName string, forceDestroy bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = "true"
  }
}
resource "google_storage_folder" "folder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/"
  force_destroy = %t
}
resource "google_storage_folder" "subfolder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "${google_storage_folder.folder.name}subfolder/"
} 
`, bucketName, forceDestroy)
}

func testAccStorageFolder_removeParentFolder(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "EU"
  uniform_bucket_level_access = true
  hierarchical_namespace {
	enabled = "true"
  }
}
resource "google_storage_folder" "subfolder" {
  bucket        = google_storage_bucket.bucket.name
  name          = "folder/subfolder/"
  force_destroy = false
} 
`, bucketName)
}

func testAccCheckStorageBucketUploadItem(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "folder/" + "bucketDestroyTestFile"}

		if res, err := config.NewStorageClient(config.UserAgent).Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		return nil
	}
}
