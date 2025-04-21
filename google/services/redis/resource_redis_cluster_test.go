// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package redis_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/redis"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccRedisCluster_createUpdateClusterWithNodeType(t *testing.T) {

	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with node type "REDIS_STANDARD_SMALL"
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 3, deletionProtectionEnabled: true, nodeType: "REDIS_STANDARD_SMALL", zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "TUESDAY", maintenanceHours: 2, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with node type "REDIS_HIGHMEM_MEDIUM"
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 3, deletionProtectionEnabled: true, nodeType: "REDIS_HIGHMEM_MEDIUM", zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "TUESDAY", maintenanceHours: 2, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 3, deletionProtectionEnabled: false, nodeType: "REDIS_HIGHMEM_MEDIUM", zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "TUESDAY", maintenanceHours: 2, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
		},
	})
}

// Validate zone distribution for the cluster.
func TestAccRedisCluster_createClusterWithZoneDistribution(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with replica count 1
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "SINGLE_ZONE", zone: "us-central1-b"}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "SINGLE_ZONE", zone: "us-central1-b"}),
			},
		},
	})
}

// Validate that replica count is updated for the cluster
func TestAccRedisCluster_updateReplicaCount(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with replica count 1
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update replica count to 2
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 2, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update replica count to 0
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
		},
	})
}

// Validate that shard count is updated for the cluster
func TestAccRedisCluster_updateShardCount(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with shard count 3
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update shard count to 5
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 5, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 1, shardCount: 5, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
		},
	})
}

// Validate that redisConfigs is updated for the cluster
func TestAccRedisCluster_updateRedisConfigs(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster
				Config: createOrUpdateRedisCluster(&ClusterParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					maintenanceDay:       "MONDAY",
					maintenanceHours:     1,
					maintenanceMinutes:   0,
					maintenanceSeconds:   0,
					maintenanceNanos:     0,
					redisConfigs: map[string]string{
						"maxmemory-policy": "volatile-ttl",
					}}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// add a new redis config key-value pair and update existing redis config
				Config: createOrUpdateRedisCluster(&ClusterParams{
					name:                 name,
					shardCount:           3,
					zoneDistributionMode: "MULTI_ZONE",
					maintenanceDay:       "MONDAY",
					maintenanceHours:     1,
					maintenanceMinutes:   0,
					maintenanceSeconds:   0,
					maintenanceNanos:     0,
					redisConfigs: map[string]string{
						"maxmemory-policy":  "allkeys-lru",
						"maxmemory-clients": "90%",
					}}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// remove all redis configs
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, shardCount: 3, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0}),
			},
		},
	})
}

// Validate that deletion protection enabled/disabled cluster is created updated
func TestAccRedisCluster_createUpdateDeletionProtection(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with deletion protection set to false
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE"}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update deletion protection to true
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE"}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update deletion protection to false and delete the cluster
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE"}),
			},
		},
	})
}
func TestAccRedisCluster_automatedBackupConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisCluster_automatedBackupConfig(context),
			},
			{
				ResourceName:      "google_redis_cluster.cluster_abc",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisCluster_automatedBackupConfigWithout(context),
			},
			{
				ResourceName:      "google_redis_cluster.cluster_abc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedisCluster_automatedBackupConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster_abc" {
  name                           = "tf-test-redis-abc-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false
  automated_backup_config {
   retention                     = "259200s"
   fixed_frequency_schedule {
    start_time {
      hours                      = 20
    }
   }
  }

}   
`, context)
}

func testAccRedisCluster_automatedBackupConfigWithout(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster_abc" {
  name                           = "tf-test-redis-abc-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false 
  
}   
`, context)
}

