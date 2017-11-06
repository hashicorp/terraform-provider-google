package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
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
				Config: testAccComputeAddress_basic(acctest.RandString(10)),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeAddress_internal(acctest.RandString(10)),
			},

			resource.TestStep{
				ResourceName:      "google_compute_address.internal",
				ImportState:       true,
				ImportStateVerify: true,
			},

			resource.TestStep{
				ResourceName:      "google_compute_address.internal_with_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},

			resource.TestStep{
				ResourceName:      "google_compute_address.internal_with_subnet_and_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
