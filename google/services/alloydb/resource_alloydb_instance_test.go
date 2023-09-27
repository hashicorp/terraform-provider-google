// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccAlloydbInstance_update(t *testing.T) {
	t.Parallel()

	random_suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-instance-update-1"),
		"random_suffix": random_suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_alloydbInstanceBasic(context),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbInstance_update(context),
			},
			{
				ResourceName:            "google_alloydb_instance.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "instance_id", "reconciling", "update_time", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccAlloydbInstance_alloydbInstanceBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

func testAccAlloydbInstance_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 4
  }

  labels = {
	test = "tf-test-alloydb-instance%{random_suffix}"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id

  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a primary instance with minimal number of fields
func TestAccAlloydbInstance_createInstanceWithMandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-instance-mandatory-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithMandatoryFields(context),
			},
		},
	})
}

func testAccAlloydbInstance_createInstanceWithMandatoryFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a primary instance with maximum number of fields
/* func TestAccAlloydbInstance_createInstanceWithMaximumFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-instance-maximum-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithMaximumFields(context),
			},
		},
	})
}

func testAccAlloydbInstance_createInstanceWithMaximumFields(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  labels = {
    test_label = "test-alloydb-label"
  }
  annotations = {
    test_annotation = "test-alloydb-annotation"
  }
  gce_zone = "us-east1-a"
  database_flags = {
	  "alloydb.enable_auto_explain" = "true"
  }
  availability_type = "REGIONAL"
  machine_config {
	  cpu_count = 4
  }
  query_insights_config {
    query_string_length = 300
    record_application_tags = "false"
    record_client_address = "true"
    query_plans_per_minute = 10
  }
  depends_on = [google_service_networking_connection.vpc_connection]
  lifecycle {
    ignore_changes = [
      gce_zone,
      annotations
    ]
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}*/

// This test passes if we are able to create a primary instance with an associated read-pool instance
func TestAccAlloydbInstance_createPrimaryAndReadPoolInstance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-instance-readpool-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createPrimaryAndReadPoolInstance(context),
			},
		},
	})
}

func testAccAlloydbInstance_createPrimaryAndReadPoolInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_instance" "read_pool" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}-read"
  instance_type = "READ_POOL"
  read_pool_config {
    node_count = 4
  }
  depends_on = [google_alloydb_instance.primary]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to update a database flag in primary instance
/*func TestAccAlloydbInstance_updateDatabaseFlagInPrimaryInstance(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  acctest.BootstrapSharedServiceNetworkingConnection(t, "alloydb-instance-updatedb-1"),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_autoExplainEnabledInPrimaryInstance(context),
			},
			{
				ResourceName:      "google_alloydb_instance.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlloydbInstance_autoExplainDisabledInPrimaryInstance(context),
			},
		},
	})
}

func testAccAlloydbInstance_autoExplainEnabledInPrimaryInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  database_flags = {
	  "alloydb.enable_auto_explain" = "true"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}*/

func testAccAlloydbInstance_autoExplainDisabledInPrimaryInstance(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "primary" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
  database_flags = {
	  "alloydb.enable_auto_explain" = "false"
  }
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

data "google_project" "project" {}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// This test passes if we are able to create a primary instance by specifying network_config.network and network_config.allocated_ip_range
func TestAccAlloydbInstance_createInstanceWithNetworkConfigAndAllocatedIPRange(t *testing.T) {
	t.Parallel()

	projectNumber := envvar.GetTestProjectNumberFromEnv()
	testId := "alloydbinstance-network-config-1"
	networkName := acctest.BootstrapSharedTestNetwork(t, testId)
	networkId := fmt.Sprintf("projects/%v/global/networks/%v", projectNumber, networkName)
	addressName := acctest.BootstrapSharedTestGlobalAddress(t, testId, networkId)
	acctest.BootstrapSharedServiceNetworkingConnection(t, testId)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"network_name":  networkName,
		"address_name":  addressName,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbInstance_createInstanceWithNetworkConfigAndAllocatedIPRange(context),
			},
		},
	})
}

func testAccAlloydbInstance_createInstanceWithNetworkConfigAndAllocatedIPRange(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network_config {
    network    = data.google_compute_network.default.id
    allocated_ip_range = data.google_compute_global_address.private_ip_alloc.name
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_compute_global_address" "private_ip_alloc" {
  name =  "%{address_name}"
}
`, context)
}
