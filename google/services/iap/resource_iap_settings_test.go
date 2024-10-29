// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iap_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccIapSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
		"role":            "roles/iap.settingsAdmin",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckIapSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIapSettings_basic(context),
			},
			{
				ResourceName:            "google_iap_settings.iap_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_settings.0.workforce_identity_settings.0.oauth2.0.client_secret"},
			},
			{
				Config: testAccIapSettings_update(context),
			},
			{
				ResourceName:            "google_iap_settings.iap_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_settings.0.workforce_identity_settings.0.oauth2.0.client_secret"},
			},
		},
	})
}

func testAccIapSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_compute_region_backend_service" "default" {
  name                            = "tf-test-iap-settings-tf%{random_suffix}"
  region                          = "us-central1"
  health_checks                   = [google_compute_health_check.default.id]
  connection_draining_timeout_sec = 10
  session_affinity                = "CLIENT_IP"
}

resource "google_compute_health_check" "default" {
  name               = "tf-test-iap-bs-health-check%{random_suffix}"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "80"
  }
}

resource "google_iap_settings" "iap_settings" {
  name = "projects/${data.google_project.project.number}/iap_web/compute-us-central1/services/${google_compute_region_backend_service.default.name}"
  access_settings {
    identity_sources = ["WORKFORCE_IDENTITY_FEDERATION"]
    allowed_domains_settings {
      domains = ["test.abc.com"]
      enable  = true
    }
    cors_settings {
      allow_http_options = true
    }
    reauth_settings {
      method = "SECURE_KEY"
      max_age = "305s"
      policy_type = "MINIMUM"
    }
    gcip_settings {
      login_page_uri = "https://test.com/?apiKey=abc"
    }
    oauth_settings {
      login_hint = "test"
    }
    workforce_identity_settings {
      workforce_pools = ["wif-pool"]
      oauth2 {
        client_id = "test-client-id"
	client_secret = "test-client-secret"
      }
    }
  }
  application_settings {
    cookie_domain = "test.abc.com"
    csm_settings {
      rctoken_aud = "test-aud-set"
    }
    access_denied_page_settings {
      access_denied_page_uri = "test-uri"
      generate_troubleshooting_uri = true
      remediation_token_generation_enabled = false
    }
    attribute_propagation_settings {
      output_credentials = ["HEADER"]
      expression = "attributes.saml_attributes.filter(attribute, attribute.name in [\"test1\", \"test2\"])"
      enable = false
    }
  }
}
`, context)
}

func testAccIapSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "my_project" {
  name            = "tf-test-%{random_suffix}"
  project_id      = "tf-test-%{random_suffix}"
  org_id = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_project.my_project]

  create_duration = "60s"
}

resource "google_project_service" "project_service" {
  project = google_project.my_project.project_id
  service = "iap.googleapis.com"

  # Needed for CI tests for permissions to propagate, should not be needed for actual usage
  depends_on = [time_sleep.wait_60_seconds]
}

resource "google_app_engine_application" "app" {
  project     = google_project_service.project_service.project
  location_id = "us-central"
}

resource "google_iap_web_type_app_engine_iam_member" "foo" {
  project = google_app_engine_application.app.project
  app_id = google_app_engine_application.app.app_id
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}

resource "google_iap_settings" "iap_settings" {
  name = "projects/${google_project.my_project.project_id}/iap_web/appengine-${google_app_engine_application.app.app_id}"
  access_settings {
    allowed_domains_settings {
      domains = ["appengine.abc.com"]
      enable  = true
    }
    cors_settings {
      allow_http_options = true
    }
  }
  application_settings {
    cookie_domain = "appengine.abc.com"
    attribute_propagation_settings {
      output_credentials = ["JWT"]
      expression = "attributes.saml_attributes.filter(attribute, attribute.name in [\"test1\", \"test2\"])"
      enable = false
    }
  }
  depends_on = [google_iap_web_type_app_engine_iam_member.foo]
}
`, context)
}
