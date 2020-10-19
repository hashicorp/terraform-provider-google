package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGoogleProjectServiceAccountsBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_service_accounts.acceptance"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectServiceAccountBasic(getTestProjectFromEnv()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "id", getTestProjectFromEnv()),
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "action"),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectServiceAccountBasic(project string) string {
	return fmt.Sprintf(`
resource "google_project_default_service_accounts" "acceptance" {
	project = "%s"
}
`, project)
}
