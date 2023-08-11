// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeInstance_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeInstanceConfig(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeInstanceCheck("data.google_compute_instance.bar", "google_compute_instance.foo"),
					resource.TestCheckResourceAttr("data.google_compute_instance.bar", "network_interface.#", "1"),
					resource.TestCheckResourceAttr("data.google_compute_instance.bar", "boot_disk.0.initialize_params.0.size", "10"),
					resource.TestCheckResourceAttr("data.google_compute_instance.bar", "boot_disk.0.initialize_params.0.type", "pd-standard"),
					resource.TestCheckResourceAttr("data.google_compute_instance.bar", "scratch_disk.0.interface", "SCSI"),
					resource.TestCheckResourceAttr("data.google_compute_instance.bar", "network_interface.0.access_config.0.network_tier", "PREMIUM"),
					resource.TestCheckResourceAttr("data.google_compute_instance.bar", "enable_display", "true"),
				),
			},
		},
	})
}

func testAccDataSourceComputeInstanceCheck(datasourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		datasourceAttributes := ds.Primary.Attributes
		resourceAttributes := rs.Primary.Attributes

		instanceAttrsToTest := []string{
			"name",
			"machine_type",
			"current_status",
			"can_ip_forward",
			"description",
			"deletion_protection",
			"labels",
			"metadata",
			"min_cpu_platform",
			"project",
			"tags",
			"zone",
			"cpu_platform",
			"instance_id",
			"label_fingerprint",
			"metadata_fingerprint",
			"self_link",
			"tags_fingerprint",
		}

		for _, attrToCheck := range instanceAttrsToTest {
			if datasourceAttributes[attrToCheck] != resourceAttributes[attrToCheck] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attrToCheck,
					datasourceAttributes[attrToCheck],
					resourceAttributes[attrToCheck],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceComputeInstanceConfig(instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "foo" {
  name           = "%s"
  machine_type   = "n1-standard-1"   // can't be e2 because of local-ssd
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = "debian-8-jessie-v20160803"
    }
  }

  scratch_disk {
	interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    foo            = "bar"
    baz            = "qux"
    startup-script = "echo Hello"
  }

  labels = {
    my_key       = "my_value"
    my_other_key = "my_other_value"
  }

  enable_display = true
}

data "google_compute_instance" "bar" {
  name = google_compute_instance.foo.name
  zone = "us-central1-a"
}

data "google_compute_instance" "baz" {
  self_link = google_compute_instance.foo.self_link
}
`, instanceName)
}
