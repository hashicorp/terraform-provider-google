// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accessapproval_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceAccessApprovalProjectServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id": envvar.GetTestProjectFromEnv(),
	}

	resourceName := "data.google_access_approval_project_service_account.aa_account"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessApprovalProjectServiceAccount_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_email"),
				),
			},
		},
	})
}

func testAccDataSourceAccessApprovalProjectServiceAccount_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_access_approval_project_service_account" "aa_account" {
  project_id = "%{project_id}"
}
`, context)
}
