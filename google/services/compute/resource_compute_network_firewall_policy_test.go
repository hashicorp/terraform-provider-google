// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeNetworkFirewallPolicy_GlobalHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkFirewallPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkFirewallPolicy_GlobalHandWritten(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeNetworkFirewallPolicy_GlobalHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeNetworkFirewallPolicy_LongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeNetworkFirewallPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkFirewallPolicy_LongForm(context),
			},
			{
				ResourceName:      "google_compute_network_firewall_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Import won't get the long form of any URL parameter
				ImportStateVerifyIgnore: []string{"project", "region"},
			},
		},
	})
}

func testAccComputeNetworkFirewallPolicy_GlobalHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Sample global network firewall policy"
}

`, context)
}

func testAccComputeNetworkFirewallPolicy_GlobalHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Updated global network firewall policy"
}

`, context)
}

func testAccComputeNetworkFirewallPolicy_LongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "projects/%{project_name}"
  description = "Sample global network firewall policy"
}

`, context)
}
