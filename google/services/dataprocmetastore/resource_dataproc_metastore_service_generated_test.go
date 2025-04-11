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

package dataprocmetastore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDataprocMetastoreService_dataprocMetastoreServiceBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceBasicExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "default" {
  service_id = "tf-test-metastore-srv%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }

  hive_metastore_config {
    version = "2.3.6"
  }

  labels = {
    env = "test"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceDeletionProtectionExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceDeletionProtectionExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceDeletionProtectionExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "default" {
    service_id          = "tf-test-metastore-srv%{random_suffix}"
    location            = "us-central1"
    port                = 9080
    tier                = "DEVELOPER"
    deletion_protection = %{deletion_protection}
  
    maintenance_window {
      hour_of_day = 2
      day_of_week = "SUNDAY"
    }
  
    hive_metastore_config {
      version = "2.3.6"
    }
  
    labels = {
      env = "test"
    }
  }
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceCmekTestExample(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceCmekTestExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceCmekTestExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

data "google_storage_project_service_account" "gcs_account" {}


resource "google_dataproc_metastore_service" "default" {
  service_id = "tf-test-example-service%{random_suffix}"
  location   = "us-central1"

  encryption_config {
    kms_key = "tf-test-acctest.BootstrapKMSKeyWithPurposeInLocationAn%{random_suffix}"
  }

  hive_metastore_config {
    version = "3.1.2"
  }

  depends_on = [
    google_kms_crypto_key_iam_member.crypto_key_member_1,
    google_kms_crypto_key_iam_member.crypto_key_member_2,
  ]
}

resource "google_kms_crypto_key_iam_member" "crypto_key_member_1" {
  crypto_key_id = "tf-test-acctest.BootstrapKMSKeyWithPurposeInLocationAn%{random_suffix}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-metastore.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "crypto_key_member_2" {
  crypto_key_id = "tf-test-acctest.BootstrapKMSKeyWithPurposeInLocationAn%{random_suffix}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceEndpointExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceEndpointExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.endpoint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceEndpointExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "endpoint" {
  service_id = "tf-test-metastore-endpoint%{random_suffix}"
  location   = "us-central1"
  tier       = "DEVELOPER"

  hive_metastore_config {
    version           = "3.1.2"
    endpoint_protocol = "GRPC"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceAuxExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceAuxExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.aux",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceAuxExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "aux" {
  service_id = "tf-test-metastore-aux%{random_suffix}"
  location   = "us-central1"
  tier       = "DEVELOPER"

  hive_metastore_config {
    version = "3.1.2"
    auxiliary_versions {
      key     = "aux-test"
      version = "2.3.6"
    }
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceMetadataExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceMetadataExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.metadata",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceMetadataExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "metadata" {
  service_id = "tf-test-metastore-metadata%{random_suffix}"
  location   = "us-central1"
  tier       = "DEVELOPER"

  metadata_integration {
    data_catalog_config {
      enabled = true
    }
  }

  hive_metastore_config {
    version = "3.1.2"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceTelemetryExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceTelemetryExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.telemetry",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceTelemetryExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "telemetry" {
  service_id = "tf-test-ms-telemetry%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  hive_metastore_config {
    version = "3.1.2"
  }

  telemetry_config {
    log_format = "LEGACY"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceDpms2Example(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceDpms2Example(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.dpms2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceDpms2Example(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "dpms2" {
  service_id = "tf-test-ms-dpms2%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    instance_size = "EXTRA_SMALL"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceDpms2ScalingFactorExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceDpms2ScalingFactorExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.dpms2_scaling_factor",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceDpms2ScalingFactorExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "dpms2_scaling_factor" {
  service_id = "tf-test-ms-dpms2sf%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    scaling_factor = "2"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceDpms2ScalingFactorLt1Example(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceDpms2ScalingFactorLt1Example(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.dpms2_scaling_factor_lt1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceDpms2ScalingFactorLt1Example(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "dpms2_scaling_factor_lt1" {
  service_id = "tf-test-ms-dpms2sflt1%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    scaling_factor = "0.1"
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "backup" {
  service_id = "backup%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }

  hive_metastore_config {
    version = "2.3.6"
  }

  scheduled_backup {
    enabled         = true
    cron_schedule   = "0 0 * * *"
    time_zone       = "UTC"
    backup_location = "gs://${google_storage_bucket.bucket.name}"
  }

  labels = {
    env = "test"
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "backup%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMaxScalingFactorExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMaxScalingFactorExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.test_resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMaxScalingFactorExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "test_resource" {
  service_id = "tf-test-test-service%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    autoscaling_config {
      autoscaling_enabled = true
      limit_config {
        max_scaling_factor = 1.0
      }
    }
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMinAndMaxScalingFactorExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMinAndMaxScalingFactorExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.test_resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMinAndMaxScalingFactorExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "test_resource" {
  service_id = "tf-test-test-service%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    autoscaling_config {
      autoscaling_enabled = true
      limit_config {
        min_scaling_factor = 0.1
        max_scaling_factor = 1.0
      }
    }
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMinScalingFactorExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMinScalingFactorExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.test_resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingMinScalingFactorExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "test_resource" {
  service_id = "tf-test-test-service%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    autoscaling_config {
      autoscaling_enabled = true
      limit_config {
        min_scaling_factor = 0.1
      }
    }
  }
}
`, context)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingNoLimitConfigExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingNoLimitConfigExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.test_resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service_id", "tags", "terraform_labels"},
			},
		},
	})
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceAutoscalingNoLimitConfigExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "test_resource" {
  service_id = "tf-test-test-service%{random_suffix}"
  location   = "us-central1"

  # DPMS 2 requires SPANNER database type, and does not require
  # a maintenance window.
  database_type = "SPANNER"

  hive_metastore_config {
    version           = "3.1.2"
  }

  scaling_config {
    autoscaling_config {
      autoscaling_enabled = true
    }
  }
}
`, context)
}

func testAccCheckDataprocMetastoreServiceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_metastore_service" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DataprocMetastoreBasePath}}projects/{{project}}/locations/{{location}}/services/{{service_id}}")
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
				return fmt.Errorf("DataprocMetastoreService still exists at %s", url)
			}
		}

		return nil
	}
}
