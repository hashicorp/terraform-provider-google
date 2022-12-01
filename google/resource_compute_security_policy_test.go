package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeSecurityPolicy_basic(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_basic(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRule(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRule(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRuleExpr(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRuleExpr(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_update(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRule(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      testAccComputeSecurityPolicy_updateSamePriority(spName),
				ExpectError: regexp.MustCompile("Two rules have the same priority, please update one of the priorities to be different."),
			},

			{
				Config: testAccComputeSecurityPolicy_update(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccComputeSecurityPolicy_withRule(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withAdvancedOptionsConfig(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_basic(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withAdvancedOptionsConfig(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_basic(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withAdaptiveProtection(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withAdaptiveProtection(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withAdaptiveProtectionUpdate(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRateLimitOptions(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRateLimitOptions(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRateLimitWithRedirectOptions(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRateLimitWithRedirectOptions(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRecaptchaOptionsConfig(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_basic(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withRecaptchaOptionsConfig(project, spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withRedirectSiteKeyUpdate(project, spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withEmptyRedirectSiteKey(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withHeadAction(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	headerName := fmt.Sprintf("tf-test-header-name-%s", randString(t, 10))
	headerNameUpdate := fmt.Sprintf("tf-test-header-name-update-%s", randString(t, 10))
	headerValue := fmt.Sprintf("tf-test-header-value-%s", randString(t, 10))
	headerValueUpdate := fmt.Sprintf("tf-test-header-value-update-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withoutHeadAction(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withHeadAction(spName, headerName, headerValue),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withHeadAction(spName, headerNameUpdate, headerValueUpdate),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withMultipleHeaders(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withoutHeadAction(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccComputeSecurityPolicy_withRecaptchaOptionsConfig(project, spName string) string {
	return fmt.Sprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "test"

  labels = {
    label-one = "value-one"
   }

  project = "%s"

  web_settings {
    integration_type  = "INVISIBLE"
    allow_all_domains = true
    allowed_domains   = ["localhost"]
  }
}

resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type        = "CLOUD_ARMOR"

  recaptcha_options_config {
    redirect_site_key = google_recaptcha_enterprise_key.primary.name
  }
}
`, project, spName)
}

func testAccComputeSecurityPolicy_withRedirectSiteKeyUpdate(project, spName string) string {
	return fmt.Sprintf(`
resource "google_recaptcha_enterprise_key" "primary1" {
  display_name = "test"

  labels = {
    label-one = "value-one"
   }

  project = "%s"

  web_settings {
    integration_type  = "INVISIBLE"
    allow_all_domains = true
    allowed_domains   = ["localhost"]
  }
}

resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type        = "CLOUD_ARMOR"

  recaptcha_options_config {
    redirect_site_key = google_recaptcha_enterprise_key.primary1.name
  }
}
`, project, spName)
}

func testAccComputeSecurityPolicy_withEmptyRedirectSiteKey(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type        = "CLOUD_ARMOR"

  recaptcha_options_config {
    redirect_site_key = ""
  }
}
`, spName)
}

func testAccCheckComputeSecurityPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_security_policy" {
				continue
			}

			pol := rs.Primary.Attributes["name"]

			_, err := config.NewComputeClient(config.userAgent).SecurityPolicies.Get(config.Project, pol).Do()
			if err == nil {
				return fmt.Errorf("Security policy %q still exists", pol)
			}
		}

		return nil
	}
}

func testAccComputeSecurityPolicy_basic(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "basic security policy"
  type        = "CLOUD_ARMOR"
}
`, spName)
}

func testAccComputeSecurityPolicy_withRule(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name = "%s"

  rule {
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

  rule {
    action   = "allow"
    priority = "2000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["10.0.0.0/24"]
      }
    }
    preview = true
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_updateSamePriority(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "updated description"

  // keep this
  rule {
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

  // add this
  rule {
    action   = "deny(403)"
    priority = "2000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["10.0.1.0/24"]
      }
    }
  }

  rule {
    action   = "allow"
    priority = "2000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["10.0.0.0/24"]
      }
    }
    preview = true
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_update(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "updated description"

  // keep this
  rule {
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

  // add this
  rule {
    action   = "deny(403)"
    priority = "1000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["10.0.1.0/24"]
      }
    }
  }

  // update this
  rule {
    action   = "allow"
    priority = "2000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["10.0.0.0/24"]
      }
    }
    description = "updated description"
    preview     = false
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withRuleExpr(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
	name = "%s"

	rule {
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

	rule {
		action   = "allow"
		priority = "2000"
		match {
			expr {
				// These fields are not yet supported (Issue hashicorp/terraform-provider-google#4497: mbang)
				// title = "Has User"
				// description = "Determines whether the request has a user account"
				expression = "evaluatePreconfiguredExpr('xss-canary')"
			}
		}
		preview = true
	}
}
`, spName)
}

func testAccComputeSecurityPolicy_withoutHeadAction(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name = "%s"

  rule {
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

  rule {
    action   = "allow"
    priority = "1000"
    match {
      expr {
        expression = "request.path.matches(\"/login.html\") && token.recaptcha_session.score < 0.2"
      }
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withHeadAction(spName, headerName, headerValue string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name = "%s"

  rule {
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

  rule {
    action   = "allow"
    priority = "1000"
    match {
      expr {
        expression = "request.path.matches(\"/login.html\") && token.recaptcha_session.score < 0.2"
      }
    }

    header_action {
      request_headers_to_adds {
        header_name  = "%s"
        header_value = "%s"
      }
    }
  }
}
`, spName, headerName, headerValue)
}

func testAccComputeSecurityPolicy_withMultipleHeaders(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name = "%s"

  rule {
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

  rule {
    action   = "allow"
    priority = "1000"
    match {
      expr {
        expression = "request.path.matches(\"/login.html\") && token.recaptcha_session.score < 0.2"
      }
    }

    header_action {
      request_headers_to_adds {
        header_name  = "reCAPTCHA-Warning"
        header_value = "high"
      }

      request_headers_to_adds {
        header_name  = "X-Hello"
        header_value = "World"
      }

      request_headers_to_adds {
        header_name  = "X-Resource"
        header_value = "test"
      }
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withAdvancedOptionsConfig(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "updated description"

  advanced_options_config {
    json_parsing = "STANDARD"
    json_custom_config {
      content_types = [
        "application/json",
        "application/vnd.api+json",
        "application/vnd.collection+json",
        "application/vnd.hyper+json"
      ]
    }
    log_level    = "VERBOSE"
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withAdaptiveProtection(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "updated description"

  adaptive_protection_config {
    layer_7_ddos_defense_config {
      enable = true
      rule_visibility = "STANDARD"
	}
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withAdaptiveProtectionUpdate(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
  name        = "%s"
  description = "updated description"

  adaptive_protection_config {
    layer_7_ddos_defense_config {
      enable = false
      rule_visibility = "STANDARD"
    }
  }
}
`, spName)
}

func testAccComputeSecurityPolicy_withRateLimitOptions(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
	name        = "%s"
	description = "updated description"

	rule {
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

	rule {
		action = "throttle"
		priority = 100
		match {
			versioned_expr = "SRC_IPS_V1"
			config {
				src_ip_ranges = [
					"0.0.0.0/32",
				]
			}
		}
		rate_limit_options {
			conform_action = "allow"
			exceed_action = "deny(403)"
			enforce_on_key = "IP"
			rate_limit_threshold {
				count = 100
				interval_sec = 60
			}
		}
	}
}
`, spName)
}

func testAccComputeSecurityPolicy_withRateLimitWithRedirectOptions(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
	name        = "%s"
	description = "updated description"

	rule {
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

	rule {
		action = "throttle"
		priority = 100
		match {
			versioned_expr = "SRC_IPS_V1"
			config {
				src_ip_ranges = [
					"0.0.0.0/32",
				]
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
				count = 100
				interval_sec = 60
			}
		}
	}
}
`, spName)
}

func TestAccComputeSecurityPolicy_withRedirectOptionsRecaptcha(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRedirectOptionsRecaptcha(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRedirectOptionsUpdate(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRedirectOptionsRecaptcha(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSecurityPolicy_withRedirectOptionsExternal(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSecurityPolicy_withRedirectOptionsExternal(t *testing.T) {
	t.Parallel()

	spName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSecurityPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSecurityPolicy_withRedirectOptionsExternal(spName),
			},
			{
				ResourceName:      "google_compute_security_policy.policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeSecurityPolicy_withRedirectOptionsRecaptcha(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
	name        = "%s"

	rule {
		action   = "redirect"
		priority = "2147483647"
		match {
			versioned_expr = "SRC_IPS_V1"
			config {
				src_ip_ranges = ["*"]
			}
		}
		description = "default rule"
		redirect_options {
			type = "GOOGLE_RECAPTCHA"
		}
	}
}
`, spName)
}

func testAccComputeSecurityPolicy_withRedirectOptionsExternal(spName string) string {
	return fmt.Sprintf(`
resource "google_compute_security_policy" "policy" {
	name        = "%s"

	rule {
		action   = "redirect"
		priority = "2147483647"
		match {
			versioned_expr = "SRC_IPS_V1"
			config {
				src_ip_ranges = ["*"]
			}
		}
		description = "default rule"
		redirect_options {
			type = "EXTERNAL_302"
			target = "https://example.com"
		}
	}
}
`, spName)
}
