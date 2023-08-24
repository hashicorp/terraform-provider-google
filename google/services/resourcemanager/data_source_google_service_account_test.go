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

func TestAccDatasourceGoogleServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account.acceptance"
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccount_basic(account),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "id", fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", envvar.GetTestProjectFromEnv(), account, envvar.GetTestProjectFromEnv())),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "unique_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "member"),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccount_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%s"
  display_name = "Testing Account"
}

data "google_service_account" "acceptance" {
  account_id = google_service_account.acceptance.account_id
}
`, account)
}
