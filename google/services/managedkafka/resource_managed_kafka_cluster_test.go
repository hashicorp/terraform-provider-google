// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package managedkafka_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccManagedKafkaCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckManagedKafkaClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKafkaCluster_basic(context),
			},
			{
				ResourceName:            "google_managed_kafka_cluster.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccManagedKafkaCluster_update(context),
			},
			{
				ResourceName:            "google_managed_kafka_cluster.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccManagedKafkaCluster_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_managed_kafka_cluster" "example" {
  cluster_id = "tf-test-my-cluster%{random_suffix}"
  location = "us-central1"
  capacity_config {
    vcpu_count = 3
    memory_bytes = 3221225472
  }
  gcp_config {
    access_config {
      network_configs {
        subnet = "projects/${data.google_project.project.number}/regions/us-central1/subnetworks/default"
      }
    }
  }
  rebalance_config {
    mode = "NO_REBALANCE"
  }
  labels = {
    key = "value"
  }
}

data "google_project" "project" {
}
`, context)
}

func testAccManagedKafkaCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_managed_kafka_cluster" "example" {
  cluster_id = "tf-test-my-cluster%{random_suffix}"
  location = "us-central1"
  capacity_config {
    vcpu_count = 4
    memory_bytes = 4512135122
  }
  gcp_config {
    access_config {
      network_configs {
        subnet = "projects/${data.google_project.project.number}/regions/us-central1/subnetworks/default"
      }
    }
  }
  rebalance_config {
    mode = "AUTO_REBALANCE_ON_SCALE_UP"
  }
  labels = {
    key = "new-value"
  }
}

data "google_project" "project" {
}
`, context)
}
