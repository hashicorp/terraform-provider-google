// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGoogleFolder_byFullName(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFolder_byFullNameConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleFolderCheck("data.google_folder.folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_byShortName(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFolder_byShortNameConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleFolderCheck("data.google_folder.folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_lookupOrganization(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFolder_lookupOrganizationConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleFolderCheck("data.google_folder.folder", "google_folder.foobar"),
					resource.TestCheckResourceAttr("data.google_folder.folder", "organization", parent),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_byFullNameNotFound(t *testing.T) {
	name := "folders/" + acctest.RandString(t, 16)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGoogleFolder_byFullNameNotFoundConfig(name),
				ExpectError: regexp.MustCompile("Folder Not Found : " + name),
			},
		},
	})
}

func testAccDataSourceGoogleFolderCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

func testAccCheckGoogleFolder_byFullNameConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
  parent       = "%s"
  display_name = "%s"
}

data "google_folder" "folder" {
  folder = google_folder.foobar.name
}
`, parent, displayName)
}

func testAccCheckGoogleFolder_byShortNameConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
  parent       = "%s"
  display_name = "%s"
}

data "google_folder" "folder" {
  folder = replace(google_folder.foobar.name, "folders/", "")
}
`, parent, displayName)
}

func testAccCheckGoogleFolder_lookupOrganizationConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
  parent       = "%s"
  display_name = "%s"
}

data "google_folder" "folder" {
  folder              = google_folder.foobar.name
  lookup_organization = true
}
`, parent, displayName)
}

func testAccCheckGoogleFolder_byFullNameNotFoundConfig(name string) string {
	return fmt.Sprintf(`
data "google_folder" "folder" {
  folder = "%s"
}
`, name)
}
