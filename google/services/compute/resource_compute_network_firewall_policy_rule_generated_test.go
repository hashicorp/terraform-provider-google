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

func TestAccComputeNetworkFirewallPolicyRule_networkFirewallPolicyRuleExample(t *testing.T) {
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
		CheckDestroy:             testAccCheckComputeNetworkFirewallPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNetworkFirewallPolicyRule_networkFirewallPolicyRuleExample(context),
			},
			{
				ResourceName:            "google_compute_network_firewall_policy_rule.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
		},
	})
}

func testAccComputeNetworkFirewallPolicyRule_networkFirewallPolicyRuleExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_security_address_group" "basic_global_networksecurity_address_group" {
  name        = "address%{random_suffix}"
  parent      = "projects/%{project_name}"
  description = "Sample global networksecurity_address_group"
  location    = "global"
  items       = ["208.80.154.224/32"]
  type        = "IPV4"
  capacity    = 100
}

resource "google_compute_network_firewall_policy" "basic_network_firewall_policy" {
  name        = "policy%{random_suffix}"
  description = "Sample global network firewall policy"
  project     = "%{project_name}"
}

resource "google_compute_network_firewall_policy_rule" "primary" {
  action                  = "allow"
  description             = "This is a simple rule description"
  direction               = "INGRESS"
  disabled                = false
  enable_logging          = true
  firewall_policy         = google_compute_network_firewall_policy.basic_network_firewall_policy.name
  priority                = 1000
  rule_name               = "test-rule"
  target_service_accounts = ["%{service_acct}"]

  match {
    src_ip_ranges = ["10.100.0.1/32"]
    src_fqdns = ["google.com"]
    src_region_codes = ["US"]
    src_threat_intelligences = ["iplist-known-malicious-ips"]

    src_secure_tags {
      name = google_tags_tag_value.basic_value.id
    }

    layer4_configs {
      ip_protocol = "all"
    }

    src_address_groups = [google_network_security_address_group.basic_global_networksecurity_address_group.id]
  }
}

resource "google_compute_network" "basic_network" {
  name = "network%{random_suffix}"
}

resource "google_tags_tag_key" "basic_key" {
  description = "For keyname resources."
  parent      = "organizations/%{org_id}"
  purpose     = "GCE_FIREWALL"
  short_name  = "tagkey%{random_suffix}"
  purpose_data = {
    network = "%{project_name}/${google_compute_network.basic_network.name}"
  }
}

resource "google_tags_tag_value" "basic_value" {
  description = "For valuename resources."
  parent      = google_tags_tag_key.basic_key.id
  short_name  = "tagvalue"
}
`, context)
}

func testAccCheckComputeNetworkFirewallPolicyRuleDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_network_firewall_policy_rule" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/getRule?priority={{priority}}")
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
				return fmt.Errorf("ComputeNetworkFirewallPolicyRule still exists at %s", url)
			}
		}

		return nil
	}
}
