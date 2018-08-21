package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	compute "google.golang.org/api/compute/v1"
)

func TestAccComputeInstanceFromTemplate_basic(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))
	templateName := fmt.Sprintf("terraform-test-%s", acctest.RandString(10))
	resourceName := "google_compute_instance_from_template.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceFromTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceFromTemplate_basic(instanceName, templateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),

					// Check that fields were set based on the template
					resource.TestCheckResourceAttr(resourceName, "machine_type", "n1-standard-1"),
					resource.TestCheckResourceAttr(resourceName, "attached_disk.#", "1"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceFromTemplateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_instance_from_template" {
			continue
		}

		_, err := config.clientCompute.Instances.Get(
			config.Project, rs.Primary.Attributes["zone"], rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Instance still exists")
		}
	}

	return nil
}

func testAccComputeInstanceFromTemplate_basic(instance, template string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
	family  = "debian-9"
	project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
	name = "%s"
	image = "${data.google_compute_image.my_image.self_link}"
	size = 10
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
	name = "%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "${data.google_compute_image.my_image.self_link}"
		auto_delete = true
		disk_size_gb = 100
		boot = true
	}

	disk {
		source = "${google_compute_disk.foobar.name}"
		auto_delete = false
		boot = false
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}

	can_ip_forward = true
}

resource "google_compute_instance_from_template" "foobar" {
	name           = "%s"
	zone           = "us-central1-a"

	source_instance_template = "${google_compute_instance_template.foobar.self_link}"

	// Overrides
	can_ip_forward = false
	labels {
		my_key       = "my_value"
	}
}
`, template, template, instance)
}
