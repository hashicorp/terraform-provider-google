// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigqueryconnection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigqueryConnectionConnection_bigqueryConnectionBasic(t *testing.T) {
	// Uses random provider
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryConnectionConnectionDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryConnectionConnection_bigqueryConnectionBasic(context),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_sql.0.credential.0.password", "cloud_sql.0.credential.0.username"},
				ResourceName:            "google_bigquery_connection.connection",
			},
			{
				Config: testAccBigqueryConnectionConnection_bigqueryConnectionBasicUpdate(context),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_sql.0.credential.0.password", "cloud_sql.0.credential.0.username"},
				ResourceName:            "google_bigquery_connection.connection",
			},
		},
	})
}

func testAccBigqueryConnectionConnection_bigqueryConnectionBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
    name             = "tf-test-pg-database-instance%{random_suffix}"
    database_version = "POSTGRES_11"
    region           = "us-central1"
    settings {
		tier = "db-f1-micro"
	}

    deletion_protection = false
}

resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "db"
}

resource "random_password" "pwd" {
    length = 16
    special = false
}

resource "google_sql_user" "user" {
    name = "username"
    instance = google_sql_database_instance.instance.name
    password = random_password.pwd.result
}

resource "google_bigquery_connection" "connection" {
    connection_id = "tf-test-my-connection%{random_suffix}"
    location      = "US"
    friendly_name = "👋"
    description   = "a riveting description"
    cloud_sql {
        instance_id = google_sql_database_instance.instance.connection_name
        database    = google_sql_database.db.name
        type        = "POSTGRES"
        credential {
            username = google_sql_user.user.name
            password = google_sql_user.user.password
        }
    }
}
`, context)
}

func testAccBigqueryConnectionConnection_bigqueryConnectionBasicUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
    name             = "tf-test-mysql-database-instance%{random_suffix}"
    database_version = "MYSQL_5_6"
    region           = "us-central1"
    settings {
		tier = "db-f1-micro"
	}

    deletion_protection = false
}

resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "db2"
}

resource "random_password" "pwd" {
    length = 16
    special = false
}

resource "google_sql_user" "user" {
    name = "username"
    instance = google_sql_database_instance.instance.name
    password = random_password.pwd.result
}

resource "google_bigquery_connection" "connection" {
    connection_id = "tf-test-my-connection%{random_suffix}"
    location      = "US"
    friendly_name = "👋👋"
    description   = "a very riveting description"
    cloud_sql {
        instance_id = google_sql_database_instance.instance.connection_name
        database    = google_sql_database.db.name
        type        = "MYSQL"
        credential {
            username = google_sql_user.user.name
            password = google_sql_user.user.password
        }
    }
}
`, context)
}

func TestAccBigqueryConnectionConnection_bigqueryConnectionAwsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckBigqueryConnectionConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryConnectionConnection_bigqueryConnectionAws(context),
			},
			{
				ResourceName:            "google_bigquery_connection.connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryConnectionConnection_bigqueryConnectionAwsUpdate(context),
			},
			{
				ResourceName:            "google_bigquery_connection.connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryConnectionConnection_bigqueryConnectionAws(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_connection" "connection" {
   connection_id = "tf-test-my-connection%{random_suffix}"
   location      = "aws-us-east-1"
   friendly_name = "👋"
   description   = "a riveting description"
   aws {
      access_role {
         iam_role_id =  "arn:aws:iam::999999999999:role/omnirole%{random_suffix}"
      }
   }
}
`, context)
}

func testAccBigqueryConnectionConnection_bigqueryConnectionAwsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_connection" "connection" {
   connection_id = "tf-test-my-connection%{random_suffix}"
   location      = "aws-us-east-1"
   friendly_name = "👋"
   description   = "a riveting description"
   aws {
      access_role {
         iam_role_id =  "arn:aws:iam::999999999999:role/omnirole%{random_suffix}update"
      }
   }
}
`, context)
}
