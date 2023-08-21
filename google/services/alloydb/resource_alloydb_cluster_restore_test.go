// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

/*
 * Restore tests are kept separate from other cluster tests because they require an instance and a backup to exist
 */

// Restore tests depend on instances and backups being taken, which can take up to 10 minutes. Since the instance doesn't change in between tests,
// we condense everything into individual test cases.
func TestAccAlloydbCluster_restore(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedTestNetwork(t, "alloydbinstance-restore"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbClusterAndInstanceAndBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.source",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				// Invalid input check - cannot pass in both sources
				Config:      testAccAlloydbClusterAndInstanceAndBackup_OnlyOneSourceAllowed(context),
				ExpectError: regexp.MustCompile("\"restore_continuous_backup_source\": conflicts with restore_backup_source"),
			},
			{
				// Invalid input check - both source cluster and point in time are required
				Config:      testAccAlloydbClusterAndInstanceAndBackup_SourceClusterAndPointInTimeRequired(context),
				ExpectError: regexp.MustCompile("The argument \"point_in_time\" is required, but no definition was found."),
			},
			{
				// Validate backup restore succeeds
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.restored_from_backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "restore_backup_source"},
			},
			{
				// Validate PITR succeeds
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.restored_from_point_in_time",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "restore_continuous_backup_source"},
			},
			{
				// Make sure updates work without recreating the clusters
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_AllowedUpdate(context),
			},
			{
				Config:      testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_NotAllowedUpdate(context),
				ExpectError: regexp.MustCompile("the plan calls for this resource to be\ndestroyed"),
			},
			{
				Config: testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_AllowDestroy(context),
			},
		},
	})
}

func testAccAlloydbClusterAndInstanceAndBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

// Cannot restore if multiple sources are present
func testAccAlloydbClusterAndInstanceAndBackup_OnlyOneSourceAllowed(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored" {
  cluster_id             = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_backup_source {
    backup_name = google_alloydb_backup.default.name
  }
  restore_continuous_backup_source {
    cluster = google_alloydb_cluster.source.name
    point_in_time = google_alloydb_backup.default.update_time
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

// Cannot restore if multiple sources are present
func testAccAlloydbClusterAndInstanceAndBackup_SourceClusterAndPointInTimeRequired(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored" {
  cluster_id             = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id

  restore_continuous_backup_source {
    cluster = google_alloydb_cluster.source.name
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_backup_source {
    backup_name = google_alloydb_backup.default.name
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

// The source cluster, instance, and backup should all exist prior to this being invoked. Otherwise the PITR restore will not succeed
// due to the time being too early.
func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_backup_source {
    backup_name = google_alloydb_backup.default.name
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_alloydb_cluster" "restored_from_point_in_time" {
  cluster_id             = "tf-test-alloydb-pitr-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_continuous_backup_source {
    cluster = google_alloydb_cluster.source.name
    point_in_time = google_alloydb_backup.default.update_time
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

// This updates the PITR and backup restored clusters by adding a continuous backup configuration. This can be done in place
// and does not require re-creating the cluster from scratch.
func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_AllowedUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_backup_source {
    backup_name = google_alloydb_backup.default.name
  }

  continuous_backup_config {
    enabled              = true
    recovery_window_days = 20
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_alloydb_cluster" "restored_from_point_in_time" {
  cluster_id             = "tf-test-alloydb-pitr-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_continuous_backup_source {
    cluster = google_alloydb_cluster.source.name
    point_in_time = google_alloydb_backup.default.update_time
  }

  continuous_backup_config {
    enabled              = true
    recovery_window_days = 20
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

// Updating the source cluster, point in time, or source backup are not allowed. This type of operation would
// require deleting and recreating the cluster
func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_NotAllowedUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_backup" "default2" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup2-%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_backup_source {
    backup_name = google_alloydb_backup.default2.name
  }

  continuous_backup_config {
    enabled              = true
    recovery_window_days = 20
  }

  lifecycle {
    prevent_destroy = true
  }

  depends_on = [google_alloydb_backup.default2]
}

resource "google_alloydb_cluster" "restored_from_point_in_time" {
  cluster_id             = "tf-test-alloydb-pitr-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_continuous_backup_source {
    cluster = google_alloydb_cluster.restored_from_backup.name
    point_in_time = google_alloydb_backup.default.update_time
  }

  continuous_backup_config {
    enabled              = true
    recovery_window_days = 20
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}

// The source cluster, instance, and backup should all exist prior to this being invoked. Otherwise the PITR restore will not succeed
// due to the time being too early.
func testAccAlloydbClusterAndInstanceAndBackup_RestoredFromBackupAndRestoredFromPointInTime_AllowDestroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "source" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "source" {
  cluster       = google_alloydb_cluster.source.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.source.name

  depends_on = [google_alloydb_instance.source]
}

resource "google_alloydb_cluster" "restored_from_backup" {
  cluster_id            = "tf-test-alloydb-backup-restored-cluster-%{random_suffix}"
  location              = "us-central1"
  network               = data.google_compute_network.default.id
  restore_backup_source {
    backup_name = google_alloydb_backup.default.name
  }
}

resource "google_alloydb_cluster" "restored_from_point_in_time" {
  cluster_id             = "tf-test-alloydb-pitr-restored-cluster-%{random_suffix}"
  location               = "us-central1"
  network                = data.google_compute_network.default.id
  restore_continuous_backup_source {
    cluster = google_alloydb_cluster.source.name
    point_in_time = google_alloydb_backup.default.update_time
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
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
`, context)
}
