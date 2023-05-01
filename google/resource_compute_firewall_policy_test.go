package google

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeFirewallPolicy_update(t *testing.T) {
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	policyName := fmt.Sprintf("tf-test-firewall-policy-%s", RandString(t, 10))
	folderName := fmt.Sprintf("tf-test-folder-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
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
