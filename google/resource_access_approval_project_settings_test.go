package google

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since access approval settings are hierarchical, and only one can exist per folder/project/org,
// and all refer to the same organization, they need to be run serially.
// See AccessApprovalOrganizationSettings for the test runner.
func testAccAccessApprovalProjectSettings(t *testing.T) {
	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessApprovalProjectSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessApprovalProjectSettings_basic(context),
			},
			{
				ResourceName:            "google_project_access_approval_settings.project_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
			{
				Config: testAccAccessApprovalProjectSettings_update(context),
			},
			{
				ResourceName:            "google_project_access_approval_settings.project_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccAccessApprovalProjectSettings_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project_access_approval_settings" "project_access_approval" {
  project_id          = "%{project}"

  enrolled_services {
    cloud_product = "all"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccAccessApprovalProjectSettings_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project_access_approval_settings" "project_access_approval" {
  project_id          = "%{project}"
  notification_emails = ["testuser@example.com", "example.user@example.com"]

  enrolled_services {
    cloud_product = "all"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccCheckAccessApprovalProjectSettingsDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_project_access_approval_settings" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			log.Printf("[DEBUG] Ignoring destroy during test")
		}

		return nil
	}
}
