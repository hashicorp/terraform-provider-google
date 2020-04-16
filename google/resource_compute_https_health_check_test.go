package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeHttpsHealthCheck_update(t *testing.T) {
	t.Parallel()

	hhckName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHttpsHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHttpsHealthCheck_update1(hhckName),
			},
			{
				ResourceName:      "google_compute_https_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeHttpsHealthCheck_update2(hhckName),
			},
			{
				ResourceName:      "google_compute_https_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeHttpsHealthCheck_update1(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
  name         = "%s"
  description  = "Resource created for Terraform acceptance testing"
  request_path = "/not_default"
}
`, hhckName)
}

func testAccComputeHttpsHealthCheck_update2(hhckName string) string {
	return fmt.Sprintf(`
resource "google_compute_https_health_check" "foobar" {
  name                = "%s"
  description         = "Resource updated for Terraform acceptance testing"
  healthy_threshold   = 10
  unhealthy_threshold = 10
}
`, hhckName)
}
