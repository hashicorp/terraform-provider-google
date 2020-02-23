package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceRegionInstanceGroup(t *testing.T) {
	t.Parallel()
	name := "acctest-" + acctest.RandString(6)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRegionInstanceGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group.data_source", "name", name),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group.data_source", "project", getTestProjectFromEnv()),
					resource.TestCheckResourceAttr("data.google_compute_region_instance_group.data_source", "instances.#", "1")),
			},
		},
	})
}

func testAccDataSourceRegionInstanceGroup_basic(instanceManagerName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo" {
  name = "%s"
}

data "google_compute_image" "debian" {
  project = "debian-cloud"
  name    = "debian-9-stretch-v20171129"
}

resource "google_compute_instance_template" "foo" {
  machine_type = "n1-standard-1"
  disk {
    source_image = data.google_compute_image.debian.self_link
  }
  network_interface {
    access_config {
    }
    network = "default"
  }
}

resource "google_compute_region_instance_group_manager" "foo" {
  name               = "%s"
  base_instance_name = "foo"
  version {
    instance_template = google_compute_instance_template.foo.self_link
    name              = "primary"
  }
  region       = "us-central1"
  target_pools = [google_compute_target_pool.foo.self_link]
  target_size  = 1

  named_port {
    name = "web"
    port = 80
  }
  wait_for_instances = true
}

data "google_compute_region_instance_group" "data_source" {
  self_link = google_compute_region_instance_group_manager.foo.instance_group
}
`, acctest.RandomWithPrefix("test-rigm-"), instanceManagerName)
}
