package google

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/googleapi"
)

func TestAccSpannerDatabase_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	rnd := acctest.RandString(10)
	instanceName := fmt.Sprintf("my-instance-%s", rnd)
	databaseName := fmt.Sprintf("mydb_%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
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

func TestAccSpannerDatabase_basicWithInitialDDL(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	instanceName := fmt.Sprintf("my-instance-%s", rnd)
	databaseName := fmt.Sprintf("mydb-%s", rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_basicWithInitialDDL(instanceName, databaseName),
			},
			{
				ResourceName:      "google_spanner_database.basic",
				ImportState:       true,
				ImportStateVerify: true,
				// DDL statements get issued at the time the create/update
				// occurs, which means storing them in state isn't really
				// necessary.
				ImportStateVerifyIgnore: []string{"ddl"},
			},
		},
	})
}

func testAccCheckSpannerDatabaseDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_spanner_database" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Unable to verify delete of spanner database, ID is empty")
		}

		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		id := spannerDatabaseId{
			Project:  project,
			Instance: rs.Primary.Attributes["instance"],
			Database: rs.Primary.Attributes["name"],
		}
		_, err = config.clientSpanner.Projects.Instances.Databases.Get(
			id.databaseUri()).Do()

		if err == nil {
			return fmt.Errorf("Spanner database still exists")
		}

		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return nil
		}
		return errwrap.Wrapf("Error verifying spanner database deleted: {{err}}", err)
	}

	return nil
}

func testAccSpannerDatabase_basic(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "display-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "%s"
}
`, instanceName, instanceName, databaseName)
}

func testAccSpannerDatabase_basicWithInitialDDL(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "display-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "%s"
  ddl           =  [
     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
     "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)" ]
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