// Validate that Import managedBackupSource can be used to create the cluster
func TestAccRedisCluster_managedBackupSource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"back_up":       "back_me_up",
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisCluster_managedBackupSourceSetup(context),
				Check: resource.ComposeTestCheckFunc(
					// Create an on-demand backup
					testAccCheckRedisClusterOnDemandBackup(t, "google_redis_cluster.cluster_mbs_main", context["back_up"].(string)),
				),
			},
			{
				ResourceName:      "google_redis_cluster.cluster_mbs_main",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisCluster_managedBackupSourceImport(context),
			},
			{
				ResourceName:      "google_redis_cluster.cluster_mbs_main",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedisCluster_managedBackupSourceSetup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster_mbs_main" {
  name                           = "tf-test-mbs-main-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false
}
`, context)
}

func testAccRedisCluster_managedBackupSourceImport(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster_mbs_main" {
  name                           = "tf-test-mbs-main-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false
}

resource "google_redis_cluster" "cluster_mb_copy" {
  name                           = "tf-test-mbs-copy-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false
   managed_backup_source {
    backup                       = join("", [google_redis_cluster.cluster_mbs_main.backup_collection , "/backups/%{back_up}"])
  }
}   
`, context)
}

// Takes the backup and verifies that the backup operation was successful.
func testAccCheckRedisClusterOnDemandBackup(t *testing.T, resourceName string, backupId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for Redis cluster")
		}

		config := acctest.GoogleProviderConfig(t)

		// Extract the cluster name, project, and region from the resource
		project, err := acctest.GetTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		region := rs.Primary.Attributes["region"]
		name := rs.Primary.Attributes["name"]

		// Construct the backup request
		backupRequest := map[string]interface{}{
			"backupId": backupId,
		}

		// Make the API call to create an on-demand backup
		url := fmt.Sprintf("https://redis.googleapis.com/v1/projects/%s/locations/%s/clusters/%s:backup", project, region, name)

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      backupRequest,
		})

		if err != nil {
			return fmt.Errorf("Error creating on-demand backup for Redis cluster %s: %s", name, err)
		}

		// Wait for the operation to complete
		err = redis.RedisOperationWaitTime(
			config, res, project, "Creating Redis Cluster Backup", config.UserAgent,
			time.Minute*20)

		// Check if the operation was successful
		if res == nil {
			return fmt.Errorf("Empty response when creating on-demand backup for Redis cluster %s", name)
		}

		return nil
	}
}

// Validate that Import gcsSource can be used to create the cluster
func TestAccRedisCluster_gcsSource(t *testing.T) {
	t.Parallel()
	randomSuffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix": randomSuffix,
		"back_up":       "back_me_up",
		"gcs_bucket":    fmt.Sprintf("tf-test-redis-backup-%s", randomSuffix),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisCluster_gcsSourceSetup(context),
				Check: resource.ComposeTestCheckFunc(
					// Create an on-demand backup
					testAccCheckRedisClusterOnDemandBackup(t, "google_redis_cluster.cluster_gbs_main", context["back_up"].(string)),
					// Export the backup to GCS
					testAccCheckRedisClusterExportBackup(t, "google_redis_cluster.cluster_gbs_main", context["back_up"].(string), context["gcs_bucket"].(string)),
				),
			},
			{
				ResourceName:      "google_redis_cluster.cluster_gbs_main",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisCluster_gcsSource(context),
			},
			{
				ResourceName:      "google_redis_cluster.cluster_gbs_main",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccRedisCluster_gcsSourceSetup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster_gbs_main" {
  name                           = "tf-test-gbs-main-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false
}

# Create a GCS bucket for exporting Redis backups
resource "google_storage_bucket" "redis_backup_bucket" {
  name                           = "%{gcs_bucket}"
  location                       = "us-central1"
  uniform_bucket_level_access    = true
  force_destroy                  = true
}

# Grant the Redis service account permission to access the bucket
# The Memorystore service account has the format:
# service-{project_number}@cloud-redis.iam.gserviceaccount.com
data "google_project" "project" {}

