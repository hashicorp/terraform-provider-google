// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccDataSourceGoogleSubnetwork(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleSubnetwork(fmt.Sprintf("tf-test-subnetwork-ds-%d", acctest.RandInt(t)), fmt.Sprintf("tf-test-subnetwork-ds-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSubnetworkCheck("data.google_compute_subnetwork.my_subnetwork", "google_compute_subnetwork.foobar"),
					testAccDataSourceGoogleSubnetworkCheck("data.google_compute_subnetwork.my_subnetwork_self_link", "google_compute_subnetwork.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSubnetworkCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

		subnetwork_attrs_to_test := []string{
			"id",
			"name",
			"description",
			"ip_cidr_range",
			"private_ip_google_access",
			"secondary_ip_range",
		}

		for _, attr_to_check := range subnetwork_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if !tpgresource.CompareSelfLinkOrResourceName("", ds_attr["network"], rs_attr["network"], nil) && ds_attr["network"] != rs_attr["network"] {
			return fmt.Errorf("network does not match: %s vs %s", ds_attr["network"], rs_attr["network"])
		}

		if !tpgresource.CompareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		return nil
	}
}

func testAccDataSourceGoogleSubnetwork(networkName, subnetName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name        = "%s"
  description = "my-description"
}

resource "google_compute_subnetwork" "foobar" {
  name                     = "%s"
  description              = "my-description"
  ip_cidr_range            = "10.0.0.0/24"
  network                  = google_compute_network.foobar.self_link
  private_ip_google_access = true
  secondary_ip_range {
    range_name    = "tf-test-secondary-range"
    ip_cidr_range = "192.168.1.0/24"
  }
}

data "google_compute_subnetwork" "my_subnetwork" {
  name = google_compute_subnetwork.foobar.name
}

data "google_compute_subnetwork" "my_subnetwork_self_link" {
  self_link = google_compute_subnetwork.foobar.self_link
}
`, networkName, subnetName)
}
