// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatastreamStream_update(t *testing.T) {
	// this test uses the random provider
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":       acctest.RandString(t, 10),
		"deletion_protection": false,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckDatastreamStreamDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamStream_datastreamStreamBasicExample(context),
				Check:  resource.TestCheckResourceAttr("google_datastream_stream.default", "state", "NOT_STARTED"),
			},
			{
				ResourceName:            "google_datastream_stream.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stream_id", "location", "desired_state"},
			},
			{
				Config: testAccDatastreamStream_datastreamStreamBasicUpdate(context, "RUNNING", true),
				Check:  resource.TestCheckResourceAttr("google_datastream_stream.default", "state", "RUNNING"),
			},
			{
				ResourceName:            "google_datastream_stream.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stream_id", "location", "desired_state"},
			},
			{
				Config: testAccDatastreamStream_datastreamStreamBasicUpdate(context, "PAUSED", true),
				Check:  resource.TestCheckResourceAttr("google_datastream_stream.default", "state", "PAUSED"),
			},
			{
				ResourceName:            "google_datastream_stream.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stream_id", "location", "desired_state"},
			},
			{
				Config: testAccDatastreamStream_datastreamStreamBasicUpdate(context, "RUNNING", true),
				Check:  resource.TestCheckResourceAttr("google_datastream_stream.default", "state", "RUNNING"),
			},
			{
				ResourceName:            "google_datastream_stream.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stream_id", "location", "desired_state"},
			},
			{
				// Disable prevent_destroy
				Config: testAccDatastreamStream_datastreamStreamBasicUpdate(context, "RUNNING", false),
			},
		},
	})
}

func testAccDatastreamStream_datastreamStreamBasicUpdate(context map[string]interface{}, desiredState string, preventDestroy bool) string {
	context["lifecycle_block"] = ""
	if preventDestroy {
		context["lifecycle_block"] = `
        lifecycle {
            prevent_destroy = true
        }`
	}
	context["desired_state"] = desiredState
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_sql_database_instance" "instance" {
    name             = "tf-test-my-instance%{random_suffix}"
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

    deletion_protection  = %{deletion_protection}
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
    name     = "user"
    instance = google_sql_database_instance.instance.name
    host     = "%"
    password = random_password.pwd.result
}

resource "google_datastream_connection_profile" "source_connection_profile" {
    display_name          = "Source connection profile"
    location              = "us-central1"
    connection_profile_id = "tf-test-source-profile%{random_suffix}"

    mysql_profile {
        hostname = google_sql_database_instance.instance.public_ip_address
        username = google_sql_user.user.name
        password = google_sql_user.user.password
    }
}

resource "google_storage_bucket" "bucket" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "viewer" {
    bucket = google_storage_bucket.bucket.name
    role   = "roles/storage.objectViewer"
    member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-datastream.iam.gserviceaccount.com"
}

resource "google_storage_bucket_iam_member" "creator" {
    bucket = google_storage_bucket.bucket.name
    role   = "roles/storage.objectCreator"
    member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-datastream.iam.gserviceaccount.com"
}

resource "google_storage_bucket_iam_member" "reader" {
    bucket = google_storage_bucket.bucket.name
    role   = "roles/storage.legacyBucketReader"
    member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-datastream.iam.gserviceaccount.com"
}

resource "google_datastream_connection_profile" "destination_connection_profile" {
    display_name          = "Connection profile"
    location              = "us-central1"
    connection_profile_id = "tf-test-destination-profile%{random_suffix}"

    gcs_profile {
        bucket    = google_storage_bucket.bucket.name
        root_path = "/path"
    }
}

resource "google_datastream_stream" "default" {
    stream_id = "tf-test-my-stream%{random_suffix}"
    location = "us-central1"
    display_name = "my stream update"
    desired_state = "%{desired_state}"

    labels = {
    	key = "updated"
    }

    source_config {
        source_connection_profile = google_datastream_connection_profile.source_connection_profile.id

        mysql_source_config {}
    }
    destination_config {
        destination_connection_profile = google_datastream_connection_profile.destination_connection_profile.id
        gcs_destination_config {
            path = "mydata"
            file_rotation_mb = 200
            file_rotation_interval = "60s"
            json_file_format {
                schema_file_format = "NO_SCHEMA_FILE"
                compression = "GZIP"
            }
        }
    }

    backfill_all {
    }
	%{lifecycle_block}
}
`, context)
}
