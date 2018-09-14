package google

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/googleapi"
)

// Unit Tests

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

func TestImportSpannerDatabaseId_InstanceDB(t *testing.T) {
	id, e := importSpannerDatabaseId("instance456/database789")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "", id.Project)
	expectEquals(t, "instance456", id.Instance)
	expectEquals(t, "database789", id.Database)
}

func TestImportSpannerDatabaseId_ProjectInstanceDB(t *testing.T) {
	id, e := importSpannerDatabaseId("project123/instance456/database789")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "project123", id.Project)
	expectEquals(t, "instance456", id.Instance)
	expectEquals(t, "database789", id.Database)
}

func TestImportSpannerDatabaseId_projectId(t *testing.T) {
	shouldPass := []string{
		"project-id/instance/database",
		"123123/instance/123",
		"hashicorptest.net:project-123/instance/123",
		"123/456/789",
	}

	shouldFail := []string{
		"project-id#/instance/database",
		"project-id/instance#/database",
		"project-id/instance/database#",
		"hashicorptest.net:project-123:invalid:project/instance/123",
		"hashicorptest.net:/instance/123",
	}

	for _, element := range shouldPass {
		_, e := importSpannerDatabaseId(element)
		if e != nil {
			t.Error("importSpannerDatabaseId should pass on '" + element + "' but doesn't")
		}
	}

	for _, element := range shouldFail {
		_, e := importSpannerDatabaseId(element)
		if e == nil {
			t.Error("importSpannerDatabaseId should fail on '" + element + "' but doesn't")
		}
	}
}

func TestImportSpannerDatabaseId_invalidLeadingSlash(t *testing.T) {
	id, e := importSpannerDatabaseId("/instance456/database789")
	expectInvalidSpannerDbImportId(t, id, e)
}

func TestImportSpannerDatabaseId_invalidTrailingSlash(t *testing.T) {
	id, e := importSpannerDatabaseId("instance456/database789/")
	expectInvalidSpannerDbImportId(t, id, e)
}

func TestImportSpannerDatabaseId_invalidSingleSlash(t *testing.T) {
	id, e := importSpannerDatabaseId("/")
	expectInvalidSpannerDbImportId(t, id, e)
}

func TestImportSpannerDatabaseId_invalidMultiSlash(t *testing.T) {
	id, e := importSpannerDatabaseId("project123/instance456/db789/next")
	expectInvalidSpannerDbImportId(t, id, e)
}

func expectInvalidSpannerDbImportId(t *testing.T, id *spannerDatabaseId, e error) {
	if id != nil {
		t.Errorf("Expected spannerDatabaseId to be nil")
		return
	}
	if e == nil {
		t.Errorf("Expected an Error but did not get one")
		return
	}
	if !strings.HasPrefix(e.Error(), "Invalid spanner database specifier") {
		t.Errorf("Expecting Error starting with 'Invalid spanner database specifier'")
	}
}

// Acceptance Tests

func TestAccSpannerDatabase_basic(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_basic(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_database.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerDatabase_basicWithInitialDDL(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_basicWithInitialDDL(rnd),
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

func testAccSpannerDatabase_basic(rnd string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "my-instance-%s"
  config        = "regional-us-central1"
  display_name  = "my-displayname-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "my-db-%s"
}
`, rnd, rnd, rnd)
}

func testAccSpannerDatabase_basicWithInitialDDL(rnd string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "my-instance-%s"
  config        = "regional-us-central1"
  display_name  = "my-displayname-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "my-db-%s"
  ddl           =  [
     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
     "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)" ]
}
`, rnd, rnd, rnd)
}
