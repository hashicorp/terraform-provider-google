package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeHealthCheckDatasource_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheckDatasourceConfig(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_compute_health_check.hc", "google_compute_health_check.hc"),
				),
			},
		},
	})
}

func testAccComputeHealthCheckDatasourceConfig(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "hc" {
  name        = "tf-test-%s"
  description = "Health check via tcp"

  timeout_sec         = 1
  check_interval_sec  = 1
  healthy_threshold   = 4
  unhealthy_threshold = 5

  tcp_health_check {
    port_name          = "health-check-port"
    port_specification = "USE_NAMED_PORT"
    request            = "ARE YOU HEALTHY?"
    proxy_header       = "NONE"
    response           = "I AM HEALTHY"
  }
}

data "google_compute_health_check" "hc" {
  name = google_compute_health_check.hc.name
}
`, suffix)
}
