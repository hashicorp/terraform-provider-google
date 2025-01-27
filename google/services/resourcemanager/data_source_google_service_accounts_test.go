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

func TestAccDataSourceGoogleServiceAccounts_basic(t *testing.T) {
	t.Parallel()

	// Common resource configuration
	static_prefix := "tf-test"
	random_suffix := acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()

	// Configuration of network resources
	sa_1 := static_prefix + "-sa-1-" + random_suffix
	sa_2 := static_prefix + "-sa-2-" + random_suffix

	// Configuration map used in test deployment
	context := map[string]interface{}{
		"project": project,
		"sa_1":    sa_1,
		"sa_2":    sa_2,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccountsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// We can't guarantee that no more service accounts are in the project, so we'll check set-ness rather than correctness
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.0.account_id"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.0.disabled"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.0.email"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.0.member"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.0.name"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.0.unique_id"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.1.account_id"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.1.disabled"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.1.email"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.1.member"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.1.name"),
					resource.TestCheckResourceAttrSet("data.google_service_accounts.all", "accounts.1.unique_id"),

					// Check for prefix on account id
					resource.TestCheckResourceAttr("data.google_service_accounts.with_prefix", "accounts.0.account_id", sa_1),

					// Check for regex on email
					resource.TestCheckResourceAttr("data.google_service_accounts.with_regex", "accounts.0.email", fmt.Sprintf("%s@%s.iam.gserviceaccount.com", sa_1, project)),

					// Check if the account_id matches the prefix
					resource.TestCheckResourceAttr("data.google_service_accounts.with_prefix_and_regex", "accounts.0.account_id", fmt.Sprintf(sa_1)),

					// Check if the email matches the regex
					resource.TestCheckResourceAttr("data.google_service_accounts.with_prefix_and_regex", "accounts.0.email", fmt.Sprintf("%s@%s.iam.gserviceaccount.com", sa_1, project)),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountsConfig(context map[string]interface{}) string {
	return fmt.Sprintf(`
locals {
  project_id = "%s"
  sa_one     = "%s"
  sa_two     = "%s"
}

resource "google_service_account" "sa_one" {
  account_id   = local.sa_one
  description  = local.sa_one
  display_name = local.sa_one
}

resource "google_service_account" "sa_two" {
  account_id   = local.sa_two
  description  = local.sa_two
  display_name = local.sa_two
}

data "google_service_accounts" "all" {
  project = local.project_id

  depends_on = [
    google_service_account.sa_one,
    google_service_account.sa_two,
  ]
}

data "google_service_accounts" "with_prefix" {
  prefix  = google_service_account.sa_one.account_id
  project = local.project_id
}

data "google_service_accounts" "with_regex" {
  project = local.project_id
  regex   = ".*${google_service_account.sa_one.account_id}.*@.*\\.gserviceaccount\\.com"
}

data "google_service_accounts" "with_prefix_and_regex" {
  prefix  = google_service_account.sa_one.account_id
  project = local.project_id
  regex   = ".*${google_service_account.sa_one.account_id}.*@.*\\.gserviceaccount\\.com"
}
`,
		context["project"].(string),
		context["sa_1"].(string),
		context["sa_2"].(string),
	)
}
