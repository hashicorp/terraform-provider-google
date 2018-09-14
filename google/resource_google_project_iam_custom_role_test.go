package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccProjectIamCustomRole_basic(t *testing.T) {
	t.Parallel()

	roleId := "tfIamCustomRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamCustomRole_basic(roleId),
				Check:  resource.TestCheckResourceAttr("google_project_iam_custom_role.foo", "stage", "GA"),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckGoogleProjectIamCustomRole_update(roleId),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectIamCustomRole_undelete(t *testing.T) {
	t.Parallel()

	roleId := "tfIamCustomRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamCustomRole_basic(roleId),
				Check:  resource.TestCheckResourceAttr("google_project_iam_custom_role.foo", "deleted", "false"),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Soft-delete
			{
				Config: testAccCheckGoogleProjectIamCustomRole_deleted(roleId),
				Check:  resource.TestCheckResourceAttr("google_project_iam_custom_role.foo", "deleted", "true"),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Undelete
			{
				Config: testAccCheckGoogleProjectIamCustomRole_basic(roleId),
				Check:  resource.TestCheckResourceAttr("google_project_iam_custom_role.foo", "deleted", "false"),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectIamCustomRole_createAfterDestroy(t *testing.T) {
	t.Parallel()

	roleId := "tfIamCustomRole" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamCustomRole_basic(roleId),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Destroy resources
			{
				Config:  " ",
				Destroy: true,
			},
			// Re-create with no existing state
			{
				Config: testAccCheckGoogleProjectIamCustomRole_basic(roleId),
			},
			{
				ResourceName:      "google_project_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleProjectIamCustomRoleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_project_iam_custom_role" {
			continue
		}

		role, err := config.clientIAM.Projects.Roles.Get(rs.Primary.ID).Do()

		if err != nil {
			return err
		}

		if !role.Deleted {
			return fmt.Errorf("Iam custom role still exists")
		}

	}

	return nil
}

func testAccCheckGoogleProjectIamCustomRole_basic(roleId string) string {
	return fmt.Sprintf(`
resource "google_project_iam_custom_role" "foo" {
  role_id     = "%s"
  title       = "My Custom Role"
  description = "foo"
  permissions = ["iam.roles.list"]
}
`, roleId)
}

func testAccCheckGoogleProjectIamCustomRole_deleted(roleId string) string {
	return fmt.Sprintf(`
resource "google_project_iam_custom_role" "foo" {
  role_id     = "%s"
  title       = "My Custom Role"
  description = "foo"
  permissions = ["iam.roles.list"]
  deleted     = true
}
`, roleId)
}

func testAccCheckGoogleProjectIamCustomRole_update(roleId string) string {
	return fmt.Sprintf(`
resource "google_project_iam_custom_role" "foo" {
  role_id     = "%s"
  title       = "My Custom Role Updated"
  description = "bar"
  permissions = ["iam.roles.list", "iam.roles.create", "iam.roles.delete"]
  stage       = "BETA"
}
`, roleId)
}
