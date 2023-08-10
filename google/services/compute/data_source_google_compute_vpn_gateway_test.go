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

func TestAccDataSourceGoogleVpnGateway(t *testing.T) {
	t.Parallel()

	vpnGatewayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleVpnGatewayConfig(vpnGatewayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleVpnGatewayCheck("data.google_compute_vpn_gateway.my_vpn_gateway", "google_compute_vpn_gateway.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleVpnGatewayCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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
		vpn_gateway_attrs_to_test := []string{
			"id",
			"name",
			"description",
			"network",
		}

		for _, attr_to_check := range vpn_gateway_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if !tpgresource.CompareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		return nil
	}
}

func testAccDataSourceGoogleVpnGatewayConfig(name string) string {
	return fmt.Sprintf(`
resource "google_compute_vpn_gateway" "foobar" {
  name        = "%s"
  description = "my-description"
  network     = "default"
}

data "google_compute_vpn_gateway" "my_vpn_gateway" {
  name = google_compute_vpn_gateway.foobar.name
}
`, name)
}
