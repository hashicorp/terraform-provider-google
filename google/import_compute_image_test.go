package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeImage_importFromRawDisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeImage_basic,
			},
			resource.TestStep{
				ResourceName:            "google_compute_image.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"raw_disk", "create_timeout"},
			},
		},
	})
}

func TestAccComputeImage_importFromSourceDisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeImage_basedondisk,
			},
			resource.TestStep{
				ResourceName:      "google_compute_image.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
