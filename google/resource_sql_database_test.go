package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/sqladmin/v1beta4"
)

func TestAccGoogleSqlDatabase_basic(t *testing.T) {
	var database sqladmin.Database

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabase_basic, acctest.RandString(10), acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseExists(
						"google_sql_database.database", &database),
					testAccCheckGoogleSqlDatabaseEquals(
						"google_sql_database.database", &database),
				),
			},
		},
	})
}

func TestAccGoogleSqlDatabase_update(t *testing.T) {
	var database sqladmin.Database

	instance_name := acctest.RandString(10)
	database_name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlDatabaseInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabase_basic, instance_name, database_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseExists(
						"google_sql_database.database", &database),
					testAccCheckGoogleSqlDatabaseEquals(
						"google_sql_database.database", &database),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(
					testGoogleSqlDatabase_latin1, instance_name, database_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseExists(
						"google_sql_database.database", &database),
					testAccCheckGoogleSqlDatabaseEquals(
						"google_sql_database.database", &database),
				),
			},
		},
	})
}

func testAccCheckGoogleSqlDatabaseEquals(n string,
	database *sqladmin.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		database_name := rs.Primary.Attributes["name"]
		instance_name := rs.Primary.Attributes["instance"]
		charset := rs.Primary.Attributes["charset"]
		collation := rs.Primary.Attributes["collation"]

		if database_name != database.Name {
			return fmt.Errorf("Error name mismatch, (%s, %s)", database_name, database.Name)
		}

		if instance_name != database.Instance {
			return fmt.Errorf("Error instance_name mismatch, (%s, %s)", instance_name, database.Instance)
		}

		if charset != database.Charset {
			return fmt.Errorf("Error charset mismatch, (%s, %s)", charset, database.Charset)
		}

		if collation != database.Collation {
			return fmt.Errorf("Error collation mismatch, (%s, %s)", collation, database.Collation)
		}

		return nil
	}
}

func testAccCheckGoogleSqlDatabaseExists(n string,
	database *sqladmin.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		database_name := rs.Primary.Attributes["name"]
		instance_name := rs.Primary.Attributes["instance"]
		found, err := config.clientSqlAdmin.Databases.Get(config.Project,
			instance_name, database_name).Do()

		if err != nil {
			return fmt.Errorf("Not found: %s: %s", n, err)
		}

		*database = *found

		return nil
	}
}

func testAccGoogleSqlDatabaseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*Config)
		if rs.Type != "google_sql_database" {
			continue
		}

		database_name := rs.Primary.Attributes["name"]
		instance_name := rs.Primary.Attributes["instance"]
		_, err := config.clientSqlAdmin.Databases.Get(config.Project,
			instance_name, database_name).Do()

		if err == nil {
			return fmt.Errorf("Database resource still exists")
		}
	}

	return nil
}

var testGoogleSqlDatabase_basic = `
resource "google_sql_database_instance" "instance" {
	name = "sqldatabasetest%s"
	region = "us-central"
	settings {
		tier = "D0"
	}
}

resource "google_sql_database" "database" {
	name = "sqldatabasetest%s"
	instance = "${google_sql_database_instance.instance.name}"
}
`
var testGoogleSqlDatabase_latin1 = `
resource "google_sql_database_instance" "instance" {
	name = "sqldatabasetest%s"
	region = "us-central"
	settings {
		tier = "D0"
	}
}

resource "google_sql_database" "database" {
	name = "sqldatabasetest%s"
	instance = "${google_sql_database_instance.instance.name}"
	charset = "latin1"
	collation = "latin1_swedish_ci"
}
`
