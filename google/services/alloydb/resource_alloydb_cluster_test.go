// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAlloydbCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_update(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

func testAccAlloydbCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  labels = {
	foo = "bar" 
  }

  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// Test if adding automatedBackupPolicy AND initialUser re-creates the cluster.
// Ideally, cluster shouldn't be re-created. This test will only pass if the cluster
// isn't re-created but updated in-place.
func TestAccAlloydbCluster_addAutomatedBackupPolicyAndInitialUser(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"hour":          23,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withoutInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_withInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

// Test if deleting automatedBackupPolicy AND initialUser re-creates the cluster.
// Ideally, cluster shouldn't be re-created. This test will only pass if the cluster
// isn't re-created but updated in-place.
func TestAccAlloydbCluster_deleteAutomatedBackupPolicyAndInitialUser(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"hour":          23,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_withoutInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

// Test if automatedBackupPolicy properly handles a startTime of 0 (aka midnight). Calling terraform plan
// after creating the cluster should not bring anything up.
func TestAccAlloydbCluster_AutomatedBackupPolicyHandlesMidnight(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"hour":          0,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withInitialUserAndAutomatedBackupPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

func testAccAlloydbCluster_withInitialUserAndAutomatedBackupPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id   = "tf-test-alloydb-cluster%{random_suffix}"
  location     = "us-central1"
  network      = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  initial_user {
    user     = "tf-test-alloydb-cluster%{random_suffix}"
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }

  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true

    weekly_schedule {
      days_of_week = ["MONDAY"]

      start_times {
        hours   = %{hour}
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }

    quantity_based_retention {
      count = 1
    }

    labels = {
      test = "tf-test-alloydb-cluster%{random_suffix}"
    }
  }
  lifecycle {
    prevent_destroy = true
  }  
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func testAccAlloydbCluster_withoutInitialUserAndAutomatedBackupPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  lifecycle {
    prevent_destroy = true
  }  
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// The cluster creation should work fine even without a weekly schedule.
func TestAccAlloydbCluster_missingWeeklySchedule(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_missingWeeklySchedule(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlloydbCluster_missingWeeklySchedule(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    quantity_based_retention {
	  count = 1
	}
    labels = {
	  test = "tf-test-alloydb-cluster%{random_suffix}"
	}
  }
}
data "google_project" "project" {}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// The cluster creation should succeed with minimal number of arguments.
func TestAccAlloydbCluster_mandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

// The cluster creation should succeed with maximal number of arguments.
func TestAccAlloydbCluster_maximumFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_alloydbClusterFullExample(context),
			},
		},
	})
}

// Deletion of time-based retention policy should be an in-place operation
func TestAccAlloydbCluster_deleteTimeBasedRetentionPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withTimeBasedRetentionPolicy(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlloydbCluster_withoutTimeBasedRetentionPolicy(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

func testAccAlloydbCluster_withTimeBasedRetentionPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true

    weekly_schedule {
      days_of_week = ["MONDAY"]

      start_times {
        hours   = 23
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
    time_based_retention {
      retention_period = "4.5s"
    }
  }
  lifecycle {
    ignore_changes = [
      automated_backup_policy[0].time_based_retention
    ]
    prevent_destroy = true
  }
}

data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func testAccAlloydbCluster_withoutTimeBasedRetentionPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true

    weekly_schedule {
      days_of_week = ["MONDAY"]

      start_times {
        hours   = 23
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
  }
  lifecycle {
    ignore_changes = [
      automated_backup_policy[0].time_based_retention
    ]
    prevent_destroy = true
  }
}

data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}
func TestAccAlloydbCluster_usingCMEK(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_usingCMEK(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
		},
	})
}

func testAccAlloydbCluster_usingCMEK(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
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

func TestAccAlloydbCluster_CMEKInAutomatedBackupIsUpdatable(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_usingCMEKInClusterAndAutomatedBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_updateCMEKInAutomatedBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_usingCMEKallowDeletion(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
		},
	})
}

