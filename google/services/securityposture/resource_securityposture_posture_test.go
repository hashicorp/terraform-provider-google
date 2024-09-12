// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securityposture_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityposturePosture_securityposturePosture_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgTargetFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityposturePosture_securityposturePosture_full(context),
			},
			{
				ResourceName:            "google_securityposture_posture.posture_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "posture_id", "annotations"},
			},
			{
				Config: testAccSecurityposturePosture_securityposturePosture_update(context),
			},
			{
				ResourceName:            "google_securityposture_posture.posture_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "posture_id", "annotations"},
			},
		},
	})
}

func testAccSecurityposturePosture_securityposturePosture_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture_test" {
	posture_id          = "posture_test"
	parent = "organizations/%{org_id}"
	location = "global"
	state = "ACTIVE"
    description = "a new posture"
    policy_sets {
    	policy_set_id = "org_policy_set"
    	description = "set of org policies"
    	policies {
    		policy_id = "policy_1"
    		constraint {
    			org_policy_constraint {
    				canned_constraint_id = "storage.uniformBucketLevelAccess"
    				policy_rules {
    					enforce = true
    				}
    			}
    		}
    	}
	}
}
`, context)
}

func testAccSecurityposturePosture_securityposturePosture_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture_test" {
	posture_id          = "posture_test"
	parent = "organizations/%{org_id}"
	location = "global"
	state = "ACTIVE"
    description = "a new posture"
    policy_sets {
    	policy_set_id = "org_policy_set"
    	description = "set of org policies"
    	policies {
    		policy_id = "policy_1"
    		constraint {
    			org_policy_constraint_custom {
    				custom_constraint {
    					name         = "organizations/%{org_id}/customConstraints/custom.disableGkeAutoUpgrade"
					  	display_name = "Disable GKE auto upgrade"
					  	description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

					  	action_type    = "ALLOW"
					  	condition      = "resource.management.autoUpgrade == false"
					  	method_types   = ["CREATE", "UPDATE"]
					  	resource_types = ["container.googleapis.com/NodePool"]
    				}
    				policy_rules {
    					enforce = true
    				}
    			}
    		}
		}
	}
}
`, context)
}
