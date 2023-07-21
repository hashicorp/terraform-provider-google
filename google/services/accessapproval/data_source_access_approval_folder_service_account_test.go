// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accessapproval_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceAccessApprovalFolderServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	resourceName := "data.google_access_approval_folder_service_account.aa_account"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessApprovalFolderServiceAccount_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_email"),
				),
			},
		},
	})
}

func testAccDataSourceAccessApprovalFolderServiceAccount_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "my_folder" {
  display_name = "tf-test-my-folder%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

# Wait after folder creation to limit eventual consistency errors.
resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.my_folder]

  create_duration = "120s"
}

data "google_access_approval_folder_service_account" "aa_account" {
  folder_id = google_folder.my_folder.folder_id

  depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}
