package google

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGoogleOrganizationPolicy_import(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_list_allowAll(org),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
