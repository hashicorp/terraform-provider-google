package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeHttpsHealthCheck_basic(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HttpsHealthCheck

	hhckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHttpsHealthCheckDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeHttpsHealthCheck_basic(hhckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHttpsHealthCheckExists(
						"google_compute_https_health_check.foobar", &healthCheck),
					testAccCheckComputeHttpsHealthCheckRequestPath(
						"/health_check", &healthCheck),
					testAccCheckComputeHttpsHealthCheckThresholds(
						3, 3, &healthCheck),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_https_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHttpsHealthCheck_update(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HttpsHealthCheck

	hhckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHttpsHealthCheckDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeHttpsHealthCheck_update1(hhckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHttpsHealthCheckExists(
						"google_compute_https_health_check.foobar", &healthCheck),
					testAccCheckComputeHttpsHealthCheckRequestPath(
						"/not_default", &healthCheck),
					testAccCheckComputeHttpsHealthCheckThresholds(
						2, 2, &healthCheck),
				),
			},
			resource.TestStep{
				Config: testAccComputeHttpsHealthCheck_update2(hhckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHttpsHealthCheckExists(
						"google_compute_https_health_check.foobar", &healthCheck),
					testAccCheckComputeHttpsHealthCheckRequestPath(
						"/", &healthCheck),
					testAccCheckComputeHttpsHealthCheckThresholds(
						10, 10, &healthCheck),
				),
			},
		},
	})
}

func testAccCheckComputeHttpsHealthCheckExists(n string, healthCheck *compute.HttpsHealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.HttpsHealthChecks.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("HttpsHealthCheck not found")
		}

		*healthCheck = *found

		return nil
	}
}

func testAccCheckComputeHttpsHealthCheckRequestPath(path string, healthCheck *compute.HttpsHealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if healthCheck.RequestPath != path {
			return fmt.Errorf("RequestPath doesn't match: expected %s, got %s", path, healthCheck.RequestPath)
		}

		return nil
	}
}

func testAccCheckComputeHttpsHealthCheckThresholds(healthy, unhealthy int64, healthCheck *compute.HttpsHealthCheck) resource.TestCheckFunc {
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

func testAccComputeHttpsHealthCheck_basic(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
	check_interval_sec = 3
	description = "Resource created for Terraform acceptance testing"
	healthy_threshold = 3
	host = "foobar"
	name = "%s"
	port = "80"
	request_path = "/health_check"
	timeout_sec = 2
	unhealthy_threshold = 3
}
`, hhckName)
}

func testAccComputeHttpsHealthCheck_update1(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
	name = "%s"
	description = "Resource created for Terraform acceptance testing"
	request_path = "/not_default"
}
`, hhckName)
}

func testAccComputeHttpsHealthCheck_update2(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
	name = "%s"
	description = "Resource updated for Terraform acceptance testing"
	healthy_threshold = 10
	unhealthy_threshold = 10
}
`, hhckName)
}
