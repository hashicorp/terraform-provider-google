// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	compute "google.golang.org/api/compute/v1"
)

func TestAccDataSourceComputeInstanceGuestAttributes_basic(t *testing.T) {
	t.Parallel()

	var instance compute.Instance
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				//need to create the guest_attributes metadata from startup script first
				Config: testAccDataSourceComputeInstanceGuestAttributesInitialConfig(instanceName),
				Check:  testAccCheckComputeInstanceExists(t, "google_compute_instance.foo", &instance),
			},
			{
				Config: testAccDataSourceComputeInstanceGuestAttributesConfig_variableKey(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_instance_guest_attributes.bar", "variable_key", "testing/key2"),
					resource.TestCheckResourceAttr("data.google_compute_instance_guest_attributes.bar", "variable_value", "test2"),
				),
			},
			{
				Config: testAccDataSourceComputeInstanceGuestAttributesConfig_queryPath(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_instance_guest_attributes.bar", "query_path", "testing/"),
					resource.TestCheckResourceAttr("data.google_compute_instance_guest_attributes.bar", "query_value.0.value", "test1"),
					resource.TestCheckResourceAttr("data.google_compute_instance_guest_attributes.bar", "query_value.1.value", "test2"),
				),
			},
		},
	})
}

func testAccDataSourceComputeInstanceGuestAttributesInitialConfig(instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "foo" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  zone           = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-8-jessie-v20160803"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    enable-guest-attributes = "TRUE"
  }

  metadata_startup_script = <<-EOF
  curl -X PUT --data "test1" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/testing/key1 -H "Metadata-Flavor: Google"
  curl -X PUT --data "test2" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/testing/key2 -H "Metadata-Flavor: Google"
  EOF
}
`, instanceName)
}

func testAccDataSourceComputeInstanceGuestAttributesConfig_queryPath(instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "foo" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  zone           = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-8-jessie-v20160803"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    enable-guest-attributes = "TRUE"
  }

  metadata_startup_script = <<-EOF
  curl -X PUT --data "test1" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/testing/key1 -H "Metadata-Flavor: Google"
  curl -X PUT --data "test2" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/testing/key2 -H "Metadata-Flavor: Google"
  EOF
}

data "google_compute_instance_guest_attributes" "bar" {
  name = google_compute_instance.foo.name
  zone = "us-central1-a"
  query_path = "testing/"
}
`, instanceName)
}

func testAccDataSourceComputeInstanceGuestAttributesConfig_variableKey(instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "foo" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  zone           = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-8-jessie-v20160803"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    enable-guest-attributes = "TRUE"
  }

  metadata_startup_script = <<-EOF
  curl -X PUT --data "test1" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/testing/key1 -H "Metadata-Flavor: Google"
  curl -X PUT --data "test2" http://metadata.google.internal/computeMetadata/v1/instance/guest-attributes/testing/key2 -H "Metadata-Flavor: Google"
  EOF
}

data "google_compute_instance_guest_attributes" "bar" {
  name = google_compute_instance.foo.name
  zone = "us-central1-a"
  variable_key = "testing/key2"
}
`, instanceName)
}
