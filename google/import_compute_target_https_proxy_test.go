package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeTargetHttpsProxy_import(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("thttps-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetHttpsProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetHttpsProxy_basic1(id),
			},
			resource.TestStep{
				ResourceName:      "google_compute_target_https_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
