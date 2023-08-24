// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComposerImageVersions_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleComposerImageVersionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleComposerImageVersionsMeta("data.google_composer_image_versions.versions"),
				),
			},
		},
	})
}

func testAccCheckGoogleComposerImageVersionsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find versions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("versions data source ID not set.")
		}

		versionCountStr, ok := rs.Primary.Attributes["image_versions.#"]
		if !ok {
			return errors.New("can't find 'image_versions' attribute")
		}

		versionCount, err := strconv.Atoi(versionCountStr)
		if err != nil {
			return errors.New("failed to read number of valid image versions")
		}
		if versionCount < 1 {
			return fmt.Errorf("expected at least 1 valid image versions, received %d, this is most likely a bug",
				versionCount)
		}

		for i := 0; i < versionCount; i++ {
			idx := "image_versions." + strconv.Itoa(i)
			if v, ok := rs.Primary.Attributes[idx+".image_version_id"]; !ok || v == "" {
				return fmt.Errorf("image_version %v is missing image_version_id", i)
			}
			if v, ok := rs.Primary.Attributes[idx+".supported_python_versions.#"]; !ok || v == "" || v == "0" {
				return fmt.Errorf("image_version %v is missing supported_python_versions", i)
			}
		}

		return nil
	}
}

var testAccCheckGoogleComposerImageVersionsConfig = `
data "google_composer_image_versions" "versions" {
}
`
