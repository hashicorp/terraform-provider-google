package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeHttpHealthCheck_update(t *testing.T) {
	t.Parallel()

	var healthCheck compute.HttpHealthCheck

	hhckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHttpHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHttpHealthCheck_update1(hhckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHttpHealthCheckExists(
						"google_compute_http_health_check.foobar", &healthCheck),
					testAccCheckComputeHttpHealthCheckRequestPath(
						"/not_default", &healthCheck),
					testAccCheckComputeHttpHealthCheckThresholds(
						2, 2, &healthCheck),
				),
			},
			{
				Config: testAccComputeHttpHealthCheck_update2(hhckName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHttpHealthCheckExists(
						"google_compute_http_health_check.foobar", &healthCheck),
					testAccCheckComputeHttpHealthCheckRequestPath(
						"/", &healthCheck),
					testAccCheckComputeHttpHealthCheckThresholds(
						10, 10, &healthCheck),
				),
			},
		},
	})
}

func testAccCheckComputeHttpHealthCheckExists(n string, healthCheck *compute.HttpHealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No name is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.HttpHealthChecks.Get(
			config.Project, rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("HttpHealthCheck not found")
		}

		*healthCheck = *found

		return nil
	}
}

func testAccCheckComputeHttpHealthCheckRequestPath(path string, healthCheck *compute.HttpHealthCheck) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if healthCheck.RequestPath != path {
			return fmt.Errorf("RequestPath doesn't match: expected %s, got %s", path, healthCheck.RequestPath)
		}

		return nil
	}
}

func testAccCheckComputeHttpHealthCheckThresholds(healthy, unhealthy int64, healthCheck *compute.HttpHealthCheck) resource.TestCheckFunc {
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

func testAccComputeHttpHealthCheck_update1(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_http_health_check" "foobar" {
  name         = "%s"
  description  = "Resource created for Terraform acceptance testing"
  request_path = "/not_default"
}
`, hhckName)
}

func testAccComputeHttpHealthCheck_update2(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_http_health_check" "foobar" {
  name                = "%s"
  description         = "Resource updated for Terraform acceptance testing"
  healthy_threshold   = 10
  unhealthy_threshold = 10
}
`, hhckName)
}
