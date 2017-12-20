package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleBillingAccount_basic(t *testing.T) {
	billingId := getTestBillingAccountFromEnv(t)
	name := "billingAccounts/" + billingId

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleBillingAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "id", billingId),
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "name", name),
					resource.TestCheckResourceAttr("data.google_billing_account.acct", "open", "true"),
				),
			},
		},
	})
}

func testAccCheckGoogleBillingAccount_basic(billingId string) string {
	return fmt.Sprintf(`
data "google_billing_account" "acct" {
  name = "%s"
}`, billingId)
}
