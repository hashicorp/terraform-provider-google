// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgdns "github.com/hashicorp/terraform-provider-google/google/services/dns"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestIpv6AddressDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New       []string
		ShouldSuppress bool
	}{
		"compact form should suppress diff": {
			Old:            []string{"2a03:b0c0:1:e0::29b:8001"},
			New:            []string{"2a03:b0c0:0001:00e0:0000:0000:029b:8001"},
			ShouldSuppress: true,
		},
		"different address should not suppress diff": {
			Old:            []string{"2a03:b0c0:1:e00::29b:8001"},
			New:            []string{"2a03:b0c0:0001:00e0:0000:0000:029b:8001"},
			ShouldSuppress: false,
		},
		"increase address should not suppress diff": {
			Old:            []string{""},
			New:            []string{"2a03:b0c0:0001:00e0:0000:0000:029b:8001"},
			ShouldSuppress: false,
		},
		"decrease address should not suppress diff": {
			Old:            []string{"2a03:b0c0:1:e00::29b:8001"},
			New:            []string{""},
			ShouldSuppress: false,
		},
		"switch address positions should suppress diff": {
			Old:            []string{"2a03:b0c0:1:e00::28b:8001", "2a03:b0c0:1:e0::29b:8001"},
			New:            []string{"2a03:b0c0:1:e0::29b:8001", "2a03:b0c0:1:e00::28b:8001"},
			ShouldSuppress: true,
		},
	}

	parseFunc := func(x string) string {
		return net.ParseIP(x).String()
	}

	for tn, tc := range cases {
		shouldSuppress := tpgdns.RrdatasListDiffSuppress(tc.Old, tc.New, parseFunc, nil)
		if shouldSuppress != tc.ShouldSuppress {
			t.Errorf("%s: expected %t", tn, tc.ShouldSuppress)
		}
	}
}

