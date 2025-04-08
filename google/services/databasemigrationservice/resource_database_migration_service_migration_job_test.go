// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package databasemigrationservice_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatabaseMigrationServiceMigrationJob_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatabaseMigrationServiceMigrationJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceMigrationJob_full(context),
			},
			{
				ResourceName:            "google_database_migration_service_migration_job.mysqltomysql",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "migration_job_id", "terraform_labels"},
			},
			{
				Config: testAccDatabaseMigrationServiceMigrationJob_update(context),
			},
			{
				ResourceName:            "google_database_migration_service_migration_job.mysqltomysql",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "migration_job_id", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceMigrationJob_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_sql_database_instance" "source_csql" {
  name             = "tf-test-source-csql%{random_suffix}"
  database_version = "MYSQL_5_7"
  settings {
    tier = "db-n1-standard-1"
    deletion_protection_enabled = false
  }
  deletion_protection = false
}

resource "google_sql_ssl_cert" "source_sql_client_cert" {
  common_name = "cert%{random_suffix}"
  instance    = google_sql_database_instance.source_csql.name

  depends_on = [google_sql_database_instance.source_csql]
}

resource "google_sql_user" "source_sqldb_user" {
  name     = "username%{random_suffix}"
  instance = google_sql_database_instance.source_csql.name
  password = "password%{random_suffix}"

  depends_on = [google_sql_ssl_cert.source_sql_client_cert]
}

resource "google_database_migration_service_connection_profile" "source_cp" {
  location              = "us-central1"
  connection_profile_id = "tf-test-source-cp%{random_suffix}"
  display_name          = "tf-test-source-cp%{random_suffix}_display"
  labels = {
    foo = "bar"
  }
  mysql {
    host     = google_sql_database_instance.source_csql.ip_address.0.ip_address
    port     = 3306
    username = google_sql_user.source_sqldb_user.name
    password = google_sql_user.source_sqldb_user.password
    ssl {
      client_key         = google_sql_ssl_cert.source_sql_client_cert.private_key
      client_certificate = google_sql_ssl_cert.source_sql_client_cert.cert
      ca_certificate     = google_sql_ssl_cert.source_sql_client_cert.server_ca_cert
      type = "SERVER_CLIENT"
    }
    cloud_sql_id = "tf-test-source-csql%{random_suffix}"
  }

  depends_on = [google_sql_user.source_sqldb_user]
}

resource "google_sql_database_instance" "destination_csql" {
  name             = "tf-test-destination-csql%{random_suffix}"
  database_version = "MYSQL_5_7"
  settings {
    tier = "db-n1-standard-1"
    deletion_protection_enabled = false
  }
  deletion_protection = false
}

resource "google_database_migration_service_connection_profile" "destination_cp" {
  location              = "us-central1"
  connection_profile_id = "tf-test-destination-cp%{random_suffix}"
  display_name          = "tf-test-destination-cp%{random_suffix}_display"
  labels = {
    foo = "bar"
  }
  mysql {
    cloud_sql_id = "tf-test-destination-csql%{random_suffix}"
  }
  depends_on = [google_sql_database_instance.destination_csql]
}

resource "google_compute_network" "default" {
  name = "tf-test-destination-csql%{random_suffix}"
}

resource "google_database_migration_service_migration_job" "mysqltomysql" {
  location              = "us-central1"
  migration_job_id = "tf-test-my-migrationid%{random_suffix}"
  display_name = "tf-test-my-migrationid%{random_suffix}_display"
  labels = {
    foo = "bar"
  }
  performance_config {
    dump_parallel_level = "MAX"
  }
  vpc_peering_connectivity {
    vpc = google_compute_network.default.id
  }
  dump_type = "LOGICAL"
  dump_flags {
    dump_flags {
      name = "max-allowed-packet"
      value = "1073741824"
    }
  }
  source          = google_database_migration_service_connection_profile.source_cp.name
  destination     = google_database_migration_service_connection_profile.destination_cp.name
  type            = "CONTINUOUS"
}
`, context)
}

func testAccDatabaseMigrationServiceMigrationJob_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_sql_database_instance" "source_csql" {
  name             = "tf-test-source-csql%{random_suffix}"
  database_version = "MYSQL_5_7"
  settings {
    tier = "db-n1-standard-1"
    deletion_protection_enabled = false
  }
  deletion_protection = false
}

resource "google_sql_ssl_cert" "source_sql_client_cert" {
  common_name = "cert%{random_suffix}"
  instance    = google_sql_database_instance.source_csql.name

  depends_on = [google_sql_database_instance.source_csql]
}

resource "google_sql_user" "source_sqldb_user" {
  name     = "username%{random_suffix}"
  instance = google_sql_database_instance.source_csql.name
  password = "password%{random_suffix}"

  depends_on = [google_sql_ssl_cert.source_sql_client_cert]
}

resource "google_database_migration_service_connection_profile" "source_cp" {
  location              = "us-central1"
  connection_profile_id = "tf-test-source-cp%{random_suffix}"
  display_name          = "tf-test-source-cp%{random_suffix}_display"
  labels = {
    foo = "bar"
  }
  mysql {
    host     = google_sql_database_instance.source_csql.ip_address.0.ip_address
    port     = 3306
    username = google_sql_user.source_sqldb_user.name
    password = google_sql_user.source_sqldb_user.password
    ssl {
      client_key         = google_sql_ssl_cert.source_sql_client_cert.private_key
      client_certificate = google_sql_ssl_cert.source_sql_client_cert.cert
      ca_certificate     = google_sql_ssl_cert.source_sql_client_cert.server_ca_cert
      type = "SERVER_CLIENT"
    }
    cloud_sql_id = "tf-test-source-csql%{random_suffix}"
  }

  depends_on = [google_sql_user.source_sqldb_user]
}

resource "google_sql_database_instance" "destination_csql" {
  name             = "tf-test-destination-csql%{random_suffix}"
  database_version = "MYSQL_5_7"
  settings {
    tier = "db-n1-standard-1"
    deletion_protection_enabled = false
  }
  deletion_protection = false
}

resource "google_database_migration_service_connection_profile" "destination_cp" {
  location              = "us-central1"
  connection_profile_id = "tf-test-destination-cp%{random_suffix}"
  display_name          = "tf-test-destination-cp%{random_suffix}_display"
  labels = {
    foo = "bar"
  }
  mysql {
    cloud_sql_id = "tf-test-destination-csql%{random_suffix}"
  }
  depends_on = [google_sql_database_instance.destination_csql]
}

resource "google_compute_network" "default" {
  name = "tf-test-destination-csql%{random_suffix}"
}

resource "google_database_migration_service_migration_job" "mysqltomysql" {
  location              = "us-central1"
  migration_job_id = "tf-test-my-migrationid%{random_suffix}"
  display_name = "tf-test-my-migrationid%{random_suffix}_display"
  labels = {
    foo = "bar"
  }
  performance_config {
    dump_parallel_level = "MIN"
  }
  static_ip_connectivity {
  }
  dump_type = "LOGICAL"
  dump_flags {
    dump_flags {
      name = "max-allowed-packet"
      value = "1231231234"
    }
  }
  source          = google_database_migration_service_connection_profile.source_cp.name
  destination     = google_database_migration_service_connection_profile.destination_cp.name
  type            = "ONE_TIME"
}
`, context)
}