func testAccAlloydbCluster_usingCMEKInClusterAndAutomatedBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = google_kms_crypto_key.key.id
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  lifecycle {
	prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
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

func testAccAlloydbCluster_updateCMEKInAutomatedBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = google_kms_crypto_key.key2.id
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  lifecycle {
	prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
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

resource "google_kms_crypto_key" "key2" {
	name     = "%{key_name}-2"
	key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key2" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	members = [
	  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
	]
}
`, context)
}

func testAccAlloydbCluster_usingCMEKallowDeletion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = google_kms_crypto_key.key2.id
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
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

resource "google_kms_crypto_key" "key2" {
	name     = "%{key_name}-2"
	key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
	"serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key2" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	members = [
	  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
	]
}
`, context)
}

// Validates continuous backups defaults to being enabled with 14d retention, even if not explicitly configured.
func TestAccAlloydbCluster_continuousBackup_enabledByDefault(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withoutContinuousBackupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

// Continuous backups defaults to being enabled with 14d retention. If the same configuration is set explicitly, terraform plan
// should return no changes.
func TestAccAlloydbCluster_continuousBackup_update_noChangeIfDefaultsSet(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"enabled":              true,
		"recovery_window_days": 14,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withoutContinuousBackupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_continuousBackupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

// This test ensures that if you start with a terraform configuration where continuous backups are explicitly set to the default configuration
// and then remove continuous backups and call terraform plan, no changes would be found.
func TestAccAlloydbCluster_continuousBackup_noChangeIfRemoved(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"enabled":              true,
		"recovery_window_days": 14,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_continuousBackupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
		},
	})
}

// Ensures changes to the continuous backup config properly applies
func TestAccAlloydbCluster_continuousBackup_update(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix":        suffix,
		"enabled":              true,
		"recovery_window_days": 15,
	}
	context2 := map[string]interface{}{
		"random_suffix":        suffix,
		"enabled":              false,
		"recovery_window_days": 14,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withoutContinuousBackupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_continuousBackupConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "15"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_continuousBackupConfig(context2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.enabled", "false"),
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "continuous_backup_config.0.recovery_window_days", "14"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

func testAccAlloydbCluster_withoutContinuousBackupConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func testAccAlloydbCluster_continuousBackupConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  continuous_backup_config {
    enabled              = %{enabled}
    recovery_window_days = %{recovery_window_days}
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func TestAccAlloydbCluster_continuousBackup_CMEKIsUpdatable(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	kms := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-alloydb-key1")
	context := map[string]interface{}{
		"random_suffix": suffix,
		"key_ring":      kms.KeyRing.Name,
		"key_name":      kms.CryptoKey.Name,
	}

	kms2 := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-alloydb-key2")
	context2 := map[string]interface{}{
		"random_suffix": suffix,
		"key_ring":      kms2.KeyRing.Name,
		"key_name":      kms2.CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_usingCMEKInClusterAndContinuousBackup(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_usingCMEKInClusterAndContinuousBackup(context2),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_continuousBackupUsingCMEKAllowDeletion(context2),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "location"},
			},
		},
	})
}

func testAccAlloydbCluster_usingCMEKInClusterAndContinuousBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  continuous_backup_config {
    enabled       		 = true
	recovery_window_days = 20
    encryption_config {
      kms_key_name = "%{key_name}"
    }
  }
  lifecycle {
	prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
  crypto_key_id = "%{key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
  ]
}
`, context)
}

func testAccAlloydbCluster_continuousBackupUsingCMEKAllowDeletion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  continuous_backup_config {
    enabled       		 = true
	recovery_window_days = 20
    encryption_config {
      kms_key_name = "%{key_name}"
    }
  }
  depends_on = [google_kms_crypto_key_iam_binding.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_binding" "crypto_key" {
	crypto_key_id = "%{key_name}"
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	members = [
	  "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com",
	]
  }
`, context)
}