func TestAccDNSRecordSet_basic(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/test-record.%s.hashicorptest.com./A", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check both import formats
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_Update(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.11", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.11", 600),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_changeType(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_bigChange(zoneName, 600),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./CNAME", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_nestedNS(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", acctest.RandString(t, 10))
	recordSetName := fmt.Sprintf("\"nested.%s.hashicorptest.com.\"", zoneName)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_NS(zoneName, recordSetName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/nested.%s.hashicorptest.com./NS", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_secondaryNS(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", acctest.RandString(t, 10))
	recordSetName := "google_dns_managed_zone.parent-zone.dns_name"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_NS(zoneName, recordSetName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s.hashicorptest.com./NS", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_quotedTXT(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-txt-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_quotedTXT(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/test-record.%s.hashicorptest.com./TXT", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_uppercaseMX(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-txt-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_uppercaseMX(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./MX", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_routingPolicy(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-network-%s", acctest.RandString(t, 10))
	backendSubnetName := fmt.Sprintf("tf-test-backend-subnet-%s", acctest.RandString(t, 10))
	proxySubnetName := fmt.Sprintf("tf-test-proxy-subnet-%s", acctest.RandString(t, 10))
	httpHealthCheckName := fmt.Sprintf("tf-test-http-health-check-%s", acctest.RandString(t, 10))
	backendName := fmt.Sprintf("tf-test-backend-%s", acctest.RandString(t, 10))
	urlMapName := fmt.Sprintf("tf-test-url-map-%s", acctest.RandString(t, 10))
	httpProxyName := fmt.Sprintf("tf-test-http-proxy-%s", acctest.RandString(t, 10))
	forwardingRuleName := fmt.Sprintf("tf-test-forwarding-rule-%s", acctest.RandString(t, 10))
	zoneName := fmt.Sprintf("dnszone-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_routingPolicyWRR(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyGEO(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyPrimaryBackup(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyRegionalL7PrimaryBackup(networkName, proxySubnetName, httpHealthCheckName, backendName, urlMapName, httpProxyName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyCrossRegionL7PrimaryBackup(networkName, backendSubnetName, proxySubnetName, httpHealthCheckName, backendName, urlMapName, httpProxyName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_changeRouting(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicy(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", envvar.GetTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Tracks fix for https://github.com/hashicorp/terraform-provider-google/issues/12043
func TestAccDNSRecordSet_interpolated(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_interpolated(zoneName),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/test-record.%s.hashicorptest.com./TXT", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDnsRecordSetDestroyProducer(t *testing.T) func(s *terraform.State) error {

	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dns_record_set" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DNSResourceDnsRecordSet still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccDnsRecordSet_basic(zoneName string, addr2 string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  rrdatas      = ["127.0.0.1", "%s"]
  ttl          = %d
}
`, zoneName, zoneName, zoneName, addr2, ttl)
}

func testAccDnsRecordSet_routingPolicy(zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  ttl          = %d
  routing_policy {
    wrr {
      weight  = 0
      rrdatas = ["1.2.3.4", "4.3.2.1"]
    }

    wrr {
      weight  = 0
      rrdatas = ["2.3.4.5", "5.4.3.2"]
    }
  }
}
`, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_NS(name string, recordSetName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = %s
  type         = "NS"
  rrdatas      = ["ns.hashicorp.services.", "ns2.hashicorp.services."]
  ttl          = %d
}
`, name, name, recordSetName, ttl)
}

func testAccDnsRecordSet_bigChange(zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "CNAME"
  rrdatas      = ["www.terraform.io."]
  ttl          = %d
}
`, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_quotedTXT(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "TXT"
  rrdatas      = ["test", "\"quoted test\""]
  ttl          = %d
}
`, name, name, name, ttl)
}

func testAccDnsRecordSet_uppercaseMX(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "MX"
  rrdatas = [
    "1 ASPMX.L.GOOGLE.COM.",
    "5 ALT1.ASPMX.L.GOOGLE.COM.",
    "5 ALT2.ASPMX.L.GOOGLE.COM.",
    "10 ASPMX2.GOOGLEMAIL.COM.",
    "10 ASPMX3.GOOGLEMAIL.COM.",
  ]
  ttl = %d
}
`, name, name, name, ttl)
}

func testAccDnsRecordSet_routingPolicyWRR(networkName, backendName, forwardingRuleName, zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name = "%s"
}

resource "google_compute_region_backend_service" "backend" {
  name   = "%s"
  region = "us-central1"
}

resource "google_compute_forwarding_rule" "default" {
  name   = "%s"
  region = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.backend.id
  all_ports             = true
  allow_global_access   = true
  network               = google_compute_network.default.name
}

resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
  visibility = "private"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  ttl          = %d

  routing_policy {
    wrr {
      weight  = 0
      rrdatas = ["1.2.3.4", "4.3.2.1"]
    }

    wrr {
      weight  = 0
      rrdatas = ["2.3.4.5", "5.4.3.2"]
    }

    wrr {
      weight = 1.0

      health_checked_targets {
        internal_load_balancers {
          load_balancer_type = "regionalL4ilb"
          ip_address         = google_compute_forwarding_rule.default.ip_address
          port               = "80"
          ip_protocol        = "tcp"
          network_url        = google_compute_network.default.id
          project            = google_compute_forwarding_rule.default.project
          region             = google_compute_forwarding_rule.default.region
        }
      }
    }
  }
}
`, networkName, backendName, forwardingRuleName, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_routingPolicyGEO(networkName, backendName, forwardingRuleName, zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name = "%s"
}

resource "google_compute_region_backend_service" "backend" {
  name   = "%s"
  region = "us-central1"
}

resource "google_compute_forwarding_rule" "default" {
  name   = "%s"
  region = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.backend.id
  all_ports             = true
  allow_global_access   = true
  network               = google_compute_network.default.name
}

resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
  visibility = "private"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  ttl          = %d

  routing_policy {
    enable_geo_fencing = true

    geo {
      location = "us-east4"
      rrdatas  = ["1.2.3.4", "4.3.2.1"]
    }

    geo {
      location = "asia-east1"
      rrdatas  = ["2.3.4.5", "5.4.3.2"]
    }

    geo {
      location = "us-central1"

      health_checked_targets {
        internal_load_balancers {
          load_balancer_type = "regionalL4ilb"
          ip_address         = google_compute_forwarding_rule.default.ip_address
          port               = "80"
          ip_protocol        = "tcp"
          network_url        = google_compute_network.default.id
          project            = google_compute_forwarding_rule.default.project
          region             = google_compute_forwarding_rule.default.region
        }
      }
    }
  }
}
`, networkName, backendName, forwardingRuleName, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_routingPolicyPrimaryBackup(networkName, backendName, forwardingRuleName, zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name = "%s"
}

resource "google_compute_region_backend_service" "backend" {
  name   = "%s"
  region = "us-central1"
}

resource "google_compute_forwarding_rule" "default" {
  name   = "%s"
  region = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.backend.id
  all_ports             = true
  allow_global_access   = true
  network               = google_compute_network.default.name
}

resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
  visibility = "private"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  ttl          = %d

  routing_policy {
    primary_backup {
      trickle_ratio                  = 0.1
      enable_geo_fencing_for_backups = true

      primary {
        internal_load_balancers {
          load_balancer_type = "regionalL4ilb"
          ip_address         = google_compute_forwarding_rule.default.ip_address
          port               = "80"
          ip_protocol        = "tcp"
          network_url        = google_compute_network.default.id
          project            = google_compute_forwarding_rule.default.project
          region             = google_compute_forwarding_rule.default.region
        }
      }

      backup_geo {
        location = "us-west1"
        rrdatas  = ["1.2.3.4"]
      }

      backup_geo {
        location = "asia-east1"
        rrdatas  = ["5.6.7.8"]
      }
    }
  }
}
`, networkName, backendName, forwardingRuleName, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_routingPolicyRegionalL7PrimaryBackup(networkName, proxySubnetName, healthCheckName, backendName, urlMapName, httpProxyName, forwardingRuleName, zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name = "%s"
}

resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "%s"
  ip_cidr_range = "10.100.0.0/24"
  region        = "us-central1"
  purpose       = "INTERNAL_HTTPS_LOAD_BALANCER"
  role          = "ACTIVE"
  network       = google_compute_network.default.id
}

resource "google_compute_region_health_check" "health_check" {
  name   = "%s"
  region = "us-central1"

  http_health_check {
    port = 80
  }
}

resource "google_compute_region_backend_service" "backend" {
  name                  = "%s"
  region                = "us-central1"
  load_balancing_scheme = "INTERNAL_MANAGED"
  protocol              = "HTTP"
  health_checks         = [google_compute_region_health_check.health_check.id]
}

resource "google_compute_region_url_map" "url_map" {
  name            = "%s"
  region          = "us-central1"
  default_service = google_compute_region_backend_service.backend.id
}

resource "google_compute_region_target_http_proxy" "http_proxy" {
  name    = "%s"
  region  = "us-central1"
  url_map = google_compute_region_url_map.url_map.id
}

resource "google_compute_forwarding_rule" "default" {
  name                  = "%s"
  region                = "us-central1"
  depends_on            = [google_compute_subnetwork.proxy_subnet]
  load_balancing_scheme = "INTERNAL_MANAGED"
  target                = google_compute_region_target_http_proxy.http_proxy.id
  port_range            = "80"
  allow_global_access   = true
  network               = google_compute_network.default.name
  ip_protocol           = "TCP"
}

resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
  visibility = "private"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  ttl          = %d

  routing_policy {
    primary_backup {
      trickle_ratio                  = 0.1
      enable_geo_fencing_for_backups = true

      primary {
        internal_load_balancers {
          load_balancer_type = "regionalL7ilb"
          ip_address         = google_compute_forwarding_rule.default.ip_address
          port               = "80"
          ip_protocol        = "tcp"
          network_url        = google_compute_network.default.id
          project            = google_compute_forwarding_rule.default.project
          region             = google_compute_forwarding_rule.default.region
        }
      }

      backup_geo {
        location = "us-west1"
        rrdatas  = ["1.2.3.4"]
      }

      backup_geo {
        location = "asia-east1"
        rrdatas  = ["5.6.7.8"]
      }
    }
  }
}
`, networkName, proxySubnetName, healthCheckName, backendName, urlMapName, httpProxyName, forwardingRuleName, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_routingPolicyCrossRegionL7PrimaryBackup(networkName, backendSubnetName, proxySubnetName, healthCheckName, backendName, urlMapName, httpProxyName, forwardingRuleName, zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name = "%s"
}

