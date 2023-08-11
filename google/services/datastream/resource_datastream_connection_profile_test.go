// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatastreamConnectionProfile_update(t *testing.T) {
	// this test uses the random provider
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

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
