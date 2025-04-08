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
