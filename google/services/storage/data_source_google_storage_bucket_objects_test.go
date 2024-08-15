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

func TestAccDataSourceGoogleStorageBucketObjects_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	bucket := "tf-bucket-object-test-" + acctest.RandString(t, 10)

	context := map[string]interface{}{
		"bucket":        bucket,
		"project":       project,
		"object_0_name": "bee",
		"object_1_name": "fly",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleStorageBucketObjectsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.0.content_type"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.0.media_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.0.name"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.0.self_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.0.storage_class"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.1.content_type"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.1.media_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.1.name"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.1.self_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_objects.my_insects", "bucket_objects.1.storage_class"),
					// Test content
					resource.TestCheckResourceAttr("data.google_storage_bucket_objects.my_insects", "bucket", context["bucket"].(string)),
					resource.TestCheckResourceAttr("data.google_storage_bucket_objects.my_insects", "bucket_objects.0.name", context["object_0_name"].(string)),
					resource.TestCheckResourceAttr("data.google_storage_bucket_objects.my_insects", "bucket_objects.1.name", context["object_1_name"].(string)),
					// Test match_glob
					resource.TestCheckResourceAttr("data.google_storage_bucket_objects.my_bee_glob", "bucket_objects.0.name", context["object_0_name"].(string)),
					// Test prefix
					resource.TestCheckResourceAttr("data.google_storage_bucket_objects.my_fly_prefix", "bucket_objects.0.name", context["object_1_name"].(string)),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketObjectsConfig(context map[string]interface{}) string {
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

resource "google_storage_bucket_object" "fly" {
  bucket  = google_storage_bucket.my_insect_cage.name
  content = "zzzzzt"
  name    = "%s"
}

data "google_storage_bucket_objects" "my_insects" {
  bucket = google_storage_bucket.my_insect_cage.name

  depends_on = [
    google_storage_bucket_object.bee,
	google_storage_bucket_object.fly,
  ]
}

data "google_storage_bucket_objects" "my_bee_glob" {
  bucket     = google_storage_bucket.my_insect_cage.name
  match_glob = "b*"

  depends_on = [
    google_storage_bucket_object.bee,
  ]
}

data "google_storage_bucket_objects" "my_fly_prefix" {
  bucket = google_storage_bucket.my_insect_cage.name
  prefix = "f"

  depends_on = [
    google_storage_bucket_object.fly,
  ]
}`,
		context["bucket"].(string),
		context["project"].(string),
		context["object_0_name"].(string),
		context["object_1_name"].(string),
	)
}
