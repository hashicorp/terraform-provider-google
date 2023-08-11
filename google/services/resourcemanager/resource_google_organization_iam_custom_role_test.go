// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOrganizationIamCustomRole_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationIamCustomRoleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					t,
					"google_organization_iam_custom_role.foo",
					"My Custom Role",
					"foo",
					"GA",
					[]string{"resourcemanager.projects.list"}),
			},
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_update(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					t,
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

	org := envvar.GetTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationIamCustomRoleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check:  testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus(t, "google_organization_iam_custom_role.foo", false),
			},
			// Soft-delete
			{
				Config:  testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check:   testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus(t, "google_organization_iam_custom_role.foo", true),
				Destroy: true,
			},
			// Undelete
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check:  testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus(t, "google_organization_iam_custom_role.foo", false),
			},
		},
	})
}

func TestAccOrganizationIamCustomRole_createAfterDestroy(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationIamCustomRoleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_basic(org, roleId),
				Check: testAccCheckGoogleOrganizationIamCustomRole(
					t,
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
					t,
					"google_organization_iam_custom_role.foo",
					"My Custom Role",
					"foo",
					"GA",
					[]string{"resourcemanager.projects.list"}),
			},
		},
	})
}

func testAccCheckGoogleOrganizationIamCustomRoleDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_organization_iam_custom_role" {
				continue
			}

			role, err := config.NewIamClient(config.UserAgent).Organizations.Roles.Get(rs.Primary.ID).Do()

			if err != nil {
				return err
			}

			if !role.Deleted {
				return fmt.Errorf("Iam custom role still exists")
			}

		}

		return nil
	}
}

func testAccCheckGoogleOrganizationIamCustomRole(t *testing.T, n, title, description, stage string, permissions []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)
		role, err := config.NewIamClient(config.UserAgent).Organizations.Roles.Get(rs.Primary.ID).Do()

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

func testAccCheckGoogleOrganizationIamCustomRoleDeletionStatus(t *testing.T, n string, deleted bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)
		role, err := config.NewIamClient(config.UserAgent).Organizations.Roles.Get(rs.Primary.ID).Do()

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
