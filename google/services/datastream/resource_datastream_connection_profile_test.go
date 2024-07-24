// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatastreamConnectionProfile_update(t *testing.T) {
	// this test uses the random provider
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	random_pass_1 := acctest.RandString(t, 10)
	random_pass_2 := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckDatastreamConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamConnectionProfile_update(context),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location"},
			},
			{
				Config: testAccDatastreamConnectionProfile_update2(context, true),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "postgresql_profile.0.password"},
			},
			{
				// Disable prevent_destroy
				Config: testAccDatastreamConnectionProfile_update2(context, false),
			},
			{
				Config: testAccDatastreamConnectionProfile_mySQLUpdate(context, true, random_pass_1),
			},
			{
				ResourceName:            "google_datastream_connection_profile.mysql_con_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql_profile.0.password"},
			},
			{
				// run once more to update the password. it should update it in-place
				Config: testAccDatastreamConnectionProfile_mySQLUpdate(context, true, random_pass_2),
			},
			{
				ResourceName:            "google_datastream_connection_profile.mysql_con_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "mysql_profile.0.password"},
			},
			{
				// Disable prevent_destroy
				Config: testAccDatastreamConnectionProfile_mySQLUpdate(context, false, random_pass_2),
			},
		},
	})
}

func testAccDatastreamConnectionProfile_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_datastream_connection_profile" "default" {
	display_name          = "Connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-my-profile%{random_suffix}"

	gcs_profile {
		bucket    = "my-bucket"
		root_path = "/path"
	}
	lifecycle {
		prevent_destroy = true
	}
}
`, context)
}

func testAccDatastreamConnectionProfile_update2(context map[string]interface{}, preventDestroy bool) string {
	context["lifecycle_block"] = ""
	if preventDestroy {
		context["lifecycle_block"] = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
    name             = "tf-test-my-database-instance%{random_suffix}"
    database_version = "POSTGRES_14"
    region           = "us-central1"
    settings {
      tier = "db-f1-micro"

      ip_configuration {

        // Datastream IPs will vary by region.
        authorized_networks {
            value = "34.71.242.81"
        }

        authorized_networks {
            value = "34.72.28.29"
        }

        authorized_networks {
            value = "34.67.6.157"
        }

        authorized_networks {
            value = "34.67.234.134"
        }

        authorized_networks {
            value = "34.72.239.218"
        }
      }
    }

    deletion_protection  = "false"
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
    name = "user"
    instance = google_sql_database_instance.instance.name
    password = random_password.pwd.result
}

resource "google_datastream_connection_profile" "default" {
	display_name          = "Connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-my-profile%{random_suffix}"

	postgresql_profile {
		hostname = google_sql_database_instance.instance.public_ip_address
		username = google_sql_user.user.name
		password = google_sql_user.user.password
		database = google_sql_database.db.name
	}
	%{lifecycle_block}
}
`, context)
}

func testAccDatastreamConnectionProfile_mySQLUpdate(context map[string]interface{}, preventDestroy bool, password string) string {
	context["lifecycle_block"] = ""
	if preventDestroy {
		context["lifecycle_block"] = `
		lifecycle {
			prevent_destroy = true
		}`
	}

	context["password"] = password

	return acctest.Nprintf(`
resource "google_sql_database_instance" "mysql_instance" {
    name             = "tf-test-mysql-database-instance%{random_suffix}"
    database_version = "MYSQL_8_0"
    region           = "us-central1"
    settings {
      tier = "db-f1-micro"
        backup_configuration {
            enabled            = true
            binary_log_enabled = true
        }

      ip_configuration {

        // Datastream IPs will vary by region.
        authorized_networks {
            value = "34.71.242.81"
        }

        authorized_networks {
            value = "34.72.28.29"
        }

        authorized_networks {
            value = "34.67.6.157"
        }

        authorized_networks {
            value = "34.67.234.134"
        }

        authorized_networks {
            value = "34.72.239.218"
        }
      }
    }

    deletion_protection  = "false"
}

resource "google_sql_database" "mysql_db" {
    instance = google_sql_database_instance.mysql_instance.name
    name     = "db"
}

resource "google_sql_user" "mysql_user" {
    name = "user"
    instance = google_sql_database_instance.mysql_instance.name
    host     = "%"
    password = "%{password}"
}

resource "google_datastream_connection_profile" "mysql_con_profile" {
    display_name          = "Source connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-mysql-profile%{random_suffix}"

    mysql_profile {
		hostname = google_sql_database_instance.mysql_instance.public_ip_address
		username = google_sql_user.mysql_user.name
		password = google_sql_user.mysql_user.password
	}
	%{lifecycle_block}
}
`, context)
}