resource "google_storage_bucket_iam_member" "redis_backup_writer" {
  bucket 						= google_storage_bucket.redis_backup_bucket.name
  role   						= "roles/storage.admin"
  member 						= "serviceAccount:service-${data.google_project.project.number}@cloud-redis.iam.gserviceaccount.com"
}
`, context)
}

func testAccRedisCluster_gcsSource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster_gbs_main" {
  name                            = "tf-test-gbs-main-%{random_suffix}"
  shard_count                     = 1
  region                          = "us-central1"
  deletion_protection_enabled     = false
}

# Reference the bucket created in the setup
resource "google_storage_bucket" "redis_backup_bucket" {
  name                        	  = "%{gcs_bucket}"
  location                    	  = "us-central1"
  uniform_bucket_level_access 	  = true
  force_destroy               	  = true
}

# Grant the Redis service account permission to access the bucket
data "google_project" "project" {}

data "google_storage_bucket_objects" "backup" {
  bucket 					      = "%{gcs_bucket}"
}

resource "google_storage_bucket_iam_member" "redis_backup_writer" {
  bucket 						  = google_storage_bucket.redis_backup_bucket.name
  role   						  = "roles/storage.admin"
  member 						  = "serviceAccount:service-${data.google_project.project.number}@cloud-redis.iam.gserviceaccount.com"
}

# Create a Redis cluster that imports from the GCS bucket
resource "google_redis_cluster" "cluster_gbs" {
  name                           = "tf-test-gbs-copy-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false
  gcs_source {
    uris                         = [join("", ["gs://%{gcs_bucket}/" , data.google_storage_bucket_objects.backup.bucket_objects[0]["name"]])]
  }
  depends_on                     = [google_storage_bucket_iam_member.redis_backup_writer]
}
`, context)
}

func testAccCheckRedisClusterExportBackup(t *testing.T, resourceName string, backupId string, gcsDestination string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] Starting Redis Cluster backup export for resource %s, backup %s to %s", resourceName, backupId, gcsDestination)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for Redis cluster")
		}

		log.Printf("[DEBUG] Resource state: %#v", rs.Primary)

		config := acctest.GoogleProviderConfig(t)

		// Extract the cluster name, project, and region from the resource
		project, err := acctest.GetTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		region := rs.Primary.Attributes["region"]
		name := rs.Primary.Attributes["name"]

		log.Printf("[DEBUG] Exporting backup for cluster: project=%s, region=%s, name=%s", project, region, name)

		// Step 1: Find the backup collection ID for this cluster
		// First, list all backup collections in this region
		listCollectionsUrl := fmt.Sprintf("https://redis.googleapis.com/v1/projects/%s/locations/%s/backupCollections",
			project, region)

		log.Printf("[DEBUG] Listing backup collections URL: %s", listCollectionsUrl)

		collectionsRes, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    listCollectionsUrl,
			UserAgent: config.UserAgent,
		})

		if err != nil {
			log.Printf("[ERROR] Error listing backup collections: %s", err)
			return fmt.Errorf("Error listing backup collections: %s", err)
		}

		log.Printf("[DEBUG] Backup collections response: %#v", collectionsRes)

		// Find the backup collection that belongs to our cluster
		backupCollectionId := ""
		if collections, ok := collectionsRes["backupCollections"].([]interface{}); ok {
			log.Printf("[DEBUG] Found %d backup collections", len(collections))
			for i, collection := range collections {
				if collectionMap, ok := collection.(map[string]interface{}); ok {
					// The backup collection name format is projects/{project}/locations/{location}/backupCollections/{backupCollection}
					cluster := collectionMap["cluster"].(string)
					log.Printf("[DEBUG] CLuster %d Long name: %s", i, cluster)

					parts := strings.Split(cluster, "/")
					cluster_name := parts[len(parts)-1]

					log.Printf("[DEBUG] CLuster %d name: %s", i, cluster_name)
					log.Printf("[DEBUG] Provided Cluster Name  name: %s", name)

					if strings.Contains(cluster_name, name) {
						collection_id_long := collectionMap["name"].(string)
						parts := strings.Split(collection_id_long, "/")
						backupCollectionId = parts[len(parts)-1]
						log.Printf("[DEBUG] Found collection ID: %s for cluster %s ", backupCollectionId, cluster_name)
						break
					}

				}
			}

		} else {
			log.Printf("[DEBUG] No 'backupCollections' field found in response or it's not a slice")
		}

		if backupCollectionId == "" {
			return fmt.Errorf("Could not find backup collection for cluster %s", name)
		}

		exportRequest := map[string]interface{}{
			"gcsBucket": gcsDestination,
		}

		log.Printf("[DEBUG] Export request: %#v", exportRequest)

		exportUrl := fmt.Sprintf("https://redis.googleapis.com/v1/projects/%s/locations/%s/backupCollections/%s/backups/%s:export",
			project, region, backupCollectionId, backupId)

		log.Printf("[DEBUG] Export URL: %s", exportUrl)

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    exportUrl,
			UserAgent: config.UserAgent,
			Body:      exportRequest,
		})

		if err != nil {
			log.Printf("[ERROR] Error initiating backup export: %s", err)
			return fmt.Errorf("Error initiating backup export: %s", err)
		}

		log.Printf("[DEBUG] Export response: %#v", res)

		// Wait for the export operation to complete
		log.Printf("[DEBUG] Waiting for export operation to complete")
		err = redis.RedisOperationWaitTime(
			config, res, project, "Exporting Redis Cluster Backup", config.UserAgent,
			time.Minute*20)

		if err != nil {
			log.Printf("[ERROR] Error during backup export operation: %s", err)
			return fmt.Errorf("Error during backup export operation: %s", err)
		}

		return nil
	}
}

