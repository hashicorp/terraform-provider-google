// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleActiveFolder_default(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleActiveFolderConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleActiveFolderCheck("data.google_active_folder.my_folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleActiveFolder_space(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "terraform test " + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleActiveFolderConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleActiveFolderCheck("data.google_active_folder.my_folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleActiveFolder_dash(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "terraform - test " + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleActiveFolderConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleActiveFolderCheck("data.google_active_folder.my_folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleActiveFolderCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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
		folder_attrs_to_test := []string{"parent", "display_name", "name"}

		for _, attr_to_check := range folder_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}
		return nil
	}
}

func testAccDataSourceGoogleActiveFolderConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
  parent       = "%s"
  display_name = "%s"
}

data "google_active_folder" "my_folder" {
  parent       = google_folder.foobar.parent
  display_name = google_folder.foobar.display_name
}
`, parent, displayName)
}
