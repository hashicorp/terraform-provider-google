package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeProjectMetadataItem_basic(t *testing.T) {
	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic(key, "myValue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata(key, "myValue"),
				),
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicMultiple(t *testing.T) {
	// Generate a config of two config keys
	config := testAccProjectMetadataItem_basic("myKey", "myValue") +
		testAccProjectMetadataItem_basic("myOtherKey", "myOtherValue")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata("myKey", "myValue"),
					testAccCheckProjectMetadataItem_hasMetadata("myOtherKey", "myOtherValue"),
				),
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicWithEmptyVal(t *testing.T) {
	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic(key, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata(key, ""),
				),
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicUpdate(t *testing.T) {
	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic(key, "myValue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata(key, "myValue"),
				),
			},
			{
				Config: testAccProjectMetadataItem_basic(key, "myUpdatedValue"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectMetadataItem_hasMetadata(key, "myUpdatedValue"),
				),
			},
		},
	})
}

func testAccCheckProjectMetadataItem_hasMetadata(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		project, err := config.clientCompute.Projects.Get(config.Project).Do()
		if err != nil {
			return err
		}

		metadata := flattenComputeMetadata(project.CommonInstanceMetadata.Items)

		val, ok := metadata[key]
		if !ok {
			return fmt.Errorf("Unable to find a value for key '%s'", key)
		}
		if val != value {
			return fmt.Errorf("Value for key '%s' does not match. Expected '%s' but found '%s'", key, value, val)
		}
		return nil
	}
}

func testAccCheckProjectMetadataItemDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	project, err := config.clientCompute.Projects.Get(config.Project).Do()
	if err != nil {
		return err
	}

	metadata := flattenComputeMetadata(project.CommonInstanceMetadata.Items)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_project_metadata_item" {
			continue
		}

		_, ok := metadata[rs.Primary.ID]
		if ok {
			return fmt.Errorf("Metadata key/value '%s': '%s' still exist", rs.Primary.Attributes["key"], rs.Primary.Attributes["value"])
		}
	}

	return nil
}

func testAccProjectMetadataItem_basic(key, val string) string {
	return fmt.Sprintf(`
resource "google_compute_project_metadata_item" "foobar-%s" {
  key   = "%s"
  value = "%s"
}
`, acctest.RandString(10), key, val)
}
