// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testAccAccessContextManagerAccessPolicyIamBinding(t *testing.T) {
	acctest.SkipIfVcr(t)

	org := envvar.GetTestOrgFromEnv(t)
	account := "tf-test-" + acctest.RandString(t, 10)
	role := "roles/accesscontextmanager.policyAdmin"
	policy := createScopedPolicy(t, org)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

func testAccAccessContextManagerAccessPolicyIamMember(t *testing.T) {
	acctest.SkipIfVcr(t)

	org := envvar.GetTestOrgFromEnv(t)
	account := "tf-test-" + acctest.RandString(t, 10)
	role := "roles/accesscontextmanager.policyAdmin"
	policy := createScopedPolicy(t, org)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccAccessContextManagerAccessPolicyIamMember_basic(policy, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_access_context_manager_access_policy_iam_member.member", "role", role),
					resource.TestCheckResourceAttr(
						"google_access_context_manager_access_policy_iam_member.member", "member", "serviceAccount:"+envvar.ServiceAccountCanonicalEmail(account)),
				),
			},
		},
	})
}

func testAccAccessContextManagerAccessPolicyIamPolicy(t *testing.T) {
	acctest.SkipIfVcr(t)

	org := envvar.GetTestOrgFromEnv(t)
	account := "tf-test-" + acctest.RandString(t, 10)
	role := "roles/accesscontextmanager.policyAdmin"
	policy := createScopedPolicy(t, org)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccAccessContextManagerAccessPolicyIamPolicy_basic(policy, account, role),
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

func testAccAccessContextManagerAccessPolicyIamMember_basic(policy, account, role string) string {
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

func testAccAccessContextManagerAccessPolicyIamPolicy_basic(policy, account, role string) string {
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
	rand := acctest.RandString(t, 10)
	return fmt.Sprintf(`
		resource "google_project" "project" {
		project_id      = "acm-tf-test-%s"
		name            = "acm-tf-test-%s"
		org_id          = "%s"
		deletion_policy = "DELETE"
		}

		resource "google_access_context_manager_access_policy" "access-policy" {
			parent = "organizations/%s"
			title  = "test policy"
			scopes = ["projects/${google_project.project.number}"]
		}
	`, rand, rand, org, org)
}
