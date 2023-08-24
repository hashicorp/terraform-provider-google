// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeZones_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeZones_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComputeZonesMeta("data.google_compute_zones.available"),
				),
			},
		},
	})
}

func TestAccComputeZones_filter(t *testing.T) {
	t.Parallel()
	region := "us-central1"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeZones_filter(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComputeZonesRegion("data.google_compute_zones.available", region),
				),
			},
		},
	})
}

func testAccCheckGoogleComputeZonesMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find zones data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("zones data source ID not set.")
		}

		count, ok := rs.Primary.Attributes["names.#"]
		if !ok {
			return errors.New("can't find 'names' attribute")
		}

		noOfNames, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of zones")
		}
		if noOfNames < 2 {
			return fmt.Errorf("expected at least 2 zones, received %d, this is most likely a bug",
				noOfNames)
		}

		for i := 0; i < noOfNames; i++ {
			idx := "names." + strconv.Itoa(i)
			v, ok := rs.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("zone list is corrupt (%q not found), this is definitely a bug", idx)
			}
			if len(v) < 1 {
				return fmt.Errorf("Empty zone name (%q), this is definitely a bug", idx)
			}
		}

		return nil
	}
}

func testAccCheckGoogleComputeZonesRegion(n, region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find zones data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("zones data source ID not set.")
		}

		count, ok := rs.Primary.Attributes["names.#"]
		if !ok {
			return errors.New("can't find 'names' attribute")
		}

		noOfNames, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of zones")
		}

		for i := 0; i < noOfNames; i++ {
			idx := "names." + strconv.Itoa(i)
			v, ok := rs.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("zone list is corrupt (%q not found), this is definitely a bug", idx)
			}
			if !strings.Contains(v, region) {
				return fmt.Errorf("zone name %q does not contain region %q", v, region)
			}
		}

		return nil
	}
}

var testAccComputeZones_basic = `
data "google_compute_zones" "available" {}
`

func testAccComputeZones_filter(region string) string {
	return fmt.Sprintf(`
data "google_compute_zones" "available" {
  region = "%s"
  status = "UP"
}
`, region)
}
