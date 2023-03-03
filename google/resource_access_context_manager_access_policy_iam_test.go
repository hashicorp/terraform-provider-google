package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAccessContextManagerAccessPolicyIamBinding(t *testing.T) {
	SkipIfVcr(t)

	org := GetTestOrgFromEnv(t)
	account := "tf-acm-iam-" + RandString(t, 10)
	role := "roles/accesscontextmanager.policyAdmin"
	policy := createScopedPolicy(t, org)
	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccAccessContextManagerAccessPolicyIamBinding_basic(policy, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_access_context_manager_access_policy_iam_binding.binding", "role", role),
				),
			},
		},
	})
}

func TestAccAccessContextManagerAccessPolicyIamMember(t *testing.T) {
	SkipIfVcr(t)

	org := GetTestOrgFromEnv(t)
	account := "tf-acm-iam-" + RandString(t, 10)
	role := "roles/accesscontextmanager.policyAdmin"
	policy := createScopedPolicy(t, org)
	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccAccessContextManagerAccessPolicyIamMember(policy, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_access_context_manager_access_policy_iam_member.member", "role", role),
					resource.TestCheckResourceAttr(
						"google_access_context_manager_access_policy_iam_member.member", "member", "serviceAccount:"+serviceAccountCanonicalEmail(account)),
				),
			},
		},
	})
}

func TestAccAccessContextManagerAccessPolicyIamPolicy(t *testing.T) {
	SkipIfVcr(t)

	org := GetTestOrgFromEnv(t)
	account := "tf-acm-iam-" + RandString(t, 10)
	role := "roles/accesscontextmanager.policyAdmin"
	policy := createScopedPolicy(t, org)
	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccAccessContextManagerAccessPolicyIamPolicy(policy, account, role),
			},
		},
	})
}

func testAccAccessContextManagerAccessPolicyIamBinding_basic(policy, account, role string) string {
	return fmt.Sprintf(policy+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Access Context Manager IAM Testing Account"
}

resource google_access_context_manager_access_policy_iam_binding binding {
	name = google_access_context_manager_access_policy.access-policy.name
	role = "%s" 
	members = [
		"serviceAccount:${google_service_account.test-account1.email}",
	]
}
`, account, role)
}

func testAccAccessContextManagerAccessPolicyIamMember(policy, account, role string) string {
	return fmt.Sprintf(policy+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Access Context Manager IAM Testing Account"
}

resource google_access_context_manager_access_policy_iam_member member {
	name = google_access_context_manager_access_policy.access-policy.name
	role = "%s" 
    member = "serviceAccount:${google_service_account.test-account.email}"
}

`, account, role)
}

func testAccAccessContextManagerAccessPolicyIamPolicy(policy, account, role string) string {
	return fmt.Sprintf(policy+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Access Context Manager IAM Testing Account"
}

data google_iam_policy admin {
	binding {
		role = "%s"
    	members = ["serviceAccount:${google_service_account.test-account.email}"]
	}
}
   
resource google_access_context_manager_access_policy_iam_policy policy {
	name = google_access_context_manager_access_policy.access-policy.name
	policy_data = data.google_iam_policy.admin.policy_data
}

`, account, role)
}

func createScopedPolicy(t *testing.T, org string) string {
	rand := RandString(t, 10)
	return fmt.Sprintf(`
		resource "google_project" "project" {
		project_id      = "acm-tf-test-%s"
		name            = "acm-tf-test-%s"
		org_id          = "%s"
		}

		resource "google_access_context_manager_access_policy" "access-policy" {
			parent = "organizations/%s"
			title  = "test policy"
			scopes = ["projects/${google_project.project.number}"]
		}
	`, rand, rand, org, org)
}
