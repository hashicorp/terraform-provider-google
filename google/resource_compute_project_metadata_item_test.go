package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeProjectMetadataItem_basic(t *testing.T) {
	t.Parallel()

	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic("foobar", key, "myValue"),
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicMultiple(t *testing.T) {
	t.Parallel()

	// Generate a config of two config keys
	key1 := "myKey" + randString(t, 10)
	key2 := "myKey" + randString(t, 10)
	config := testAccProjectMetadataItem_basic("foobar", key1, "myValue") +
		testAccProjectMetadataItem_basic("foobar2", key2, "myOtherValue")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicWithEmptyVal(t *testing.T) {
	t.Parallel()

	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic("foobar", key, ""),
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_basicUpdate(t *testing.T) {
	t.Parallel()

	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectMetadataItem_basic("foobar", key, "myValue"),
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProjectMetadataItem_basic("foobar", key, "myUpdatedValue"),
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectMetadataItem_exists(t *testing.T) {
	t.Parallel()

	// Key must be unique to avoid concurrent tests interfering with each other
	key := "myKey" + randString(t, 10)
	originalConfig := testAccProjectMetadataItem_basic("foobar", key, "myValue")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectMetadataItemDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: originalConfig,
			},
			{
				ResourceName:      "google_compute_project_metadata_item.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Add a second resource with the same key
			{
				Config:      originalConfig + testAccProjectMetadataItem_basic("foobar2", key, "myValue"),
				ExpectError: regexp.MustCompile("already present in metadata for project"),
			},
		},
	})
}

func testAccCheckProjectMetadataItemDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		project, err := config.clientCompute.Projects.Get(config.Project).Do()
		if err != nil {
			return err
		}

		metadata := flattenMetadata(project.CommonInstanceMetadata)

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
}

func testAccProjectMetadataItem_basic(resourceName, key, val string) string {
	return fmt.Sprintf(`
resource "google_compute_project_metadata_item" "%s" {
  key   = "%s"
  value = "%s"
}
`, resourceName, key, val)
}
