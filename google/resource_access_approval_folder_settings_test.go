package google

import (
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAccessApprovalFolderSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessApprovalFolderSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessApprovalFolderSettings_full(context),
			},
			{
				ResourceName:            "google_folder_access_approval_settings.folder_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder_id"},
			},
			{
				Config: testAccAccessApprovalFolderSettings_update(context),
			},
			{
				ResourceName:            "google_folder_access_approval_settings.folder_access_approval",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder_id"},
			},
		},
	})
}

func testAccAccessApprovalFolderSettings_full(context map[string]interface{}) string {
	return Nprintf(`
resource "google_folder" "my_folder" {
  display_name = "tf-test-my-folder%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

resource "google_folder_access_approval_settings" "folder_access_approval" {
  folder_id           = google_folder.my_folder.folder_id
  notification_emails = ["testuser@example.com"]

  enrolled_services {
    cloud_product = "all"
  }
}
`, context)
}

func testAccAccessApprovalFolderSettings_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_folder" "my_folder" {
  display_name = "tf-test-my-folder%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

resource "google_folder_access_approval_settings" "folder_access_approval" {
  folder_id           = google_folder.my_folder.folder_id
  notification_emails = ["testuser@example.com", "example.user@example.com"]

  enrolled_services {
    cloud_product = "all"
  }
}
`, context)
}

func testAccCheckAccessApprovalFolderSettingsDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_folder_access_approval_settings" {
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
