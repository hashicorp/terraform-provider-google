// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccComputeRegionTargetHttpsProxy_update(t *testing.T) {
	t.Parallel()

	resourceSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetHttpsProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionTargetHttpsProxy_basic1(resourceSuffix),
			},
			{
				ResourceName:      "google_compute_region_target_https_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionTargetHttpsProxy_basic2(resourceSuffix),
			},
			{
				ResourceName:      "google_compute_region_target_https_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionTargetHttpsProxy_basic3(resourceSuffix),
			},
			{
				ResourceName:      "google_compute_region_target_https_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionTargetHttpsProxy_basic1(id string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_https_proxy" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "httpsproxy-test-%s"
  url_map          = google_compute_region_url_map.foobar1.self_link
  ssl_certificates = [google_compute_region_ssl_certificate.foobar1.self_link]
}

resource "google_compute_region_backend_service" "foobar1" {
  name          = "httpsproxy-test-backend1-%s"
  health_checks = [google_compute_region_health_check.zero.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_backend_service" "foobar2" {
  name          = "httpsproxy-test-backend2-%s"
  health_checks = [google_compute_region_health_check.one.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name     = "httpsproxy-test-health-check1-%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_health_check" "one" {
  name     = "httpsproxy-test-health-check2-%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_url_map" "foobar1" {
  name            = "httpsproxy-test-url-map1-%s"
  default_service = google_compute_region_backend_service.foobar1.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar1.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar1.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar1.self_link
  }
}

resource "google_compute_region_url_map" "foobar2" {
  name            = "httpsproxy-test-url-map2-%s"
  default_service = google_compute_region_backend_service.foobar2.self_link
  host_rule {
    hosts        = ["mysite2.com", "myothersite2.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar2.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar2.self_link
    }
  }
  test {
    host    = "mysite2.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar2.self_link
  }
}

resource "google_compute_region_ssl_certificate" "foobar1" {
  name        = "httpsproxy-test-cert1-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_region_ssl_certificate" "foobar2" {
  name        = "httpsproxy-test-cert2-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}
`, id, id, id, id, id, id, id, id, id)
}

func testAccComputeRegionTargetHttpsProxy_basic2(id string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_https_proxy" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "httpsproxy-test-%s"
  url_map     = google_compute_region_url_map.foobar2.self_link
  ssl_certificates = [
    google_compute_region_ssl_certificate.foobar1.self_link,
    google_compute_region_ssl_certificate.foobar2.self_link,
  ]
}

resource "google_compute_region_backend_service" "foobar1" {
  name          = "httpsproxy-test-backend1-%s"
  health_checks = [google_compute_region_health_check.zero.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_backend_service" "foobar2" {
  name          = "httpsproxy-test-backend2-%s"
  health_checks = [google_compute_region_health_check.one.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name     = "httpsproxy-test-health-check1-%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_health_check" "one" {
  name     = "httpsproxy-test-health-check2-%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_url_map" "foobar1" {
  name            = "httpsproxy-test-url-map1-%s"
  default_service = google_compute_region_backend_service.foobar1.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar1.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar1.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar1.self_link
  }
}

resource "google_compute_region_url_map" "foobar2" {
  name            = "httpsproxy-test-url-map2-%s"
  default_service = google_compute_region_backend_service.foobar2.self_link
  host_rule {
    hosts        = ["mysite2.com", "myothersite2.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar2.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar2.self_link
    }
  }
  test {
    host    = "mysite2.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar2.self_link
  }
}

resource "google_compute_region_ssl_certificate" "foobar1" {
  name        = "httpsproxy-test-cert1-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_region_ssl_certificate" "foobar2" {
  name        = "httpsproxy-test-cert2-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}
`, id, id, id, id, id, id, id, id, id)
}

func testAccComputeRegionTargetHttpsProxy_basic3(id string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_https_proxy" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "httpsproxy-test-%s"
  url_map          = google_compute_region_url_map.foobar2.self_link
  ssl_certificates = [google_compute_region_ssl_certificate.foobar2.self_link]
  ssl_policy       = google_compute_region_ssl_policy.foobar.self_link
}

resource "google_compute_region_backend_service" "foobar1" {
  name          = "httpsproxy-test-backend1-%s"
  health_checks = [google_compute_region_health_check.zero.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_backend_service" "foobar2" {
  name          = "httpsproxy-test-backend2-%s"
  health_checks = [google_compute_region_health_check.one.self_link]
  protocol      = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name     = "httpsproxy-test-health-check1-%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_health_check" "one" {
  name     = "httpsproxy-test-health-check2-%s"
  http_health_check {
    port = 443
  }
}

resource "google_compute_region_url_map" "foobar1" {
  name            = "httpsproxy-test-url-map1-%s"
  default_service = google_compute_region_backend_service.foobar1.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar1.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar1.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar1.self_link
  }
}

resource "google_compute_region_url_map" "foobar2" {
  name            = "httpsproxy-test-url-map2-%s"
  default_service = google_compute_region_backend_service.foobar2.self_link
  host_rule {
    hosts        = ["mysite2.com", "myothersite2.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_region_backend_service.foobar2.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_region_backend_service.foobar2.self_link
    }
  }
  test {
    host    = "mysite2.com"
    path    = "/*"
    service = google_compute_region_backend_service.foobar2.self_link
  }
}

resource "google_compute_region_ssl_policy" "foobar" {
  name            = "sslproxy-test-%s"
  description     = "my-description"
  min_tls_version = "TLS_1_2"
  profile         = "MODERN"
  region          = "us-central1"
}

resource "google_compute_region_ssl_certificate" "foobar1" {
  name        = "httpsproxy-test-cert1-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_region_ssl_certificate" "foobar2" {
  name        = "httpsproxy-test-cert2-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}
`, id, id, id, id, id, id, id, id, id, id)
}

func TestAccComputeRegionTargetHttpsProxy_addSslPolicy_withForwardingRule(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"resource_suffix": acctest.RandString(t, 10),
		"project_id":      envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetHttpsProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionTargetHttpsProxy_withForwardingRule(context),
			},
			{
				ResourceName:      "google_compute_region_target_https_proxy.default-https",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionTargetHttpsProxy_withForwardingRule_withSslPolicy(context),
			},
			{
				ResourceName:      "google_compute_region_target_https_proxy.default-https",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRegionTargetHttpsProxy_withForwardingRule(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_forwarding_rule" "default-https" {
  project               = "%{project_id}"
  region                = "us-central1"
  name                  = "https-frwd-rule-%{resource_suffix}"
  load_balancing_scheme = "INTERNAL_MANAGED"
  target                = google_compute_region_target_https_proxy.default-https.self_link
  network               = google_compute_network.ilb_network.name
  subnetwork            = google_compute_subnetwork.ilb_subnet.name
  ip_address            = google_compute_address.consumer_address.id
  ip_protocol           = "TCP"
  port_range            = "443"
  allow_global_access   = "true"
  depends_on            = [google_compute_subnetwork.ilb_subnet2]
}

resource "google_compute_region_backend_service" "default" {
  project               = "%{project_id}"
  region                = "us-central1"
  name                  = "backend-service-%{resource_suffix}"
  protocol              = "HTTPS"
  port_name             = "https-server"
  load_balancing_scheme = "INTERNAL_MANAGED"
  session_affinity      = "HTTP_COOKIE"
  health_checks         = [google_compute_region_health_check.default.self_link]
  locality_lb_policy    = "RING_HASH"

  # websocket handling: https://stackoverflow.com/questions/63822612/websocket-connection-being-closed-on-google-compute-engine
  timeout_sec = 600

  consistent_hash {
    http_cookie {
      ttl {
        # 24hr cookie ttl
        seconds = 86400
        nanos   = null
      }
      name = "X-CLIENT-SESSION"
      path = null
    }
    http_header_name  = null
    minimum_ring_size = 1024
  }

  log_config {
    enable      = true
    sample_rate = 1.0
  }
}

resource "google_compute_region_health_check" "default" {
  project             = "%{project_id}"
  region              = "us-central1"
  name                = "hc-%{resource_suffix}"
  timeout_sec         = 5
  check_interval_sec  = 30
  healthy_threshold   = 3
  unhealthy_threshold = 3

  https_health_check {
    port         = 443
    request_path = "/health"
  }
}

resource "google_compute_region_target_https_proxy" "default-https" {
  project          = "%{project_id}"
  region           = "us-central1"
  name             = "https-proxy-%{resource_suffix}"
  url_map          = google_compute_region_url_map.default-https.self_link
  ssl_certificates = [google_compute_region_ssl_certificate.foobar0.self_link]
}

resource "google_compute_region_url_map" "default-https" {
  project         = "%{project_id}"
  region          = "us-central1"
  name            = "lb-%{resource_suffix}"
  default_service = google_compute_region_backend_service.default.id
}

resource "google_compute_region_ssl_certificate" "foobar0" {
  name        = "httpsproxy-test-cert0-%{resource_suffix}"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l4-ilb-network-%{resource_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l4-ilb-subnet-%{resource_suffix}"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = google_compute_network.ilb_network.id
}

resource "google_compute_subnetwork" "ilb_subnet2" {
  name          = "tf-test-l4-ilb-subnet2-%{resource_suffix}"
	ip_cidr_range = "10.142.0.0/20"
  region        = "us-central1"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.ilb_network.id
}

resource "google_compute_address" "consumer_address" {
  name         = "tf-test-website-ip-%{resource_suffix}-1"
  region       = "us-central1"
  subnetwork   = google_compute_subnetwork.ilb_subnet.id
  address_type = "INTERNAL"
}
`, context)
}

func testAccComputeRegionTargetHttpsProxy_withForwardingRule_withSslPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_forwarding_rule" "default-https" {
  project               = "%{project_id}"
  region                = "us-central1"
  name                  = "https-frwd-rule-%{resource_suffix}"
  load_balancing_scheme = "INTERNAL_MANAGED"
  target                = google_compute_region_target_https_proxy.default-https.self_link
  network               = google_compute_network.ilb_network.name
  subnetwork            = google_compute_subnetwork.ilb_subnet.name
  ip_address            = google_compute_address.consumer_address.id
  ip_protocol           = "TCP"
  port_range            = "443"
  allow_global_access   = "true"
  depends_on            = [google_compute_subnetwork.ilb_subnet2]
}

resource "google_compute_region_backend_service" "default" {
  project               = "%{project_id}"
  region                = "us-central1"
  name                  = "backend-service-%{resource_suffix}"
  protocol              = "HTTPS"
  port_name             = "https-server"
  load_balancing_scheme = "INTERNAL_MANAGED"
  session_affinity      = "HTTP_COOKIE"
  health_checks         = [google_compute_region_health_check.default.self_link]
  locality_lb_policy    = "RING_HASH"

  # websocket handling: https://stackoverflow.com/questions/63822612/websocket-connection-being-closed-on-google-compute-engine
  timeout_sec = 600

  consistent_hash {
    http_cookie {
      ttl {
        # 24hr cookie ttl
        seconds = 86400
        nanos   = null
      }
      name = "X-CLIENT-SESSION"
      path = null
    }
    http_header_name  = null
    minimum_ring_size = 1024
  }

  log_config {
    enable      = true
    sample_rate = 1.0
  }
}

resource "google_compute_region_health_check" "default" {
  project             = "%{project_id}"
  region              = "us-central1"
  name                = "hc-%{resource_suffix}"
  timeout_sec         = 5
  check_interval_sec  = 30
  healthy_threshold   = 3
  unhealthy_threshold = 3

  https_health_check {
    port         = 443
    request_path = "/health"
  }
}

resource "google_compute_region_target_https_proxy" "default-https" {
  project          = "%{project_id}"
  region           = "us-central1"
  name             = "https-proxy-%{resource_suffix}"
  url_map          = google_compute_region_url_map.default-https.self_link
  ssl_certificates = [google_compute_region_ssl_certificate.foobar0.self_link]
  ssl_policy       = google_compute_region_ssl_policy.default.id
}

resource "google_compute_region_url_map" "default-https" {
  project         = "%{project_id}"
  region          = "us-central1"
  name            = "lb-%{resource_suffix}"
  default_service = google_compute_region_backend_service.default.id
}

resource "google_compute_region_ssl_policy" "default" {
  project = "%{project_id}"
  region  = "us-central1"
  name    = "ssl-policy-%{resource_suffix}"

  profile         = "RESTRICTED"
  min_tls_version = "TLS_1_2"
}

resource "google_compute_region_ssl_certificate" "foobar0" {
  name        = "httpsproxy-test-cert0-%{resource_suffix}"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l4-ilb-network-%{resource_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l4-ilb-subnet-%{resource_suffix}"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = google_compute_network.ilb_network.id
}

resource "google_compute_subnetwork" "ilb_subnet2" {
  name          = "tf-test-l4-ilb-subnet2-%{resource_suffix}"
	ip_cidr_range = "10.142.0.0/20"
  region        = "us-central1"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.ilb_network.id
}

resource "google_compute_address" "consumer_address" {
  name         = "tf-test-website-ip-%{resource_suffix}-1"
  region       = "us-central1"
  subnetwork   = google_compute_subnetwork.ilb_subnet.id
  address_type = "INTERNAL"
}
`, context)
}
