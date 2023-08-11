// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeGlobalForwardingRule_updateTarget(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	proxy := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	proxyUpdated := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	backend := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	hc := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	urlmap := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalForwardingRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalForwardingRule_httpProxy(fr, "proxy", proxy, proxyUpdated, backend, hc, urlmap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"google_compute_global_forwarding_rule.forwarding_rule", "target", regexp.MustCompile(proxy+"$")),
				),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target"},
			},
			{
				Config: testAccComputeGlobalForwardingRule_httpProxy(fr, "proxy2", proxy, proxyUpdated, backend, hc, urlmap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"google_compute_global_forwarding_rule.forwarding_rule", "target", regexp.MustCompile(proxyUpdated+"$")),
				),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target"},
			},
		},
	})
}

func TestAccComputeGlobalForwardingRule_ipv6(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	proxy := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	backend := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	hc := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))
	urlmap := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalForwardingRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalForwardingRule_ipv6(fr, proxy, backend, hc, urlmap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_compute_global_forwarding_rule.forwarding_rule", "ip_version", "IPV6"),
				),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target"},
			},
		},
	})
}

func testAccComputeGlobalForwardingRule_httpProxy(fr, targetProxy, proxy, proxy2, backend, hc, urlmap string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "TCP"
  name        = "%s"
  port_range  = "80"
  target      = google_compute_target_http_proxy.%s.self_link
}

resource "google_compute_target_http_proxy" "proxy" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_url_map.url_map.self_link
}

resource "google_compute_target_http_proxy" "proxy2" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_url_map.url_map.self_link
}

resource "google_compute_backend_service" "backend" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "url_map" {
  name            = "%s"
  default_service = google_compute_backend_service.backend.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.backend.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.backend.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.backend.self_link
  }
}
`, fr, targetProxy, proxy, proxy2, backend, hc, urlmap)
}

func testAccComputeGlobalForwardingRule_ipv6(fr, proxy, backend, hc, urlmap string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "TCP"
  name        = "%s"
  port_range  = "80"
  target      = google_compute_target_http_proxy.proxy.self_link
  ip_version  = "IPV6"
}

resource "google_compute_target_http_proxy" "proxy" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  url_map     = google_compute_url_map.urlmap.self_link
}

resource "google_compute_backend_service" "backend" {
  name          = "%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "urlmap" {
  name            = "%s"
  default_service = google_compute_backend_service.backend.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.backend.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.backend.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.backend.self_link
  }
}
`, fr, proxy, backend, hc, urlmap)
}
