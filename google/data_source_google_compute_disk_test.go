package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleComputeDisk_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeDiskDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeDisk_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_compute_disk.foo", "google_compute_disk.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeDisk_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_disk" "foo" {
  name     = "tf-test-compute-disk-%{random_suffix}"
}

data "google_compute_disk" "foo" {
  name     = google_compute_disk.foo.name
  project  = google_compute_disk.foo.project
}
`, context)
}
