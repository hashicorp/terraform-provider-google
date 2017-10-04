package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeTargeTcpProxy_basic(t *testing.T) {
	target := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetTcpProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetTcpProxy_basic1(target, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetTcpProxyExists(
						"google_compute_target_tcp_proxy.foobar"),
				),
			},
		},
	})
}

func TestAccComputeTargetTcpProxy_update(t *testing.T) {
	target := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("ttcp-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetTcpProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetTcpProxy_basic1(target, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetTcpProxyExists(
						"google_compute_target_tcp_proxy.foobar"),
				),
			},
		},
	})
}

func testAccCheckComputeTargetTcpProxyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_target_tcp_proxy" {
			continue
		}

		_, err := config.clientCompute.TargetTcpProxies.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("TargetTcpProxy still exists")
		}
	}

	return nil
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

		found, err := config.clientCompute.TargetTcpProxies.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("TargetTcpProxy not found")
		}

		return nil
	}
}

func testAccComputeTargetTcpProxy_basic1(target, backend, hc string) string {
	return fmt.Sprintf(`
	resource "google_compute_target_tcp_proxy" "foobar" {
		description = "Resource created for Terraform acceptance testing"
		name = "%s"
		backend_service = "${google_compute_backend_service.foobar.self_link}"
	}

	resource "google_compute_backend_service" "foobar" {
		name = "%s"
		health_checks = ["${google_compute_http_health_check.zero.self_link}"]
	}

	resource "google_compute_http_health_check" "zero" {
		name = "%s"
		request_path = "/"
		check_interval_sec = 1
		timeout_sec = 1
	}
	`, target, backend, hc)
}
