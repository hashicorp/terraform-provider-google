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

func TestAccDataSourceGoogleProjectIamCustomRole_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	roleId := "tfIamCustomRole" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamCustomRoleConfig(project, roleId),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_project_iam_custom_role.this",
						"google_project_iam_custom_role.this",
					),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectIamCustomRoleConfig(project string, roleId string) string {
	return fmt.Sprintf(`
locals {
  project = "%s"
  role_id = "%s"
}

resource "google_project_iam_custom_role" "this" {
  project = local.project
  role_id = local.role_id
  title   = "Terraform Test"

  permissions = [
	"iam.roles.create",
	"iam.roles.delete",
    "iam.roles.list",
  ]
}

data "google_project_iam_custom_role" "this" {
  project = google_project_iam_custom_role.this.project
  role_id = google_project_iam_custom_role.this.role_id
}
`, project, roleId)
}
