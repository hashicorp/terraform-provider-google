// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccComputeRegionSecurityPolicy_regionSecurityPolicyBasicUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_basic(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_update(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "basic region security policy"
  type        = "CLOUD_ARMOR_NETWORK"

  ddos_protection_config {
    ddos_protection = "STANDARD"
  }
}
`, context)
}

func testAccComputeRegionSecurityPolicy_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "basic update region security policy"
  type        = "CLOUD_ARMOR_NETWORK"

  ddos_protection_config {
    ddos_protection = "ADVANCED"
  }
}
`, context)
}

func TestAccComputeRegionSecurityPolicy_regionSecurityPolicyUserDefinedFieldsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_withoutUserDefinedFields(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withUserDefinedFields(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withUserDefinedFieldsUpdate(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withoutUserDefinedFields(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_withoutUserDefinedFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "basic region security policy"
  type        = "CLOUD_ARMOR_NETWORK"
}
`, context)
}

func testAccComputeRegionSecurityPolicy_withUserDefinedFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "basic update region security policy"
  type        = "CLOUD_ARMOR_NETWORK"
  user_defined_fields {
    name = "SIG1_AT_0"
    base = "TCP"
    offset = 8
    size = 2
    mask = "0x8F00"
  }
}
`, context)
}

func testAccComputeRegionSecurityPolicy_withUserDefinedFieldsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_region_security_policy" "policy" {
  name        = "tf-test%{random_suffix}"
  description = "basic update region security policy"
  type        = "CLOUD_ARMOR_NETWORK"
  user_defined_fields {
    name = "SIG1_AT_0"
    base = "UDP"
    offset = 4
    size = 4
    mask = "0xFFFF"
  }
  user_defined_fields {
    name = "SIG2_AT_8"
    base = "TCP"
    offset = 8
    size = 2
    mask = "0x8700"
  }
}
`, context)
}

func TestAccComputeRegionSecurityPolicy_regionSecurityPolicyWithRulesBasicUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_withRules(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withRulesUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_region_security_policy.policy", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_withRules(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
		name        = "tf-test%{random_suffix}"
		description = "region security policy with rules updated"
		type        = "CLOUD_ARMOR"

		rules {
			action   = "deny"
			priority = "1000"
			match {
				expr {
					expression = "request.path.matches(\"/login.html\") && token.recaptcha_session.score < 0.2"
				}
			}
		}


		rules {
			action   = "deny"
			priority = "2147483647"
			match {
				versioned_expr = "SRC_IPS_V1"
				config {
					src_ip_ranges = ["*"]
				}
			}
			description = "default rule"
		}

	}
	`, context)
}

func testAccComputeRegionSecurityPolicy_withRulesUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_compute_region_security_policy" "policy" {
		name        = "tf-test%{random_suffix}"
		description = "region security policy with rules updated"
		type        = "CLOUD_ARMOR"

		rules {
			action   = "allow"
			priority = "2147483647"
			match {
				versioned_expr = "SRC_IPS_V1"
				config {
					src_ip_ranges = ["*"]
				}
			}
			description = "default rule updated"
		}
	}
	`, context)
}

func TestAccComputeRegionSecurityPolicy_regionSecurityPolicyWithRulesPreconfiguredWafConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_withPreconfiguredWafConfig(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withPreconfiguredWafConfig_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_region_security_policy.policy", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_withPreconfiguredWafConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
			name	= "tf-test%{random_suffix}"
			type	= "CLOUD_ARMOR"

			rules {
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
	
			rules {
				action   = "allow"
				priority = "2147483647"
				match {
					versioned_expr = "SRC_IPS_V1"
					config {
						src_ip_ranges = ["*"]
					}
				}
				description = "default rule"
			}

		}
	`, context)
}

