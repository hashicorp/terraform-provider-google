// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeImages_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"image":         "debian-cloud/debian-11",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleComputeImagesConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Test schema
					resource.TestCheckResourceAttrSet("data.google_compute_images.all", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.google_compute_images.all", "images.1.name"),
					resource.TestCheckResourceAttrSet("data.google_compute_images.all", "images.0.self_link"),
					resource.TestCheckResourceAttrSet("data.google_compute_images.all", "images.1.self_link"),
					resource.TestCheckResourceAttrSet("data.google_compute_images.all", "images.0.image_id"),
					resource.TestCheckResourceAttrSet("data.google_compute_images.all", "images.1.image_id"),
				),
			},
		},
	})
}

func testAccCheckGoogleComputeImagesConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_disk" "test-disk" {
  name  = "tf-test-disk-%{random_suffix}"
  type  = "pd-standard"
  image = "%{image}"
  size  = 10
}

resource "google_compute_image" "foo" {
  name = "tf-test-image1-%{random_suffix}"
  source_disk = google_compute_disk.test-disk.self_link
}

resource "google_compute_image" "bar" {
  name = "tf-test-image2-%{random_suffix}"
  source_image = google_compute_image.foo.self_link
}

data "google_compute_images" "all" {
  depends_on = [
    google_compute_image.foo,
    google_compute_image.bar,
  ]
}
`, context)
}
