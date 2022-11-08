// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccComputeRegionNetworkFirewallPolicy_RegionalHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"region":        getTestRegionFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionNetworkFirewallPolicyDestroyProducer(t),
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

func testAccComputeRegionNetworkFirewallPolicy_RegionalHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_region_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Sample regional network firewall policy"
  region = "%{region}"
}


`, context)
}

func testAccComputeRegionNetworkFirewallPolicy_RegionalHandWrittenUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_region_network_firewall_policy" "primary" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Updated regional network firewall policy"
  region = "%{region}"
}


`, context)
}

func testAccCheckComputeRegionNetworkFirewallPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_compute_region_network_firewall_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &compute.NetworkFirewallPolicy{
				Name:              dcl.String(rs.Primary.Attributes["name"]),
				Description:       dcl.String(rs.Primary.Attributes["description"]),
				Project:           dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Location:          dcl.StringOrNil(rs.Primary.Attributes["region"]),
				CreationTimestamp: dcl.StringOrNil(rs.Primary.Attributes["creation_timestamp"]),
				Fingerprint:       dcl.StringOrNil(rs.Primary.Attributes["fingerprint"]),
				Id:                dcl.StringOrNil(rs.Primary.Attributes["region_network_firewall_policy_id"]),
				SelfLink:          dcl.StringOrNil(rs.Primary.Attributes["self_link"]),
				SelfLinkWithId:    dcl.StringOrNil(rs.Primary.Attributes["self_link_with_id"]),
			}

			client := NewDCLComputeClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetNetworkFirewallPolicy(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_compute_region_network_firewall_policy still exists %v", obj)
			}
		}
		return nil
	}
}
