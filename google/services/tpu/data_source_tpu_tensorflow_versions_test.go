// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpu_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccTPUTensorflowVersions_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccTpuTensorFlowVersionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleTpuTensorflowVersions("data.google_tpu_tensorflow_versions.available"),
				),
			},
		},
	})
}

func testAccCheckGoogleTpuTensorflowVersions(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find TPU Tensorflow versions data source: %s", n)
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
			return errors.New("failed to read number of version")
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

var testAccTpuTensorFlowVersionsConfig = `
data "google_tpu_tensorflow_versions" "available" {}
`
