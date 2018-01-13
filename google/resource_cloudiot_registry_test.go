package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudiotRegistryCreate(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudiotRegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudiotRegistry(),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudiotRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
		},
	})
}

func testAccCheckCloudiotRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudiot_registry" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		registry, _ := config.clientCloudiot.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if registry != nil {
			return fmt.Errorf("Registry still present")
		}
	}

	return nil
}

func testAccCloudiotRegistryExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientCloudiot.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Topic does not exist")
		}

		return nil
	}
}

func testAccCloudiotRegistry() string {
	return fmt.Sprintf(`
resource "google_cloudiot_registry" "foobar" {
	name = "psregistry-test-%s"
}`, acctest.RandString(10))
}
