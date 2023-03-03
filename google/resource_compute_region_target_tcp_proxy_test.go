package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeRegionTargetTcpProxy_update(t *testing.T) {
	t.Parallel()

	target := fmt.Sprintf("trtcp-test-%s", RandString(t, 10))
	backend := fmt.Sprintf("trtcp-test-%s", RandString(t, 10))
	hc := fmt.Sprintf("trtcp-test-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckComputeRegionTargetTcpProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionTargetTcpProxy_basic1(target, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionTargetTcpProxyExists(
						t, "google_compute_region_target_tcp_proxy.foobar"),
				),
			},
			{
				ResourceName:      "google_compute_region_target_tcp_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_region_backend_service.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_region_health_check.zero",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeRegionTargetTcpProxy_update2(target, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeRegionTargetTcpProxyExists(
						t, "google_compute_region_target_tcp_proxy.foobar"),
				),
			},
		},
	})
}

func testAccCheckComputeRegionTargetTcpProxyExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := GoogleProviderConfig(t)
		name := rs.Primary.Attributes["name"]
		region := rs.Primary.Attributes["region"]

		found, err := config.NewComputeClient(config.UserAgent).RegionTargetTcpProxies.Get(
			config.Project, region, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("RegionTargetTcpProxy not found")
		}

		return nil
	}
}

func testAccComputeRegionTargetTcpProxy_basic1(target, backend, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_tcp_proxy" "foobar" {
  description     = "Resource created for Terraform acceptance testing"
  name            = "%s"
  backend_service = google_compute_region_backend_service.foobar.self_link
  proxy_header    = "NONE"
  region          = "us-central1"
}

resource "google_compute_region_backend_service" "foobar" {
  name          = "%s"
  protocol      = "TCP"
  health_checks = [google_compute_region_health_check.zero.self_link]
  region        = "us-central1"

  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
  region = "us-central1"
}
`, target, backend, hc)
}

func testAccComputeRegionTargetTcpProxy_update2(target, backend, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_region_target_tcp_proxy" "foobar" {
  description     = "Resource created for Terraform acceptance testing"
  name            = "%s"
  backend_service = google_compute_region_backend_service.foobar2.self_link
  proxy_header    = "PROXY_V1"
  region          = "us-central1"
}

resource "google_compute_region_backend_service" "foobar" { 
  name          = "%s"
  protocol      = "TCP"
  health_checks = [google_compute_region_health_check.zero.self_link]
  region        = "us-central1"

  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_backend_service" "foobar2" { 
  name          = "%s-2"
  protocol      = "TCP"
  health_checks = [google_compute_region_health_check.zero.self_link]
  region        = "us-central1"

  load_balancing_scheme = "INTERNAL_MANAGED"
}

resource "google_compute_region_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
  region = "us-central1"
}
`, target, backend, backend, hc)
}
