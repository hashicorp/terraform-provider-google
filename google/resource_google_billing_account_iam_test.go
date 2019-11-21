package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBillingAccountIam(t *testing.T) {
	t.Parallel()

	billing := getTestBillingAccountFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/billing.viewer"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccBillingAccountIamBinding_basic(account, billing, role),
				Check: testAccCheckGoogleBillingAccountIamBindingExists("foo", role, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_billing_account_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s roles/billing.viewer", billing),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccBillingAccountIamBinding_update(account, billing, role),
				Check: testAccCheckGoogleBillingAccountIamBindingExists("foo", role, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_billing_account_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s roles/billing.viewer", billing),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccBillingAccountIamMember_basic(account, billing, role),
				Check: testAccCheckGoogleBillingAccountIamMemberExists("foo", "roles/billing.viewer",
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				),
			},
			{
				ResourceName:      "google_billing_account_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s roles/billing.viewer serviceAccount:%s@%s.iam.gserviceaccount.com", billing, account, getTestProjectFromEnv()),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleBillingAccountIamBindingExists(bindingResourceName, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources["google_billing_account_iam_binding."+bindingResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := testAccProvider.Meta().(*Config)
		p, err := config.clientBilling.BillingAccounts.GetIamPolicy("billingAccounts/" + bindingRs.Primary.Attributes["billing_account_id"]).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

func testAccCheckGoogleBillingAccountIamMemberExists(n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_billing_account_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := testAccProvider.Meta().(*Config)
		p, err := config.clientBilling.BillingAccounts.GetIamPolicy("billingAccounts/" + rs.Primary.Attributes["billing_account_id"]).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				for _, m := range binding.Members {
					if m == member {
						return nil
					}
				}

				return fmt.Errorf("Missing member %q, got %v", member, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

func testAccBillingAccountIamBinding_basic(account, billingAccountId, role string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Billing Account Iam Testing Account"
}

resource "google_billing_account_iam_binding" "foo" {
  billing_account_id = "%s"
  role               = "%s"
  members            = ["serviceAccount:${google_service_account.test-account.email}"]
}
`, account, billingAccountId, role)
}

func testAccBillingAccountIamBinding_update(account, billingAccountId, role string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Billing Account Iam Testing Account"
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Billing Account Iam Testing Account"
}

resource "google_billing_account_iam_binding" "foo" {
  billing_account_id = "%s"
  role               = "%s"
  members = [
    "serviceAccount:${google_service_account.test-account.email}",
    "serviceAccount:${google_service_account.test-account-2.email}",
  ]
}
`, account, account, billingAccountId, role)
}

func testAccBillingAccountIamMember_basic(account, billingAccountId, role string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Billing Account Iam Testing Account"
}

resource "google_billing_account_iam_member" "foo" {
  billing_account_id = "%s"
  role               = "%s"
  member             = "serviceAccount:${google_service_account.test-account.email}"
}
`, account, billingAccountId, role)
}
