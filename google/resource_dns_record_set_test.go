package google

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
