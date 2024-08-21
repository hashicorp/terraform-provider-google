// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// This test passes if secondary instance's machine config can be updated
func TestAccAlloydbInstance_secondaryInstanceUpdateMachineConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceInitial(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceUpdateMachineConfig(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceInitial(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
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
    network = data.google_compute_network.default.id
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
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbInstance_secondaryInstanceUpdateMachineConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
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
    network = data.google_compute_network.default.id
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
    cpu_count = 4
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a secondary instance with an associated read-pool instance
func TestAccAlloydbInstance_secondaryInstanceWithReadPoolInstance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceWithReadPoolInstance(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceWithReadPoolInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
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
  location     = "us-west1"
  network_config {
    network = data.google_compute_network.default.id
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
}

resource "google_alloydb_instance" "read_pool" {
  cluster       = google_alloydb_cluster.secondary.name
  instance_id   = "tf-test-alloydb-read-instance%{random_suffix}-read"
  instance_type = "READ_POOL"
  read_pool_config {
    node_count = 4
  }
  depends_on = [google_alloydb_instance.secondary]
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a secondary instance by specifying network_config.network and network_config.allocated_ip_range
func TestAccAlloydbCluster_secondaryInstanceWithNetworkConfigAndAllocatedIPRange(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
		"address_name":  acctest.BootstrapSharedTestGlobalAddress(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_secondaryInstanceWithNetworkConfigAndAllocatedIPRange(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbCluster_secondaryInstanceWithNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network    = data.google_compute_network.default.id
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
  location     = "us-west1"
  network_config {
    network    = data.google_compute_network.default.id
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
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name =  "%{address_name}"
}
`, context)
}

// This test passes if secondary instance's database flag config can be updated
func TestAccAlloydbInstance_secondaryInstanceUpdateDatabaseFlag(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceInitial(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceUpdateDatabaseFlag(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceUpdateDatabaseFlag(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
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
    network = data.google_compute_network.default.id
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

  database_flags = {
	  "alloydb.enable_auto_explain" = "true"
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if secondary instance's query insight config can be updated
func TestAccAlloydbInstance_secondaryInstanceUpdateQueryInsightConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceInitial(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
			{
				Config: testAccAlloydbInstance_secondaryInstanceUpdateQueryInsightConfig(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceUpdateQueryInsightConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
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
    network = data.google_compute_network.default.id
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

  query_insights_config {
      query_plans_per_minute = 10
      query_string_length = 2048
      record_application_tags = true
      record_client_address = true
  }
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a secondary instance with maximum fields
func TestAccAlloydbInstance_secondaryInstanceMaximumFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydbinstance-network-config-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_secondaryInstanceMaximumFields(context),
			},
			{
				ResourceName:            "google_alloydb_instance.secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_secondaryInstanceMaximumFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_cluster" "primary" {
  cluster_id = "tf-test-alloydb-primary-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network = data.google_compute_network.default.id
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
  location     = "us-west1"
  network_config {
    network = data.google_compute_network.default.id
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

  labels = {
    test_label = "test-alloydb-label"
  }

  query_insights_config {
      query_plans_per_minute = 10
      query_string_length = 2048
      record_application_tags = true
      record_client_address = true
  }

  availability_type = "REGIONAL"
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}
