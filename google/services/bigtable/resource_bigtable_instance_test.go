// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBigtableInstance_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigtableInstance_invalid(instanceName),
				ExpectError: regexp.MustCompile("config is invalid: Too few cluster blocks: Should have at least 1 \"cluster\" block"),
			},
			{
				Config: testAccBigtableInstance(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance(instanceName, 4),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_cluster(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_cluster(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "cluster"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterReordered(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				PlanOnly:                true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "cluster"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterModified(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "cluster", "labels", "terraform_labels"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterReordered(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				PlanOnly:                true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "cluster"}, // we don't read instance type back
			},
			{
				Config: testAccBigtableInstance_clusterModifiedAgain(instanceName, 5),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "cluster", "labels", "terraform_labels"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_development(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_development(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_allowDestroy(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_noAllowDestroy(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				Config:      testAccBigtableInstance_noAllowDestroy(instanceName, 3),
				Destroy:     true,
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccBigtableInstance(instanceName, 3),
			},
		},
	})
}

func TestAccBigtableInstance_kms(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	kms1 := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	kms2 := acctest.BootstrapKMSKeyInLocation(t, "us-east1")
	pid := envvar.GetTestProjectFromEnv()
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_kms(pid, instanceName, kms1.CryptoKey.Name, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			// TODO(kevinsi4508): Verify that the instance can be recreated due to `kms_key_name` change.
			{
				Config:             testAccBigtableInstance_kms(pid, instanceName, kms2.CryptoKey.Name, 3),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccBigtableInstance_createWithAutoscalingAndUpdate(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create Autoscaling config with 2 nodes. Default storage_target is set by service based on storage type.
				Config: testAccBigtableInstance_autoscalingCluster(instanceName, 2, 5, 70),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "cluster.0.num_nodes", "2"),
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "cluster.0.autoscaling_config.0.storage_target", "8192"),
				),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				// Update Autoscaling configs. storage_target is unchanged.
				Config: testAccBigtableInstance_autoscalingCluster(instanceName, 1, 5, 80),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "cluster.0.autoscaling_config.0.storage_target", "8192"),
				),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_createWithAutoscalingAndUpdateWithStorageTarget(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create Autoscaling config with 2 nodes. Set storage_target to a non-default value.
				Config: autoscalingClusterConfigWithStorageTarget(instanceName, 2, 5, 70, 9000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "cluster.0.num_nodes", "2"),
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "cluster.0.autoscaling_config.0.storage_target", "9000"),
				),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				// Update Autoscaling configs.
				Config: autoscalingClusterConfigWithStorageTarget(instanceName, 1, 5, 80, 10000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "cluster.0.autoscaling_config.0.storage_target", "10000"),
				),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_enableAndDisableAutoscaling(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance(instanceName, 2),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
			},
			{
				// Enable Autoscaling.
				Config: testAccBigtableInstance_autoscalingCluster(instanceName, 2, 5, 70),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("google_bigtable_instance.instance",
					"cluster.0.num_nodes", "2")),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				// Disable Autoscaling specifying num_nodes=1 and node count becomes 1.
				Config: testAccBigtableInstance(instanceName, 1),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
			},
		},
	})
}

func TestAccBigtableInstance_enableAndDisableAutoscalingWithoutNumNodes(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create Autoscaling cluster with 2 nodes.
				Config: testAccBigtableInstance_autoscalingCluster(instanceName, 2, 5, 70),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("google_bigtable_instance.instance",
					"cluster.0.num_nodes", "2")),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"}, // we don't read instance type back
			},
			{
				// Disable Autoscaling without specifying num_nodes, it should use the current node count, which is 2.
				Config: testAccBigtableInstance_noNumNodes(instanceName),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("google_bigtable_instance.instance",
					"cluster.0.num_nodes", "2")),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
			},
		},
	})
}

func testAccCheckBigtableInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_instance" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			c, err := config.BigTableClientFactory(config.UserAgent).NewInstanceAdminClient(config.Project)
			if err != nil {
				return fmt.Errorf("Error starting instance admin client. %s", err)
			}

			defer c.Close()

			_, err = c.InstanceInfo(ctx, rs.Primary.Attributes["name"])
			if err == nil {
				return fmt.Errorf("Instance %s still exists.", rs.Primary.Attributes["name"])
			}
		}

		return nil
	}
}

func TestAccBigtableInstance_MultipleClustersSameID(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigtableInstance_multipleClustersSameID(instanceName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("duplicated cluster_id: %q", instanceName)),
			},
			{
				Config: testAccBigtableInstance(instanceName, 3),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
			},
			{
				Config:      testAccBigtableInstance_multipleClustersSameID(instanceName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("duplicated cluster_id: %q", instanceName)),
			},
		},
	})
}

