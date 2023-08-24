// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Bindings and members are tested serially to avoid concurrent updates of the org's IAM policy.
// When concurrent changes happen, the behavior is to abort and ask the user to retry allowing
// them to see the new diff instead of blindly overriding the policy stored in GCP. This desired
// behavior however induces flakiness in our acceptance tests, hence the need for running them
// serially.
// Policies are *not tested*, because testing them will ruin changes made to the test org.
func TestAccOrganizationIamMembersAndBindings(t *testing.T) {
	if os.Getenv("TF_RUN_ORG_IAM") != "true" {
		t.Skip("Environment variable TF_RUN_ORG_IAM is not set, skipping.")
	}
	t.Parallel()

	testCases := map[string]func(t *testing.T){
		"bindingBasic":     testAccOrganizationIamBinding_basic,
		"bindingCondition": testAccOrganizationIamBinding_condition,
		"memberBasic":      testAccOrganizationIamMember_basic,
		"memberCondition":  testAccOrganizationIamMember_condition,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccOrganizationIamBinding_basic(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "tfIamTest" + acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccOrganizationIamBinding_basicConfig(account, roleId, org),
				Check: testAccCheckGoogleOrganizationIamBindingExists(t, "foo", "test-role", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_organization_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s organizations/%s/roles/%s", org, org, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccOrganizationIamBinding_update(account, roleId, org),
				Check: testAccCheckGoogleOrganizationIamBindingExists(t, "foo", "test-role", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_organization_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s organizations/%s/roles/%s", org, org, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationIamBinding_condition(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	roleId := "tfIamTest" + acctest.RandString(t, 10)
	conditionTitle := "expires_after_2019_12_31"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccOrganizationIamBinding_conditionConfig(account, roleId, org, conditionTitle),
				Check: testAccCheckGoogleOrganizationIamBindingExists(t, "foo", "test-role", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_organization_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s organizations/%s/roles/%s %s", org, org, roleId, conditionTitle),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationIamMember_basic(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccOrganizationIamMember_basicConfig(account, org),
				Check: testAccCheckGoogleOrganizationIamMemberExists(t, "foo", "roles/browser",
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				),
			},
			{
				ResourceName:      "google_organization_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s roles/browser serviceAccount:%s@%s.iam.gserviceaccount.com", org, account, envvar.GetTestProjectFromEnv()),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationIamMember_condition(t *testing.T) {
	org := envvar.GetTestOrgFromEnv(t)
	account := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	conditionTitle := "expires_after_2019_12_31"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccOrganizationIamMember_conditionConfig(account, org, conditionTitle),
				Check: testAccCheckGoogleOrganizationIamMemberExists(t, "foo", "roles/browser",
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, envvar.GetTestProjectFromEnv()),
				),
			},
			{
				ResourceName: "google_organization_iam_member.foo",
				ImportStateId: fmt.Sprintf(
					"%s roles/browser serviceAccount:%s@%s.iam.gserviceaccount.com %s",
					org,
					account,
					envvar.GetTestProjectFromEnv(),
					conditionTitle,
				),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleOrganizationIamBindingExists(t *testing.T, bindingResourceName, roleResourceName string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources["google_organization_iam_binding."+bindingResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		roleRs, ok := s.RootModule().Resources["google_organization_iam_custom_role."+roleResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", roleResourceName)
		}

		config := acctest.GoogleProviderConfig(t)
		p, err := config.NewResourceManagerClient(config.UserAgent).Organizations.GetIamPolicy(
			"organizations/"+bindingRs.Primary.Attributes["org_id"],
			&cloudresourcemanager.GetIamPolicyRequest{
				Options: &cloudresourcemanager.GetPolicyOptions{
					RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
				},
			},
		).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == roleRs.Primary.ID {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", roleRs.Primary.ID)
	}
}

func testAccCheckGoogleOrganizationIamMemberExists(t *testing.T, n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_organization_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := acctest.GoogleProviderConfig(t)
		p, err := config.NewResourceManagerClient(config.UserAgent).Organizations.GetIamPolicy(
			"organizations/"+rs.Primary.Attributes["org_id"],
			&cloudresourcemanager.GetIamPolicyRequest{
				Options: &cloudresourcemanager.GetPolicyOptions{
					RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
				},
			},
		).Do()
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

// We are using a custom role since iam_binding is authoritative on the member list and
// we want to avoid removing members from an existing role to prevent unwanted side effects.
func testAccOrganizationIamBinding_basicConfig(account, role, org string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Organization Iam Testing Account"
}

resource "google_organization_iam_custom_role" "test-role" {
  role_id     = "%s"
  org_id      = "%s"
  title       = "Iam Testing Role"
  permissions = ["genomics.datasets.get"]
}

resource "google_organization_iam_binding" "foo" {
  org_id  = "%s"
  role    = google_organization_iam_custom_role.test-role.id
  members = ["serviceAccount:${google_service_account.test-account.email}"]
}
`, account, role, org, org)
}

func testAccOrganizationIamBinding_update(account, role, org string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Organization Iam Testing Account"
}

resource "google_organization_iam_custom_role" "test-role" {
  role_id     = "%s"
  org_id      = "%s"
  title       = "Iam Testing Role"
  permissions = ["genomics.datasets.get"]
}

resource "google_service_account" "test-account-2" {
  account_id   = "%s-2"
  display_name = "Organization Iam Testing Account"
}

resource "google_organization_iam_binding" "foo" {
  org_id = "%s"
  role   = google_organization_iam_custom_role.test-role.id
  members = [
    "serviceAccount:${google_service_account.test-account.email}",
    "serviceAccount:${google_service_account.test-account-2.email}",
  ]
}
`, account, role, org, account, org)
}

func testAccOrganizationIamBinding_conditionConfig(account, role, org, conditionTitle string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Organization Iam Testing Account"
}

resource "google_organization_iam_custom_role" "test-role" {
  role_id     = "%s"
  org_id      = "%s"
  title       = "Iam Testing Role"
  permissions = ["genomics.datasets.get"]
}

resource "google_organization_iam_binding" "foo" {
  org_id  = "%s"
  role    = google_organization_iam_custom_role.test-role.id
  members = ["serviceAccount:${google_service_account.test-account.email}"]
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, account, role, org, org, conditionTitle)
}

func testAccOrganizationIamMember_basicConfig(account, org string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Organization Iam Testing Account"
}

resource "google_organization_iam_member" "foo" {
  org_id = "%s"
  role   = "roles/browser"
  member = "serviceAccount:${google_service_account.test-account.email}"
}
`, account, org)
}

func testAccOrganizationIamMember_conditionConfig(account, org, conditionTitle string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Organization Iam Testing Account"
}

resource "google_organization_iam_member" "foo" {
  org_id = "%s"
  role   = "roles/browser"
  member = "serviceAccount:${google_service_account.test-account.email}"
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, account, org, conditionTitle)
}
