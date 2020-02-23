package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOrganizationIamCustomRole_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					"google_organization_iam_custom_role.foo",
					"My Custom Role",
					"foo",
					"GA",
					[]string{"resourcemanager.projects.list"}),
			},
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_update(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					"google_organization_iam_custom_role.foo",
					"My Custom Role Updated",
					"bar",
					"BETA",
					[]string{"resourcemanager.projects.list", "resourcemanager.organizations.get"}),
			},
			{
				ResourceName:      "google_organization_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationIamCustomRole_undelete(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check:  testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus("google_organization_iam_custom_role.foo", false),
			},
			// Soft-delete
			{
				Config:  testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check:   testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus("google_organization_iam_custom_role.foo", true),
				Destroy: true,
			},
			// Undelete
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check:  testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus("google_organization_iam_custom_role.foo", false),
			},
		},
	})
}

func TestAccOrganizationIamCustomRole_createAfterDestroy(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					"google_organization_iam_custom_role.foo",
					"My Custom Role",
					"foo",
					"GA",
					[]string{"resourcemanager.projects.list"}),
			},
			// Destroy resources
			{
				Config:  " ",
				Destroy: true,
			},
			// Re-create with no existing state
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					"google_organization_iam_custom_role.foo",
					"My Custom Role",
					"foo",
					"GA",
					[]string{"resourcemanager.projects.list"}),
			},
		},
	})
}

func testAccCheckGoogleOrganizationIamCustomRoleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_organization_iam_custom_role" {
			continue
		}

		role, err := config.clientIAM.Organizations.Roles.Get(rs.Primary.ID).Do()

		if err != nil {
			return err
		}

		if !role.Deleted {
			return fmt.Errorf("Iam custom role still exists")
		}

	}

	return nil
}

func testAccCheckGoogleOrganizationIamCustomRole(n, title, description, stage string, permissions []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		role, err := config.clientIAM.Organizations.Roles.Get(rs.Primary.ID).Do()

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

func testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus(n string, deleted bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		role, err := config.clientIAM.Organizations.Roles.Get(rs.Primary.ID).Do()

		if err != nil {
			return err
		}

		if deleted != role.Deleted {
			return fmt.Errorf("Incorrect deletion status. Expected %t, got %t", deleted, role.Deleted)
		}

		return nil
	}
}

func testAccCheckGoogleOrganizationIamCustomRole_basic(orgId, roleId string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_custom_role" "foo" {
  role_id     = "%s"
  org_id      = "%s"
  title       = "My Custom Role"
  description = "foo"
  permissions = ["resourcemanager.projects.list"]
}
`, roleId, orgId)
}

func testAccCheckGoogleOrganizationIamCustomRole_update(orgId, roleId string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_custom_role" "foo" {
  role_id     = "%s"
  org_id      = "%s"
  title       = "My Custom Role Updated"
  description = "bar"
  permissions = ["resourcemanager.projects.list", "resourcemanager.organizations.get"]
  stage       = "BETA"
}
`, roleId, orgId)
}
