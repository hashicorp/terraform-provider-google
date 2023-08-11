// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeResourcePolicy_attached(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeResourcePolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeResourcePolicy_attached(acctest.RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_resource_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeResourcePolicy_attached(suffix string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance" "foobar" {
  name           = "tf-test-%s"
  machine_type   = "e2-medium"
  zone           = "us-central1-a"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  //deletion_protection = false is implicit in this config due to default value

  boot_disk {
    initialize_params {
      image = data.google_compute_image.my_image.self_link
    }
  }

  network_interface {
    network = "default"
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

  resource_policies = [google_compute_resource_policy.foo.self_link]
}

resource "google_compute_resource_policy" "foo" {
  name   = "tf-test-policy-%s"
  region = "us-central1"
  group_placement_policy {
    availability_domain_count = 2
  }
}

`, suffix, suffix)
}
