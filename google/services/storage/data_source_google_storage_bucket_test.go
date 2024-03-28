// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageBucket_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"bucket_name": "tf-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageBucketConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_storage_bucket.bar", "google_storage_bucket.foo", map[string]struct{}{"force_destroy": {}}),
				),
			},
		},
	})
}

// Test that the data source can take a project argument, which is used as a way to avoid using Compute API to
// get project id for the project number returned from the Storage API.
func TestAccDataSourceGoogleStorageBucket_avoidComputeAPI(t *testing.T) {
	// Cannot use t.Parallel() if using t.Setenv

	project := envvar.GetTestProjectFromEnv()

	context := map[string]interface{}{
		"bucket_name":          "tf-bucket-" + acctest.RandString(t, 10),
		"real_project_id":      project,
		"incorrect_project_id": "foobar",
	}

	// Unset ENV so no provider default is available to the data source
	t.Setenv("GOOGLE_PROJECT", "")

	acctest.VcrTest(t, resource.TestCase{
		// Removed PreCheck because it wants to enforce GOOGLE_PROJECT being set
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageBucketConfig_setProjectInConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// We ignore project to show that the project argument on the data source is retained and isn't impacted
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_storage_bucket.bar", "google_storage_bucket.foo", map[string]struct{}{"force_destroy": {}, "project": {}}),

					resource.TestCheckResourceAttrSet(
						"google_storage_bucket.foo", "project_number"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.foo", "project", context["real_project_id"].(string)),

					resource.TestCheckResourceAttrSet(
						"data.google_storage_bucket.bar", "project_number"),
					resource.TestCheckResourceAttr(
						"data.google_storage_bucket.bar", "project", context["incorrect_project_id"].(string)),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageBucketConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "foo" {
  name     = "%{bucket_name}"
  location = "US"
}

data "google_storage_bucket" "bar" {
  name = google_storage_bucket.foo.name
  depends_on = [
    google_storage_bucket.foo,
  ]
}
`, context)
}

func testAccDataSourceGoogleStorageBucketConfig_setProjectInConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "foo" {
  project = "%{real_project_id}"
  name     = "%{bucket_name}"
  location = "US"
}

// The project argument here doesn't help the provider retrieve data about the bucket
// It only serves to stop the data source using the compute API to convert the project number to an id
data "google_storage_bucket" "bar" {
  project = "%{incorrect_project_id}"
  name = google_storage_bucket.foo.name
  depends_on = [
    google_storage_bucket.foo,
  ]
}
`, context)
}
