package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestIpv6AddressDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New       string
		ShouldSuppress bool
	}{
		"compact form should suppress diff": {
			Old:            "2a03:b0c0:1:e0::29b:8001",
			New:            "2a03:b0c0:0001:00e0:0000:0000:029b:8001",
			ShouldSuppress: true,
		},
		"different address should not suppress diff": {
			Old:            "2a03:b0c0:1:e00::29b:8001",
			New:            "2a03:b0c0:0001:00e0:0000:0000:029b:8001",
			ShouldSuppress: false,
		},
	}

	for tn, tc := range cases {
		shouldSuppress := ipv6AddressDiffSuppress("", tc.Old, tc.New, nil)
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
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

func TestAccDNSRecordSet_modify(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.10", 300),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
			},
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.11", 300),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
			},
			{
				Config: testAccDnsRecordSet_basic(zoneName, "127.0.0.11", 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
			},
			{
				Config: testAccDnsRecordSet_bigChange(zoneName, 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
			},
		},
	})
}

func TestAccDNSRecordSet_ns(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_ns(zoneName, 300),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
			},
			{
				ResourceName:      "google_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s.hashicorptest.com./NS", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSRecordSet_nestedNS(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSet_nestedNS(zoneName, 300),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsRecordSetExists(
						t, "google_dns_record_set.foobar", zoneName),
				),
			},
		},
	})
}

func testAccCheckDnsRecordSetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			// Deletion of the managed_zone implies everything is gone
			if rs.Type == "google_dns_managed_zone" {
				_, err := config.clientDns.ManagedZones.Get(
					config.Project, rs.Primary.ID).Do()
				if err == nil {
					return fmt.Errorf("DNS ManagedZone still exists")
				}
			}
		}

		return nil
	}
}

func testAccCheckDnsRecordSetExists(t *testing.T, resourceType, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceType]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dnsName := rs.Primary.Attributes["name"]
		dnsType := rs.Primary.Attributes["type"]

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		resp, err := config.clientDns.ResourceRecordSets.List(
			config.Project, resourceName).Name(dnsName).Type(dnsType).Do()
		if err != nil {
			return fmt.Errorf("Error confirming DNS RecordSet existence: %#v", err)
		}
		switch len(resp.Rrsets) {
		case 0:
			// The resource doesn't exist anymore
			return fmt.Errorf("DNS RecordSet not found")
		case 1:
			return nil
		default:
			return fmt.Errorf("Only expected 1 record set, got %d", len(resp.Rrsets))
		}
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

func testAccDnsRecordSet_ns(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "%s.hashicorptest.com."
  type         = "NS"
  rrdatas      = ["ns.hashicorp.services.", "ns2.hashicorp.services."]
  ttl          = %d
}
`, name, name, name, ttl)
}

func testAccDnsRecordSet_nestedNS(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "nested.%s.hashicorptest.com."
  type         = "NS"
  rrdatas      = ["ns.hashicorp.services.", "ns2.hashicorp.services."]
  ttl          = %d
}
`, name, name, name, ttl)
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
