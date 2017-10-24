package google

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleKmsKeyRing_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_kms_key_ring.key_ring"
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsKeyRing_import(name),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleKmsKeyRing_import(keyRingName string) string {
	return fmt.Sprintf(`
resource "google_kms_key_ring" "key_ring" {
	name     = "%s"
	location = "us-central1"
}
	`, keyRingName)
}
