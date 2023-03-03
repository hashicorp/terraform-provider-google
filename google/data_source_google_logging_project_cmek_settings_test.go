package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLoggingProjectCmekSettings_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    "tf-test-" + RandString(t, 10),
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
	}
	resourceName := "data.google_logging_project_cmek_settings.cmek_settings"

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectCmekSettings_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "id", fmt.Sprintf("projects/%s/cmekSettings", context["project_name"])),
					resource.TestCheckResourceAttr(
						resourceName, "name", fmt.Sprintf("projects/%s/cmekSettings", context["project_name"])),
					resource.TestCheckResourceAttrSet(resourceName, "service_account_id"),
				),
			},
		},
	})
}

func testAccLoggingProjectCmekSettings_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "default" {
	project_id      = "%{project_name}"
	name            = "%{project_name}"
	org_id          = "%{org_id}"
	billing_account = "%{billing_account}"
}

data "google_logging_project_cmek_settings" "cmek_settings" {
	project = google_project.default.name
}
`, context)
}
