package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudIdentityGroups_basic(t *testing.T) {

	context := map[string]interface{}{
		"org_domain":    getTestOrgDomainFromEnv(t),
		"cust_id":       getTestCustIdFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	return testAccCloudIdentityGroup_cloudIdentityGroupsBasicExample(context) + Nprintf(`

data "google_cloud_identity_groups" "groups" {
  parent = google_cloud_identity_group.cloud_identity_group_basic.parent
}
`, context)
}