// Validate that persistence is updated for the cluster
func TestAccRedisCluster_persistenceUpdate(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with AOF enabled
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, nodeType: "REDIS_STANDARD_SMALL", zoneDistributionMode: "MULTI_ZONE", persistenceBlock: "persistence_config {\nmode = \"AOF\"\naof_config{\nappend_fsync = \"EVERYSEC\"\n}\n}"}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// disable AOF
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, nodeType: "REDIS_STANDARD_SMALL", zoneDistributionMode: "MULTI_ZONE", persistenceBlock: "persistence_config {\nmode = \"DISABLED\"\n}"}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			}, {
				// update persistence to RDB
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, nodeType: "REDIS_STANDARD_SMALL", zoneDistributionMode: "MULTI_ZONE", persistenceBlock: "persistence_config {\nmode = \"RDB\"\nrdb_config {\nrdb_snapshot_period = \"ONE_HOUR\"\nrdb_snapshot_start_time = \"2024-10-02T15:01:23Z\"\n}\n}"}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, nodeType: "REDIS_STANDARD_SMALL", zoneDistributionMode: "MULTI_ZONE", persistenceBlock: "persistence_config {\nmode = \"RDB\"\nrdb_config {\nrdb_snapshot_period = \"ONE_HOUR\"\nrdb_snapshot_start_time = \"2024-10-02T15:01:23Z\"\n}\n}"}),
			},
		},
	})
}

