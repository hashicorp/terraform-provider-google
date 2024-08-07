// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappbackup_netappBackupFull_update(t *testing.T) {
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-1", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappbackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappbackup_netappBackupFromVolumeSnapshot(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
			{
				Config: testAccNetappbackup_netappBackupFromVolumeSnapshot_update(context),
			},
			{
				ResourceName:            "google_netapp_backup.test_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels", "vault_name"},
			},
		},
	})
}

func testAccNetappbackup_netappBackupFromVolumeSnapshot(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-central1"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}

resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = google_netapp_storage_pool.default.location
}

resource "google_netapp_volume_snapshot" "default" {
	depends_on = [google_netapp_volume.default]
	location = google_netapp_volume.default.location
	volume_name = google_netapp_volume.default.name
	description = "This is a test description"
	name = "testvolumesnap%{random_suffix}"
	labels = {
	  key= "test"
	  value= "snapshot"
	}
  }

resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test backup"
  source_volume = google_netapp_volume.default.id
  location = google_netapp_backup_vault.default.location
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
	key= "test"
	value= "backup"
  }
}
`, context)
}

func testAccNetappbackup_netappBackupFromVolumeSnapshot_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "default" {
  name = "tf-test-backup-pool%{random_suffix}"
  location = "us-central1"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "default" {
  name = "tf-test-backup-volume%{random_suffix}"
  location = google_netapp_storage_pool.default.location
  capacity_gib = "100"
  share_name = "tf-test-backup-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
  deletion_policy = "FORCE"
  backup_config {
    backup_vault = google_netapp_backup_vault.default.id
  }
}

resource "google_netapp_backup_vault" "default" {
  name = "tf-test-backup-vault%{random_suffix}"
  location = google_netapp_storage_pool.default.location
}

resource "google_netapp_volume_snapshot" "default" {
	depends_on = [google_netapp_volume.default]
	location = google_netapp_volume.default.location
	volume_name = google_netapp_volume.default.name
	description = "This is a test description"
	name = "testvolumesnap%{random_suffix}"
	labels = {
	  key= "test"
	  value= "snapshot"
	}
  }

resource "google_netapp_backup" "test_backup" {
  name = "tf-test-test-backup%{random_suffix}"
  description = "This is a test backup"
  source_volume = google_netapp_volume.default.id
  location = google_netapp_backup_vault.default.location
  vault_name = google_netapp_backup_vault.default.name
  source_snapshot = google_netapp_volume_snapshot.default.id
  labels = {
	key= "test_update"
	value= "backup_update"
  }
}
`, context)
}
