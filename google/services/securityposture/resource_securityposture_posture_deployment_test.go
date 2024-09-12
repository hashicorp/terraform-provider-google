// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package securityposture_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":         envvar.GetTestOrgTargetFromEnv(t),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_basic(context),
			},
			{
				ResourceName:            "google_securityposture_posture_deployment.postureDeployment_one",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
			{
				Config: testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context),
			},
			{
				ResourceName:            "google_securityposture_posture_deployment.postureDeployment_one",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "annotations"},
			},
		},
	})
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture_one" {
    posture_id          = "posture_one"
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

resource "google_securityposture_posture_deployment" "postureDeployment_one" {
    posture_deployment_id          = "posture_deployment_one"
    parent = "organizations/%{org_id}"
    location = "global"
    description = "a new posture deployment"
    target_resource = "projects/%{project_number}"
    posture_id = google_securityposture_posture.posture_one.name
    posture_revision_id = google_securityposture_posture.posture_one.revision_id
}
`, context)
}

func testAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_securityposture_posture" "posture_one" {
    posture_id          = "posture_one"
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
                    canned_constraint_id = "storage.publicAccessPrevention"
                    policy_rules {
                        enforce = true
                    }
                }
            }
        }
    }
}

resource "google_securityposture_posture_deployment" "postureDeployment_one" {
    posture_deployment_id          = "posture_deployment_one"
    parent = "organizations/%{org_id}"
    location = "global"
    description = "an updated posture deployment"
    target_resource = "projects/%{project_number}"
    posture_id = google_securityposture_posture.posture_one.name
    posture_revision_id = google_securityposture_posture.posture_one.revision_id
}
`, context)
}
