// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccStorageAnywhereCache_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageAnywhereCache_full(context),
			},
			{
				ResourceName:            "google_storage_anywhere_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket"},
			},
			{
				Config: testAccStorageAnywhereCache_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_storage_anywhere_cache.cache", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_storage_anywhere_cache.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket"},
			},
		},
	})
}

func testAccStorageAnywhereCache_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-bucket-name%{random_suffix}"
  location                    = "US"
}

resource "time_sleep" "destroy_wait_5000_seconds" {
  depends_on = [google_storage_bucket.bucket]
  destroy_duration = "5000s"
}

resource "google_storage_anywhere_cache" "cache" {
  bucket = google_storage_bucket.bucket.name
  zone = "us-central1-f"
  ttl = "3601s"
  depends_on = [time_sleep.destroy_wait_5000_seconds]
}
`, context)
}

func testAccStorageAnywhereCache_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-bucket-name%{random_suffix}"
  location                    = "US"
}

resource "time_sleep" "destroy_wait_5000_seconds" {
  depends_on = [google_storage_bucket.bucket]
  destroy_duration = "5000s"
}

resource "google_storage_anywhere_cache" "cache" {
  bucket = google_storage_bucket.bucket.name
  zone = "us-central1-f"
  admission_policy = "admit-on-second-miss"
  ttl = "3620s"
  depends_on = [time_sleep.destroy_wait_5000_seconds]
}
`, context)
}
