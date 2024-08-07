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

func TestAccNetappVolumeSnapshot_volumeSnapshotCreateExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-1", acctest.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappVolumeSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappVolumeSnapshot_volumeSnapshotCreateExample_full(context),
			},
			{
				ResourceName:            "google_netapp_volume_snapshot.test_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "volume_name", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappVolumeSnapshot_volumeSnapshotCreateExample_update(context),
			},
			{
				ResourceName:            "google_netapp_volume_snapshot.test_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "volume_name", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetappVolumeSnapshot_volumeSnapshotCreateExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-test-pool%{random_suffix}"
  location = "us-west2"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "default" {
  location = google_netapp_storage_pool.default.location
  name = "tf-test-test-volume%{random_suffix}"
  capacity_gib = "100"
  share_name = "tf-test-test-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
}

resource "google_netapp_volume_snapshot" "test_snapshot" {
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

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccNetappVolumeSnapshot_volumeSnapshotCreateExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_storage_pool" "default" {
  name = "tf-test-test-pool%{random_suffix}"
  location = "us-west2"
  service_level = "PREMIUM"
  capacity_gib = "2048"
  network = data.google_compute_network.default.id
}

resource "google_netapp_volume" "default" {
  location = google_netapp_storage_pool.default.location
  name = "tf-test-test-volume%{random_suffix}"
  capacity_gib = "100"
  share_name = "tf-test-test-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols = ["NFSV3"]
}

resource "google_netapp_volume_snapshot" "test_snapshot" {
  depends_on = [google_netapp_volume.default]
  location = google_netapp_volume.default.location
  volume_name = google_netapp_volume.default.name
  description = "This is a update description"
  name = "testvolumesnap%{random_suffix}"
  labels = {
	key= "test"
	value= "snapshot_update"
  }
}

resource "google_netapp_volume_snapshot" "test_snapshot2" {
	depends_on = [google_netapp_volume.default]
	location = google_netapp_volume.default.location
	volume_name = google_netapp_volume.default.name
	description = "This is a update description"
	name = "testvolumesnap2%{random_suffix}"
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}
