package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDNSManagedZone_update(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description1"),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description2"),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSManagedZone_privateUpdate(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_privateUpdate(zoneSuffix, "network-1", "network-2"),
			},
			{
				ResourceName:      "google_dns_managed_zone.private",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsManagedZone_privateUpdate(zoneSuffix, "network-2", "network-3"),
			},
			{
				ResourceName:      "google_dns_managed_zone.private",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSManagedZone_dnssec_on(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_dnssec_on(zoneSuffix),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSManagedZone_dnssec_off(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_dnssec_off(zoneSuffix),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsManagedZone_basic(suffix, description string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name        = "mzone-test-%s"
  dns_name    = "tf-acctest-%s.hashicorptest.com."
  description = "%s"
  labels = {
    foo = "bar"
  }

  visibility = "public"
}
`, suffix, suffix, description)
}

func testAccDnsManagedZone_dnssec_on(suffix string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name     = "mzone-test-%s"
  dns_name = "tf-acctest-%s.hashicorptest.com."

  dnssec_config {
    state = "on"
    default_key_specs {
      algorithm  = "rsasha256"
      key_length = "2048"
      key_type   = "zoneSigning"
    }
    default_key_specs {
      algorithm  = "rsasha256"
      key_length = "2048"
      key_type   = "keySigning"
    }
  }
}
`, suffix, suffix)
}

func testAccDnsManagedZone_dnssec_off(suffix string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name     = "mzone-test-%s"
  dns_name = "tf-acctest-%s.hashicorptest.com."

  dnssec_config {
    state = "off"
  }
}
`, suffix, suffix)
}

func testAccDnsManagedZone_privateUpdate(suffix, first_network, second_network string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "private" {
  name        = "private-zone-%s"
  dns_name    = "private.example.com."
  description = "Example private DNS zone"
  visibility  = "private"
  private_visibility_config {
    networks {
      network_url = google_compute_network.%s.self_link
    }
    networks {
      network_url = google_compute_network.%s.self_link
    }
  }
}

resource "google_compute_network" "network-1" {
  name                    = "network-1-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-2" {
  name                    = "network-2-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-3" {
  name                    = "network-3-%s"
  auto_create_subnetworks = false
}
`, suffix, first_network, second_network, suffix, suffix, suffix)
}

func TestDnsManagedZoneImport_parseImportId(t *testing.T) {
	zoneRegexes := []string{
		"projects/(?P<project>[^/]+)/managedZones/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/managedZones/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}

	cases := map[string]struct {
		ImportId             string
		IdRegexes            []string
		Config               *Config
		ExpectedSchemaValues map[string]interface{}
		ExpectError          bool
	}{
		"full self_link": {
			IdRegexes: zoneRegexes,
			ImportId:  "https://www.googleapis.com/dns/v1/projects/my-project/managedZones/my-zone",
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"name":    "my-zone",
			},
		},
		"relative self_link": {
			IdRegexes: zoneRegexes,
			ImportId:  "projects/my-project/managedZones/my-zone",
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"name":    "my-zone",
			},
		},
		"short id": {
			IdRegexes: zoneRegexes,
			ImportId:  "my-project/managedZones/my-zone",
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"name":    "my-zone",
			},
		},
		"short id with default project and region": {
			IdRegexes: zoneRegexes,
			ImportId:  "my-zone",
			Config: &Config{
				Project: "default-project",
			},
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"name":    "my-zone",
			},
		},
	}

	for tn, tc := range cases {
		d := &ResourceDataMock{
			FieldsInSchema: make(map[string]interface{}),
			id:             tc.ImportId,
		}
		config := tc.Config
		if config == nil {
			config = &Config{}
		}
		//
		if err := parseImportId(tc.IdRegexes, d, config); err == nil {
			for k, expectedValue := range tc.ExpectedSchemaValues {
				if v, ok := d.GetOk(k); ok {
					if v != expectedValue {
						t.Errorf("%s failed; Expected value %q for field %q, got %q", tn, expectedValue, k, v)
					}
				} else {
					t.Errorf("%s failed; Expected a value for field %q", tn, k)
				}
			}
		} else if !tc.ExpectError {
			t.Errorf("%s failed; unexpected error: %s", tn, err)
		}
	}
}

func TestAccDNSManagedZone_importWithProject(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)
	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSManagedZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_basicWithProject(zoneSuffix, "description1", project),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsManagedZone_basicWithProject(suffix, description, project string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name        = "mzone-test-%s"
  dns_name    = "tf-acctest-%s.hashicorptest.com."
  description = "%s"
  project     = "%s"
}
`, suffix, suffix, description, project)
}
