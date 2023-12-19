// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkebackup_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGKEBackupBackupPlan_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":         envvar.GetTestProjectFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
		"network_name":    acctest.BootstrapSharedTestNetwork(t, "gke-cluster"),
		"subnetwork_name": acctest.BootstrapSubnet(t, "gke-cluster", acctest.BootstrapSharedTestNetwork(t, "gke-cluster")),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEBackupBackupPlanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupBackupPlan_basic(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccGKEBackupBackupPlan_full(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_plan.backupplan",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccGKEBackupBackupPlan_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
  location           = "us-central1"
  initial_node_count = 1
  workload_identity_config {
    workload_pool = "%{project}.svc.id.goog"
  }
  addons_config {
    gke_backup_agent_config {
      enabled = true
    }
  }
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}

resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  backup_config {
    include_volume_data = false
    include_secrets = false
    all_namespaces = true
  }
  labels = {
    "some-key-1": "some-value-1"
  }
}
`, context)
}

func testAccGKEBackupBackupPlan_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test-testcluster%{random_suffix}"
  location           = "us-central1"
  initial_node_count = 1
  workload_identity_config {
    workload_pool = "%{project}.svc.id.goog"
  }
  addons_config {
    gke_backup_agent_config {
      enabled = true
    }
  }
  deletion_protection = false
  network       = "%{network_name}"
  subnetwork    = "%{subnetwork_name}"
}
	
resource "google_gke_backup_backup_plan" "backupplan" {
  name = "tf-test-testplan%{random_suffix}"
  cluster = google_container_cluster.primary.id
  location = "us-central1"
  retention_policy {
    backup_delete_lock_days = 30
    backup_retain_days = 180
  }
  backup_schedule {
    cron_schedule = "0 9 * * 1"
  }
  backup_config {
    include_volume_data = true
    include_secrets = true
    selected_applications {
      namespaced_names {
        name = "app1"
        namespace = "ns1"
      }
      namespaced_names {
        name = "app2"
        namespace = "ns2"
      }
    }
  }
  labels = {
    "some-key-2": "some-value-2"
  }
}
`, context)
}
