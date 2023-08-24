// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgdns "github.com/hashicorp/terraform-provider-google/google/services/dns"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/dns/v1"
)

func TestAccDNSManagedZone_update(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description1", map[string]string{"foo": "bar", "ping": "pong"}),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description2", map[string]string{"foo": "bar"}),
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

	zoneSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
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

func TestAccDNSManagedZone_dnssec_update(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_dnssec_on(zoneSuffix),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
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

func TestAccDNSManagedZone_dnssec_empty(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_dnssec_empty(zoneSuffix),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSManagedZone_privateForwardingUpdate(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_privateForwardingUpdate(zoneSuffix, "172.16.1.10", "172.16.1.20", "default", "private"),
			},
			{
				ResourceName:      "google_dns_managed_zone.private",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsManagedZone_privateForwardingUpdate(zoneSuffix, "172.16.1.10", "192.168.1.1", "private", "default"),
			},
			{
				ResourceName:      "google_dns_managed_zone.private",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSManagedZone_cloudLoggingConfigUpdate(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZone_cloudLoggingConfig_basic(zoneSuffix),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsManagedZone_cloudLoggingConfig_update(zoneSuffix, true),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDnsManagedZone_cloudLoggingConfig_update(zoneSuffix, false),
			},
			{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSManagedZone_forceDestroy(t *testing.T) {
	//t.Parallel()

	zoneSuffix := acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSManagedZone_forceDestroy(zoneSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedZoneCreateRRs(t, zoneSuffix, project),
				),
			},
		},
	})
}

func testAccCheckManagedZoneCreateRRs(t *testing.T, zoneSuffix string, project string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		zone := fmt.Sprintf("mzone-test-%s", zoneSuffix)
		// Build the change
		chg := &dns.Change{
			Additions: []*dns.ResourceRecordSet{
				{
					Name:    fmt.Sprintf("cname.%s.hashicorptest.com.", zoneSuffix),
					Type:    "CNAME",
					Ttl:     300,
					Rrdatas: []string{"foo.example.com."},
				},
				{
					Name:    fmt.Sprintf("a.%s.hashicorptest.com.", zoneSuffix),
					Type:    "A",
					Ttl:     300,
					Rrdatas: []string{"1.1.1.1"},
				},
				{
					Name:    fmt.Sprintf("nested.%s.hashicorptest.com.", zoneSuffix),
					Type:    "NS",
					Ttl:     300,
					Rrdatas: []string{"ns.hashicorp.services.", "ns2.hashicorp.services."},
				},
			},
		}

		chg, err := config.NewDnsClient(config.UserAgent).Changes.Create(project, zone, chg).Do()
		if err != nil {
			return fmt.Errorf("Error creating DNS RecordSet: %s", err)
		}

		w := &tpgdns.DnsChangeWaiter{
			Service:     config.NewDnsClient(config.UserAgent),
			Change:      chg,
			Project:     project,
			ManagedZone: zone,
		}
		_, err = w.Conf().WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for Google DNS change: %s", err)
		}

		return nil
	}
}

func testAccDNSManagedZone_forceDestroy(suffix string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name        = "mzone-test-%s"
  dns_name    = "%s.hashicorptest.com."
  labels = {
    foo = "bar"
  }
  force_destroy = true
  visibility = "public"
}
`, suffix, suffix)
}

func testAccDnsManagedZone_basic(suffix, description string, labels map[string]string) string {
	labelsRep := ""
	for k, v := range labels {
		labelsRep += fmt.Sprintf("%s = %q, ", k, v)
	}
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name        = "mzone-test-%s"
  dns_name    = "tf-acctest-%s.hashicorptest.com."
  description = "%s"
  labels 	  = {%s}
  visibility = "public"
}
`, suffix, suffix, description, labelsRep)
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

    non_existence = "nsec"
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

    non_existence = "nsec3"
  }
}
`, suffix, suffix)
}

func testAccDnsManagedZone_dnssec_empty(suffix string) string {
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
    gke_clusters {
		gke_cluster_name = google_container_cluster.cluster-1.id
	}
  }
}

resource "google_compute_network" "network-1" {
  name                    = "tf-test-net-1-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-2" {
  name                    = "tf-test-net-2-%s"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network-3" {
  name                    = "tf-test-network-3-%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork-1" {
  name                     = google_compute_network.network-1.name
  network                  = google_compute_network.network-1.name
  ip_cidr_range            = "10.0.36.0/24"
  region                   = "us-central1"
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pod"
    ip_cidr_range = "10.0.0.0/19"
  }

  secondary_ip_range {
    range_name    = "svc"
    ip_cidr_range = "10.0.32.0/22"
  }
}

resource "google_container_cluster" "cluster-1" {
  name               = "tf-test-cluster-1-%s"
  location           = "us-central1-c"
  initial_node_count = 1

  networking_mode = "VPC_NATIVE"
  default_snat_status {
    disabled = true
  }
  network    = google_compute_network.network-1.name
  subnetwork = google_compute_subnetwork.subnetwork-1.name

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "10.42.0.0/28"
    master_global_access_config {
      enabled = true
	}
  }
  master_authorized_networks_config {
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.subnetwork-1.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.subnetwork-1.secondary_ip_range[1].range_name
  }
}
`, suffix, first_network, second_network, suffix, suffix, suffix, suffix)
}

func testAccDnsManagedZone_privateForwardingUpdate(suffix, first_nameserver, second_nameserver, first_forwarding_path, second_forwarding_path string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "private" {
  name        = "private-zone-%s"
  dns_name    = "private.example.com."
  description = "Example private DNS zone"
  visibility  = "private"
  private_visibility_config {
    networks {
      network_url = google_compute_network.network-1.self_link
    }
  }

  forwarding_config {
    target_name_servers {
      ipv4_address = "%s"
      forwarding_path = "%s"
    }
    target_name_servers {
      ipv4_address = "%s"
      forwarding_path = "%s"
    }
  }
}

resource "google_compute_network" "network-1" {
  name                    = "tf-test-net-1-%s"
  auto_create_subnetworks = false
}
`, suffix, first_nameserver, first_forwarding_path, second_nameserver, second_forwarding_path, suffix)
}

func testAccDnsManagedZone_cloudLoggingConfig_basic(suffix string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name        = "mzone-test-%s"
  dns_name    = "tf-acctest-%s.hashicorptest.com."
  description = "Example DNS zone"
  labels = {
    foo = "bar"
  }
}
`, suffix, suffix)
}

func testAccDnsManagedZone_cloudLoggingConfig_update(suffix string, enableCloudLogging bool) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name        = "mzone-test-%s"
  dns_name    = "tf-acctest-%s.hashicorptest.com."
  description = "Example DNS zone"
  labels = {
    foo = "bar"
  }

  cloud_logging_config {
    enable_logging = %t
  }
}
`, suffix, suffix, enableCloudLogging)
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
		Config               *transport_tpg.Config
		ExpectedSchemaValues map[string]interface{}
		ExpectError          bool
	}{
		"full self_link": {
			IdRegexes: zoneRegexes,
			ImportId:  "https://dns.googleapis.com/dns/v1/projects/my-project/managedZones/my-zone",
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
			Config: &transport_tpg.Config{
				Project: "default-project",
			},
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"name":    "my-zone",
			},
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDataMock{
			FieldsInSchema: make(map[string]interface{}),
		}
		d.SetId(tc.ImportId)
		config := tc.Config
		if config == nil {
			config = &transport_tpg.Config{}
		}
		//
		if err := tpgresource.ParseImportId(tc.IdRegexes, d, config); err == nil {
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

	zoneSuffix := acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
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
