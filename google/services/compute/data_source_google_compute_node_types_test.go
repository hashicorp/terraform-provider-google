// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeNodeTypes_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeNodeTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComputeNodeTypes("data.google_compute_node_types.available"),
				),
			},
		},
	})
}

func testAccCheckGoogleComputeNodeTypes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find node types data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("node types data source ID not set.")
		}

		count, ok := rs.Primary.Attributes["names.#"]
		if !ok {
			return errors.New("can't find 'names' attribute")
		}

		cnt, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of node types")
		}
		if cnt < 1 {
			return fmt.Errorf("expected at least one node type, got %d", cnt)
		}

		for i := 0; i < cnt; i++ {
			idx := fmt.Sprintf("names.%d", i)
			v, ok := rs.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("expected %q, version not found", idx)
			}

			if !regexp.MustCompile(`-[0-9]+-[0-9]+$`).MatchString(v) {
				return fmt.Errorf("unexpected type format for %q, value is %v", idx, v)
			}
		}
		return nil
	}
}

var testAccDataSourceComputeNodeTypes_basic = `
data "google_compute_node_types" "available" {
	zone = "us-central1-a"
}
`
