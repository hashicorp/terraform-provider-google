// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeFirewallPolicy_update(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	policyName := fmt.Sprintf("tf-test-firewall-policy-%s", acctest.RandString(t, 10))
	folderName := fmt.Sprintf("tf-test-folder-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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
