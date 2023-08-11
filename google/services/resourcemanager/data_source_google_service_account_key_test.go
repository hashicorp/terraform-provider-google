// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceGoogleServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_key.acceptance"
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	serviceAccountName := fmt.Sprintf(
		"projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com",
		envvar.GetTestProjectFromEnv(),
		account,
		envvar.GetTestProjectFromEnv(),
	)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceGoogleServiceAccountKey(account),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleServiceAccountKeyExists(t, resourceName),
					// Check that the 'name' starts with the service account name
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(serviceAccountName)),
					resource.TestCheckResourceAttrSet(resourceName, "key_algorithm"),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
				),
			},
		},
	})
}

func testAccDatasourceGoogleServiceAccountKey(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id = "%s"
}

resource "google_service_account_key" "acceptance" {
  service_account_id = google_service_account.acceptance.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

data "google_service_account_key" "acceptance" {
  name = google_service_account_key.acceptance.name
}
`, account)
}
