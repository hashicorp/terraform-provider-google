// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkSecurityAuthzPolicy_networkServicesAuthzPolicyHttpRules(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityAuthzPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityAuthzPolicy_networkServicesAuthzPolicyHttpRules(context),
			},
			{
				ResourceName:            "google_network_security_authz_policy.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkSecurityAuthzPolicy_networkServicesAuthzPolicyHttpRules(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "lb-network-%{random_suffix}"
  project                 = "%{project}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "backend-subnet-%{random_suffix}"
  project       = "%{project}"
  region        = "us-west1"
  ip_cidr_range = "10.1.2.0/24"
  network       = google_compute_network.default.id
}

resource "google_compute_subnetwork" "proxy_only" {
  name          = "proxy-only-subnet-%{random_suffix}"
  project       = "%{project}"
  region        = "us-west1"
  ip_cidr_range = "10.129.0.0/23"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.default.id
}

resource "google_compute_address" "default" {
  name         = "l7-ilb-ip-address-%{random_suffix}"
  project      = "%{project}"
  region       = "us-west1"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
}

resource "google_compute_region_health_check" "default" {
  name    = "l7-ilb-basic-check-%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}

resource "google_compute_region_backend_service" "url_map" {
  name                  = "l7-ilb-backend-service-%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"

  health_checks = [google_compute_region_health_check.default.id]
}

resource "google_compute_region_url_map" "default" {
  name            = "l7-ilb-map-%{random_suffix}"
  project         = "%{project}"
  region          = "us-west1"
  default_service = google_compute_region_backend_service.url_map.id
}

resource "google_compute_region_target_http_proxy" "default" {
  name    = "l7-ilb-proxy-%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"
  url_map = google_compute_region_url_map.default.id
}

resource "google_compute_forwarding_rule" "default" {
  name                  = "l7-ilb-forwarding-rule-%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"
  network               = google_compute_network.default.id
  subnetwork            = google_compute_subnetwork.default.id
  ip_protocol           = "TCP"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  ip_address            = google_compute_address.default.id

  depends_on = [google_compute_subnetwork.proxy_only]
}

resource "google_compute_region_backend_service" "authz_extension" {
  name    = "authz-service-%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"

  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_name             = "grpc"
}

resource "google_network_services_authz_extension" "default" {
  name     = "my-authz-ext-%{random_suffix}"
  project  = "%{project}"
  location = "us-west1"

  description           = "my description"
  load_balancing_scheme = "INTERNAL_MANAGED"
  authority             = "ext11.com"
  service               = google_compute_region_backend_service.authz_extension.self_link
  timeout               = "0.1s"
  fail_open             = false
  forward_headers       = ["Authorization"]
}

resource "google_network_security_authz_policy" "default" {
  name        = "tf-test-my-authz-policy-%{random_suffix}"
  project     = "%{project}"
  location    = "us-west1"
  description = "my description"

  target {
    load_balancing_scheme = "INTERNAL_MANAGED"
    resources = [ google_compute_forwarding_rule.default.self_link ]
  }

  action = "CUSTOM"
  custom_provider {
	authz_extension {
      resources = [ google_network_services_authz_extension.default.id ]
    }
  }

  http_rules {
    from {
	  not_sources {
        # Prefix
		principals {
          ignore_case = false
          prefix      = "prefix"
        }
        resources {
          iam_service_account {
            ignore_case = false
          	prefix      = "prefix"
          }
          tag_value_id_set {
            ids = ["1"]
          }
        }
		# Suffix / Ignore case
		principals {
		  ignore_case = true
		  suffix      = "suffix"
		}
		resources {
		  iam_service_account {
		    ignore_case = true
			  suffix      = "suffix"
		  }
		  tag_value_id_set {
		    ids = ["2"]
		  }
		}
		# Exact
		principals {
		  ignore_case = true
		  exact       = "exact"
		}
		resources {
		  iam_service_account {
		    ignore_case = true
			exact       = "exact"
		  }
		  tag_value_id_set {
		    ids = ["3"]
		  }
		}
		# Contains / Ignore case
		principals {
		  ignore_case = true
		  contains    = "contains"
		}
		resources {
		  iam_service_account {
		    ignore_case = true
			contains    = "contains"
		  }
		  tag_value_id_set {
		    ids = ["4"]
		  }
		}
      }
      sources {
		# Prefix
        principals {
          ignore_case = false
          prefix      = "prefix"
        }
        resources {
          iam_service_account {
            ignore_case = false
          	prefix      = "prefix"
          }
          tag_value_id_set {
            ids = ["1"]
          }
        }
		# Suffix / Ignore case
		principals {
			ignore_case = true
			suffix      = "suffix"
        }
        resources {
          iam_service_account {
            ignore_case = true
          	suffix      = "suffix"
          }
          tag_value_id_set {
            ids = ["2"]
          }
        }
		# Exact
		principals {
          exact       = "exact"
          ignore_case = false
        }
        resources {
          iam_service_account {
            exact       = "exact"
          	ignore_case = false
          }
          tag_value_id_set {
            ids = ["3"]
          }
        }
		# Contains / Ignore case
		principals {
          contains    = "contains"
          ignore_case = true
        }
        resources {
          iam_service_account {
            contains    = "contains"
          	ignore_case = true
          }
          tag_value_id_set {
            ids = ["4"]
          }
        }
      }
    }
    to {
      operations {
        methods = ["GET", "PUT", "POST", "HEAD", "PATCH", "DELETE", "OPTIONS"]
		header_set {
          # Prefix
		  headers {
            name = "PrefixHeader"
            value {
			  ignore_case = false
			  prefix      = "prefix"
            }
          }
		  # Suffix / Ignore case
		  headers {
			name = "SuffixHeader"
			value {
			  ignore_case = true
			  suffix      = "suffix"
			}
		  }
		  # Exact
		  headers {
            name = "ExactHeader"
            value {
              exact       = "exact"
          	  ignore_case = false
            }
          }
		  # Contains / Ignore case
		  headers {
            name = "ContainsHeader"
            value {
              contains    = "contains"
          	  ignore_case = true
            }
          }
        }
        # Prefix
		hosts {
			ignore_case = false
			prefix      = "prefix"
        }
		paths {
          ignore_case = false
          prefix      = "prefix"
        }
		# Suffix / Ignore case
		hosts {
          ignore_case = true
          suffix      = "suffix"
        }
        paths {
          ignore_case = true
          suffix      = "suffix"
        }
		# Exact
		hosts {
          exact       = "exact"
          ignore_case = false
        }
		paths {
		  exact       = "exact"
		  ignore_case = false
        }
		# Contains / Ignore case
		hosts {
          contains    = "contains"
          ignore_case = true
        }
		paths {
          contains    = "contains"
          ignore_case = true
        }
      }
    }
	when = "request.host.endsWith('.example.com')"
  }

  labels = {
    foo = "bar"
  }
}
`, context)
}
