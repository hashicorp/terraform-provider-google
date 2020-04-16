package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleBillingAccount_byFullName(t *testing.T) {
	billingId := getTestBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	vcrTest(t, resource.TestCase{
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

func TestAccDataSourceGoogleBillingAccount_byShortName(t *testing.T) {
	billingId := getTestBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	billingId := getTestBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGoogleBillingAccount_byNameClosed(name),
				ExpectError: regexp.MustCompile("Billing account not found: " + name),
			},
		},
	})
}

func TestAccDataSourceGoogleBillingAccount_byDisplayName(t *testing.T) {
	name := randString(t, 16)

	vcrTest(t, resource.TestCase{
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
  billing_account = "%s"
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
