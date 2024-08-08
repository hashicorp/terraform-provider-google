// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testAccDataSourceCloudIdentityGroupTransitiveMemberships_basicTest(t *testing.T) {

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"identity_user": envvar.GetTestIdentityUserFromEnv(t),
		"random_suffix": randString,
		"group_b_id":    fmt.Sprintf("tf-test-group-b-%s@%s", randString, envvar.GetTestOrgDomainFromEnv(t)),
	}

	memberId := acctest.Nprintf("%{identity_user}@%{org_domain}", context)
	groupBId := context["group_b_id"].(string)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupTransitiveMembershipConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// Finds two members of Group A (1 direct, 1 indirect)
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_transitive_memberships.members",
						"memberships.#", "2"),
					// Group B is a member of Group A; DIRECT membership to A
					checkGroupTransitiveMembershipRelationship("data.google_cloud_identity_group_transitive_memberships.members", groupBId, "DIRECT"),
					// User is a member of Group B; INDIRECT membership to A
					checkGroupTransitiveMembershipRelationship("data.google_cloud_identity_group_transitive_memberships.members", memberId, "INDIRECT"),
				),
			},
		},
	})
}

// Create Group A, Group B
// Make Group B a member of Group A
// Make identity user a member of Group B; is a transitive member of Group A
func testAccCloudIdentityGroupTransitiveMembershipConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_cloud_identity_group" "group_a" {
  display_name = "tf-test-group-a-%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-group-a-%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group" "group_b" {
  display_name = "tf-test-group-b-%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "%{group_b_id}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}


resource "google_cloud_identity_group_membership" "group_b_membership_in_group_a" {
  group    = google_cloud_identity_group.group_a.id

  preferred_member_key {
    id = "%{group_b_id}"
  }

  roles {
    name = "MEMBER"
  }
}

// By putting the user in group B, they are also a member of group A via B
resource "google_cloud_identity_group_membership" "user_in_group_b" {
  group    = google_cloud_identity_group.group_b.id

  preferred_member_key {
    id = "%{identity_user}@%{org_domain}"
  }

  roles {
    name = "MEMBER"
  }

  roles {
    name = "MANAGER"
  }
}

# Wait after adding user to group B to handle eventual consistency errors.
resource "time_sleep" "wait_15_seconds" {
  depends_on = [google_cloud_identity_group_membership.user_in_group_b]

  create_duration = "15s"
}

// Look for all members of Group A. This should return Group B and the user.
data "google_cloud_identity_group_transitive_memberships" "members" {
  group = google_cloud_identity_group.group_a.id

  depends_on = [
    google_cloud_identity_group_membership.user_in_group_b,
	time_sleep.wait_15_seconds
  ]
}
`, context)
}

func checkGroupTransitiveMembershipRelationship(datasourceName, memberId, expectedRelationType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		}

		if ds.Primary.Attributes["memberships.#"] == "0" {
			return fmt.Errorf("no memberships found in %s", datasourceName)
		}

		membersCount, err := strconv.Atoi(ds.Primary.Attributes["memberships.#"])
		if err != nil {
			return fmt.Errorf("error getting number of members, %v", err)
		}
		found := false
		for i := 0; i < membersCount; i++ {
			id := ds.Primary.Attributes[fmt.Sprintf("memberships.%d.preferred_member_key.0.id", i)]
			relType := ds.Primary.Attributes[fmt.Sprintf("memberships.%d.relation_type", i)]
			found = (id == memberId) && (relType == expectedRelationType)

			if found {
				break
			}
		}

		if !found {
			return fmt.Errorf("did not find a user with id %s and relation type %s in the memberships list", memberId, expectedRelationType)
		}

		return nil
	}
}
