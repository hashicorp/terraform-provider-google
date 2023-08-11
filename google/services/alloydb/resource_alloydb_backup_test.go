// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAlloydbBackup_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "alloydb-update"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbBackup_alloydbBackupFullExample(context),
			},
			{
				ResourceName:            "google_alloydb_backup.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_id", "location", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbBackup_update(context),
			},
			{
				ResourceName:            "google_alloydb_backup.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_id", "location", "reconciling", "update_time"},
			},
		},
	})
}

// Updates "label" field from testAccAlloydbBackup_alloydbBackupFullExample
func testAccAlloydbBackup_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.default.name

  description = "example description"
  labels = {
    "label" = "updated_key"
    "label2" = "updated_key2"
  }
  depends_on = [google_alloydb_instance.default]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// Test to create on-demand backup with mandatory fields.
func TestAccAlloydbBackup_createBackupWithMandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "alloydbbackup-mandatory"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbBackup_createBackupWithMandatoryFields(context),
			},
		},
	})
}

func testAccAlloydbBackup_createBackupWithMandatoryFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_backup" "default" {
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  location = "us-central1"
  cluster_name = google_alloydb_cluster.default.name
  depends_on = [google_alloydb_instance.default]
}

resource "google_alloydb_cluster" "default" {
  location = "us-central1"
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  network    = data.google_compute_network.default.id
}

data "google_project" "project" { }

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
  lifecycle {
	ignore_changes = [
		address,
		creation_timestamp,
		id,
		network,
		project,
		self_link
	]
  }
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

func TestAccAlloydbBackup_usingCMEK(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "alloydb-cmek"),
		"random_suffix": acctest.RandString(t, 10),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbBackup_usingCMEK(context),
			},
			{
				ResourceName:            "google_alloydb_backup.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_id", "location", "reconciling", "update_time"},
			},
		},
	})
}

func testAccAlloydbBackup_usingCMEK(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_backup" "default" {
	location     = "us-central1"
	backup_id    = "tf-test-alloydb-backup%{random_suffix}"
	cluster_name = google_alloydb_cluster.default.name
	description = "example description"
	labels = {
		"label" = "updated_key"
		"label2" = "updated_key2"
	}
	encryption_config {
		kms_key_name = google_kms_crypto_key.key.id
	}
	depends_on = [google_alloydb_instance.default]
}
	  
resource "google_alloydb_cluster" "default" {
	cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
	location   = "us-central1"
	network    = data.google_compute_network.default.id
}
	  
resource "google_alloydb_instance" "default" {
	cluster       = google_alloydb_cluster.default.name
	instance_id   = "tf-test-alloydb-instance%{random_suffix}"
	instance_type = "PRIMARY"
	  
	depends_on = [google_service_networking_connection.vpc_connection]
}
	  
resource "google_compute_global_address" "private_ip_alloc" {
	name          =  "tf-test-alloydb-cluster%{random_suffix}"
	address_type  = "INTERNAL"
	purpose       = "VPC_PEERING"
	prefix_length = 16
	network       = data.google_compute_network.default.id
}
	  
resource "google_service_networking_connection" "vpc_connection" {
	network                 = data.google_compute_network.default.id
	service                 = "servicenetworking.googleapis.com"
	reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
	  
data "google_compute_network" "default" {
	name = "%{network_name}"
}
data "google_project" "project" {}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}
`, context)
}
