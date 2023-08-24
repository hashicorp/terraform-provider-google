// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceAlloydbLocations_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlloydbLocations_basic(context),
				Check: resource.ComposeTestCheckFunc(
					validateAlloydbLocationsResult(
						"data.google_alloydb_locations.qa",
					),
				),
			},
		},
	})
}

func testAccDataSourceAlloydbLocations_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_alloydb_locations" "qa" {
}
`, context)
}

func validateAlloydbLocationsResult(dataSourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}
		var dsAttr map[string]string
		dsAttr = ds.Primary.Attributes

		totalFlags, err := strconv.Atoi(dsAttr["locations.#"])
		if err != nil {
			return errors.New("Couldn't convert length of flags list to integer")
		}
		if totalFlags == 0 {
			return errors.New("No locations are fetched")
		}
		for i := 0; i < totalFlags; i++ {
			if dsAttr["locations."+strconv.Itoa(i)+".name"] == "" {
				return errors.New("name parameter is not set for the location")
			}
			if dsAttr["locations."+strconv.Itoa(i)+".location_id"] == "" {
				return errors.New("location_id parameter is not set for the location")
			}
		}
		return nil
	}
}
