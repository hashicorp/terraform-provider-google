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

func testAccDataSourceCloudIdentityGroups_basicTest(t *testing.T) {

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
				Config: testAccCloudIdentityGroupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_cloud_identity_groups.groups",
						"groups.#"),
					resource.TestMatchResourceAttr("data.google_cloud_identity_groups.groups",
						"groups.0.name", regexp.MustCompile("^groups/.*$")),
				),
			},
		},
	})
}

func testAccCloudIdentityGroupConfig(context map[string]interface{}) string {
	return testAccCloudIdentityGroup_cloudIdentityGroupsBasicExample(context) + acctest.Nprintf(`

data "google_cloud_identity_groups" "groups" {
  parent = google_cloud_identity_group.cloud_identity_group_basic.parent
}
`, context)
}
