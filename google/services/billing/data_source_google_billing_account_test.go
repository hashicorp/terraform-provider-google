// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package billing_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleBillingAccount_byFullName(t *testing.T) {
	billingId := envvar.GetTestMasterBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleBillingAccount_byName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "id", billingId),
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "name", name),
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "open", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBillingAccount_byShortName(t *testing.T) {
	billingId := envvar.GetTestMasterBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleBillingAccount_byName(billingId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "id", billingId),
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "name", name),
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "open", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBillingAccount_byFullNameClosed(t *testing.T) {
	billingId := envvar.GetTestMasterBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGoogleBillingAccount_byNameClosed(name),
				ExpectError: regexp.MustCompile("Billing account not found: " + name),
			},
		},
	})
}

func TestAccDataSourceGoogleBillingAccount_byDisplayName(t *testing.T) {
	name := acctest.RandString(t, 16)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGoogleBillingAccount_byDisplayName(name),
				ExpectError: regexp.MustCompile("Billing account not found: " + name),
			},
		},
	})
}

func testAccCheckGoogleBillingAccount_byName(name string) string {
	return fmt.Sprintf(`
data "google_billing_account" "acct" {
  billing_account = "%s"
  lookup_projects = false
}
`, name)
}

func testAccCheckGoogleBillingAccount_byNameClosed(name string) string {
	return fmt.Sprintf(`
data "google_billing_account" "acct" {
  billing_account = "%s"
  open            = false
}
`, name)
}

func testAccCheckGoogleBillingAccount_byDisplayName(name string) string {
	return fmt.Sprintf(`
data "google_billing_account" "acct" {
  display_name = "%s"
}
`, name)
}
