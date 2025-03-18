// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleOrganizationIamCustomRole_basic(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	roleId := "tfIamCustomRole" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRoleConfig(orgId, roleId),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_organization_iam_custom_role.this",
						"google_organization_iam_custom_role.this",
					),
				),
			},
		},
	})
}

func testAccCheckGoogleOrganizationIamCustomRoleConfig(orgId string, roleId string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_custom_role" "this" {
  org_id      = "%s"
  role_id     = "%s"
  title       = "Terraform Test"

  permissions = [
    "iam.roles.create",
    "iam.roles.delete",
    "iam.roles.list",
  ]
}

data "google_organization_iam_custom_role" "this" {
  org_id      = google_organization_iam_custom_role.this.org_id
  role_id     = google_organization_iam_custom_role.this.role_id
}
`, orgId, roleId)
}
