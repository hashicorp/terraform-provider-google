package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleProject(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
		}...,
	)

	name := "foobar"
	randSuffix := acctest.RandString(10)
	pid := fmt.Sprintf("%s-%s", name, randSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleProjectConfig(pid, name, org),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleProjectCheck("data.google_project.my_project", "google_project.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleProjectCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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

		project_attrs_to_test := []string{
			"project_id",
			"number",
			"name",
		}

		for _, attr_to_check := range project_attrs_to_test {
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

func testAccDataSourceGoogleProjectConfig(pid string, name string, org string) string {
	return fmt.Sprintf(`
resource "google_project" "foobar" {
	project_id = "%s"
	name = "%s"
	org_id = "%s"
}

data "google_project" "my_project" {
	name = "${google_project.foobar.name}"
}`, pid, name, org)
}
