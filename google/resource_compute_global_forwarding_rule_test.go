package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeGlobalForwardingRule_updateTarget(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	proxy := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	proxyUpdated := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	urlmap := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalForwardingRule_httpProxy(fr, "proxy", proxy, proxyUpdated, backend, hc, urlmap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"google_compute_global_forwarding_rule.forwarding_rule", "target", regexp.MustCompile(proxy+"$")),
				),
			},
			{
				ResourceName:      "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeGlobalForwardingRule_httpProxy(fr, "proxy2", proxy, proxyUpdated, backend, hc, urlmap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"google_compute_global_forwarding_rule.forwarding_rule", "target", regexp.MustCompile(proxyUpdated+"$")),
				),
			},
			{
				ResourceName:      "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeGlobalForwardingRule_ipv6(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	proxy := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))
	urlmap := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalForwardingRule_ipv6(fr, proxy, backend, hc, urlmap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_compute_global_forwarding_rule.forwarding_rule", "ip_version", "IPV6"),
				),
			},
			{
				ResourceName:      "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:       true,
				ImportStateVerify: true,
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