func TestAccBigtableInstance_forceDestroyBackups(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	randomString := acctest.RandString(t, 10)
	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"instance_name":  fmt.Sprintf("tf-test-instance-%s", randomString),
		"cluster_name_1": fmt.Sprintf("tf-test-cluster-%s-1", randomString),
		"cluster_name_2": fmt.Sprintf("tf-test-cluster-%s-2", randomString),
		"cluster_zone_1": fmt.Sprintf("%s-a", region),
		"cluster_zone_2": fmt.Sprintf("%s-b", region),
		"table_name":     fmt.Sprintf("tf-test-table-%s", randomString),
		"force_destroy":  true, // Overridden in test steps
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"http": {},
			"time": {},
		},
		CheckDestroy: testAccCheckBigtableInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Create force_destroy = false
				Config: testAccBigtableInstance_forceDestroy(context, false),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type", "labels", "terraform_labels"}, // we don't read instance type back
				Check: resource.ComposeTestCheckFunc(
					// Make sure field is set, and is set to false after import
					resource.TestCheckResourceAttr("google_bigtable_instance.instance", "force_destroy", "false"),
				),
			},
			{
				// Try to delete the instance after force_destroy = false was set before
				Config:      testAccBigtableInstance_forceDestroy_deleteInstance(),
				ExpectError: regexp.MustCompile("until all user backups have been deleted"),
			},
			{
				// Update force_destroy = true
				Config: testAccBigtableInstance_forceDestroy(context, true),
			},
			{
				// Try to delete the instance after force_destroy = true was set before
				Config: testAccBigtableInstance_forceDestroy_deleteInstance(),
			},
		},
	})
}

func testAccBigtableInstance_multipleClustersSameID(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s"
    num_nodes    = 2
    storage_type = "SSD"
  }

  deletion_protection = false

  labels = {
    env = "default"
  }
}
`, instanceName, instanceName, instanceName)
}

func testAccBigtableInstance(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    num_nodes    = %d
    storage_type = "HDD"
  }

  deletion_protection = false

  labels = {
    env = "default"
  }
}
`, instanceName, instanceName, numNodes)
}

func testAccBigtableInstance_noNumNodes(instanceName string) string {
	return fmt.Sprintf(`
		resource "google_bigtable_instance" "instance" {
			name = "%s"
			cluster {
				cluster_id   = "%s"
				storage_type = "HDD"
			}
			deletion_protection = false
			labels = {
				env = "default"
			}
		}`, instanceName, instanceName)
}

func testAccBigtableInstance_invalid(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
}
`, instanceName)
}

func testAccBigtableInstance_cluster(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s-a"
    zone         = "us-central1-a"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-b"
    zone         = "us-central1-b"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-c"
    zone         = "us-central1-c"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-d"
    zone         = "us-central1-f"
    num_nodes    = %d
    storage_type = "HDD"
  }

  deletion_protection = false
}
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

func testAccBigtableInstance_clusterReordered(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s-c"
    zone         = "us-central1-c"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-d"
    zone         = "us-central1-f"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-a"
    zone         = "us-central1-a"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-b"
    zone         = "us-central1-b"
    num_nodes    = %d
    storage_type = "HDD"
  }

  deletion_protection = false
}
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

func testAccBigtableInstance_clusterModified(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s-c"
    zone         = "us-central1-c"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-a"
    zone         = "us-central1-a"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-b"
    zone         = "us-central1-b"
    num_nodes    = %d
    storage_type = "HDD"
  }

  deletion_protection = false

  labels = {
    env = "default"
  }
}
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

// Add two clusters after testAccBigtableInstance_clusterModified.
func testAccBigtableInstance_clusterModifiedAgain(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s-c"
    zone         = "us-central1-c"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-a"
    zone         = "us-central1-a"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-b"
    zone         = "us-central1-b"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-asia-a"
    zone         = "asia-northeast1-a"
    num_nodes    = %d
    storage_type = "HDD"
  }
  cluster {
    cluster_id   = "%s-asia-b"
    zone         = "asia-northeast1-b"
    num_nodes    = %d
    storage_type = "HDD"
  }

  deletion_protection = false

  labels = {
    env = "test"
  }
}
`, instanceName, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes, instanceName, numNodes)
}

func testAccBigtableInstance_development(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }
  instance_type = "DEVELOPMENT"

  deletion_protection = false
}
`, instanceName, instanceName)
}

func testAccBigtableInstance_noAllowDestroy(instanceName string, numNodes int) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = %d
    storage_type = "HDD"
  }
}
`, instanceName, instanceName, numNodes)
}

func testAccBigtableInstance_kms(pid, instanceName, kmsKey string, numNodes int) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}


resource "google_project_iam_member" "kms_project_binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-bigtable.iam.gserviceaccount.com"
}

resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = %d
    storage_type = "HDD"
	kms_key_name = "%s"
  }
  depends_on = [google_project_iam_member.kms_project_binding]
  deletion_protection = false

}
`, pid, instanceName, instanceName, numNodes, kmsKey)
}

func testAccBigtableInstance_autoscalingCluster(instanceName string, min int, max int, cpuTarget int) string {
	return fmt.Sprintf(`resource "google_bigtable_instance" "instance" {
		name = "%s"
		cluster {
			cluster_id   = "%s"
			storage_type = "HDD"
			autoscaling_config {
				min_nodes = %d
				max_nodes = %d
				cpu_target = %d
			}
		}
	  deletion_protection = false

	}`, instanceName, instanceName, min, max, cpuTarget)
}

func autoscalingClusterConfigWithStorageTarget(instanceName string, min int, max int, cpuTarget int, storageTarget int) string {
	return fmt.Sprintf(`resource "google_bigtable_instance" "instance" {
		name = "%s"
		cluster {
			cluster_id   = "%s"
			storage_type = "HDD"
			autoscaling_config {
				min_nodes = %d
				max_nodes = %d
				cpu_target = %d
				storage_target = %d
			}
		}
	  deletion_protection = false

	}`, instanceName, instanceName, min, max, cpuTarget, storageTarget)
}

func testAccBigtableInstance_forceDestroy(context map[string]interface{}, forceDestroy bool) string {
	context["force_destroy"] = forceDestroy

	return acctest.Nprintf(`
provider "google" {
  alias = "http_auth"
}

resource "google_bigtable_instance" "instance" {
  name = "%{instance_name}"
  cluster {
    cluster_id   = "%{cluster_name_1}"
    num_nodes    = 1
    storage_type = "HDD"
    zone         = "%{cluster_zone_1}"
  }
  cluster {
    cluster_id   = "%{cluster_name_2}"
    num_nodes    = 1
    storage_type = "HDD"
    zone         = "%{cluster_zone_2}"
  }
  force_destroy = %{force_destroy}
  deletion_protection = false
  labels = {
    env = "default"
  }
}

resource "google_bigtable_table" "table" {
  name          = "%{table_name}"
  instance_name = google_bigtable_instance.instance.id
  split_keys    = ["a", "b", "c"]
}

data "google_client_config" "current" {
	provider = google.http_auth
}

locals {
  project = google_bigtable_instance.instance.project
  instance = google_bigtable_instance.instance.name
  cluster_1 = google_bigtable_instance.instance.cluster[0].cluster_id
  cluster_2 = google_bigtable_instance.instance.cluster[1].cluster_id
  backup = "backup-1"
}

data "http" "make_backup_1" {
  url    = "https://bigtableadmin.googleapis.com/v2/projects/${local.project}/instances/${local.instance}/clusters/${local.cluster_1}/backups?backupId=${local.backup}"
  method = "POST"

  request_headers = {
    Content-Type  = "application/json"
    Authorization = "Bearer ${data.google_client_config.current.access_token}"
  }

  request_body = <<EOT
{
  "sourceTable" : "${google_bigtable_table.table.id}",
  "expireTime" : "${time_offset.week-in-future.rfc3339}"
}
EOT

  depends_on = [
    google_bigtable_table.table // Needs to exist for backup to be made
  ]
}

data "http" "make_backup_2" {
  url    = "https://bigtableadmin.googleapis.com/v2/projects/${local.project}/instances/${local.instance}/clusters/${local.cluster_2}/backups?backupId=${local.backup}"
  method = "POST"

  request_headers = {
    Content-Type  = "application/json"
    Authorization = "Bearer ${data.google_client_config.current.access_token}"
  }

  request_body = <<EOT
{
  "sourceTable" : "${google_bigtable_table.table.id}",
  "expireTime" : "${time_offset.week-in-future.rfc3339}"
}
EOT

  depends_on = [
    google_bigtable_table.table // Needs to exist for backup to be made
  ]
}

check "health_check_1" {
  assert {
    condition     = data.http.make_backup_1.status_code == 200
    error_message = "HTTP request to create a backup returned a non-200 status code"
  }
}

check "health_check_2" {
  assert {
    condition     = data.http.make_backup_2.status_code == 200
    error_message = "HTTP request to create a backup returned a non-200 status code"
  }
}

# Expiration time for the backup being created
resource "time_offset" "week-in-future" {
  offset_days = 7
}

resource "time_sleep" "wait_30sec_1" {
  depends_on = [data.http.make_backup_1]
  create_duration = "30s"
}
resource "time_sleep" "wait_30sec_2" {
  depends_on = [data.http.make_backup_2]
  create_duration = "30s"
}
`, context)
}

func testAccBigtableInstance_forceDestroy_deleteInstance() string {
	// A version of the config in testAccBigtableInstance_forceDestroy missing all the BigTable resources + related things
	// This allows attempting to delete the instance when backups are present inside.
	return `
provider "google" {
  alias = "http_auth"
}
resource "time_offset" "week-in-future" {
  offset_days = 7
}
`
}
