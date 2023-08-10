// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package containerattached_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleContainerAttachedVersions(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleContainerAttachedVersionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleContainerAttachedVersionsCheck("data.google_container_attached_versions.versions"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleContainerAttachedVersionsCheck(data_source_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		count, ok := ds.Primary.Attributes["valid_versions.#"]
		if !ok {
			return fmt.Errorf("cannot find 'valid_versions' attribute")
		}
		noOfVersions, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of versions")
		}
		if noOfVersions < 1 {
			return fmt.Errorf("expected at least 1 version, received %d", noOfVersions)
		}

		for i := 0; i < noOfVersions; i++ {
			idx := "valid_versions." + strconv.Itoa(i)
			v, ok := ds.Primary.Attributes[idx]
			if !ok {
				return fmt.Errorf("versions list is corrupt (%q not found)", idx)
			}
			if v == "" {
				return fmt.Errorf("empty version returned for %q", idx)
			}
		}
		return nil
	}
}

func testAccDataSourceGoogleContainerAttachedVersionsConfig() string {
	return `
data "google_project" "project" {
}

data "google_container_attached_versions" "versions" {
	location       = "us-west1"
	project        = data.google_project.project.project_id
}
`
}
