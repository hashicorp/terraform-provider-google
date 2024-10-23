// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeRegionNetworkFirewallPolicyAssociation_regionNetworkFirewallPolicyAssociationExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionNetworkFirewallPolicyAssociationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicyAssociation_regionNetworkFirewallPolicyAssociationExample(context),
			},
			{
				ResourceName:            "google_compute_region_network_firewall_policy_association.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_policy", "region"},
			},
		},
	})
}

func testAccComputeRegionNetworkFirewallPolicyAssociation_regionNetworkFirewallPolicyAssociationExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_network_firewall_policy" "policy" {
  name = "tf-test-my-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Sample global network firewall policy"
  region = "%{region}"
}

resource "google_compute_network" "network" {
  name = "tf-test-my-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_region_network_firewall_policy_association" "default" {
  name = "tf-test-my-association%{random_suffix}"
  project = "%{project_name}"
  attachment_target = google_compute_network.network.id
  firewall_policy =  google_compute_region_network_firewall_policy.policy.id
  region = "%{region}"
}
`, context)
}

func testAccCheckComputeRegionNetworkFirewallPolicyAssociationDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_region_network_firewall_policy_association" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/firewallPolicies/{{firewall_policy}}/getAssociation?name={{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeRegionNetworkFirewallPolicyAssociation still exists at %s", url)
			}
		}

		return nil
	}
}
