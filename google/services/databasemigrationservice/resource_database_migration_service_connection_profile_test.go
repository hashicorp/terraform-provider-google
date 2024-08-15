// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package databasemigrationservice_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseMigrationServiceConnectionProfile_update(t *testing.T) {
	t.Parallel()

	suffix := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_basic(suffix),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql.0.password", "labels", "terraform_labels"},
			},
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_update(suffix),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql.0.password", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConnectionProfile_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_connection_profile" "default" {
	location = "us-central1"
	connection_profile_id = "tf-test-dbms-connection-profile%{random_suffix}"
	display_name          = "tf-test-dbms-connection-profile-display%{random_suffix}"
	labels	= { 
		foo = "bar" 
	}
	mysql {
	  host = "10.20.30.40"
	  port = 3306
	  username = "tf-test-dbms-test-user%{random_suffix}"
	  password = "tf-test-dbms-test-pass%{random_suffix}"
	}
}
`, context)
}

func testAccDatabaseMigrationServiceConnectionProfile_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_database_migration_service_connection_profile" "default" {
	location = "us-central1"
	connection_profile_id = "tf-test-dbms-connection-profile%{random_suffix}"
	display_name          = "tf-test-dbms-connection-profile-updated-display%{random_suffix}"
	labels	= { 
		bar = "foo" 
	}
	mysql {
	  host = "10.20.30.50"
	  port = 3306
	  username = "tf-test-update-dbms-test-user%{random_suffix}"
	  password = "tf-test-update-dbms-test-pass%{random_suffix}"
	}
}
`, context)
}

func TestAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileAlloydb(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "vpc-network-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatabaseMigrationServiceConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileAlloydb(context),
			},
			{
				ResourceName:            "google_database_migration_service_connection_profile.alloydbprofile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "alloydb.0.settings.0.initial_user.0.password", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccDatabaseMigrationServiceConnectionProfile_databaseMigrationServiceConnectionProfileAlloydb(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_database_migration_service_connection_profile" "alloydbprofile" {
  location = "us-central1"
  connection_profile_id = "tf-test-my-profileid%{random_suffix}"
  display_name = "tf-test-my-profileid%{random_suffix}_display"
  labels = { 
    foo = "bar" 
  }
  alloydb {
    cluster_id = "tf-test-dbmsalloycluster%{random_suffix}"
    settings {
      initial_user {
        user = "alloyuser%{random_suffix}"
        password = "alloypass%{random_suffix}"
      }
      vpc_network = data.google_compute_network.default.id
      labels  = { 
        alloyfoo = "alloybar" 
      }
      primary_instance_settings {
        id = "priminstid"
        machine_config {
          cpu_count = 2
        }
        database_flags = { 
        }
        labels = { 
          alloysinstfoo = "allowinstbar" 
        }
      }
    }
  }
}
`, context)
}
