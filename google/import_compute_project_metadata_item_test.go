package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeProjectMetadataItem_importBasic(t *testing.T) {
	t.Parallel()

	key := "myKey" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basicWithResourceName("foobar", key, "myValue"),
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
