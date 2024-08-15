// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeRegionNetworkFirewallPolicy_RegionalHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionNetworkFirewallPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicy_RegionalHandWritten(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicy_RegionalHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeRegionNetworkFirewallPolicy_LongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionNetworkFirewallPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicy_LongForm(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Import won't get the long form of any URL parameter
				ImportStateVerifyIgnore: []string{"project", "region"},
			},
		},
	})
}

func testAccComputeRegionNetworkFirewallPolicy_RegionalHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Sample regional network firewall policy"
  region = "%{region}"
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicy_RegionalHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Updated regional network firewall policy"
  region = "%{region}"
}
`, context)
}

func testAccComputeRegionNetworkFirewallPolicy_LongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "projects/%{project_name}"
  description = "Sample regional network firewall policy"
  region = "regions/%{region}"
}
`, context)
}
