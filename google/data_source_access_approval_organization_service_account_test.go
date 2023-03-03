package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAccessApprovalOrganizationServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": GetTestOrgFromEnv(t),
	}

	resourceName := "data.google_access_approval_organization_service_account.aa_account"

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessApprovalOrganizationServiceAccount_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_email"),
				),
			},
		},
	})
}

func testAccDataSourceAccessApprovalOrganizationServiceAccount_basic(context map[string]interface{}) string {
	return Nprintf(`
data "google_access_approval_organization_service_account" "aa_account" {
  organization_id = "%{org_id}"
}
`, context)
}
