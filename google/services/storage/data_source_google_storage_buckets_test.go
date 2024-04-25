// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageBuckets_basic(t *testing.T) {
	t.Parallel()

	static_prefix := "tf-bucket-test"
	random_suffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"bucket1":         static_prefix + "-1-" + random_suffix,
		"bucket2":         static_prefix + "-2-" + random_suffix,
		"project_id":      static_prefix + "-" + random_suffix,
		"organization":    envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleStorageBucketsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.0.location"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.0.name"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.0.self_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.0.storage_class"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.1.location"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.1.name"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.1.self_link"),
					resource.TestCheckResourceAttrSet("data.google_storage_buckets.all", "buckets.1.storage_class"),
					// Test content
					resource.TestCheckResourceAttr("data.google_storage_buckets.all", "project", context["project_id"].(string)),
					resource.TestCheckResourceAttr("data.google_storage_buckets.all", "buckets.0.name", context["bucket1"].(string)),
					resource.TestCheckResourceAttr("data.google_storage_buckets.all", "buckets.1.name", context["bucket2"].(string)),
					// Test with project
					resource.TestCheckResourceAttr("data.google_storage_buckets.one", "buckets.0.name", context["bucket1"].(string)),
					// Test prefix
					resource.TestCheckResourceAttr("data.google_storage_buckets.two", "buckets.0.name", context["bucket2"].(string)),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketsConfig(context map[string]interface{}) string {
	return fmt.Sprintf(`
locals {
  billing_account = "%s"
  bucket_one      = "%s"
  bucket_two      = "%s"
  organization    = "%s"
  project_id      = "%s"
}

resource "google_project" "acceptance" {
  name            = local.project_id
  project_id      = local.project_id
  org_id          = local.organization
  billing_account = local.billing_account
}

resource "google_storage_bucket" "one" {
  force_destroy               = true
  location                    = "EU"
  name                        = local.bucket_one
  project                     = google_project.acceptance.project_id
  uniform_bucket_level_access = true
}

resource "google_storage_bucket" "two" {
  force_destroy               = true
  location                    = "EU"
  name                        = local.bucket_two
  project                     = google_project.acceptance.project_id
  uniform_bucket_level_access = true
}

data "google_storage_buckets" "all" {
  project = google_project.acceptance.project_id
  
  depends_on = [
    google_storage_bucket.one,
    google_storage_bucket.two,
  ]
}

data "google_storage_buckets" "one" {
  prefix  = "tf-bucket-test-1"
  project = google_project.acceptance.project_id

  depends_on = [
    google_storage_bucket.one,
  ]
}

data "google_storage_buckets" "two" {
  prefix  = "tf-bucket-test-2"
  project = google_project.acceptance.project_id

  depends_on = [
    google_storage_bucket.two,
  ]
}`,
		context["billing_account"].(string),
		context["bucket1"].(string),
		context["bucket2"].(string),
		context["organization"].(string),
		context["project_id"].(string),
	)
}
