// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "subscription_type", "STANDARD"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccAlloydbCluster_update(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "labels", "terraform_labels"},
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

func TestAccAlloydbCluster_upgrade(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-instance-upgrade-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_beforeUpgrade(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "labels", "terraform_labels", "skip_await_major_version_upgrade"},
			},
			{
				Config: testAccAlloydbCluster_afterUpgrade(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location", "labels", "terraform_labels", "skip_await_major_version_upgrade"},
			},
		},
	})
}

func testAccAlloydbCluster_beforeUpgrade(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  skip_await_major_version_upgrade = false
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  database_version = "POSTGRES_14"
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 8
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbCluster_afterUpgrade(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  skip_await_major_version_upgrade = false
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
  }
  database_version = "POSTGRES_15"
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 8
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// Trial cluster creation should succeed with subscription type field set to Trial.
func TestAccAlloydbCluster_withSubscriptionTypeTrial(t *testing.T) {
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
				Config: testAccAlloydbCluster_withSubscriptionTypeTrial(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "subscription_type", "TRIAL"),
					resource.TestMatchResourceAttr("google_alloydb_cluster.default", "trial_metadata.0.start_time", regexp.MustCompile(".+")),
					resource.TestMatchResourceAttr("google_alloydb_cluster.default", "trial_metadata.0.end_time", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func testAccAlloydbCluster_withSubscriptionTypeTrial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  subscription_type = "TRIAL"
  network_config {
  	network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// Standard cluster creation should succeed with subscription type field set to Standard.
func TestAccAlloydbCluster_withSubscriptionTypeStandard(t *testing.T) {
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
				Config: testAccAlloydbCluster_withSubscriptionTypeStandard(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "subscription_type", "STANDARD"),
				),
			},
		},
	})
}

func testAccAlloydbCluster_withSubscriptionTypeStandard(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  subscription_type = "STANDARD"
  network_config {
  	network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }

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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
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
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-alloydb-cluster-key1").CryptoKey.Name,
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  encryption_config {
    kms_key_name = "%{kms_key_name}"
  }
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
data "google_project" "project" {}
resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}
`, context)
}

func TestAccAlloydbCluster_CMEKInAutomatedBackupIsUpdatable(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"kms_key_name1": acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-alloydb-backup-update-key1").CryptoKey.Name,
		"kms_key_name2": acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bootstrap-alloydb-backup-update-key2").CryptoKey.Name,
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  encryption_config {
    kms_key_name = "%{kms_key_name1}"
  }
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = "%{kms_key_name1}"
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  lifecycle {
	prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}
`, context)
}

func testAccAlloydbCluster_updateCMEKInAutomatedBackup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  encryption_config {
    kms_key_name = "%{kms_key_name1}"
  }
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = "%{kms_key_name2}"
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  lifecycle {
    prevent_destroy = true
  }
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "crypto_key2" {
  crypto_key_id = "%{kms_key_name2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}
`, context)
}

func testAccAlloydbCluster_usingCMEKallowDeletion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  encryption_config {
    kms_key_name = "%{kms_key_name1}"
  }
  automated_backup_policy {
    location      = "us-central1"
    backup_window = "1800s"
    enabled       = true
    encryption_config {
      kms_key_name = "%{kms_key_name2}"
    }
    time_based_retention {
      retention_period = "510s"
    }
  }
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "crypto_key2" {
  crypto_key_id = "%{kms_key_name2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
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

func testAccAlloydbCluster_continuousBackupConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }

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
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
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
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}
`, context)
}

func testAccAlloydbCluster_continuousBackupUsingCMEKAllowDeletion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  continuous_backup_config {
    enabled       		 = true
	recovery_window_days = 20
    encryption_config {
      kms_key_name = "%{key_name}"
    }
  }
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
	crypto_key_id = "%{key_name}"
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
  }
