// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iam3_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIAM3OrganizationsPolicyBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_full(context),
			},
			{
				ResourceName:            "google_iam_organizations_policy_binding.my_org_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "organization", "policy_binding_id"},
			},

			{
				Config: testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update(context),
			},
			{
				ResourceName:            "google_iam_organizations_policy_binding.my_org_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "location", "organization", "policy_binding_id"},
			},
		},
	})
}

func testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_principal_access_boundary_policy" "pab_policy" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  principal_access_boundary_policy_id = "tf-test-my-pab-policy%{random_suffix}"
}

resource "google_iam_organizations_policy_binding" "my_org_binding" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  policy_kind    = "PRINCIPAL_ACCESS_BOUNDARY"
  policy_binding_id = "tf-test-test-org-binding%{random_suffix}"
  policy         = "organizations/%{org_id}/locations/global/principalAccessBoundaryPolicies/${google_iam_principal_access_boundary_policy.pab_policy.principal_access_boundary_policy_id}"
  target {
    principal_set = "//cloudresourcemanager.googleapis.com/organizations/%{org_id}"
  }
}
`, context)
}

func testAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_iam_principal_access_boundary_policy" "pab_policy" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  principal_access_boundary_policy_id = "tf-test-my-pab-policy%{random_suffix}"
}

resource "google_iam_organizations_policy_binding" "my_org_binding" {
  organization   = "%{org_id}"
  location       = "global"
  display_name   = "test org binding%{random_suffix}"
  policy_kind    = "PRINCIPAL_ACCESS_BOUNDARY"
  policy_binding_id = "tf-test-test-org-binding%{random_suffix}"
  policy         = "organizations/%{org_id}/locations/global/principalAccessBoundaryPolicies/${google_iam_principal_access_boundary_policy.pab_policy.principal_access_boundary_policy_id}"
  annotations    = {"foo": "bar"}
  target {
    principal_set = "//cloudresourcemanager.googleapis.com/organizations/%{org_id}"
  }
  condition {
    description  = "test condition"
    expression   = "principal.subject == 'al@a.com'"
    location     = "test location"
    title        = "test title"
  }
}
`, context)
}
