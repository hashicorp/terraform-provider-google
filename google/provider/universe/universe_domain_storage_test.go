// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package universe_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccUniverseDomainStorage(t *testing.T) {
	// Skip this test in all env since this can only run in specific test project.
	// Location field from `google_storage_bucket` needs to be changed depending on the universe.
	t.Skip()

	universeDomain := envvar.GetTestUniverseDomainFromEnv(t)
	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccUniverseDomain_bucket(universeDomain, bucketName),
			},
		},
	})
}

func testAccUniverseDomain_bucket(universeDomain string, bucketName string) string {
	return fmt.Sprintf(`
provider "google" {
  universe_domain = "%s"
}
	  
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
`, universeDomain, bucketName)
}

func testAccStorageBucketDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_bucket" {
				continue
			}

			_, err := config.NewStorageClient(config.UserAgent).Buckets.Get(rs.Primary.ID).Do()
			if err == nil {
				return fmt.Errorf("Bucket still exists")
			}
		}

		return nil
	}
}
