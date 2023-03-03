package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAccessApprovalProjectServiceAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id": GetTestProjectFromEnv(),
	}

	resourceName := "data.google_access_approval_project_service_account.aa_account"

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessApprovalProjectServiceAccount_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "account_email"),
				),
			},
		},
	})
}

func testAccDataSourceAccessApprovalProjectServiceAccount_basic(context map[string]interface{}) string {
	return Nprintf(`
data "google_access_approval_project_service_account" "aa_account" {
  project_id = "%{project_id}"
}
`, context)
}
