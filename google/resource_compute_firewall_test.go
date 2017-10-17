package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeFirewall_basic(t *testing.T) {
	t.Parallel()

	var firewall compute.Firewall
	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFirewallExists(
						"google_compute_firewall.foobar", &firewall),
					testAccCheckComputeFirewallApiVersion(&firewall),
				),
			},
		},
	})
}

func TestAccComputeFirewall_update(t *testing.T) {
	t.Parallel()

	var firewall compute.Firewall
	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFirewallExists(
						"google_compute_firewall.foobar", &firewall),
					testAccCheckComputeFirewallApiVersion(&firewall),
				),
			},
			resource.TestStep{
				Config: testAccComputeFirewall_update(networkName, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFirewallExists(
						"google_compute_firewall.foobar", &firewall),
					testAccCheckComputeFirewallPorts(
						&firewall, "80-255"),
					testAccCheckComputeFirewallApiVersion(&firewall),
				),
			},
		},
	})
}

func TestAccComputeFirewall_priority(t *testing.T) {
	t.Parallel()

	var firewall computeBeta.Firewall
	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{{
			Config: testAccComputeFirewall_priority(networkName, firewallName, 1001),
			Check: resource.ComposeTestCheckFunc(
				testAccCheckComputeBetaFirewallExists(
					"google_compute_firewall.foobar", &firewall),
				testAccCheckComputeFirewallHasPriority(&firewall, 1001),
				testAccCheckComputeFirewallBetaApiVersion(&firewall),
			),
		}},
	})
}

func TestAccComputeFirewall_noSource(t *testing.T) {
	t.Parallel()

	var firewall compute.Firewall
	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_noSource(networkName, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeFirewallExists(
						"google_compute_firewall.foobar", &firewall),
					testAccCheckComputeFirewallApiVersion(&firewall),
				),
			},
		},
	})
}

func TestAccComputeFirewall_denied(t *testing.T) {
	t.Parallel()

	var firewall computeBeta.Firewall
	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_denied(networkName, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBetaFirewallExists("google_compute_firewall.foobar", &firewall),
					testAccCheckComputeBetaFirewallDenyPorts(&firewall, "22"),
					testAccCheckComputeFirewallBetaApiVersion(&firewall),
				),
			},
		},
	})
}

func TestAccComputeFirewall_egress(t *testing.T) {
	t.Parallel()

	var firewall computeBeta.Firewall
	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_egress(networkName, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBetaFirewallExists("google_compute_firewall.foobar", &firewall),
					testAccCheckComputeBetaFirewallEgress(&firewall),
					testAccCheckComputeFirewallBetaApiVersion(&firewall),
				),
			},
		},
	})
}

func testAccCheckComputeFirewallDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_firewall" {
			continue
		}

		_, err := config.clientCompute.Firewalls.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Firewall still exists")
		}
	}

	return nil
}

func testAccCheckComputeFirewallExists(n string, firewall *compute.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Firewalls.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Firewall not found")
		}

		*firewall = *found

		return nil
	}
}

func testAccCheckComputeFirewallHasPriority(firewall *computeBeta.Firewall, priority int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if firewall.Priority != int64(priority) {
			return fmt.Errorf("Priority for firewall does not match: expected %d, found %d", priority, firewall.Priority)
		}
		return nil
	}
}

func testAccCheckComputeBetaFirewallExists(n string, firewall *computeBeta.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientComputeBeta.Firewalls.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Firewall not found")
		}

		*firewall = *found

		return nil
	}
}

func testAccCheckComputeFirewallPorts(
	firewall *compute.Firewall, ports string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall.Allowed) == 0 {
			return fmt.Errorf("no allowed rules")
		}

		if firewall.Allowed[0].Ports[0] != ports {
			return fmt.Errorf("bad: %#v", firewall.Allowed[0].Ports)
		}

		return nil
	}
}

func testAccCheckComputeBetaFirewallDenyPorts(firewall *computeBeta.Firewall, ports string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(firewall.Denied) == 0 {
			return fmt.Errorf("no denied rules")
		}

		if firewall.Denied[0].Ports[0] != ports {
			return fmt.Errorf("bad: %#v", firewall.Denied[0].Ports)
		}

		return nil
	}
}

func testAccCheckComputeBetaFirewallEgress(firewall *computeBeta.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if firewall.Direction != "EGRESS" {
			return fmt.Errorf("firewall not EGRESS")
		}

		return nil
	}
}

func testAccCheckComputeFirewallBetaApiVersion(firewall *computeBeta.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// The self-link of the network field is used to determine which API was used when fetching
		// the state from the API.
		if !strings.Contains(firewall.Network, "compute/beta") {
			return fmt.Errorf("firewall beta API was not used")
		}

		return nil
	}
}

func testAccCheckComputeFirewallApiVersion(firewall *compute.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// The self-link of the network field is used to determine which API was used when fetching
		// the state from the API.
		if !strings.Contains(firewall.Network, "compute/v1") {
			return fmt.Errorf("firewall beta API was not used")
		}

		return nil
	}
}

func testAccComputeFirewall_basic(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "firewall-test-%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		allow {
			protocol = "icmp"
		}
	}`, network, firewall)
}

func testAccComputeFirewall_update(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "firewall-test-%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.self_link}"
		source_tags = ["foo"]

		allow {
			protocol = "tcp"
			ports = ["80-255"]
		}
	}`, network, firewall)
}

func testAccComputeFirewall_priority(network, firewall string, priority int) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "firewall-test-%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		allow {
			protocol = "icmp"
		}
		priority = %d
	}`, network, firewall, priority)
}

func testAccComputeFirewall_noSource(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "firewall-test-%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"

		allow {
			protocol = "tcp"
			ports    = [22]
		}
	}`, network, firewall)
}

func testAccComputeFirewall_denied(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "firewall-test-%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		deny {
			protocol = "tcp"
			ports    = [22]
		}
	}`, network, firewall)
}

func testAccComputeFirewall_egress(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "firewall-test-%s"
		direction = "EGRESS"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"

		allow {
			protocol = "tcp"
			ports    = [22]
		}
	}`, network, firewall)
}
