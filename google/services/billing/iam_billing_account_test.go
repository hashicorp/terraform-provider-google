// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package billing_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBillingAccountIam(t *testing.T) {
	// Deletes two fine-grained resources in same step
	acctest.SkipIfVcr(t)
	t.Parallel()

	billing := envvar.GetTestMasterBillingAccountFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/billing.viewer"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccBillingAccountIamBinding_basic(account, billing, role),
				Check: testAccCheckGoogleBillingAccountIamBindingExists(t, "foo", role, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
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
				Check: testAccCheckGoogleBillingAccountIamBindingExists(t, "foo", role, []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_billing_account_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s roles/billing.viewer", billing),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove the binding from state before adding a member.
				// Otherwise, we'll process the delete and create in an arbitrary order
				// and may have inconsistent results
				Config: testAccBillingAccountNoBindings(account),
			},
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccBillingAccountIamMember_basic(account, billing, role),
				Check: testAccCheckGoogleBillingAccountIamMemberExists(t, "foo", "roles/billing.viewer",
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				),
			},
			{
				ResourceName:      "google_billing_account_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s roles/billing.viewer serviceAccount:%s@%s.iam.gserviceaccount.com", billing, account, envvar.GetTestProjectFromEnv()),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleBillingAccountIamBindingExists(t *testing.T, bindingResourceName, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources["google_billing_account_iam_binding."+bindingResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		p, err := config.NewBillingClient(config.UserAgent).BillingAccounts.GetIamPolicy("billingAccounts/" + bindingRs.Primary.Attributes["billing_account_id"]).Do()
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

func testAccCheckGoogleBillingAccountIamMemberExists(t *testing.T, n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_billing_account_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		p, err := config.NewBillingClient(config.UserAgent).BillingAccounts.GetIamPolicy("billingAccounts/" + rs.Primary.Attributes["billing_account_id"]).Do()
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

func testAccBillingAccountNoBindings(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Billing Account Iam Testing Account"
}
`, account)
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
