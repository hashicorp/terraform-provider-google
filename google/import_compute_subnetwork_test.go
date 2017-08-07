package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeSubnetwork_importBasic(t *testing.T) {
	cnName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	subnetwork1Name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	subnetwork2Name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	subnetwork3Name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSubnetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSubnetwork_basic(cnName, subnetwork1Name, subnetwork2Name, subnetwork3Name),
			},
			resource.TestStep{
				ResourceName:      "google_compute_subnetwork.network-ref-by-url",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				ResourceName:      "google_compute_subnetwork.network-with-private-google-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
