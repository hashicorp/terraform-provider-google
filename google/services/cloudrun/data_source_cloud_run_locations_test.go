// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrun_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGoogleCloudRunLocations_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudRunLocationsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleCloudRunLocations("data.google_cloud_run_locations.default"),
				),
			},
		},
	})
}

func testAccCheckGoogleCloudRunLocations(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find cloud run locations data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("data source id not set")
		}

		count, ok := rs.Primary.Attributes["locations.#"]
		if !ok {
			return errors.New("can't find 'locations' attribute")
		}

		cnt, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of locations")
		}
		if cnt < 5 {
			return fmt.Errorf("expected at least 5 locations, received %d, this is most likely a bug", cnt)
		}

		for i := 0; i < cnt; i++ {
			idx := fmt.Sprintf("locations.%d", i)
			_, ok := rs.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("expected %q, location not found", idx)
			}
		}
		return nil
	}
}

const testAccDataSourceGoogleCloudRunLocationsBasic = `
data "google_cloud_run_locations" "default" {}
`
