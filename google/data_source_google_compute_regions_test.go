package google

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeRegions_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleComputeRegionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComputeRegionsMeta("data.google_compute_regions.available"),
				),
			},
		},
	})
}

func testAccCheckGoogleComputeRegionsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find regions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("regions data source ID not set.")
		}

		count, ok := rs.Primary.Attributes["names.#"]
		if !ok {
			return errors.New("can't find 'names' attribute")
		}

		noOfNames, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of regions")
		}
		if noOfNames < 2 {
			return fmt.Errorf("expected at least 2 regions, received %d, this is most likely a bug",
				noOfNames)
		}

		for i := 0; i < noOfNames; i++ {
			idx := "names." + strconv.Itoa(i)
			v, ok := rs.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("region list is corrupt (%q not found), this is definitely a bug", idx)
			}
			if len(v) < 1 {
				return fmt.Errorf("Empty region name (%q), this is definitely a bug", idx)
			}
		}

		return nil
	}
}

var testAccCheckGoogleComputeRegionsConfig = `
data "google_compute_regions" "available" {}
`
