package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeTargetTcpProxy_update(t *testing.T) {
	t.Parallel()

	target := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetTcpProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetTcpProxy_basic1(target, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetTcpProxyExists(
						"google_compute_target_tcp_proxy.foobar"),
				),
			},
			{
				Config: testAccComputeTargetTcpProxy_basic2(target, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetTcpProxyExists(
						"google_compute_target_tcp_proxy.foobar"),
				),
			},
		},
	})
}

func testAccCheckComputeTargetTcpProxyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]

		found, err := config.clientCompute.TargetTcpProxies.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("TargetTcpProxy not found")
		}

		return nil
	}
}

func testAccComputeTargetTcpProxy_basic1(target, backend, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_target_tcp_proxy" "foobar" {
  description     = "Resource created for Terraform acceptance testing"
  name            = "%s"
  backend_service = google_compute_backend_service.foobar.self_link
  proxy_header    = "NONE"
}

resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  protocol      = "TCP"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}
`, target, backend, hc)
}

func testAccComputeTargetTcpProxy_basic2(target, backend, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_target_tcp_proxy" "foobar" {
  description     = "Resource created for Terraform acceptance testing"
  name            = "%s"
  backend_service = google_compute_backend_service.foobar.self_link
  proxy_header    = "PROXY_V1"
}

resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  protocol      = "TCP"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}
`, target, backend, hc)
}
