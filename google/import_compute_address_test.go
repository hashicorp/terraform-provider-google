package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeAddress_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_compute_address.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeAddress_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeAddress_importInternal(t *testing.T) {
	t.Parallel()

	resourceName := "google_compute_address.internal_with_subnet"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeAddress_internal,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
