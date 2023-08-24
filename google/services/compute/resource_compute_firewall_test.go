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

func TestAccComputeFirewall_update(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_update(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_localRanges(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_localRanges(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_localRangesUpdate(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_localRanges(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_priority(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_priority(networkName, firewallName, 1001),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_noSource(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeFirewall_noSource(networkName, firewallName),
				ExpectError: regexp.MustCompile("one of source_tags, source_ranges, or source_service_accounts must be defined"),
			},
		},
	})
}

func TestAccComputeFirewall_denied(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_denied(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_egress(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_egress(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_serviceAccounts(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	sourceSa := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	targetSa := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_serviceAccounts(sourceSa, targetSa, networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_disabled(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_disabled(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_enableLogging(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_enableLogging(networkName, firewallName, ""),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_enableLogging(networkName, firewallName, "INCLUDE_ALL_METADATA"),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_enableLogging(networkName, firewallName, "EXCLUDE_ALL_METADATA"),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFirewall_enableLogging(networkName, firewallName, ""),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_moduleOutput(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))
	firewallName := fmt.Sprintf("tf-test-firewall-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeFirewallDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_moduleOutput(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeFirewall_basic(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]

  allow {
    protocol = "icmp"
  }
}
`, network, firewall)
}

func testAccComputeFirewall_localRanges(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]

  source_ranges      = ["10.0.0.0/8"]
  destination_ranges = ["192.168.1.0/24"]

  allow {
    protocol = "icmp"
  }
}
`, network, firewall)
}

func testAccComputeFirewall_localRangesUpdate(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]

  source_ranges      = ["192.168.1.0/24"]
  destination_ranges = ["10.0.0.0/8"]

  allow {
    protocol = "icmp"
  }
}
`, network, firewall)
}

func testAccComputeFirewall_update(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.self_link
  source_tags = ["foo"]
  target_tags = ["bar"]

  allow {
    protocol = "tcp"
    ports    = ["80-255"]
  }
}
`, network, firewall)
}

func testAccComputeFirewall_priority(network, firewall string, priority int) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]

  allow {
    protocol = "icmp"
  }
  priority = %d
}
`, network, firewall, priority)
}

func testAccComputeFirewall_noSource(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name

  allow {
    protocol = "tcp"
    ports    = [22]
  }
}
`, network, firewall)
}

func testAccComputeFirewall_denied(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]

  deny {
    protocol = "tcp"
    ports    = [22]
  }
}
`, network, firewall)
}

func testAccComputeFirewall_egress(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  direction   = "EGRESS"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name

  allow {
    protocol = "tcp"
    ports    = [22]
  }
}
`, network, firewall)
}

func testAccComputeFirewall_serviceAccounts(sourceSa, targetSa, network, firewall string) string {
	return fmt.Sprintf(`
resource "google_service_account" "source" {
  account_id = "%s"
}

resource "google_service_account" "target" {
  account_id = "%s"
}

resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name

  allow {
    protocol = "icmp"
  }

  source_service_accounts = [google_service_account.source.email]
  target_service_accounts = [google_service_account.target.email]
}
`, sourceSa, targetSa, network, firewall)
}

func testAccComputeFirewall_disabled(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  source_tags = ["foo"]

  allow {
    protocol = "icmp"
  }

  disabled = true
}
`, network, firewall)
}

func testAccComputeFirewall_enableLogging(network, firewall, logging string) string {
	enableLoggingCfg := ""
	if logging != "" {
		enableLoggingCfg = fmt.Sprintf(`log_config {
		  metadata = "%s"
		}
		`, logging)
	}
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "foobar" {
  name = "%s"
  description = "Resource created for Terraform acceptance testing"
  network = google_compute_network.foobar.name
  source_tags = ["foo"]

  allow {
    protocol = "icmp"
  }

  %s
}
`, network, firewall, enableLoggingCfg)
}

func testAccComputeFirewall_moduleOutput(network, firewall string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "%s-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.foobar.name
}

resource "google_compute_address" "foobar" {
	name         = "%s-address"
	subnetwork   = google_compute_subnetwork.foobar.id
	address_type = "INTERNAL"
	region       = "us-central1"
  }

resource "google_compute_firewall" "foobar" {
  name        = "%s"
  description = "Resource created for Terraform acceptance testing"
  network     = google_compute_network.foobar.name
  direction   = "INGRESS"

  source_ranges = ["${google_compute_address.foobar.address}/32"]
  target_tags   = ["foo"]

  allow {
    protocol = "tcp"
  }
}
`, network, network, network, firewall)
}
