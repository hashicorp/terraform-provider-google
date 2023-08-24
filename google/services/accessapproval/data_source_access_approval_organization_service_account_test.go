// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accessapproval_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceAccessApprovalOrganizationServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": envvar.GetTestOrgFromEnv(t),
	}

	resourceName := "data.google_access_approval_organization_service_account.aa_account"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessApprovalOrganizationServiceAccount_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_email"),
				),
			},
		},
	})
}

func testAccDataSourceAccessApprovalOrganizationServiceAccount_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_access_approval_organization_service_account" "aa_account" {
  organization_id = "%{org_id}"
}
`, context)
}
