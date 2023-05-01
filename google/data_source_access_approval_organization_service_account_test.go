package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceAccessApprovalOrganizationServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": acctest.GetTestOrgFromEnv(t),
	}

	resourceName := "data.google_access_approval_organization_service_account.aa_account"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
