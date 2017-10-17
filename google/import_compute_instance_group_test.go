package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeInstanceGroup_import(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("instancegroup-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccComputeInstanceGroup_destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceGroup_basic(instanceName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
