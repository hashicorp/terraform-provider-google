package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeTargetSslProxy_import(t *testing.T) {
	target := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	cert := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetSslProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetSslProxy_basic1(target, cert, backend, hc),
			},
			resource.TestStep{
				ResourceName:      "google_compute_target_ssl_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