// Validate that deletion protection enabled/disabled cluster is created updated
func TestAccRedisCluster_switchoverAndDetachSecondary(t *testing.T) {
	t.Parallel()

	pcName := fmt.Sprintf("tf-test-prim-%d", acctest.RandInt(t))
	scName := fmt.Sprintf("tf-test-sec-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create primary and secondary clusters cluster
				Config: createOrUpdateRedisCluster(&ClusterParams{name: pcName, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", shouldCreateSecondary: true, secondaryClusterName: scName, ccrRole: "SECONDARY"}),
			},
			{
				ResourceName:            "google_redis_cluster.test_secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// // Switchover to secondary cluster
				Config: createOrUpdateRedisCluster(&ClusterParams{name: pcName, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", shouldCreateSecondary: true, secondaryClusterName: scName, ccrRole: "PRIMARY"}),
			},
			{
				ResourceName:            "google_redis_cluster.test_secondary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// Detach secondary cluster and delete the clusters
				Config: createOrUpdateRedisCluster(&ClusterParams{name: pcName, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", shouldCreateSecondary: true, secondaryClusterName: scName, ccrRole: "NONE"}),
			},
		},
	})
}

// Validate that cluster endpoints are updated for the cluster
func TestAccRedisCluster_updateClusterEndpoints(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with no user created connections
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// create cluster with one user created connection
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 1}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 2 endpoints
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 2}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 1 endpoint
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 1}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// update cluster with 0 endpoints
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: true, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 0}),
			},
			{
				ResourceName:            "google_redis_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psc_configs"},
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(&ClusterParams{name: name, replicaCount: 0, shardCount: 3, deletionProtectionEnabled: false, zoneDistributionMode: "MULTI_ZONE", maintenanceDay: "MONDAY", maintenanceHours: 1, maintenanceMinutes: 0, maintenanceSeconds: 0, maintenanceNanos: 0, userEndpointCount: 0}),
			},
		},
	})
}

type ClusterParams struct {
	name                      string
	replicaCount              int
	shardCount                int
	deletionProtectionEnabled bool
	nodeType                  string
	redisConfigs              map[string]string
	zoneDistributionMode      string
	zone                      string
	maintenanceDay            string
	maintenanceHours          int
	maintenanceMinutes        int
	maintenanceSeconds        int
	maintenanceNanos          int
	persistenceBlock          string
	shouldCreateSecondary     bool
	secondaryClusterName      string
	ccrRole                   string
	userEndpointCount         int
}

func createRedisClusterEndpoints(params *ClusterParams) string {
	if params.userEndpointCount == 2 {
		return createRedisClusterEndpointsWithTwoUserCreatedConnections(params)
	} else if params.userEndpointCount == 1 {
		return createRedisClusterEndpointsWithOneUserCreatedConnections(params)
	}
	return ``
}

func createRedisClusterEndpointsWithOneUserCreatedConnections(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_redis_cluster_user_created_connections" "default" {
		
		name = "%s"
		region = "us-central1"
		cluster_endpoints {
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule1_network1.psc_connection_id
					address = google_compute_address.ip1_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule1_network1.id
					network = google_compute_network.network1.id
					project_id = data.google_project.project.project_id
					service_attachment = google_redis_cluster.test.psc_service_attachments[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule2_network1.psc_connection_id
					address = google_compute_address.ip2_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule2_network1.id
					network = google_compute_network.network1.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[1].service_attachment
				}
			}
		}
		}
		data "google_project" "project" {
		}
		%s
		`,
		params.name,
		createRedisClusterUserCreatedConnection1(params),
	)

}

func createRedisClusterEndpointsWithTwoUserCreatedConnections(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_redis_cluster_user_created_connections" "default" {
		name = "%s"
		region = "us-central1"
		cluster_endpoints {
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule1_network1.psc_connection_id
					address = google_compute_address.ip1_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule1_network1.id
					network = google_compute_network.network1.id
					project_id = data.google_project.project.project_id
					service_attachment = google_redis_cluster.test.psc_service_attachments[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule2_network1.psc_connection_id
					address = google_compute_address.ip2_network1.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule2_network1.id
					network = google_compute_network.network1.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[1].service_attachment
				}
			}
		}
		cluster_endpoints {
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule1_network2.psc_connection_id
					address = google_compute_address.ip1_network2.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule1_network2.id
					network = google_compute_network.network2.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[0].service_attachment
				}
			}
			connections {
				psc_connection {
					psc_connection_id = google_compute_forwarding_rule.forwarding_rule2_network2.psc_connection_id
					address = google_compute_address.ip2_network2.address
					forwarding_rule = google_compute_forwarding_rule.forwarding_rule2_network2.id
					network = google_compute_network.network2.id
					service_attachment = google_redis_cluster.test.psc_service_attachments[1].service_attachment
				}
			}
		}
		}
		data "google_project" "project" {
		}
		%s
		%s
		`,
		params.name,
		createRedisClusterUserCreatedConnection1(params),
		createRedisClusterUserCreatedConnection2(params),
	)
}

