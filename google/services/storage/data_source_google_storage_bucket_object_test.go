// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageBucketObject_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	bucket := "tf-bucket-object-test-" + acctest.RandString(t, 10)

	context := map[string]interface{}{
		"bucket":      bucket,
		"project":     project,
		"object_name": "bee",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleStorageBucketObjectConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object.bee", "crc32c"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object.bee", "md5hash"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object.bee", "self_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object.bee", "storage_class"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object.bee", "media_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object.bee", "generation"),
					// Test content
					resource.TestCheckResourceAttr("data.google_storage_bucket_object.bee", "bucket", context["bucket"].(string)),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object.bee", "name", context["object_name"].(string)),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketObjectConfig(context map[string]interface{}) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "my_insect_cage" {
  force_destroy               = true
  location                    = "EU"
  name                        = "%s"
  project                     = "%s"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "bee" {
  bucket  = google_storage_bucket.my_insect_cage.name
  content = "bzzzzzt"
  name    = "%s"
}

data "google_storage_bucket_object" "bee" {
  bucket = google_storage_bucket.my_insect_cage.name
  name = google_storage_bucket_object.bee.name

  depends_on = [
    google_storage_bucket_object.bee,
  ]
}`,
		context["bucket"].(string),
		context["project"].(string),
		context["object_name"].(string),
	)
}
