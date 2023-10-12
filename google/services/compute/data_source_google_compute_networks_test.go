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

func TestAccDataSourceGoogleNetworks(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleNetworksConfig(networkName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleNetworksCheck("data.google_compute_networks.my_networks", "google_compute_network.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleNetworksCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

		containsNetwork := false

		for _, itm := range ds_attr {
			if string(itm) == rs_attr["name"] {
				containsNetwork = true
				break
			}
		}

		if !containsNetwork {
			return fmt.Errorf(
				"Was expecting %s in %v",
				rs_attr["name"],
				ds_attr["networks"],
			)
		}

		return nil
	}
}

func testAccDataSourceGoogleNetworksConfig(name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name        = "%s"
  description = "my-description"
}

data "google_compute_networks" "my_networks" {
	depends_on = [
		google_compute_network.foobar
	]
}
`, name)
}
