package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeHttpsHealthCheck_importBasic(t *testing.T) {
	hhckName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeHttpsHealthCheckDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeHttpsHealthCheck_basic(hhckName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_https_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
