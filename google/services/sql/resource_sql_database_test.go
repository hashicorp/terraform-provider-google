// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func TestAccSqlDatabase_basic(t *testing.T) {
	t.Parallel()

	var database sqladmin.Database

	resourceName := "google_sql_database.database"
	instanceName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	dbName := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testGoogleSqlDatabase_basic, instanceName, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseExists(t, resourceName, &database),
					testAccCheckGoogleSqlDatabaseEquals(resourceName, &database),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     fmt.Sprintf("%s/%s", instanceName, dbName),
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      resourceName,
				ImportStateId:     fmt.Sprintf("instances/%s/databases/%s", instanceName, dbName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           fmt.Sprintf("%s/%s/%s", envvar.GetTestProjectFromEnv(), instanceName, dbName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				ResourceName:            resourceName,
				ImportStateId:           fmt.Sprintf("projects/%s/instances/%s/databases/%s", envvar.GetTestProjectFromEnv(), instanceName, dbName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccSqlDatabase_update(t *testing.T) {
	t.Parallel()

	var database sqladmin.Database

	instance_name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	database_name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabase_basic, instance_name, database_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseExists(
						t, "google_sql_database.database", &database),
					testAccCheckGoogleSqlDatabaseEquals(
						"google_sql_database.database", &database),
				),
			},
			{
				Config: fmt.Sprintf(
					testGoogleSqlDatabase_latin1, instance_name, database_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlDatabaseExists(
						t, "google_sql_database.database", &database),
					testAccCheckGoogleSqlDatabaseEquals(
						"google_sql_database.database", &database),
				),
			},
		},
	})
}

func testAccCheckGoogleSqlDatabaseEquals(n string, database *sqladmin.Database) resource.TestCheckFunc {
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

func testAccCheckGoogleSqlDatabaseExists(t *testing.T, n string, database *sqladmin.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		database_name := rs.Primary.Attributes["name"]
		instance_name := rs.Primary.Attributes["instance"]
		found, err := config.NewSqlAdminClient(config.UserAgent).Databases.Get(config.Project,
			instance_name, database_name).Do()

		if err != nil {
			return fmt.Errorf("Not found: %s: %s", n, err)
		}

		*database = *found

		return nil
	}
}

func testAccSqlDatabaseDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			config := acctest.GoogleProviderConfig(t)
			if rs.Type != "google_sql_database" {
				continue
			}

			database_name := rs.Primary.Attributes["name"]
			instance_name := rs.Primary.Attributes["instance"]
			_, err := config.NewSqlAdminClient(config.UserAgent).Databases.Get(config.Project,
				instance_name, database_name).Do()

			if err == nil {
				return fmt.Errorf("Database resource still exists")
			}
		}

		return nil
	}
}

var testGoogleSqlDatabase_basic = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "database" {
  name     = "%s"
  instance = google_sql_database_instance.instance.name
}
`
var testGoogleSqlDatabase_latin1 = `
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "database" {
  name      = "%s"
  instance  = google_sql_database_instance.instance.name
  charset   = "latin1"
  collation = "latin1_swedish_ci"
}
`
