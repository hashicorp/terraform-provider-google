package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"reflect"
	"sort"
	"testing"
)

func TestAccGoogleProjectIamRole_basic(t *testing.T) {
	t.Parallel()

	roleId := "tfIamRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectIamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamRole_basic(roleId),
				Check: testAccCheckGoogleProjectIamRole(
					"google_project_iam_role.foo",
					"My Custom Role",
					"foo",
					"GA",
					[]string{"iam.roles.list"}),
			},
			{
				Config: testAccCheckGoogleProjectIamRole_update(roleId),
				Check: testAccCheckGoogleProjectIamRole(
					"google_project_iam_role.foo",
					"My Custom Role Updated",
					"bar",
					"BETA",
					[]string{"iam.roles.list", "iam.roles.create", "iam.roles.delete"}),
			},
		},
	})
}

func TestAccGoogleProjectIamRole_undelete(t *testing.T) {
	t.Parallel()

	roleId := "tfIamRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectIamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamRole_basic(roleId),
				Check:  testAccCheckGoogleProjectIamRoleDeletionStatus("google_project_iam_role.foo", false),
			},
			// Soft-delete
			{
				Config: testAccCheckGoogleProjectIamRole_deleted(roleId),
				Check:  testAccCheckGoogleProjectIamRoleDeletionStatus("google_project_iam_role.foo", true),
			},
			// Undelete
			{
				Config: testAccCheckGoogleProjectIamRole_basic(roleId),
				Check:  testAccCheckGoogleProjectIamRoleDeletionStatus("google_project_iam_role.foo", false),
			},
		},
	})
}

func testAccCheckGoogleProjectIamRoleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_project_iam_role" {
			continue
		}

		role, err := config.clientIAM.Projects.Roles.Get(rs.Primary.ID).Do()

		if err != nil {
			return err
		}

		if !role.Deleted {
			return fmt.Errorf("Iam role still exists")
		}

	}

	return nil
}

func testAccCheckGoogleProjectIamRole(n, title, description, stage string, permissions []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		role, err := config.clientIAM.Projects.Roles.Get(rs.Primary.ID).Do()

		if err != nil {
			return err
		}

		if title != role.Title {
			return fmt.Errorf("Incorrect title. Expected %q, got %q", title, role.Title)
		}

		if description != role.Description {
			return fmt.Errorf("Incorrect description. Expected %q, got %q", description, role.Description)
		}

		if stage != role.Stage {
			return fmt.Errorf("Incorrect stage. Expected %q, got %q", stage, role.Stage)
		}

		sort.Strings(permissions)
		sort.Strings(role.IncludedPermissions)
		if !reflect.DeepEqual(permissions, role.IncludedPermissions) {
			return fmt.Errorf("Incorrect permissions. Expected %q, got %q", permissions, role.IncludedPermissions)
		}

		return nil
	}
}

func testAccCheckGoogleProjectIamRoleDeletionStatus(n string, deleted bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		role, err := config.clientIAM.Projects.Roles.Get(rs.Primary.ID).Do()

		if err != nil {
			return err
		}

		if deleted != role.Deleted {
			return fmt.Errorf("Incorrect deletion status. Expected %t, got %t", deleted, role.Deleted)
		}

		return nil
	}
}

func testAccCheckGoogleProjectIamRole_basic(roleId string) string {
	return fmt.Sprintf(`
resource "google_project_iam_role" "foo" {
  role_id = "%s"
  title = "My Custom Role"
  description = "foo"
  permissions = ["iam.roles.list"]
}
`, roleId)
}

func testAccCheckGoogleProjectIamRole_deleted(roleId string) string {
	return fmt.Sprintf(`
resource "google_project_iam_role" "foo" {
  role_id = "%s"
  title = "My Custom Role"
  description = "foo"
  permissions = ["iam.roles.list"]
  deleted = true
}
`, roleId)
}

func testAccCheckGoogleProjectIamRole_update(roleId string) string {
	return fmt.Sprintf(`
resource "google_project_iam_role" "foo" {
  role_id = "%s"
  title = "My Custom Role Updated"
  description = "bar"
  permissions = ["iam.roles.list", "iam.roles.create", "iam.roles.delete"]
  stage = "BETA"
}
`, roleId)
}
