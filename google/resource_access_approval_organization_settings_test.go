package google

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Since access approval settings are hierarchical, and only one can exist per folder/project/org,
// and all refer to the same organization, they need to be run serially
func TestAccAccessApprovalSettings(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"folder":       testAccAccessApprovalFolderSettings,
		"project":      testAccAccessApprovalProjectSettings,
		"organization": testAccAccessApprovalOrganizationSettings,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccAccessApprovalOrganizationSettings(t *testing.T) {
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
    cloud_product = "App Engine"
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
