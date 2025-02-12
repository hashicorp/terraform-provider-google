// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAlloydbDatabaseInstanceDatasourceConfig(t *testing.T) {
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
				Config: testAccAlloydbDatabaseInstanceDatasourceConfig(context),
			},
		},
	})
}

func testAccAlloydbDatabaseInstanceDatasourceConfig(context map[string]interface{}) string {
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
  network_config {
    network = data.google_compute_network.default.id
  }
  initial_user {
    password = "tf-test-alloydb-cluster%{random_suffix}"
  }
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}

data "google_alloydb_instance" "default" {
  cluster_id = google_alloydb_cluster.default.cluster_id
  instance_id = google_alloydb_instance.default.instance_id
  location = "us-central1"
}
`, context)
}
