package google

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAccessApprovalOrganizationSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessApprovalOrganizationSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessApprovalOrganizationSettings_full(context),
			},
			{
				ResourceName:            "google_organization_access_approval_settings.organization_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
			{
				Config: testAccAccessApprovalOrganizationSettings_update(context),
			},
			{
				ResourceName:            "google_organization_access_approval_settings.organization_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessApprovalOrganizationSettings_full(context map[string]interface{}) string {
	return Nprintf(`
resource "google_organization_access_approval_settings" "organization_access_approval" {
  organization_id     = "%{org_id}"
  notification_emails = ["testuser@example.com"]

  enrolled_services {
    cloud_product = "appengine.googleapis.com"
  }

  enrolled_services {
    cloud_product = "dataflow.googleapis.com"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccAccessApprovalOrganizationSettings_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_organization_access_approval_settings" "organization_access_approval" {
  organization_id     = "%{org_id}"
  notification_emails = ["testuser@example.com", "example.user@example.com"]

  enrolled_services {
    cloud_product = "all"
    enrollment_level = "BLOCK_ALL"
  }
}
`, context)
}

func testAccCheckAccessApprovalOrganizationSettingsDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_organization_access_approval_settings" {
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
