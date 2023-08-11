// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
