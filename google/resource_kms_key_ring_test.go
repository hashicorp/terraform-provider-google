package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func TestAccGoogleKmsKeyRing_basic(t *testing.T) {
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testGoogleKmsKeyRing_recreate(name),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsKeyRing_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyRingExists("google_kms_key_ring.key_ring"),
				),
			},
		},
	})
}

func testGoogleKmsKeyRing_basic(name string) string {
	return fmt.Sprintf(`
	resource "google_kms_key_ring" "key_ring" {
		name = "%s"
		location = "us-central1"
	}
	`, name)
}

func testAccCheckGoogleKmsKeyRingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		location := rs.Primary.Attributes["location"]

		parent := createKmsResourceParentString(config.Project, location)
		keyRingName := createKmsResourceKeyRingName(config.Project, location, name)

		listKeyRingsResponse, err := config.clientKms.Projects.Locations.KeyRings.List(parent).Do()
		if err != nil {
			return fmt.Errorf("Error listing KeyRings: %s", err)
		}

		for _, keyRing := range listKeyRingsResponse.KeyRings {
			log.Printf("[DEBUG] Found KeyRing: %s", keyRing.Name)

			if keyRing.Name == keyRingName {
				return nil
			}
		}

		return fmt.Errorf("KeyRing not found: %s", keyRingName)
	}
}

// TODO
// KMS KeyRings cannot be deleted. This will test if the resource can be added back to state after being removed
func testGoogleKmsKeyRing_recreate(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}
