// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeSecurityPolicyRule_basicUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_preBasicUpdate(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_postBasicUpdate(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_withRuleExpr(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_withRuleExpr(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_extendedUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_extPreUpdate(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccComputeSecurityPolicyRule_extPosUpdateSamePriority(context),
				ExpectError: regexp.MustCompile("Cannot have rules with the same priorities."),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_extPosUpdate(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_withPreconfiguredWafConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_withPreconfiguredWafConfig_create(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withPreconfiguredWafConfig_update(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withPreconfiguredWafConfig_clear(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_security_policy_rule.policy_rule", "preconfigured_waf_config.0"),
				),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeSecurityPolicyRule_preBasicUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
  type        = "CLOUD_ARMOR"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic rule pre update"
  action          = "allow"
  priority        = 100
  preview         = false
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["192.168.0.0/16", "10.0.0.0/8"]
    }
  }
}
`, context)
}

func testAccComputeSecurityPolicyRule_postBasicUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
  type        = "CLOUD_ARMOR"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic rule post update"
  action          = "deny(403)"
  priority        = 100
  preview         = true
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["172.16.0.0/12"]
    }
  }
}
`, context)
}

func testAccComputeSecurityPolicyRule_withRuleExpr(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic description"
  action          = "allow"
  priority        = "2000"
  match {
    expr {
      expression = "evaluatePreconfiguredExpr('xss-canary')"
    }
  }
  preview = true
}
`, context)
}

func testAccComputeSecurityPolicyRule_extPreUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic description"
  action          = "allow"
  priority        = "2000"
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["10.0.0.0/24"]
    }
  }
  preview = true
}
`, context)
}

func testAccComputeSecurityPolicyRule_extPosUpdateSamePriority(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
}

//add this
resource "google_compute_security_policy_rule" "policy_rule2" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic description"
  action          = "deny(403)"
  priority        = "2000"
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["10.0.0.0/24"]
    }
  }
  preview = true
}

//keep this
resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic description"
  action          = "allow"
  priority        = "2000"
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["10.0.0.0/24"]
    }
  }
  preview = true
}
`, context)
}

func testAccComputeSecurityPolicyRule_extPosUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
}

//add this
resource "google_compute_security_policy_rule" "policy_rule2" {
  security_policy = google_compute_security_policy.default.name
  description     = "basic description"
  action          = "deny(403)"
  priority        = "1000"
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["10.0.0.0/24"]
    }
  }
  preview = true
}

//update this
resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "updated description"
  action          = "allow"
  priority        = "2000"
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["10.0.0.0/24"]
    }
  }
  preview = true
}
`, context)
}

func testAccComputeSecurityPolicyRule_withPreconfiguredWafConfig_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "Global security policy - create"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "Rule with preconfiguredWafConfig - create"
  action   = "deny"
  priority = "1000"
  match {
    expr {
      expression = "evaluatePreconfiguredWaf('sqli-stable')"
    }
  }
  preconfigured_waf_config {
    exclusion {
      request_cookie {
        operator = "EQUALS_ANY"
      }
      request_header {
        operator = "EQUALS"
        value    = "Referer"
      }
      request_uri {
        operator = "STARTS_WITH"
        value    = "/admin"
      }
      request_query_param {
        operator = "EQUALS"
        value    = "password"
      }
      request_query_param {
        operator = "STARTS_WITH"
        value    = "freeform"
      }
      target_rule_set = "sqli-stable"
    }
    exclusion {
      request_query_param {
        operator = "CONTAINS"
        value    = "password"
      }
      request_query_param {
        operator = "STARTS_WITH"
        value    = "freeform"
      }
      target_rule_set = "xss-stable"
    }
  }
  preview = false
}
`, context)
}

func testAccComputeSecurityPolicyRule_withPreconfiguredWafConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "Global security policy - update"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "Rule with preconfiguredWafConfig - update"
  action   = "deny"
  priority = "1000"
  match {
    expr {
      expression = "evaluatePreconfiguredWaf('rce-stable') || evaluatePreconfiguredWaf('xss-stable')"
    }
  }
  preconfigured_waf_config {
    exclusion {
      request_uri {
        operator = "STARTS_WITH"
        value    = "/admin"
      }
      target_rule_set = "rce-stable"
    }
    exclusion {
      request_query_param {
        operator = "CONTAINS"
        value    = "password"
      }
      request_query_param {
        operator = "STARTS_WITH"
        value    = "freeform"
      }
      request_query_param {
        operator = "EQUALS"
        value    = "description"
      }
      request_cookie {
        operator = "CONTAINS"
        value    = "TokenExpired"
      }
      target_rule_set = "xss-stable"
      target_rule_ids = [
        "owasp-crs-v030001-id941330-xss",
        "owasp-crs-v030001-id941340-xss",
      ]
    }
  }
  preview = false
}
`, context)
}

func testAccComputeSecurityPolicyRule_withPreconfiguredWafConfig_clear(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "Global security policy - clear"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "Rule with preconfiguredWafConfig - clear"
  action   = "deny"
  priority = "1000"
  match {
    expr {
      expression = "evaluatePreconfiguredWaf('rce-stable') || evaluatePreconfiguredWaf('xss-stable')"
    }
  }
  preview = false
}
`, context)
}
