package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeUrlMap_import(t *testing.T) {
	t.Parallel()

	bsName := fmt.Sprintf("bs-test-%s", acctest.RandString(10))
	hcName := fmt.Sprintf("hc-test-%s", acctest.RandString(10))
	umName := fmt.Sprintf("um-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeUrlMapDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeUrlMap_basic1(bsName, hcName, umName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_url_map.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}})
}
