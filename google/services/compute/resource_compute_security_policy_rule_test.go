// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func TestAccComputeSecurityPolicyRule_withRateLimitOptions(t *testing.T) {
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
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptionsCreate(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptionsUpdate(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_withRateLimit_withEnforceOnKeyConfigs(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKeyConfigs(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_withRateLimitOption_withMultipleEnforceOnKeyConfigs(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOption_withMultipleEnforceOnKeyConfigs(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOption_withMultipleEnforceOnKeyConfigs2(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_EnforceOnKeyUpdates(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withoutRateLimitOptions(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKeyName(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKey(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKeyConfigs(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKey(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKeyName(spName),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_withExprOptions(t *testing.T) {
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
				Config: testAccComputeSecurityPolicyRule_withExprOptions(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicyRule_modifyExprOptions(t *testing.T) {
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
			{
				Config: testAccComputeSecurityPolicyRule_withExprOptions(context),
			},
			{
				ResourceName:      "google_compute_security_policy_rule.policy_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicyRule_modifyExprOptions(context),
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

func testAccComputeSecurityPolicyRule_withRateLimitOptionsCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_compute_security_policy" "default" {
    name        = "tf-test%{random_suffix}"
    description = "basic global security policy"
  }

  resource "google_compute_security_policy_rule" "policy_rule" {
    security_policy = google_compute_security_policy.default.name
    description     = "rule create with rate limit"
    priority        = 101
    action          = "rate_based_ban"
    rate_limit_options {
      rate_limit_threshold {
        count = 500
        interval_sec = 10
      }
      conform_action = "allow"
      exceed_action = "deny(404)"
      enforce_on_key = "ALL"
      ban_threshold {
        count = 750
        interval_sec = 180
      }
      ban_duration_sec = 180
    }
    match {
      config {
        src_ip_ranges = [
          "*"
        ]
      }
      versioned_expr = "SRC_IPS_V1"
    }
  }
`, context)
}

func testAccComputeSecurityPolicyRule_withRateLimitOptionsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_compute_security_policy" "default" {
    name        = "tf-test%{random_suffix}"
    description = "basic global security policy"
  }

  resource "google_compute_security_policy_rule" "policy_rule" {
    security_policy = google_compute_security_policy.default.name
    description     = "rule update with rate limit update"
    priority        = 101
    action          = "rate_based_ban"
    rate_limit_options {
      rate_limit_threshold {
        count = 1000
        interval_sec = 30
      }
      conform_action = "allow"
      exceed_action = "deny(404)"
      enforce_on_key = "ALL"
      ban_threshold {
        count = 2000
        interval_sec = 180
      }
      ban_duration_sec = 300
    }
    match {
      config {
        src_ip_ranges = [
          "*"
        ]
      }
      versioned_expr = "SRC_IPS_V1"
    }
  }
`, context)
}

func testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKey(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic policy base"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "throttle rule withEnforceOnKey"
  action          = "throttle"
  priority        = "100"
  
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["*"]
    }
  }

  rate_limit_options {
    conform_action = "allow"
    exceed_action = "redirect"

    enforce_on_key = "IP"

    exceed_redirect_options {
      type = "EXTERNAL_302"
      target = "https://www.example.com"
    }

    rate_limit_threshold {
      count = 10
      interval_sec = 60
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKeyConfigs(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic policy base"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "throttle rule withEnforceOnKeyConfigs"
  action          = "throttle"
  priority        = "100"

  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["*"]
    }
  }

  rate_limit_options {
    conform_action = "allow"
    exceed_action = "redirect"

    enforce_on_key = ""

    enforce_on_key_configs {
      enforce_on_key_type = "IP"
    }
    exceed_redirect_options {
      type = "EXTERNAL_302"
      target = "https://www.example.com"
    }

    rate_limit_threshold {
      count = 10
      interval_sec = 60
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicyRule_withRateLimitOption_withMultipleEnforceOnKeyConfigs(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic policy base"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "throttle rule with withMultipleEnforceOnKeyConfigs"
  action          = "throttle"
  priority        = "100"

  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["*"]
    }
  }

  rate_limit_options {
    conform_action = "allow"
    exceed_action = "deny(429)"

    rate_limit_threshold {
      count = 10
      interval_sec = 60
    }

    enforce_on_key = ""

    enforce_on_key_configs {
      enforce_on_key_type = "HTTP_PATH"
    }

    enforce_on_key_configs {
      enforce_on_key_type = "HTTP_HEADER"
      enforce_on_key_name = "user-agent"
    }

    enforce_on_key_configs {
      enforce_on_key_type = "REGION_CODE"
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicyRule_withRateLimitOption_withMultipleEnforceOnKeyConfigs2(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic policy base"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "throttle rule withMultipleEnforceOnKeyConfigs2"
  action          = "throttle"
  priority        = "100"

  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["*"]
    }
  }

  rate_limit_options {
    conform_action = "allow"
    exceed_action = "deny(429)"

    rate_limit_threshold {
      count = 10
      interval_sec = 60
    }

    enforce_on_key = ""

    enforce_on_key_configs {
      enforce_on_key_type = "REGION_CODE"
    }

    enforce_on_key_configs {
      enforce_on_key_type = "TLS_JA3_FINGERPRINT"
    }

    enforce_on_key_configs {
      enforce_on_key_type = "USER_IP"
    }
  }
}

`, spName)
}

func testAccComputeSecurityPolicyRule_withRateLimitOptions_withoutRateLimitOptions(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic policy base"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "basic policy rule withoutRateLimitOptions"
  action          = "deny(403)"
  priority        = "100"
  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["*"]
    }
  }
}

`, spName)
}

func testAccComputeSecurityPolicyRule_withRateLimitOptions_withEnforceOnKeyName(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic policy base"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.policy.name
  description     = "throttle rule withEnforceOnKeyName"
  action          = "throttle"
  priority        = "100"

  match {
    versioned_expr = "SRC_IPS_V1"
    config {
      src_ip_ranges = ["*"]
    }
  }

  rate_limit_options {
    conform_action = "allow"
    exceed_action = "redirect"

    enforce_on_key = "HTTP_HEADER"
    enforce_on_key_name = "user-agent"

    exceed_redirect_options {
      type = "EXTERNAL_302"
      target = "https://www.example.com"
    }

    rate_limit_threshold {
      count = 10
      interval_sec = 60
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicyRule_withExprOptions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
	name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "reCAPTCHA rule"
  action          = "deny(403)"
  priority        = "2000"
  preview         = true
  match {
    expr {
      expression = "request.path.endsWith('RegisterWithEmail') && token.recaptcha_action.score >= 0.8 && (token.recaptcha_action.valid)"
    }
    expr_options {
      recaptcha_options {
        action_token_site_keys = [
          "placeholder-recaptcha-action-site-key-01",
          "placeholder-recaptcha-action-site-key-02"
        ]
        session_token_site_keys = [
          "placeholder-recaptcha-session-site-key-1",
          "placeholder-recaptcha-session-site-key-2"
        ]
      }
    }
  }
}
`, context)
}

func testAccComputeSecurityPolicyRule_modifyExprOptions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_security_policy" "default" {
  name        = "tf-test%{random_suffix}"
  description = "basic global security policy"
}

resource "google_compute_security_policy_rule" "policy_rule" {
  security_policy = google_compute_security_policy.default.name
  description     = "modified reCAPTCHA rule"
  action          = "deny(403)"
  priority        = "2000"
  preview         = true
  match {
    expr {
      expression = "request.path.endsWith('RegisterWithEmail') && token.recaptcha_action.score >= 0.8 && (token.recaptcha_action.valid)"
    }
    expr_options {
      recaptcha_options {
        action_token_site_keys = [
          "placeholder-recaptcha-action-site-key-09",
          "placeholder-recaptcha-action-site-key-08",
          "placeholder-recaptcha-action-site-key-07"
        ]
        session_token_site_keys = [
          "placeholder-recaptcha-session-site-key-1"
        ]
      }
    }
  }
}
`, context)
}
