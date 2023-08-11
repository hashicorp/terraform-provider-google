// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleProject_basic(t *testing.T) {
	t.Parallel()
	org := envvar.GetTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectConfig(project, org),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_project.project",
						"google_project.project",
						map[string]struct{}{
							// Virtual fields
							"auto_create_network": {},
							"skip_delete":         {},
						}),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectConfig(project, org string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

data "google_project" "project" {
  project_id = google_project.project.project_id
}
`, project, project, org)
}
