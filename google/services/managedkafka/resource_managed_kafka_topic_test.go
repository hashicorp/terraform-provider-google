// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package managedkafka_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccManagedKafkaTopic_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckManagedKafkaTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKafkaTopic_basic(context),
			},
			{
				ResourceName:            "google_managed_kafka_topic.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "location", "topic_id"},
			},
			{
				Config: testAccManagedKafkaTopic_update(context),
			},
			{
				ResourceName:            "google_managed_kafka_topic.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "location", "topic_id"},
			},
		},
	})
}

func testAccManagedKafkaTopic_basic(context map[string]interface{}) string {
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
}

resource "google_managed_kafka_topic" "example" {
  cluster = google_managed_kafka_cluster.example.cluster_id
  topic_id = "tf-test-my-topic%{random_suffix}"
  location = "us-central1"
  partition_count = 2
  replication_factor = 3
  configs = {
    "cleanup.policy" = "compact"
  }
}

data "google_project" "project" {
}
`, context)
}

func testAccManagedKafkaTopic_update(context map[string]interface{}) string {
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
}

resource "google_managed_kafka_topic" "example" {
  cluster = google_managed_kafka_cluster.example.cluster_id
  topic_id = "tf-test-my-topic%{random_suffix}"
  location = "us-central1"
  partition_count = 3
  replication_factor = 3
  configs = {
    "cleanup.policy" = "compact"
  }
}

data "google_project" "project" {
}
`, context)
}
