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

func TestAccComputeRegionNetworkFirewallPolicyAssociation_RegionalHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"region":        getTestRegionFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionNetworkFirewallPolicyAssociationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionNetworkFirewallPolicyAssociation_RegionalHandWritten(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_association.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionNetworkFirewallPolicyAssociation_RegionalHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_compute_region_network_firewall_policy_association.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionNetworkFirewallPolicyAssociation_RegionalHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_region_network_firewall_policy" "basic_regional_network_firewall_policy" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Sample global network firewall policy"
  region = "%{region}"
}

resource "google_compute_network" "basic_network" {
  name = "tf-test-network%{random_suffix}"
}

resource "google_compute_region_network_firewall_policy_association" "primary" {
  name = "tf-test-association%{random_suffix}"
  attachment_target = google_compute_network.basic_network.id
  firewall_policy =  google_compute_region_network_firewall_policy.basic_regional_network_firewall_policy.name
  project =  "%{project_name}"
  region = "%{region}"
}

`, context)
}

func testAccComputeRegionNetworkFirewallPolicyAssociation_RegionalHandWrittenUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_region_network_firewall_policy" "basic_regional_network_firewall_policy" {
  name = "tf-test-policy%{random_suffix}"
  project = "%{project_name}"
  description = "Sample global network firewall policy"
  region = "%{region}"
}

resource "google_compute_network" "basic_network" {
  name = "tf-test-network%{random_suffix}"
}

resource "google_compute_network" "basic_network2" {
  name = "update-tf-test-network%{random_suffix}"
}

resource "google_compute_region_network_firewall_policy_association" "primary" {
  name = "tf-test-association%{random_suffix}"
  attachment_target = google_compute_network.basic_network2.id
  firewall_policy =  google_compute_region_network_firewall_policy.basic_regional_network_firewall_policy.name
  project =  "%{project_name}"
  region = "%{region}"
}

`, context)
}

func testAccCheckComputeRegionNetworkFirewallPolicyAssociationDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_compute_region_network_firewall_policy_association" {
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

			obj := &compute.NetworkFirewallPolicyAssociation{
				AttachmentTarget: dcl.String(rs.Primary.Attributes["attachment_target"]),
				FirewallPolicy:   dcl.String(rs.Primary.Attributes["firewall_policy"]),
				Name:             dcl.String(rs.Primary.Attributes["name"]),
				Project:          dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Location:         dcl.StringOrNil(rs.Primary.Attributes["region"]),
				ShortName:        dcl.StringOrNil(rs.Primary.Attributes["short_name"]),
			}

			client := NewDCLComputeClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetNetworkFirewallPolicyAssociation(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_compute_region_network_firewall_policy_association still exists %v", obj)
			}
		}
		return nil
	}
}
