package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGoogleTagsTagKey_default(t *testing.T) {
	org := GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	shortName := "tf-test-" + RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleTagsTagKeyConfig(parent, shortName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleTagsTagKeyCheck("data.google_tags_tag_key.my_tag_key", "google_tags_tag_key.foobar"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleTagsTagKey_dot(t *testing.T) {
	org := GetTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	shortName := "terraform.test." + RandString(t, 10)

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleTagsTagKeyConfig(parent, shortName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleTagsTagKeyCheck("data.google_tags_tag_key.my_tag_key", "google_tags_tag_key.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleTagsTagKeyCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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
		tag_key_attrs_to_test := []string{"parent", "short_name", "name", "namespaced_name", "create_time", "update_time", "description"}

		for _, attr_to_check := range tag_key_attrs_to_test {
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

func testAccDataSourceGoogleTagsTagKeyConfig(parent string, shortName string) string {
	return fmt.Sprintf(`
resource "google_tags_tag_key" "foobar" {
  parent     = "%s"
  short_name = "%s"
}

data "google_tags_tag_key" "my_tag_key" {
  parent     = google_tags_tag_key.foobar.parent
  short_name = google_tags_tag_key.foobar.short_name
}
`, parent, shortName)
}
