// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package orgpolicy_test

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrgPolicyCustomConstraint_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckOrgPolicyCustomConstraintDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrgPolicyCustomConstraint_v1(context),
			},
			{
				ResourceName:            "google_org_policy_custom_constraint.constraint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccOrgPolicyCustomConstraint_v2(context),
			},
			{
				ResourceName:            "google_org_policy_custom_constraint.constraint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccOrgPolicyCustomConstraint_v1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "custom.tfTest%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Disable GKE auto upgrade"
  description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE", "UPDATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}
`, context)
}

func testAccOrgPolicyCustomConstraint_v2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "custom.tfTest%{random_suffix}"
  parent       = "organizations/%{org_id}"
  display_name = "Updated"
  description  = "Updated"

  action_type    = "DENY"
  condition      = "resource.management.autoUpgrade == true"
  method_types   = ["CREATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}
`, context)
}
