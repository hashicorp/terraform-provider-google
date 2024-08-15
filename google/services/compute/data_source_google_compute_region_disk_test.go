// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleComputeRegionDisk_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeRegionDisk_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_region_disk.foo", "google_compute_region_disk.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeRegionDisk_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_disk" "foo" {
  name          = "tf-test-compute-disk-%{random_suffix}"
  type          = "pd-standard"
  replica_zones = ["us-central1-a", "us-central1-f"]

  labels = {
    my-label = "my-label-value"
  }
}

data "google_compute_region_disk" "foo" {
  name     = google_compute_region_disk.foo.name
  project  = google_compute_region_disk.foo.project
}
`, context)
}