func createRedisClusterUserCreatedConnection1(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_compute_forwarding_rule" "forwarding_rule1_network1" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip1_network1.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network1.id
		target                 = google_redis_cluster.test.psc_service_attachments[0].service_attachment
		}

		resource "google_compute_forwarding_rule" "forwarding_rule2_network1" {	
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip2_network1.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network1.id
		target                 = google_redis_cluster.test.psc_service_attachments[1].service_attachment
		}

		resource "google_compute_address" "ip1_network1" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network1.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_address" "ip2_network1" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network1.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_subnetwork" "subnet_network1" {
		name          = "%s"
		ip_cidr_range = "10.0.0.248/29"
		region        = "us-central1"
		network       = google_compute_network.network1.id
		}

		resource "google_compute_network" "network1" {
		name                    = "%s"
		auto_create_subnetworks = false
		}
		`,
		params.name+"-11", // fwd-rule1-net1
		params.name+"-12", // fwd-rule2-net1
		params.name+"-11", // ip1-net1
		params.name+"-12", // ip2-net1
		params.name+"-1",  // subnet-net1
		params.name+"-1",  // net1
	)
}

func createRedisClusterUserCreatedConnection2(params *ClusterParams) string {
	return fmt.Sprintf(`
		resource "google_compute_forwarding_rule" "forwarding_rule1_network2" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip1_network2.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network2.id
		target                 = google_redis_cluster.test.psc_service_attachments[0].service_attachment
		}

		resource "google_compute_forwarding_rule" "forwarding_rule2_network2" {
		name                   = "%s"
		region                 = "us-central1"
		ip_address             = google_compute_address.ip2_network2.id
		load_balancing_scheme  = ""
		network                = google_compute_network.network2.id
		target                 = google_redis_cluster.test.psc_service_attachments[1].service_attachment
		}

		resource "google_compute_address" "ip1_network2" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network2.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_address" "ip2_network2" {
		name         = "%s"
		region       = "us-central1"
		subnetwork   = google_compute_subnetwork.subnet_network2.id
		address_type = "INTERNAL"
		purpose      = "GCE_ENDPOINT"
		}

		resource "google_compute_subnetwork" "subnet_network2" {
		name          = "%s"
		ip_cidr_range = "10.0.0.248/29"
		region        = "us-central1"
		network       = google_compute_network.network2.id
		}

		resource "google_compute_network" "network2" {
		name                    = "%s"
		auto_create_subnetworks = false
		}
		`,
		params.name+"-21", // fwd-rule1-net2
		params.name+"-22", // fwd-rule2-net2
		params.name+"-21", // ip1-net2
		params.name+"-22", // ip2-net2
		params.name+"-2",  // subnet-net2
		params.name+"-2",  // net2
	)

}

func createOrUpdateRedisCluster(params *ClusterParams) string {
	clusterResourceBlock := createRedisClusterResourceConfig(params /*isSecondaryCluster*/, false)
	secClusterResourceBlock := ``
	if params.shouldCreateSecondary {
		secClusterResourceBlock = createRedisClusterResourceConfig(params /*isSecondaryCluster*/, true)
	}

	endpointBlock := createRedisClusterEndpoints(params)

	return fmt.Sprintf(`
		%s
		%s
		%s
		resource "google_network_connectivity_service_connection_policy" "default" {
			name = "%s"
			location = "us-central1"
			service_class = "gcp-memorystore-redis"
			description   = "my basic service connection policy"
			network = google_compute_network.producer_net.id
			psc_config {
			subnetworks = [google_compute_subnetwork.producer_subnet.id]
			}
		}

		resource "google_compute_subnetwork" "producer_subnet" {
			name          = "%s"
			ip_cidr_range = "10.0.0.16/28"
			region        = "us-central1"
			network       = google_compute_network.producer_net.id
		}

		resource "google_compute_network" "producer_net" {
			name                    = "%s"
			auto_create_subnetworks = false
		}
		`,
		endpointBlock,
		clusterResourceBlock,
		secClusterResourceBlock,
		params.name,
		params.name,
		params.name)
}

func createRedisClusterResourceConfig(params *ClusterParams, isSecondaryCluster bool) string {
	tfClusterResourceName := "test"
	clusterName := params.name
	dependsOnBlock := "google_network_connectivity_service_connection_policy.default"

	var redsConfigsStrBuilder strings.Builder
	for key, value := range params.redisConfigs {
		redsConfigsStrBuilder.WriteString(fmt.Sprintf("%s =  \"%s\"\n", key, value))
	}

	zoneDistributionConfigBlock := ``
	if params.zoneDistributionMode != "" {
		zoneDistributionConfigBlock = fmt.Sprintf(`
		zone_distribution_config {
			mode = "%s"
			zone = "%s"
		}
		`, params.zoneDistributionMode, params.zone)
	}

	maintenancePolicyBlock := ``
	if params.maintenanceDay != "" {
		maintenancePolicyBlock = fmt.Sprintf(`
		maintenance_policy {
			weekly_maintenance_window {
				day = "%s"
				start_time {
					hours = %d
					minutes = %d
					seconds = %d
					nanos = %d
				}
			}
		}
		`, params.maintenanceDay, params.maintenanceHours, params.maintenanceMinutes, params.maintenanceSeconds, params.maintenanceNanos)
	}

	crossClusterReplicationConfigBlock := ``
	if isSecondaryCluster {
		tfClusterResourceName = "test_secondary"
		clusterName = params.secondaryClusterName
		dependsOnBlock = dependsOnBlock + ", google_redis_cluster.test"

		// Construct cross_cluster_replication_config block
		pcBlock := ``
		scsBlock := ``
		if params.ccrRole == "SECONDARY" {
			pcBlock = fmt.Sprintf(`
			primary_cluster {
				cluster = google_redis_cluster.test.id
			}
			`)
		} else if params.ccrRole == "PRIMARY" {
			scsBlock = fmt.Sprintf(`
			secondary_clusters {
				cluster = google_redis_cluster.test.id
			}
			`)
		}
		crossClusterReplicationConfigBlock = fmt.Sprintf(`
		cross_cluster_replication_config {
			cluster_role = "%s"
			%s
			%s
		}
		`, params.ccrRole, pcBlock, scsBlock)
	}

	return fmt.Sprintf(`
		resource "google_redis_cluster" "%s" {
		name           = "%s"
		replica_count = %d
		shard_count = %d
		node_type = "%s"
		deletion_protection_enabled = %v
		region         = "us-central1"
		psc_configs {
				network = google_compute_network.producer_net.id
		}
		redis_configs = {
			%s
		}
		%s
		%s
		%s
		%s
		depends_on = [
				%s
			]
		}
		`,
		tfClusterResourceName,
		clusterName,
		params.replicaCount,
		params.shardCount,
		params.nodeType,
		params.deletionProtectionEnabled,
		redsConfigsStrBuilder.String(),
		zoneDistributionConfigBlock,
		maintenancePolicyBlock,
		params.persistenceBlock,
		crossClusterReplicationConfigBlock,
		dependsOnBlock)
}
