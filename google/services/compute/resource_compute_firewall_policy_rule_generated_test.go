// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

package compute_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeFirewallPolicyRule_BasicFirSecRuleHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_name":  envvar.GetTestProjectFromEnv(),
		"service_acct":  envvar.GetTestServiceAccountFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicyRule_BasicFirSecRuleHandWritten(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewallPolicyRule_BasicFirSecRuleHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeFirewallPolicyRule_BasicFirSecRuleHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "organizations/%{org_id}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_folder" "folder" {
  display_name = "tf-test-policy%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.id
  short_name  = "tf-test-policy%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "primary" {
  firewall_policy = google_compute_firewall_policy.default.name
  description     = "Resource created for Terraform acceptance testing"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [8080]
    }
    layer4_configs {
      ip_protocol = "udp"
      ports = [22]
    }
    dest_ip_ranges = ["11.100.0.1/32"]
    dest_fqdns = []
    dest_region_codes = ["US"]
    dest_threat_intelligences = ["iplist-known-malicious-ips"]
    src_address_groups = []
    dest_address_groups = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
  target_service_accounts = ["%{service_acct}"]
}

`, context)
}

func testAccComputeFirewallPolicyRule_BasicFirSecRuleHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "tf-test-policy%{random_suffix}"
  parent      = "organizations/%{org_id}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_folder" "folder" {
  display_name = "tf-test-policy%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.id
  short_name  = "tf-test-policy%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "primary" {
  firewall_policy = google_compute_firewall_policy.default.name
  description     = "Resource created for Terraform acceptance testing - Updated"
  priority        = 9000
  enable_logging  = true
  action          = "allow"
  direction       = "EGRESS"
  disabled        = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [8080]
    }
    layer4_configs {
      ip_protocol = "udp"
      ports = [22]
    }
    dest_ip_ranges = ["11.100.0.1/32", "10.0.0.0/24"]
    dest_fqdns = ["google.com"]
    dest_region_codes = ["US"]
    dest_threat_intelligences = ["iplist-known-malicious-ips"]
    src_address_groups = []
    dest_address_groups = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
  target_service_accounts = ["%{service_acct}"]
}

`, context)
}

func testAccCheckComputeFirewallPolicyRuleDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_compute_firewall_policy_rule" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &compute.FirewallPolicyRule{
				Action:         dcl.String(rs.Primary.Attributes["action"]),
				Direction:      compute.FirewallPolicyRuleDirectionEnumRef(rs.Primary.Attributes["direction"]),
				FirewallPolicy: dcl.String(rs.Primary.Attributes["firewall_policy"]),
				Description:    dcl.String(rs.Primary.Attributes["description"]),
				Disabled:       dcl.Bool(rs.Primary.Attributes["disabled"] == "true"),
				EnableLogging:  dcl.Bool(rs.Primary.Attributes["enable_logging"] == "true"),
				Kind:           dcl.StringOrNil(rs.Primary.Attributes["kind"]),
			}

			client := transport_tpg.NewDCLComputeClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetFirewallPolicyRule(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_compute_firewall_policy_rule still exists %v", obj)
			}
		}
		return nil
	}
}
