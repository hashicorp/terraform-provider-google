package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleContainerNodePool_import(t *testing.T) {
	resourceName := "google_container_node_pool.np"
	cluster := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	np := fmt.Sprintf("tf-nodepool-test-%s", acctest.RandString(10))
	conf := testAccContainerNodePool_basic(cluster, np)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerNodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
