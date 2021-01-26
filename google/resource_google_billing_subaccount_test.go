package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBillingSubaccount_renameOnDestroy(t *testing.T) {
	t.Parallel()

	masterBilling := getTestMasterBillingAccountFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleBillingSubaccountRenameOnDestroy,
		Steps: []resource.TestStep{
			{
				// Test Billing Subaccount creation
				Config: testAccBillingSubccount_renameOnDestroy(masterBilling),
				Check:  testAccCheckGoogleBillingSubaccountExists("subaccount_with_rename_on_destroy"),
			},
		},
	})
}

func TestAccBillingSubaccount_basic(t *testing.T) {
	t.Parallel()

	masterBilling := getTestMasterBillingAccountFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Billing Subaccount creation
				Config: testAccBillingSubccount_basic(masterBilling),
				Check:  testAccCheckGoogleBillingSubaccountExists("subaccount"),
			},
			{
				ResourceName:            "google_billing_subaccount.subaccount",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_policy"},
			},
			{
				// Test Billing Subaccount update
				Config: testAccBillingSubccount_update(masterBilling),
				Check:  testAccCheckGoogleBillingSubaccountExists("subaccount"),
			},
			{
				ResourceName:            "google_billing_subaccount.subaccount",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_policy"},
			},
		},
	})
}

func testAccBillingSubccount_basic(masterBillingAccountId string) string {
	return fmt.Sprintf(`
resource "google_billing_subaccount" "subaccount" {
  display_name = "Test Billing Subaccount"
  master_billing_account  = "%s"
}
`, masterBillingAccountId)
}

func testAccBillingSubccount_update(masterBillingAccountId string) string {
	return fmt.Sprintf(`
resource "google_billing_subaccount" "subaccount" {
  display_name = "Rename Test Billing Subaccount"
  master_billing_account  = "%s"
}
`, masterBillingAccountId)
}

func testAccBillingSubccount_renameOnDestroy(masterBillingAccountId string) string {
	return fmt.Sprintf(`
resource "google_billing_subaccount" "subaccount_with_rename_on_destroy" {
  display_name = "Test Billing Subaccount (Rename on Destroy)"
  master_billing_account  = "%s"
  deletion_policy = "RENAME_ON_DESTROY"
}
`, masterBillingAccountId)
}

func testAccCheckGoogleBillingSubaccountExists(bindingResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		subaccount, ok := s.RootModule().Resources["google_billing_subaccount."+bindingResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := testAccProvider.Meta().(*Config)
		_, err := config.NewBillingClient(config.userAgent).BillingAccounts.Get(subaccount.Primary.ID).Do()
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckGoogleBillingSubaccountRenameOnDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_billing_subaccount" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		res, err := config.NewBillingClient(config.userAgent).BillingAccounts.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if !strings.HasPrefix(res.DisplayName, "Terraform Destroyed") {
			return fmt.Errorf("Billing account %s was not renamed on destroy", rs.Primary.ID)
		}
	}

	return nil
}
