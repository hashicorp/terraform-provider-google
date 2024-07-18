// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// The cluster creation should succeed with minimal number of arguments
func TestAccAlloydbCluster_secondaryClusterMandatoryFields(t *testing.T) {
	t.Parallel()
	// https://github.com/hashicorp/terraform-provider-google/issues/16231
	acctest.SkipIfVcr(t)
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterMandatoryFields(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterMandatoryFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = data.google_compute_network.default.id
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// Validation test to ensure proper error is raised if create secondary cluster is called without any secondary_config field
func TestAccAlloydbCluster_secondaryClusterMissingSecondaryConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_secondaryClusterMissingSecondaryConfig(context),
				ExpectError: regexp.MustCompile("Error creating cluster. Can not create secondary cluster without secondary_config field."),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterMissingSecondaryConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = data.google_compute_network.default.id
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// Validation test to ensure proper error is raised if secondary_config field is provided but no cluster_type is specified.
func TestAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButMissingClusterTypeSecondary(t *testing.T) {
	t.Parallel()

	// Unskip in https://github.com/hashicorp/terraform-provider-google/issues/16231
	acctest.SkipIfVcr(t)

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButMissingClusterTypeSecondary(context),
				ExpectError: regexp.MustCompile("Error creating cluster. Add {cluster_type: \"SECONDARY\"} if attempting to create a secondary cluster, otherwise remove the secondary_config."),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButMissingClusterTypeSecondary(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = data.google_compute_network.default.id

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// Validation test to ensure proper error is raised if secondary_config field is provided but cluster_type is primary
func TestAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButClusterTypeIsPrimary(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButClusterTypeIsPrimary(context),
				ExpectError: regexp.MustCompile("Error creating cluster. Add {cluster_type: \"SECONDARY\"} if attempting to create a secondary cluster, otherwise remove the secondary_config."),
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterDefinedSecondaryConfigButClusterTypeIsPrimary(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if secondary cluster can be updated
func TestAccAlloydbCluster_secondaryClusterUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterMandatoryFields(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterUpdate(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = data.google_compute_network.default.id
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  labels = {
    foo = "bar"
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func TestAccAlloydbCluster_secondaryClusterUsingCMEK(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
		"key_name":      "tf-test-key-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterUsingCMEK(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterUsingCMEK(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network      = data.google_compute_network.default.id
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  encryption_config {
    kms_key_name = google_kms_crypto_key.key.id
  }

  depends_on = [
    google_alloydb_instance.primary,
    google_kms_crypto_key_iam_member.crypto_key
  ]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_kms_key_ring" "keyring" {
  name     = "%{key_name}"
  location = "us-east1"
}

resource "google_kms_crypto_key" "key" {
  name     = "%{key_name}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = google_kms_crypto_key.key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-alloydb.iam.gserviceaccount.com"
}
`, context)
}

// Ensures secondary cluster creation works with networkConfig.
func TestAccAlloydbCluster_secondaryClusterWithNetworkConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterWithNetworkConfig(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterWithNetworkConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
	network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
  }
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network_config {
	network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
  }
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// Ensures secondary cluster creation works with networkConfig and a specified allocated IP range.
func TestAccAlloydbCluster_secondaryClusterWithNetworkConfigAndAllocatedIPRange(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"address_name":  acctest.BootstrapSharedTestGlobalAddress(t, "alloydbinstance-network-config-1"),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryClusterWithNetworkConfigAndAllocatedIPRange(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterWithNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
	network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
	allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-east1"
  network_config {
	network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
	allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  depends_on = [google_alloydb_instance.primary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name          =  "%{address_name}"
}
`, context)
}

// This test passes if secondary cluster can be promoted
func TestAccAlloydbCluster_secondaryClusterPromote(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-east1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryClusterWithInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  deletion_policy = "FORCE"

  depends_on = [google_alloydb_instance.primary]
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbCluster_secondaryClusterPromote(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if secondary cluster can be promoted and updated simultaneously
func TestAccAlloydbCluster_secondaryClusterPromoteAndSimultaneousUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-east1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteAndSimultaneousUpdate(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromoteAndSimultaneousUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = true
  }

  labels = {
    foo = "bar" 
  }  
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if secondary cluster can be promoted and the original primary can be deleted after promotion
func TestAccAlloydbCluster_secondaryClusterPromoteAndDeleteOriginalPrimary(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-east1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteAndDeleteOriginalPrimary(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromoteAndDeleteOriginalPrimary(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if the promoted secondary cluster can be updated
func TestAccAlloydbCluster_secondaryClusterPromoteAndUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-east1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteAndUpdate(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromoteAndUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }

  labels = {
    foo = "bar"
  }

}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if secondary cluster can be promoted with networkConfig and a specified allocated IP range
func TestAccAlloydbCluster_secondaryClusterPromoteWithNetworkConfigAndAllocatedIPRange(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"address_name":  acctest.BootstrapSharedTestGlobalAddress(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstanceAndNetworkConfigAndAllocatedIPRange(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteWithNetworkConfigAndAllocatedIPRange(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryClusterWithInstanceAndNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
    allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-south1"
  network_config {
    network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
    allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
  cluster_type = "SECONDARY"

  continuous_backup_config {
    enabled = false
  }

  secondary_config {
    primary_cluster_name = google_alloydb_cluster.primary.name
  }

  deletion_policy = "FORCE"

  depends_on = [google_alloydb_instance.primary]
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name          =  "%{address_name}"
}
`, context)
}

func testAccAlloydbCluster_secondaryClusterPromoteWithNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
    allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "us-south1"
  network_config {
    network    = "projects/${data.google_project.project.number}/global/networks/${data.google_compute_network.default.name}"
    allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name          =  "%{address_name}"
}
`, context)
}

// This test passes if automated backup policy and inital user can be added and deleted from the promoted secondary cluster
func TestAccAlloydbCluster_secondaryClusterPromoteAndAddAndDeleteAutomatedBackupPolicyAndInitialUser(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-south1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"hour":                       23,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteAndAddAutomatedBackupPolicyAndInitialUser(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromoteAndAddAutomatedBackupPolicyAndInitialUser(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }

  initial_user {
    user     = "tf-test-alloydb-secondary-cluster%{random_suffix}"
    password = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  }

  automated_backup_policy {
    location      = "%{secondary_cluster_location}"
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
      test = "tf-test-alloydb-secondary-cluster%{random_suffix}"
    }
  }
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if time based retention policy can be added and deleted from the promoted secondary cluster
func TestAccAlloydbCluster_secondaryClusterPromoteAndDeleteTimeBasedRetentionPolicy(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-south1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteWithTimeBasedRetentionPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteWithoutTimeBasedRetentionPolicy(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromoteWithTimeBasedRetentionPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }

  initial_user {
    user     = "tf-test-alloydb-secondary-cluster%{random_suffix}"
    password = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  }

  automated_backup_policy {
    location      = "%{secondary_cluster_location}"
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
  }
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbCluster_secondaryClusterPromoteWithoutTimeBasedRetentionPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled = false
  }

  initial_user {
    user     = "tf-test-alloydb-secondary-cluster%{random_suffix}"
    password = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  }

  automated_backup_policy {
    location      = "%{secondary_cluster_location}"
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
  }
}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if continuous backup config can be enabled in the promoted secondary cluster
func TestAccAlloydbCluster_secondaryClusterPromoteAndAddContinuousBackupConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":              acctest.RandString(t, 10),
		"secondary_cluster_location": "us-south1",
		"network_name":               acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryClusterWithInstance(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromote(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
			{
				Config: testAccAlloydbCluster_secondaryClusterPromoteAndAddContinuousBackupConfig(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "restore_backup_source", "restore_continuous_backup_source", "cluster_id", "location", "deletion_policy", "labels", "annotations", "terraform_labels", "reconciling"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryClusterPromoteAndAddContinuousBackupConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network      = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.primary.name
  instance_id   = "tf-test-alloydb-primary-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "secondary" {
  cluster_id   = "tf-test-alloydb-secondary-cluster%{random_suffix}"
  location     = "%{secondary_cluster_location}"
  network      = data.google_compute_network.default.id
  cluster_type = "PRIMARY"

  continuous_backup_config {
    enabled              = true
    recovery_window_days = 14
  }

}

resource "google_alloydb_instance" "secondary" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-secondary-instance%{random_suffix}"
  instance_type = google_alloydb_cluster.secondary.cluster_type

  machine_config {
    cpu_count = 2
  }

  lifecycle {
    ignore_changes = [instance_type]
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}
