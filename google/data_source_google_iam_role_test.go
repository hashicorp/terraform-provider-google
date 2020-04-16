package google

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceIAMRole(t *testing.T) {
	name := "roles/viewer"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleIamRoleConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleIAMRoleCheck("data.google_iam_role.role"),
				),
			},
		},
	})
}

func testAccCheckGoogleIAMRoleCheck(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find iam role data source: %s", n)
		}

		_, ok = ds.Primary.Attributes["included_permissions.#"]
		if !ok {
			return errors.New("can't find 'included_permissions' attribute")
		}

		return nil
	}
}

func testAccCheckGoogleIamRoleConfig(name string) string {
	return fmt.Sprintf(`
data "google_iam_role" "role" {
  name = "%s"
}
`, name)
}
