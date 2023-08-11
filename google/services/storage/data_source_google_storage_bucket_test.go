// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucket := "tf-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageBucketConfig(bucket),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_storage_bucket.bar", "google_storage_bucket.foo", map[string]struct{}{"force_destroy": {}}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageBucketConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "foo" {
  name     = "%s"
  location = "US"
}

data "google_storage_bucket" "bar" {
  name = google_storage_bucket.foo.name
  depends_on = [
    google_storage_bucket.foo,
  ]
}
`, bucketName)
}