func testAccComputeRegionSecurityPolicy_withPreconfiguredWafConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
			name	= "tf-test%{random_suffix}"
			type	= "CLOUD_ARMOR"

			rules {
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
						request_cookie {
							operator = "EQUALS"
							value    = "Referer"
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
						target_rule_set = "xss-stable"
						target_rule_ids = [
							"owasp-crs-v030001-id941330-xss",
							"owasp-crs-v030001-id941340-xss",
						]
					}
				}
				preview = false
			}

			rules {
				action   = "allow"
				priority = "2147483647"
				match {
					versioned_expr = "SRC_IPS_V1"
					config {
						src_ip_ranges = ["*"]
					}
				}
				description = "default rule"
			}

		}
	`, context)
}

func TestAccComputeRegionSecurityPolicy_regionSecurityPolicyWithRulesRateLimitOptions(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_withRateLimitOptions(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withRateLimitOptions_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_region_security_policy.policy", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_withRateLimitOptions(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
			name	= "tf-test%{random_suffix}"
			type	= "CLOUD_ARMOR"
			region  = "us-west2"

			rules {
				priority = "1000"
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

			rules {
				action   = "allow"
				priority = "2147483647"
				preview 	= false
				match {
					versioned_expr = "SRC_IPS_V1"
					config {
						src_ip_ranges = ["*"]
					}
				}
				description = "default rule"
			}
		}
	`, context)
}

func testAccComputeRegionSecurityPolicy_withRateLimitOptions_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
			name	= "tf-test%{random_suffix}"
			type	= "CLOUD_ARMOR"
			region  = "us-west2"

			rules {
				priority = "100"
				action          = "rate_based_ban"
				rate_limit_options {
					rate_limit_threshold {
						count = 100
						interval_sec = 30
					}
					conform_action = "allow"
					exceed_action = "deny(404)"
					enforce_on_key = "HTTP_HEADER"
					enforce_on_key_name = "user-agent"
					ban_threshold {
						count = 1000
						interval_sec = 300
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

			rules {
				action   = "allow"
				priority = "2147483647"
				preview 	= false
				match {
					versioned_expr = "SRC_IPS_V1"
					config {
						src_ip_ranges = ["*"]
					}
				}
				description = "default rule"
			}
		}
	`, context)
}

func TestAccComputeRegionSecurityPolicy_regionSecurityPolicyWithRulesMultipleEnforceOnKeyConfigs(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRegionSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionSecurityPolicy_withMultipleEnforceOnKeyConfigs(context),
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionSecurityPolicy_withMultipleEnforceOnKeyConfigs_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_region_security_policy.policy", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_region_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionSecurityPolicy_withMultipleEnforceOnKeyConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
			name	= "tf-test%{random_suffix}"
			type	= "CLOUD_ARMOR"
			region  = "us-west2"

			rules {
				priority = "1000"
				action          = "throttle"
				rate_limit_options {
					conform_action = "allow"
					exceed_action = "deny(429)"

					rate_limit_threshold {
						count = 10
						interval_sec = 60
					}

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
				match {
					config {
						src_ip_ranges = [
							"*"
						]
					}
					versioned_expr = "SRC_IPS_V1"
				}
			}

			rules {
				action   = "allow"
				priority = "2147483647"
				preview 	= false
				match {
					versioned_expr = "SRC_IPS_V1"
					config {
						src_ip_ranges = ["*"]
					}
				}
				description = "default rule"
			}
		}
	`, context)
}

func testAccComputeRegionSecurityPolicy_withMultipleEnforceOnKeyConfigs_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_compute_region_security_policy" "policy" {
			name	= "tf-test%{random_suffix}"
			type	= "CLOUD_ARMOR"
			region  = "us-west2"

			rules {
				priority = "100"
				action          = "throttle"
				rate_limit_options {
					conform_action = "allow"
					exceed_action = "deny(429)"

					rate_limit_threshold {
						count = 10
						interval_sec = 60
					}

					enforce_on_key_configs {
						enforce_on_key_type = "USER_IP"
					}

					enforce_on_key_configs {
						enforce_on_key_type = "TLS_JA3_FINGERPRINT"
					}

					enforce_on_key_configs {
						enforce_on_key_type = "REGION_CODE"
					}
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

			rules {
				action   = "allow"
				priority = "2147483647"
				preview 	= false
				match {
					versioned_expr = "SRC_IPS_V1"
					config {
						src_ip_ranges = ["*"]
					}
				}
				description = "default rule"
			}
		}
	`, context)
}
