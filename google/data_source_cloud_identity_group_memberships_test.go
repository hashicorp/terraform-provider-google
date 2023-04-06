package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceCloudIdentityGroupMemberships_basicTest(t *testing.T) {

	context := map[string]interface{}{
		"org_domain":    GetTestOrgDomainFromEnv(t),
		"cust_id":       GetTestCustIdFromEnv(t),
		"identity_user": GetTestIdentityUserFromEnv(t),
		"random_suffix": RandString(t, 10),
	}

	memberId := Nprintf("%{identity_user}@%{org_domain}", context)

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembershipConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_memberships.members",
						"memberships.#", "1"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_memberships.members",
						"memberships.0.roles.#", "2"),
					resource.TestCheckResourceAttr("data.google_cloud_identity_group_memberships.members",
						"memberships.0.preferred_member_key.0.id", memberId),
				),
			},
		},
	})
}

func testAccCloudIdentityGroupMembershipConfig(context map[string]interface{}) string {
	return testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipUserExample(context) + Nprintf(`

data "google_cloud_identity_group_memberships" "members" {
  group = google_cloud_identity_group_membership.cloud_identity_group_membership_basic.group
}
`, context)
}
