package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeForwardingRule_update(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeForwardingRule_basic(poolName, ruleName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeForwardingRule_update(poolName, ruleName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_singlePort(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeForwardingRule_singlePort(poolName, ruleName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_ip(t *testing.T) {
	t.Parallel()

	addrName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeForwardingRule_ip(addrName, poolName, ruleName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_internalLoadBalancing(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	networkName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName1 := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName2 := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeForwardingRule_internalLoadBalancing(serviceName, checkName, networkName, ruleName1, ruleName2),
			},
			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeForwardingRule_networkTier(t *testing.T) {
	t.Parallel()

	poolName := fmt.Sprintf("tf-%s", acctest.RandString(10))
	ruleName := fmt.Sprintf("tf-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeForwardingRule_networkTier(poolName, ruleName),
			},

			resource.TestStep{
				ResourceName:      "google_compute_forwarding_rule.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeForwardingRule_basic(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "foo-%s"
}
resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = "${google_compute_target_pool.foo-tp.self_link}"
  labels      = {"foo" = "bar"}
}
`, poolName, ruleName)
}

func testAccComputeForwardingRule_update(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foo-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "foo-%s"
}
resource "google_compute_target_pool" "bar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "bar-%s"
}
resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = "${google_compute_target_pool.bar-tp.self_link}"
  labels      = {"baz" = "qux"}
}
`, poolName, poolName, ruleName)
}

func testAccComputeForwardingRule_singlePort(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foobar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "%s"
}
resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80"
  target      = "${google_compute_target_pool.foobar-tp.self_link}"
}
`, poolName, ruleName)
}

func testAccComputeForwardingRule_ip(addrName, poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_address" "foo" {
  name = "%s"
}
resource "google_compute_target_pool" "foobar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "%s"
}
resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_address  = "${google_compute_address.foo.address}"
  ip_protocol = "TCP"
  name        = "%s"
  port_range  = "80-81"
  target      = "${google_compute_target_pool.foobar-tp.self_link}"
}
`, addrName, poolName, ruleName)
}

func testAccComputeForwardingRule_internalLoadBalancing(serviceName, checkName, networkName, ruleName1, ruleName2 string) string {
	return fmt.Sprintf(`
resource "google_compute_region_backend_service" "foobar-bs" {
  name                  = "%s"
  description           = "Resource created for Terraform acceptance testing"
  health_checks         = ["${google_compute_health_check.zero.self_link}"]
  region                = "us-central1"
}
resource "google_compute_health_check" "zero" {
  name               = "%s"
  description        = "Resource created for Terraform acceptance testing"
  check_interval_sec = 1
  timeout_sec        = 1

  tcp_health_check {
    port = "80"
  }
}
resource "google_compute_network" "foobar" {
  name = "%s"
  auto_create_subnetworks = true
}
resource "google_compute_forwarding_rule" "foobar" {
  description           = "Resource created for Terraform acceptance testing"
  name                  = "%s"
  load_balancing_scheme = "INTERNAL"
  backend_service       = "${google_compute_region_backend_service.foobar-bs.self_link}"
  ports                 = ["80"]
  network               = "${google_compute_network.foobar.name}"
  subnetwork            = "%s"
}
resource "google_compute_forwarding_rule" "foobar2" {
  description           = "Resource created for Terraform acceptance testing"
  name                  = "%s"
  load_balancing_scheme = "INTERNAL"
  backend_service       = "${google_compute_region_backend_service.foobar-bs.self_link}"
  ports                 = ["80"]
  network               = "${google_compute_network.foobar.self_link}"
}
`, serviceName, checkName, networkName, ruleName1, networkName, ruleName2)
}

func testAccComputeForwardingRule_networkTier(poolName, ruleName string) string {
	return fmt.Sprintf(`
resource "google_compute_target_pool" "foobar-tp" {
  description = "Resource created for Terraform acceptance testing"
  instances   = ["us-central1-a/foo", "us-central1-b/bar"]
  name        = "%s"
}
resource "google_compute_forwarding_rule" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  ip_protocol = "UDP"
  name        = "%s"
  port_range  = "80-81"
  target      = "${google_compute_target_pool.foobar-tp.self_link}"

  network_tier = "STANDARD"
}
`, poolName, ruleName)
}
