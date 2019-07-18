package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeHealthCheck_tcp(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_tcp(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckThresholds(
						3, 3, &healthCheck),
					testAccCheckComputeHealthCheckTcpPort(80, &healthCheck),
					testAccCheckComputeHealthCheckPortSpec(
						"TCP", "", &healthCheck,
					),
				),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_tcp_update(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_tcp(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckThresholds(
						3, 3, &healthCheck),
					testAccCheckComputeHealthCheckTcpPort(80, &healthCheck),
				),
			},
			{
				Config: testAccComputeHealthCheck_tcp_update(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckThresholds(
						10, 10, &healthCheck),
					testAccCheckComputeHealthCheckTcpPort(8080, &healthCheck),
				),
			},
		},
	})
}

func TestAccComputeHealthCheck_ssl(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_ssl(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckThresholds(
						3, 3, &healthCheck),
				),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_ssl_port_spec(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_ssl_fixed_port(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckPortSpec(
						"SSL", "USE_FIXED_PORT", &healthCheck),
				),
			},
		},
	})
}

func TestAccComputeHealthCheck_http(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_http(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckThresholds(
						3, 3, &healthCheck),
				),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_http_port_spec(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeHealthCheck_http_port_spec(hckName),
				ExpectError: regexp.MustCompile("Error in http_health_check: Must specify port_name when using USE_NAMED_PORT as port_specification."),
			},
			{
				Config: testAccComputeHealthCheck_http_named_port(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckPortSpec(
						"HTTP", "USE_NAMED_PORT", &healthCheck,
					),
				),
			},
		},
	})
}

func TestAccComputeHealthCheck_https(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_https(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckThresholds(
						3, 3, &healthCheck),
				),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_https_serving_port(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HealthCheck

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_https_serving_port(hckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHealthCheckExists(
						"google_compute_health_check.foobar", &healthCheck),
					testAccCheckComputeHealthCheckPortSpec(
						"HTTPS", "USE_SERVING_PORT", &healthCheck,
					),
				),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_typeTransition(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_https(hckName),
			},
			{
				Config: testAccComputeHealthCheck_http(hckName),
			},
			{
				Config: testAccComputeHealthCheck_ssl(hckName),
			},
			{
				Config: testAccComputeHealthCheck_tcp(hckName),
			},
			{
				Config: testAccComputeHealthCheck_https(hckName),
			},
		},
	})
}

func TestAccComputeHealthCheck_tcpAndSsl_shouldFail(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeHealthCheck_tcpAndSsl_shouldFail(hckName),
				ExpectError: regexp.MustCompile("conflicts with tcp_health_check"),
			},
		},
	})
}

func testAccCheckComputeHealthCheckExists(n string, healthCheck *compute.HealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.HealthChecks.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("HealthCheck not found")
		}

		*healthCheck = *found

		return nil
	}
}

func testAccCheckComputeHealthCheckThresholds(healthy, unhealthy int64, healthCheck *compute.HealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if healthCheck.HealthyThreshold != healthy {
			return fmt.Errorf("HealthyThreshold doesn't match: expected %d, got %d", healthy, healthCheck.HealthyThreshold)
		}

		if healthCheck.UnhealthyThreshold != unhealthy {
			return fmt.Errorf("UnhealthyThreshold doesn't match: expected %d, got %d", unhealthy, healthCheck.UnhealthyThreshold)
		}

		return nil
	}
}

func testAccCheckComputeHealthCheckTcpPort(port int64, healthCheck *compute.HealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if healthCheck.TcpHealthCheck.Port != port {
			return fmt.Errorf("Port doesn't match: expected %v, got %v", port, healthCheck.TcpHealthCheck.Port)
		}
		return nil
	}
}

func testAccCheckComputeHealthCheckPortSpec(blockType, portSpec string, healthCheck *compute.HealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var actualPortSpec string

		switch blockType {
		case "SSL":
			actualPortSpec = healthCheck.SslHealthCheck.PortSpecification
		case "HTTP":
			actualPortSpec = healthCheck.HttpHealthCheck.PortSpecification
		case "HTTPS":
			actualPortSpec = healthCheck.HttpsHealthCheck.PortSpecification
		case "TCP":
			actualPortSpec = healthCheck.TcpHealthCheck.PortSpecification
		}

		if actualPortSpec != portSpec {
			return fmt.Errorf("Port Specification doesn't match: expected %v, got %v", portSpec, actualPortSpec)
		}

		return nil
	}
}

func testAccComputeHealthCheck_tcp(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	tcp_health_check {
	}
}
`, hckName)
}

func testAccComputeHealthCheck_tcp_update(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource updated for Terraform acceptance testing"
	healthy_threshold = 10
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 10
	tcp_health_check {
		port = "8080"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_ssl(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	ssl_health_check {
		port = "443"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_ssl_fixed_port(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	ssl_health_check {
		port = "443"
		port_specification = "USE_FIXED_PORT"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_http(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	http_health_check {
		port = "80"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_http_port_spec(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	http_health_check {
		port_specification = "USE_NAMED_PORT"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_http_named_port(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	http_health_check {
		port_name = "http"
		port_specification = "USE_NAMED_PORT"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_https(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	https_health_check {
		port = "443"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_https_serving_port(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3
	https_health_check {
		port_specification = "USE_SERVING_PORT"
	}
}
`, hckName)
}

func testAccComputeHealthCheck_tcpAndSsl_shouldFail(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	name = "health-test-%s"
	timeout_sec = 2
	unhealthy_threshold = 3

	tcp_health_check {
	}
	ssl_health_check {
	}
}
`, hckName)
}
