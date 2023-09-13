// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpuv2_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccTpuV2RuntimeVersions_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTpuV2RuntimeVersionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTpuV2RuntimeVersions("data.google_tpu_v2_runtime_versions.available"),
				),
			},
		},
	})
}

func testAccCheckTpuV2RuntimeVersions(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TPU v2 runtime versions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("data source id not set")
		}

		count, ok := rs.Primary.Attributes["versions.#"]
		if !ok {
			return errors.New("can't find 'versions' attribute")
		}

		cnt, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of versions")
		}
		if cnt < 2 {
			return fmt.Errorf("expected at least 2 versions, received %d, this is most likely a bug", cnt)
		}

		for i := 0; i < cnt; i++ {
			idx := fmt.Sprintf("versions.%d", i)
			_, ok := rs.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("expected %q, version not found", idx)
			}
		}
		return nil
	}
}

const testAccTpuV2RuntimeVersionsConfig = `
data "google_tpu_v2_runtime_versions" "available" {}
`
