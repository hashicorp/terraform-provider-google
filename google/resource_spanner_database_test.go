package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSpannerDatabase_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	rnd := randString(t, 10)
	instanceName := fmt.Sprintf("my-instance-%s", rnd)
	databaseName := fmt.Sprintf("mydb_%s", rnd)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_basic(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
				),
			},
			{
				// Test import with default Terraform ID
				ResourceName:      "google_spanner_database.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_spanner_database.basic",
				ImportStateId:     fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instanceName, databaseName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_spanner_database.basic",
				ImportStateId:     fmt.Sprintf("instances/%s/databases/%s", instanceName, databaseName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_spanner_database.basic",
				ImportStateId:     fmt.Sprintf("%s/%s", instanceName, databaseName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerDatabase_basic(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "display-%s"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
}
`, instanceName, instanceName, databaseName)
}

// Unit Tests for type spannerDatabaseId
func TestDatabaseNameForApi(t *testing.T) {
	id := spannerDatabaseId{
		Project:  "project123",
		Instance: "instance456",
		Database: "db789",
	}
	actual := id.databaseUri()
	expected := "projects/project123/instances/instance456/databases/db789"
	expectEquals(t, expected, actual)
}
