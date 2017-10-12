package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeTargetTcpProxy_import(t *testing.T) {
	t.Parallel()

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
			},
			resource.TestStep{
				ResourceName:      "google_compute_target_tcp_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
