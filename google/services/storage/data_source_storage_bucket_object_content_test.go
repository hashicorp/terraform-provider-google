// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				),
			},
		},
	})
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
