// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"testing"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageSignedUrl_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSignedUrlConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccSignedUrlExists(t, "data.google_storage_object_signed_url.blerg"),
				),
			},
		},
	})
}

func TestAccStorageSignedUrl_accTest(t *testing.T) {
	// URL includes an expires time
	acctest.SkipIfVcr(t)
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt(t))

	headers := map[string]string{
		"x-goog-test":                    "foo",
		"x-goog-if-metageneration-match": "1",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTestGoogleStorageObjectSignedURL(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccSignedUrlRetrieval("data.google_storage_object_signed_url.story_url", nil),
					testAccSignedUrlRetrieval("data.google_storage_object_signed_url.story_url_w_headers", headers),
					testAccSignedUrlRetrieval("data.google_storage_object_signed_url.story_url_w_content_type", nil),
					testAccSignedUrlRetrieval("data.google_storage_object_signed_url.story_url_w_md5", nil),
				),
			},
		},
	})
}

func testAccSignedUrlExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["signed_url"] == "" {
			return fmt.Errorf("signed_url is empty: %v", a)
		}

		return nil
	}
}

func testAccSignedUrlRetrieval(n string, headers map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		if r == nil {
			return fmt.Errorf("Datasource not found")
		}
		a := r.Primary.Attributes

		if a["signed_url"] == "" {
			return fmt.Errorf("signed_url is empty: %v", a)
		}

		// create HTTP request
		url := a["signed_url"]
		method := a["http_method"]
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}

		// Add extension headers to request, if provided
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// content_type is optional, add to test query if provided in datasource config
		contentType := a["content_type"]
		if contentType != "" {
			req.Header.Add("Content-Type", contentType)
		}

		// content_md5 is optional, add to test query if provided in datasource config
		contentMd5 := a["content_md5"]
		if contentMd5 != "" {
			req.Header.Add("Content-MD5", contentMd5)
		}

		// send request using signed url
		client := cleanhttp.DefaultClient()
		response, err := client.Do(req)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		// check content in response, should be our test string or XML with error
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if string(body) != "once upon a time..." {
			return fmt.Errorf("Got unexpected object contents: %s\n\tURL: %s", string(body), url)
		}

		return nil
	}
}

const testGoogleSignedUrlConfig = `
data "google_storage_object_signed_url" "blerg" {
  bucket = "friedchicken"
  path   = "path/to/file"

}
`

func testAccTestGoogleStorageObjectSignedURL(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "story" {
  name   = "path/to/file"
  bucket = google_storage_bucket.bucket.name

  content = "once upon a time..."
}

data "google_storage_object_signed_url" "story_url" {
  bucket = google_storage_bucket.bucket.name
  path   = google_storage_bucket_object.story.name
}

data "google_storage_object_signed_url" "story_url_w_headers" {
  bucket = google_storage_bucket.bucket.name
  path   = google_storage_bucket_object.story.name
  extension_headers = {
    x-goog-test                = "foo"
    x-goog-if-metageneration-match = 1
  }
}

data "google_storage_object_signed_url" "story_url_w_content_type" {
  bucket = google_storage_bucket.bucket.name
  path   = google_storage_bucket_object.story.name

  content_type = "text/plain"
}

data "google_storage_object_signed_url" "story_url_w_md5" {
  bucket = google_storage_bucket.bucket.name
  path   = google_storage_bucket_object.story.name

  content_md5 = google_storage_bucket_object.story.md5hash
}
`, bucketName)
}
