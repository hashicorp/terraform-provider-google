package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Unit Tests

func TestSpannerInstanceId_instanceUri(t *testing.T) {
	id := spannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.instanceUri()
	expected := "projects/project123/instances/instance456"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_instanceConfigUri(t *testing.T) {
	id := spannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.instanceConfigUri("conf987")
	expected := "projects/project123/instanceConfigs/conf987"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_parentProjectUri(t *testing.T) {
	id := spannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.parentProjectUri()
	expected := "projects/project123"
	expectEquals(t, expected, actual)
}

func expectEquals(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

// Acceptance Tests

func TestAccSpannerInstance_basic(t *testing.T) {
	t.Parallel()

	idName := fmt.Sprintf("spanner-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(idName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutogenName(t *testing.T) {
	// Randomness
	skipIfVcr(t)
	t.Parallel()

	displayName := fmt.Sprintf("spanner-test-%s-dname", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutogenName(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "name"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_update(t *testing.T) {
	// Randomness
	skipIfVcr(t)
	t.Parallel()

	dName1 := fmt.Sprintf("spanner-dname1-%s", randString(t, 10))
	dName2 := fmt.Sprintf("spanner-dname2-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_update(dName1, 1, false),
			},
			{
				ResourceName:      "google_spanner_instance.updater",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstance_update(dName2, 2, true),
			},
			{
				ResourceName:      "google_spanner_instance.updater",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"
  num_nodes    = 1
}
`, name, name)
}

func testAccSpannerInstance_basicWithAutogenName(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = 1
}
`, name)
}

func testAccSpannerInstance_update(name string, nodes int, addLabel bool) string {
	extraLabel := ""
	if addLabel {
		extraLabel = "\"key2\" = \"value2\""
	}
	return fmt.Sprintf(`
resource "google_spanner_instance" "updater" {
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = %d

  labels = {
    "key1" = "value1"
    %s
  }
}
`, name, nodes, extraLabel)
}
