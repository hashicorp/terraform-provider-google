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

func TestAccDataSourceComputeNetworkEndpointGroup(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeNetworkEndpointGroupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeNetworkEndpointGroupCheck("data.google_compute_network_endpoint_group.bar", "google_compute_network_endpoint_group.neg"),
				),
			},
		},
	})
}

func testAccDataSourceComputeNetworkEndpointGroupCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes
		network_attrs_to_test := []string{
			"self_link",
			"name",
			"zone",
			"description",
		}

		for _, attr_to_check := range network_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}
		return nil
	}
}

func testAccDataSourceComputeNetworkEndpointGroupConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint_group" "neg" {
  name         = "tf-test-my-lb-ds-neg%{random_suffix}"
  network      = "${google_compute_network.default.self_link}"
  subnetwork   = "${google_compute_subnetwork.default.self_link}"
  default_port = "90"
  zone         = "us-central1-a"
}

resource "google_compute_network" "default" {
  name = "tf-test-ds-neg-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-ds-neg-subnetwork%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = "${google_compute_network.default.self_link}"
}

data "google_compute_network_endpoint_group" "bar" {
        name = "${google_compute_network_endpoint_group.neg.name}"
        zone = "us-central1-a"
}
`, context)
}
