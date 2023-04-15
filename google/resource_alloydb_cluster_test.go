package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlloydbCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
	return Nprintf(`
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
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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

func testAccAlloydbCluster_withInitialUserAndAutomatedBackupPolicy(context map[string]interface{}) string {
	return Nprintf(`
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
        hours   = 23
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
	return Nprintf(`
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

// We expect an error when creating a cluster without location.
// Location is a `required` field.
func TestAccAlloydbCluster_missingLocation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbCluster_missingLocation(context),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testAccAlloydbCluster_missingLocation(context map[string]interface{}) string {
	return Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}
