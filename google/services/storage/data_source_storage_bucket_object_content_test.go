// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceStorageBucketObjectContent_Basic(t *testing.T) {

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageBucketObjectContent_Basic(content, bucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content", content),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content_base64"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content_base64", base64.StdEncoding.EncodeToString([]byte(content))),
				),
			},
		},
	})
}

func TestAccDataSourceStorageBucketObjectContent_FileContentBase64(t *testing.T) {
	acctest.SkipIfVcr(t)

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	folderName := "tf-folder-" + acctest.RandString(t, 10)

	if err := os.Mkdir(folderName, 0777); err != nil {
		t.Errorf("error creating directory: %v", err)
	}

	data := []byte("data data data")
	testFile := getTmpTestFile(t, folderName, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	defer os.Remove(testFile.Name()) // clean up

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"local": resource.ExternalProvider{
				VersionConstraint: "> 2.5.0",
			},
			"archive": resource.ExternalProvider{
				VersionConstraint: "> 2.5.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageBucketObjectContent_FileContentBase64(bucket, folderName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.this", "content_base64"),
					verifyValidZip(),
				),
			},
		},
	})
}

func verifyValidZip() func(*terraform.State) error {
	return func(s *terraform.State) error {
		var outputFilePath string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "local_file" {
				outputFilePath = rs.Primary.Attributes["filename"]
				break
			}
		}
		archive, err := zip.OpenReader(outputFilePath)
		if err != nil {
			return err
		}
		defer archive.Close()
		return nil
	}
}

func testAccDataSourceStorageBucketObjectContent_Basic(content, bucket string) string {
	return fmt.Sprintf(`
data "google_storage_bucket_object_content" "default" {
	bucket = google_storage_bucket.contenttest.name
	name   = google_storage_bucket_object.object.name      
}

resource "google_storage_bucket_object" "object" {
	name    = "butterfly01"
	content = "%s"
	bucket  = google_storage_bucket.contenttest.name
}

resource "google_storage_bucket" "contenttest" {
	name          = "%s"
	location      = "US"
	force_destroy = true
}`, content, bucket)
}

func testAccDataSourceStorageBucketObjectContent_FileContentBase64(bucket, folderName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "this" {
  name                        = "%s"
  location                    = "us-east4"
  uniform_bucket_level_access = true
}

data "archive_file" "this" {
  type       = "zip"
  source_dir = "${path.cwd}/%s"
  output_path = "${path.cwd}/archive.zip"
}

resource "google_storage_bucket_object" "this" {
  name   = "archive.zip"
  bucket = google_storage_bucket.this.name
  source = data.archive_file.this.output_path
}

data "google_storage_bucket_object_content" "this" {
  name   = google_storage_bucket_object.this.name
  bucket = google_storage_bucket.this.name
}

resource "local_file" "this" {
  content_base64 = (data.google_storage_bucket_object_content.this.content_base64)
  filename = "${path.cwd}/content.zip"
}`, bucket, folderName)
}

func TestAccDataSourceStorageBucketObjectContent_Issue15717(t *testing.T) {

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	config := fmt.Sprintf(`
%s

output "output" {
	value = replace(data.google_storage_bucket_object_content.default.content, "q", "Q")
}`, testAccDataSourceStorageBucketObjectContent_Basic(content, bucket))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content", content),
				),
			},
		},
	})
}

func TestAccDataSourceStorageBucketObjectContent_Issue15717BackwardCompatibility(t *testing.T) {

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	config := fmt.Sprintf(`
%s

data "google_storage_bucket_object_content" "new" {
	bucket  = google_storage_bucket.contenttest.name
	content = "%s"
	name    = google_storage_bucket_object.object.name
}`, testAccDataSourceStorageBucketObjectContent_Basic(content, bucket), content)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.new", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.new", "content", content),
				),
			},
		},
	})
}

func getTmpTestFile(t *testing.T, folderName, prefix string) *os.File {
	testFile, err := ioutil.TempFile(folderName, prefix)
	if err != nil {
		t.Fatalf("Cannot create temp file: %s", err)
	}
	return testFile
}
