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

func TestAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-1", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappVolumeReplicationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_volume_replication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "delete_destination_volume", "replication_enabled", "force_stopping", "wait_for_mirror", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_stop(context),
			},
			{
				ResourceName:            "google_netapp_volume_replication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "delete_destination_volume", "replication_enabled", "force_stopping", "wait_for_mirror", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_resume(context),
			},
			{
				ResourceName:            "google_netapp_volume_replication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "delete_destination_volume", "replication_enabled", "force_stopping", "wait_for_mirror", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_update(context),
			},
			{
				ResourceName:            "google_netapp_volume_replication.test_replication",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"destination_volume_parameters", "location", "volume_name", "name", "delete_destination_volume", "replication_enabled", "force_stopping", "wait_for_mirror", "labels", "terraform_labels"},
			},
		},
	})
}

// Basic replication
func testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
  allow_auto_tiering = true
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
  deletion_policy = "FORCE"
}

resource "google_netapp_volume_replication" "test_replication" {
  depends_on           = [google_netapp_volume.source_volume]
  location             = google_netapp_volume.source_volume.location
  volume_name          = google_netapp_volume.source_volume.name
  name                 = "tf-test-test-replication%{random_suffix}"
  replication_schedule = "EVERY_10_MINUTES"
  destination_volume_parameters {
    storage_pool = google_netapp_storage_pool.destination_pool.id
    volume_id    = "tf-test-destination-volume%{random_suffix}"
    # Keeping the share_name of source and destination the same
    # simplifies implementing client failover concepts
    share_name  = "tf-test-source-volume%{random_suffix}"
    description = "This is a replicated volume"
    tiering_policy {
      cooling_threshold_days = 20
      tier_action = "ENABLED"
    }
  }
  delete_destination_volume = true
  wait_for_mirror = true
}
`, context)
}

// Update parameters
func testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
  allow_auto_tiering = true
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
  deletion_policy = "FORCE"
}

resource "google_netapp_volume_replication" "test_replication" {
  depends_on           = [google_netapp_volume.source_volume]
  location             = google_netapp_volume.source_volume.location
  volume_name          = google_netapp_volume.source_volume.name
  name                 = "tf-test-test-replication%{random_suffix}"
  replication_schedule = "EVERY_10_MINUTES"
  description          = "This is a replication resource"
  labels = {
    key   = "test"
    value =  "replication"
  }
  destination_volume_parameters {
    storage_pool = google_netapp_storage_pool.destination_pool.id
    volume_id    = "tf-test-destination-volume%{random_suffix}"
    # Keeping the share_name of source and destination the same
    # simplifies implementing client failover concepts
    share_name  = "tf-test-source-volume%{random_suffix}"
    description = "This is a replicated volume"
    tiering_policy {
      cooling_threshold_days = 20
      tier_action = "ENABLED"
    }
  }
  replication_enabled = true
  delete_destination_volume = true
  force_stopping = true
  wait_for_mirror = true
}
`, context)
}

// Stop replication
func testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_stop(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
  allow_auto_tiering = true
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
  deletion_policy = "FORCE"
}

resource "google_netapp_volume_replication" "test_replication" {
  depends_on           = [google_netapp_volume.source_volume]
  location             = google_netapp_volume.source_volume.location
  volume_name          = google_netapp_volume.source_volume.name
  name                 = "tf-test-test-replication%{random_suffix}"
  replication_schedule = "EVERY_10_MINUTES"
  description          = "This is a replication resource"
  labels = {
    key   = "test"
    value =  "replication2"
  }
  destination_volume_parameters {
    storage_pool = google_netapp_storage_pool.destination_pool.id
    volume_id    = "tf-test-destination-volume%{random_suffix}"
    # Keeping the share_name of source and destination the same
    # simplifies implementing client failover concepts
    share_name  = "tf-test-source-volume%{random_suffix}"
    description = "This is a replicated volume"
    tiering_policy {
      cooling_threshold_days = 20
      tier_action = "ENABLED"
    }
  }
  replication_enabled = false
  delete_destination_volume = true
  force_stopping = true
  wait_for_mirror = true
}
`, context)
}

// resume replication
func testAccNetappVolumeReplication_NetappVolumeReplicationCreateExample_resume(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "source_pool" {
  name          = "tf-test-source-pool%{random_suffix}"
  location      = "us-central1"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
}

resource "google_netapp_storage_pool" "destination_pool" {
  name          = "tf-test-destination-pool%{random_suffix}"
  location      = "us-west2"
  service_level = "PREMIUM"
  capacity_gib  = 2048
  network       = data.google_compute_network.default.id
  allow_auto_tiering = true
}

resource "google_netapp_volume" "source_volume" {
  location     = google_netapp_storage_pool.source_pool.location
  name         = "tf-test-source-volume%{random_suffix}"
  capacity_gib = 100
  share_name   = "tf-test-source-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.source_pool.name
  protocols = [
    "NFSV3"
  ]
  deletion_policy = "FORCE"
}

resource "google_netapp_volume_replication" "test_replication" {
  depends_on           = [google_netapp_volume.source_volume]
  location             = google_netapp_volume.source_volume.location
  volume_name          = google_netapp_volume.source_volume.name
  name                 = "tf-test-test-replication%{random_suffix}"
  replication_schedule = "HOURLY"
  description          = "This is a replication resource"
  labels = {
    key   = "test"
    value =  "replication2"
  }
  destination_volume_parameters {
    storage_pool = google_netapp_storage_pool.destination_pool.id
    volume_id    = "tf-test-destination-volume%{random_suffix}"
    # Keeping the share_name of source and destination the same
    # simplifies implementing client failover concepts
    share_name  = "tf-test-source-volume%{random_suffix}"
    description = "This is a replicated volume"
    tiering_policy {
      cooling_threshold_days = 20
      tier_action = "ENABLED"
    }
  }
  replication_enabled = true
  delete_destination_volume = true
  force_stopping = true
  wait_for_mirror = true
}
`, context)
}
