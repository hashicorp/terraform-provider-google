// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeFirewallPolicyAssociation_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicyAssociation_basic(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_association.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
		},
	})
}

func testAccComputeFirewallPolicyAssociation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
  deletion_protection = false
}

resource "google_folder" "target_folder" {
  display_name = "tf-test-target-%{random_suffix}"
  parent       = "%{org_name}"
  deletion_protection = false
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_association" "default" {
  firewall_policy = google_compute_firewall_policy.default.id
  attachment_target = google_folder.target_folder.name
  name = "tf-test-association-%{random_suffix}"
}
`, context)
}

func TestAccComputeFirewallPolicyAssociation_organization(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", envvar.GetTestOrgFromEnv(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicyAssociation_organization(context),
			},
			{
				ResourceName:            "google_compute_firewall_policy_association.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
		},
	})
}

func testAccComputeFirewallPolicyAssociation_organization(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  display_name = "tf-test-my-folder-%{random_suffix}"
  parent       = "%{org_name}"
  deletion_protection = false
}

resource "google_compute_firewall_policy" "policy" {
  parent      = "%{org_name}"
  short_name  = "tf-test-my-policy-%{random_suffix}"
  description = "Example Resource"
}

resource "google_compute_firewall_policy_association" "default" {
  firewall_policy = google_compute_firewall_policy.policy.id
  attachment_target = google_folder.folder.name
  name = "tf-test-my-association-%{random_suffix}"
}
`, context)
}
