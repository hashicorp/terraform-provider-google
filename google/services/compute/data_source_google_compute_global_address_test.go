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

func TestAccDataSourceComputeGlobalAddress(t *testing.T) {
	t.Parallel()

	rsName := "foobar"
	rsFullName := fmt.Sprintf("google_compute_global_address.%s", rsName)
	dsName := "my_address"
	dsFullName := fmt.Sprintf("data.google_compute_global_address.%s", dsName)
	addressName := fmt.Sprintf("tf-test-address-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeGlobalAddressConfig(rsName, dsName, addressName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeGlobalAddressCheck(dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeGlobalAddressCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

		address_attrs_to_test := []string{
			"name",
			"address",
		}

		for _, attr_to_check := range address_attrs_to_test {
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

		if ds_attr["status"] != "RESERVED" {
			return fmt.Errorf("status is %s; want RESERVED", ds_attr["status"])
		}

		return nil
	}
}

func testAccDataSourceComputeGlobalAddressConfig(rsName, dsName, addressName string) string {
	return fmt.Sprintf(`
resource "google_compute_global_address" "%s" {
  name = "%s"
}

data "google_compute_global_address" "%s" {
  name = google_compute_global_address.%s.name
}
`, rsName, addressName, dsName, rsName)
}
