package google

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Bindings and members are tested serially to avoid concurrent updates of the org's IAM policy.
// When concurrent changes happen, the behavior is to abort and ask the user to retry allowing
// them to see the new diff instead of blindly overriding the policy stored in GCP. This desired
// behavior however induces flakiness in our acceptance tests, hence the need for running them
// serially.
// Policies are *not tested*, because testing them will ruin changes made to the test org.
func TestAccOrganizationIam(t *testing.T) {
	if os.Getenv("TF_RUN_ORG_IAM") != "true" {
		t.Skip("Environment variable TF_RUN_ORG_IAM is not set, skipping.")
	}

	t.Parallel()

	org := getTestOrgFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "tfIamTest" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Binding creation
				Config: testAccOrganizationIamBinding_basic(account, roleId, org),
				Check: testAccCheckGoogleOrganizationIamBindingExists("foo", "test-role", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
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
				Check: testAccCheckGoogleOrganizationIamBindingExists("foo", "test-role", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
					fmt.Sprintf("serviceAccount:%s-2@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
			{
				ResourceName:      "google_organization_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s organizations/%s/roles/%s", org, org, roleId),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccOrganizationIamMember_basic(account, org),
				Check: testAccCheckGoogleOrganizationIamMemberExists("foo", "roles/browser",
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				),
			},
			{
				ResourceName:      "google_organization_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s roles/browser serviceAccount:%s@%s.iam.gserviceaccount.com", org, account, getTestProjectFromEnv()),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleOrganizationIamBindingExists(bindingResourceName, roleResourceName string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		bindingRs, ok := s.RootModule().Resources["google_organization_iam_binding."+bindingResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", bindingResourceName)
		}

		roleRs, ok := s.RootModule().Resources["google_organization_iam_custom_role."+roleResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", roleResourceName)
		}

		config := testAccProvider.Meta().(*Config)
		p, err := config.clientResourceManager.Organizations.GetIamPolicy("organizations/"+bindingRs.Primary.Attributes["org_id"], &cloudresourcemanager.GetIamPolicyRequest{}).Do()
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

func testAccCheckGoogleOrganizationIamMemberExists(n, role, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["google_organization_iam_member."+n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		config := testAccProvider.Meta().(*Config)
		p, err := config.clientResourceManager.Organizations.GetIamPolicy("organizations/"+rs.Primary.Attributes["org_id"], &cloudresourcemanager.GetIamPolicyRequest{}).Do()
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
func testAccOrganizationIamBinding_basic(account, role, org string) string {
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
  role    = "${google_organization_iam_custom_role.test-role.id}"
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
  org_id  = "%s"
  role    = "${google_organization_iam_custom_role.test-role.id}"
  members = [
    "serviceAccount:${google_service_account.test-account.email}",
    "serviceAccount:${google_service_account.test-account-2.email}"
  ]
}
`, account, role, org, account, org)
}

func testAccOrganizationIamMember_basic(account, org string) string {
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
