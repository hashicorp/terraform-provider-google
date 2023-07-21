// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testAccDataSourceCloudIdentityGroupMemberships_basicTest(t *testing.T) {

	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"identity_user": envvar.GetTestIdentityUserFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	memberId := acctest.Nprintf("%{identity_user}@%{org_domain}", context)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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
	return testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipUserExample(context) + acctest.Nprintf(`

data "google_cloud_identity_group_memberships" "members" {
  group = google_cloud_identity_group_membership.cloud_identity_group_membership_basic.group
}
`, context)
}
