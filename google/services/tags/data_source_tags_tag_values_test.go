// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tags_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleTagsTagValues_default(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	keyShortName := "tf-testkey-" + acctest.RandString(t, 10)
	shortName := "tf-test-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleTagsTagValuesConfig(parent, keyShortName, shortName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleTagsTagValuesCheck("data.google_tags_tag_values.my_tag_values", "google_tags_tag_value.norfqux"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleTagsTagValues_dot(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	keyShortName := "tf-testkey-" + acctest.RandString(t, 10)
	shortName := "terraform.test." + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleTagsTagValuesConfig(parent, keyShortName, shortName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleTagsTagValuesCheck("data.google_tags_tag_values.my_tag_values", "google_tags_tag_value.norfqux"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleTagsTagValuesCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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
		tag_value_attrs_to_test := []string{"parent", "name", "namespaced_name", "create_time", "update_time", "description"}

		for _, attr_to_check := range tag_value_attrs_to_test {
			data := ""
			if attr_to_check == "name" {
				data = strings.Split(ds_attr[fmt.Sprintf("values.0.%s", attr_to_check)], "/")[1]
			} else {
				data = ds_attr[fmt.Sprintf("values.0.%s", attr_to_check)]
			}
			if data != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					data,
					rs_attr[attr_to_check],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceGoogleTagsTagValuesConfig(parent string, keyShortName string, shortName string) string {
	return fmt.Sprintf(`
resource "google_tags_tag_key" "foobar" {
  parent     = "%s"
  short_name = "%s"
}

resource "google_tags_tag_value" "norfqux" {
  parent     = google_tags_tag_key.foobar.id
  short_name = "%s"
}

data "google_tags_tag_values" "my_tag_values" {
  parent     = google_tags_tag_value.norfqux.parent
}
`, parent, keyShortName, shortName)
}