resource "google_compute_subnetwork" "backend_subnet" {
  name          = "%s"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-central1"
  network       = google_compute_network.default.id
}

resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "%s"
  ip_cidr_range = "10.100.0.0/24"
  region        = "us-central1"
  purpose       = "GLOBAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.default.id
}

resource "google_compute_health_check" "health_check" {
  name   = "%s"

  http_health_check {
    port = 80
  }
}

resource "google_compute_backend_service" "backend" {
  name                  = "%s"
  load_balancing_scheme = "INTERNAL_MANAGED"
  protocol              = "HTTP"
  health_checks         = [google_compute_health_check.health_check.id]
}

resource "google_compute_url_map" "url_map" {
  name            = "%s"
  default_service = google_compute_backend_service.backend.id
}

resource "google_compute_target_http_proxy" "http_proxy" {
  name    = "%s"
  url_map = google_compute_url_map.url_map.id
}

resource "google_compute_global_forwarding_rule" "default" {
  name                  = "%s"
  depends_on            = [google_compute_subnetwork.proxy_subnet]
  load_balancing_scheme = "INTERNAL_MANAGED"
  target                = google_compute_target_http_proxy.http_proxy.id
  port_range            = "80"
  network               = google_compute_network.default.name
  subnetwork            = google_compute_subnetwork.backend_subnet.name
  ip_protocol           = "TCP"
}

resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
  visibility = "private"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "A"
  ttl          = %d

  routing_policy {
    primary_backup {
      trickle_ratio                  = 0.1
      enable_geo_fencing_for_backups = true

      primary {
        internal_load_balancers {
          load_balancer_type = "globalL7ilb"
          ip_address         = google_compute_global_forwarding_rule.default.ip_address
          port               = "80"
          ip_protocol        = "tcp"
          network_url        = google_compute_network.default.id
          project            = google_compute_global_forwarding_rule.default.project
        }
      }

      backup_geo {
        location = "us-west1"
        rrdatas  = ["1.2.3.4"]
      }

      backup_geo {
        location = "asia-east1"
        rrdatas  = ["5.6.7.8"]
      }
    }
  }
}
`, networkName, backendSubnetName, proxySubnetName, healthCheckName, backendName, urlMapName, httpProxyName, forwardingRuleName, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsRecordSet_interpolated(zoneName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "test-record.%s.hashicorptest.com."
  type         = "TXT"
  rrdatas      = ["127.0.0.1", "firebase=${google_dns_managed_zone.parent-zone.id}"]
  ttl          = 10
}
`, zoneName, zoneName, zoneName)
}