`, context)
}

// Ensures cluster creation works with networkConfig.
func TestAccAlloydbCluster_withNetworkConfig(t *testing.T) {
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
				Config: testAccAlloydbCluster_withNetworkConfig(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlloydbCluster_withNetworkConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
		network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
}
data "google_project" "project" {}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// Ensures cluster creation works with networkConfig and a specified allocated IP range.
func TestAccAlloydbCluster_withNetworkConfigAndAllocatedIPRange(t *testing.T) {
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
				Config: testAccAlloydbCluster_withNetworkConfigAndAllocatedIPRange(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlloydbCluster_withNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
		network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
		allocated_ip_range = google_compute_global_address.private_ip_alloc.name
  }
}
data "google_project" "project" {}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
resource "google_compute_global_address" "private_ip_alloc" {
	name          =  "tf-test-alloydb-cluster%{random_suffix}"
	address_type  = "INTERNAL"
	purpose       = "VPC_PEERING"
	prefix_length = 16
	network       = google_compute_network.default.id
  }
  
`, context)
}

// Ensures cluster creation works with correctly specified maintenance update policy.
func TestAccAlloydbCluster_withMaintenanceWindows(t *testing.T) {
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
				Config: testAccAlloydbCluster_withMaintenanceWindows(context),
			},
			{
				ResourceName:      "google_alloydb_cluster.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAlloydbCluster_withMaintenanceWindows(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
		network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  maintenance_update_policy {
    maintenance_windows {
      day = "WEDNESDAY"
      start_time {
        hours = 12
        minutes = 0
        seconds = 0
        nanos = 0
      }
    }
  }
}
data "google_project" "project" {}
resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}

// Ensures cluster creation throws expected errors for incorrect configurations of maintenance update policy.
func TestAccAlloydbCluster_withMaintenanceWindowsMissingFields(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_withMaintenanceWindowMissingStartTime(context),
				ExpectError: regexp.MustCompile("Error: Insufficient start_time blocks"),
			},
			{
				Config:      testAccAlloydbCluster_withMaintenanceWindowMissingDay(context),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func testAccAlloydbCluster_withMaintenanceWindowMissingStartTime(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  
  maintenance_update_policy {
    maintenance_windows {
      day = "WEDNESDAY"
    }
  }
}

resource "google_compute_network" "default" {
  name     = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}
`, context)
}

func testAccAlloydbCluster_withMaintenanceWindowMissingDay(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  }
  
  maintenance_update_policy {
    maintenance_windows {
      start_time {
        hours = 12
        minutes = 0
        seconds = 0
        nanos = 0
      }
    }
  }
}

resource "google_compute_network" "default" {
  name     = "tf-test-alloydb-cluster%{random_suffix}"
}

data "google_project" "project" {}
`, context)
}

// Ensures cluster creation succeeds for a Private Service Connect enabled cluster.
func TestAccAlloydbCluster_withPrivateServiceConnect(t *testing.T) {
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
				Config: testAccAlloydbCluster_withPrivateServiceConnect(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "psc_config.0.psc_enabled", "true"),
					resource.TestMatchResourceAttr("google_alloydb_cluster.default", "psc_config.0.service_owned_project_number", regexp.MustCompile("^[1-9]\\d*$")),
				),
			},
		},
	})
}

func testAccAlloydbCluster_withPrivateServiceConnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  psc_config {
    psc_enabled = true
  }
}
data "google_project" "project" {}
`, context)
}

// Ensures cluster update from unspecified to standard and standard to standard works with no change in config.
func TestAccAlloydbCluster_standardClusterUpdate(t *testing.T) {
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_alloydb_cluster.default", "subscription_type", "STANDARD"),
				),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_withSubscriptionTypeStandard(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_withSubscriptionTypeStandard(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
		},
	})
}

// Ensures cluster update succeeds with subscription type from trial to standard and trial to trial results in no change in config.
func TestAccAlloydbCluster_trialClusterUpdate(t *testing.T) {
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
				Config: testAccAlloydbCluster_withSubscriptionTypeTrial(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_withSubscriptionTypeTrial(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_withSubscriptionTypeStandard(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
		},
	})
}

// Ensures cluster update throws expected errors for subscription update from standard to trial.
func TestAccAlloydbCluster_standardClusterUpdateFailure(t *testing.T) {
	t.Parallel()
	errorPattern := `.*The request was invalid: invalid subscription_type update`
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_withSubscriptionTypeStandard(context),
			},
			{
				Config:      testAccAlloydbCluster_withSubscriptionTypeTrial(context),
				ExpectError: regexp.MustCompile(errorPattern),
			},
		},
	})
}
