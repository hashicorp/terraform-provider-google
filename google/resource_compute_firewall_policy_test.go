package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeFirewallPolicy_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	policyName := fmt.Sprintf("tf-test-firewall-policy-%s", randString(t, 10))
	folderName := fmt.Sprintf("tf-test-folder-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicy_basic(org, policyName, folderName),
			},
			{
				ResourceName:      "google_compute_firewall_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewallPolicy_update(org, policyName, folderName),
			},
			{
				ResourceName:      "google_compute_firewall_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewallPolicy_update(org, policyName, folderName),
			},
			{
				ResourceName:      "google_compute_firewall_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeFirewallPolicy_basic(org, policyName, folderName string) string {
	return fmt.Sprintf(`
resource "google_folder" "folder" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.name
  short_name  = "%s"
  description = "Resource created for Terraform acceptance testing"
}
`, folderName, "organizations/"+org, policyName)
}

func testAccComputeFirewallPolicy_update(org, policyName, folderName string) string {
	return fmt.Sprintf(`
resource "google_folder" "folder" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.id
  short_name  = "%s"
  description = "An updated description"
}
`, folderName, "organizations/"+org, policyName)
}
