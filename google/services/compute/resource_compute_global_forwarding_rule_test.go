// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
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
					resource.TestCheckResourceAttrSet(
						"google_compute_global_forwarding_rule.forwarding_rule", "forwarding_rule_id")),
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

func TestAccComputeGlobalForwardingRule_labels(t *testing.T) {
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
				Config: testAccComputeGlobalForwardingRule_labels(fr, proxy, backend, hc, urlmap),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target", "labels", "terraform_labels"},
			},
			{
				Config: testAccComputeGlobalForwardingRule_labelsUpdated(fr, proxy, backend, hc, urlmap),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccComputeGlobalForwardingRule_allApisLabels(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("frtest%s", acctest.RandString(t, 10))
	address := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalForwardingRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalForwardingRule_allApisLabels(fr, address),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target", "labels", "terraform_labels"},
			},
			{
				Config: testAccComputeGlobalForwardingRule_allApisLabelsUpdated(fr, address),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccComputeGlobalForwardingRule_vpcscLabels(t *testing.T) {
	t.Parallel()

	fr := fmt.Sprintf("frtest%s", acctest.RandString(t, 10))
	address := fmt.Sprintf("forwardrule-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeGlobalForwardingRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalForwardingRule_vpcscLabels(fr, address),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "target", "labels", "terraform_labels"},
			},
			{
				Config: testAccComputeGlobalForwardingRule_vpcscLabelsUpdated(fr, address),
			},
			{
				ResourceName:            "google_compute_global_forwarding_rule.forwarding_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_range", "labels", "terraform_labels"},
			},
		},
	})
}

func TestUnitComputeGlobalForwardingRule_PortRangeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"different single values": {
			Old:                "80-80",
			New:                "443",
			ExpectDiffSuppress: false,
		},
		"different ranges": {
			Old:                "80-80",
			New:                "443-444",
			ExpectDiffSuppress: false,
		},
		"same single values": {
			Old:                "80-80",
			New:                "80",
			ExpectDiffSuppress: true,
		},
		"same ranges": {
			Old:                "80-80",
			New:                "80-80",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if compute.PortRangeDiffSuppress("ports", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestUnitComputeGlobalForwardingRule_InternalIpDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"suppress - same long and short ipv6 IPs without netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8000::",
			ExpectDiffSuppress: true,
		},
		"suppress - long and short ipv6 IPs with netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "2600:1900:4020:31cd:8000::/96",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP with netmask and short ipv6 IP without netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "2600:1900:4020:31cd:8000::",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP without netmask and short ipv6 IP with netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8000::/96",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP with netmask and reference": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: true,
		},
		"suppress - long ipv6 IP without netmask and reference": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: true,
		},
		"do not suppress - ipv6 IPs different netmask": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0/96",
			New:                "2600:1900:4020:31cd:8000:0:0:0/95",
			ExpectDiffSuppress: false,
		},
		"do not suppress - reference and ipv6 IP with netmask": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "2600:1900:4020:31cd:8000:0:0:0/96",
			ExpectDiffSuppress: false,
		},
		"do not suppress - ipv6 IPs - 1": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8001::",
			ExpectDiffSuppress: false,
		},
		"do not suppress - ipv6 IPs - 2": {
			Old:                "2600:1900:4020:31cd:8000:0:0:0",
			New:                "2600:1900:4020:31cd:8000:0:0:8000",
			ExpectDiffSuppress: false,
		},
		"suppress - ipv4 IPs": {
			Old:                "1.2.3.4",
			New:                "1.2.3.4",
			ExpectDiffSuppress: true,
		},
		"suppress - ipv4 IP without netmask and ipv4 IP with netmask": {
			Old:                "1.2.3.4",
			New:                "1.2.3.4/24",
			ExpectDiffSuppress: true,
		},
		"suppress - ipv4 IP without netmask and reference": {
			Old:                "1.2.3.4",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: true,
		},
		"do not suppress - reference and ipv4 IP without netmask": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "1.2.3.4",
			ExpectDiffSuppress: false,
		},
		"do not suppress - different ipv4 IPs": {
			Old:                "1.2.3.4",
			New:                "1.2.3.5",
			ExpectDiffSuppress: false,
		},
		"do not suppress - ipv4 IPs different netmask": {
			Old:                "1.2.3.4/24",
			New:                "1.2.3.5/25",
			ExpectDiffSuppress: false,
		},
		"do not suppress - different references": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "projects/project_id/regions/region/addresses/address-name-1",
			ExpectDiffSuppress: false,
		},
		"do not suppress - same references": {
			Old:                "projects/project_id/regions/region/addresses/address-name",
			New:                "projects/project_id/regions/region/addresses/address-name",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if compute.InternalIpDiffSuppress("ipv4/v6_compare", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
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

func testAccComputeGlobalForwardingRule_labels(fr, proxy, backend, hc, urlmap string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name       = "%s"
  target     = google_compute_target_http_proxy.proxy.self_link
  port_range = "80"

  labels = {
    my-label          = "a-value"
    a-different-label = "my-second-label-value"
  }
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

func testAccComputeGlobalForwardingRule_labelsUpdated(fr, proxy, backend, hc, urlmap string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name       = "%s"
  target     = google_compute_target_http_proxy.proxy.self_link
  port_range = "80"

  labels = {
    my-label          = "a-new-value"
    a-different-label = "my-third-label-value"
  }
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

func testAccComputeGlobalForwardingRule_allApisLabels(fr, address string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name                  = "%s"
  network               = "default"
  target                = "all-apis"
  ip_address            = google_compute_global_address.default.id
  load_balancing_scheme = ""
  labels = {
    my-label          = "a-value"
    a-different-label = "my-second-label-value"
  }
}

resource "google_compute_global_address" "default" {
  name          = "%s"
  address_type  = "INTERNAL"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  network       = "default"
  address       = "100.100.100.105"
}

`, fr, address)
}

func testAccComputeGlobalForwardingRule_allApisLabelsUpdated(fr, address string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name                  = "%s"
  network               = "default"
  target                = "all-apis"
  ip_address            = google_compute_global_address.default.id
  load_balancing_scheme = ""
  labels = {
    my-label          = "a-value"
    a-different-label = "my-third-label-value"
  }
}

resource "google_compute_global_address" "default" {
  name          = "%s"
  address_type  = "INTERNAL"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  network       = "default"
  address       = "100.100.100.105"
}

`, fr, address)
}

func testAccComputeGlobalForwardingRule_vpcscLabels(fr, address string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name                  = "%s"
  network               = "default"
  target                = "vpc-sc"
  ip_address            = google_compute_global_address.default.id
  load_balancing_scheme = ""
  labels = {
    my-label          = "a-value"
    a-different-label = "my-second-label-value"
  }
}

resource "google_compute_global_address" "default" {
  name          = "%s"
  address_type  = "INTERNAL"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  network       = "default"
  address       = "100.100.100.106"
}

`, fr, address)
}

func testAccComputeGlobalForwardingRule_vpcscLabelsUpdated(fr, address string) string {
	return fmt.Sprintf(`
resource "google_compute_global_forwarding_rule" "forwarding_rule" {
  name                  = "%s"
  network               = "default"
  target                = "vpc-sc"
  ip_address            = google_compute_global_address.default.id
  load_balancing_scheme = ""
  labels = {
    my-label          = "a-value"
    a-different-label = "my-third-label-value"
  }
}

resource "google_compute_global_address" "default" {
  name          = "%s"
  address_type  = "INTERNAL"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  network       = "default"
  address       = "100.100.100.106"
}

`, fr, address)
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
