package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleBillingAccount_byName(t *testing.T) {
	billingId := getTestBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccDataSourceGoogleBillingAccount_byDisplayName(t *testing.T) {
	name := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
  name = "%s"
}`, name)
}

func testAccCheckGoogleBillingAccount_byDisplayName(name string) string {
	return fmt.Sprintf(`
data "google_billing_account" "acct" {
  display_name = "%s"
}`, name)
}
