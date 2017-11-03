package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeTargetHttpProxy_import(t *testing.T) {
	t.Parallel()

	target := fmt.Sprintf("thttp-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("thttp-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("thttp-test-%s", acctest.RandString(10))
	urlmap1 := fmt.Sprintf("thttp-test-%s", acctest.RandString(10))
	urlmap2 := fmt.Sprintf("thttp-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetHttpProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetHttpProxy_basic1(target, backend, hc, urlmap1, urlmap2),
			},
			resource.TestStep{
				ResourceName:      "google_compute_target_http_proxy.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
