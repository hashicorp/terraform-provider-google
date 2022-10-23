package google

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestValidateRecordNameTrailingDot(t *testing.T) {
	cases := []StringValidationTestCase{
		// No errors
		{TestName: "trailing dot", Value: "test-record.hashicorptest.com."},

		// With errors
		{TestName: "empty string", Value: "", ExpectError: true},
		{TestName: "no trailing dot", Value: "test-record.hashicorptest.com", ExpectError: true},
	}

	es := testStringValidationCases(cases, validateRecordNameTrailingDot)
	if len(es) > 0 {
		t.Errorf("Failed to validate DNS Record name with value: %v", es)
	}
}

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
		shouldSuppress := rrdatasListDiffSuppress(tc.Old, tc.New, parseFunc, nil)
		if shouldSuppress != tc.ShouldSuppress {
			t.Errorf("%s: expected %t", tn, tc.ShouldSuppress)
		}
	}
}

func TestAccDNSRecordSet_basic(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
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
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_Update(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.11", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.11", 600),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_changeType(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_bigChange(zoneName, 600),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./CNAME", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_nestedNS(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", randString(t, 10))
	recordSetName := fmt.Sprintf("\"nested.%s.hashicorptest.com.\"", zoneName)
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
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

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", randString(t, 10))
	recordSetName := "google_dns_managed_zone.parent-zone.dns_name"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_NS(zoneName, recordSetName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s.hashicorptest.com./NS", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_quotedTXT(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-txt-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
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

	zoneName := fmt.Sprintf("dnszone-test-txt-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_uppercaseMX(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./MX", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_routingPolicy(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-network-%s", randString(t, 10))
	backendName := fmt.Sprintf("tf-test-backend-%s", randString(t, 10))
	forwardingRuleName := fmt.Sprintf("tf-test-forwarding-rule-%s", randString(t, 10))
	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_routingPolicyWRR(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyGEO(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyPrimaryBackup(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_changeRouting(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-network-%s", randString(t, 10))
	backendName := fmt.Sprintf("tf-test-backend-%s", randString(t, 10))
	forwardingRuleName := fmt.Sprintf("tf-test-forwarding-rule-%s", randString(t, 10))
	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsRecordSet_routingPolicyGEO(networkName, backendName, forwardingRuleName, zoneName, 300),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/test-record.%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Tracks fix for https://github.com/hashicorp/terraform-provider-google/issues/12043
func TestAccDNSRecordSet_interpolated(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
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

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
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
  network               = google_compute_network.default.name
}

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
  network               = google_compute_network.default.name
}

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
  network               = google_compute_network.default.name
}

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
