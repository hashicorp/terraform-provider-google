package google

import (
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccGoogleOrganizationPolicy_import(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_list_allowAll(org),
			},
			{
				ResourceName:      "google_organization_policy.listAll",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
