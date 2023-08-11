// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSecurityGatewaySecurityPolicyRule_update(t *testing.T) {
	t.Parallel()

	gatewaySecurityPolicyName := fmt.Sprintf("tf-test-gateway-sp-%s", acctest.RandString(t, 10))
	gatewaySecurityPolicyRuleName := fmt.Sprintf("tf-test-gateway-sp-rule-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityGatewaySecurityPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityGatewaySecurityPolicyRule_basic(gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName),
			},
			{
				ResourceName:      "google_network_security_gateway_security_policy_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkSecurityGatewaySecurityPolicyRule_update(gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName),
			},
			{
				ResourceName:      "google_network_security_gateway_security_policy_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkSecurityGatewaySecurityPolicyRule_basic(gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName),
			},
			{
				ResourceName:      "google_network_security_gateway_security_policy_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkSecurityGatewaySecurityPolicyRule_basic(gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName string) string {
	return fmt.Sprintf(`
resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-central1"
  description = "gateway security policy created to be used as reference by the rule."
}
	
resource "google_network_security_gateway_security_policy_rule" "foobar" {
  name                    = "%s"
  location                = "us-central1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true  
  description             = "my description"
  priority                = 0
  session_matcher         = "host() == 'example.com'"
  application_matcher     = "request.method == 'POST'"
  basic_profile           = "ALLOW"
}
`, gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName)
}

func testAccNetworkSecurityGatewaySecurityPolicyRule_update(gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName string) string {
	return fmt.Sprintf(`
resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-central1"
  description = "gateway security policy created to be used as reference by the rule."
}
	
resource "google_network_security_gateway_security_policy_rule" "foobar" {
  name                    = "%s"
  location                = "us-central1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = false  
  description             = "my description updated"
  priority                = 1
  session_matcher         = "host() == 'update.com'"
  application_matcher     = "request.method == 'GET'"
  tls_inspection_enabled  = false
  basic_profile           = "DENY"
}
`, gatewaySecurityPolicyName, gatewaySecurityPolicyRuleName)
}
