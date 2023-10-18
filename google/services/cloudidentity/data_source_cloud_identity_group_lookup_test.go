// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testAccDataSourceCloudIdentityGroupLookup_basicTest(t *testing.T) {

	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create group and look it up via its group key, i.e. email we set
				Config: testAccCloudIdentityGroupLookupConfig_groupKeyLookup(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_cloud_identity_group_lookup.email",
						"name", regexp.MustCompile("^groups/.*$")),
					resource.TestCheckResourceAttrPair("data.google_cloud_identity_group_lookup.email", "name",
						"google_cloud_identity_group.cloud_identity_group_basic", "name"),
				),
			},
			{
				// Look up group via an API-generated 'additional group key'
				Config: testAccCloudIdentityGroupLookupConfig_additionalGroupKeyLookup(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.google_cloud_identity_group_lookup.additional-groupkey", "name",
						"google_cloud_identity_group.cloud_identity_group_basic", "name"),
				),
			},
		},
	})
}

func testAccCloudIdentityGroupLookupConfig_groupKeyLookup(context map[string]interface{}) string {
	return acctest.Nprintf(`
# config matching testAccCloudIdentityGroup_cloudIdentityGroupsBasicExample
resource "google_cloud_identity_group" "cloud_identity_group_basic" {
  display_name         = "tf-test-my-identity-group%{random_suffix}"
  initial_group_config = "WITH_INITIAL_OWNER"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

data "google_cloud_identity_group_lookup" "email" {
  group_key {
    id = google_cloud_identity_group.cloud_identity_group_basic.group_key[0].id
  }
}
`, context)
}

func testAccCloudIdentityGroupLookupConfig_additionalGroupKeyLookup(context map[string]interface{}) string {
	return acctest.Nprintf(`
# config matching testAccCloudIdentityGroup_cloudIdentityGroupsBasicExample
resource "google_cloud_identity_group" "cloud_identity_group_basic" {
  display_name         = "tf-test-my-identity-group%{random_suffix}"
  initial_group_config = "WITH_INITIAL_OWNER"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

data "google_cloud_identity_group_lookup" "additional-groupkey" {
  group_key {
    # This value is an automatically created 'additionalGroupKeys' value
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}.test-google-a.com"
  }
  depends_on = [
    google_cloud_identity_group.cloud_identity_group_basic,
  ]
}
`, context)
}
