// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package cloudrunv2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCloudRunV2Service_cloudrunv2ServiceBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceBasicExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  ingress = "INGRESS_TRAFFIC_ALL"
  
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
  }
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceLimitsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceLimitsExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceLimitsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      resources {
        limits = {
          cpu    = "2"
          memory = "1024Mi"
        }
      }
    }
  }
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceSqlExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceSqlExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceSqlExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  ingress = "INGRESS_TRAFFIC_ALL"
  
  template {
    scaling {
      max_instance_count = 2
    }
  
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.instance.connection_name]
      }
    }

    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"

      env {
        name = "FOO"
        value = "bar"
      }
      env {
        name = "SECRET_ENV_VAR"
        value_source {
          secret_key_ref {
            secret = google_secret_manager_secret.secret.secret_id
            version = "1"
          }
        }
      }
      volume_mounts {
        name = "cloudsql"
        mount_path = "/cloudsql"
      }
    }
  }

  traffic {
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
  depends_on = [google_secret_manager_secret_version.secret-version-data]
}

data "google_project" "project" {
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "tf-test-secret-1%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-data" {
  secret = google_secret_manager_secret.secret.name
  secret_data = "secret-data"
}

resource "google_secret_manager_secret_iam_member" "secret-access" {
  secret_id = google_secret_manager_secret.secret.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret.secret]
}

resource "google_sql_database_instance" "instance" {
  name             = "tf-test-cloudrun-sql%{random_suffix}"
  region           = "us-central1"
  database_version = "MYSQL_5_7"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection  = "%{deletion_protection}"
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceVpcaccessExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceVpcaccessExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceVpcaccessExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
    vpc_access{
      connector = google_vpc_access_connector.connector.id
      egress = "ALL_TRAFFIC"
    }
  }
}

resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-run-vpc%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
  region        = "us-central1"
}
resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-run-subnetwork%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}
resource "google_compute_network" "custom_test" {
  name                    = "tf-test-run-network%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceDirectvpcExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceDirectvpcExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceDirectvpcExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  launch_stage = "GA"
  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
    vpc_access{
      network_interfaces {
        network = "default"
        subnetwork = "default"
        tags = ["tag1", "tag2", "tag3"]
      }
    }
  }
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceProbesExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceProbesExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceProbesExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      startup_probe {
        initial_delay_seconds = 0
        timeout_seconds = 1
        period_seconds = 3
        failure_threshold = 1
        tcp_socket {
          port = 8080
        }
      }
      liveness_probe {
        http_get {
          path = "/"
        }
      }
    }
  }
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceSecretExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceSecretExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceSecretExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1"
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    volumes {
      name = "a-volume"
      secret {
        secret = google_secret_manager_secret.secret.secret_id
        default_mode = 292 # 0444
        items {
          version = "1"
          path = "my-secret"
        }
      }
    }
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      volume_mounts {
        name = "a-volume"
        mount_path = "/secrets"
      }
    }
  }
  depends_on = [google_secret_manager_secret_version.secret-version-data]
}

data "google_project" "project" {
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "tf-test-secret-1%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-data" {
  secret = google_secret_manager_secret.secret.name
  secret_data = "secret-data"
}

resource "google_secret_manager_secret_iam_member" "secret-access" {
  secret_id = google_secret_manager_secret.secret.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret.secret]
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceMountGcsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceMountGcsExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceMountGcsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"

  location     = "us-central1"
  launch_stage = "BETA"

  template {
    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"

    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      volume_mounts {
        name       = "bucket"
        mount_path = "/var/www"
      }
    }

    volumes {
      name = "bucket"
      gcs {
        bucket    = google_storage_bucket.default.name
        read_only = false
      }
    }
  }
}

resource "google_storage_bucket" "default" {
    name     = "tf-test-cloudrun-service%{random_suffix}"
    location = "US"
}
`, context)
}

func TestAccCloudRunV2Service_cloudrunv2ServiceMountNfsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunV2ServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunV2Service_cloudrunv2ServiceMountNfsExample(context),
			},
			{
				ResourceName:            "google_cloud_run_v2_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "labels", "annotations", "terraform_labels"},
			},
		},
	})
}

func testAccCloudRunV2Service_cloudrunv2ServiceMountNfsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_v2_service" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"

  location     = "us-central1"
  ingress      = "INGRESS_TRAFFIC_ALL"
  launch_stage = "BETA"

  template {
    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello:latest"
      volume_mounts {
        name       = "nfs"
        mount_path = "/mnt/nfs/filestore"
      }
    }
    vpc_access {
      network_interfaces {
        network    = "default"
        subnetwork = "default"
      }
    }

    volumes {
      name = "nfs"
      nfs {
        server    = google_filestore_instance.default.networks[0].ip_addresses[0]
        path      = "/share1"
        read_only = false
      }
    }
  }
}

resource "google_filestore_instance" "default" {
  name     = "tf-test-cloudrun-service%{random_suffix}"
  location = "us-central1-b"
  tier     = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
`, context)
}

func testAccCheckCloudRunV2ServiceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_cloud_run_v2_service" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{CloudRunV2BasePath}}projects/{{project}}/locations/{{location}}/services/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("CloudRunV2Service still exists at %s", url)
			}
		}

		return nil
	}
}
