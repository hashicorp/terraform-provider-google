// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
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
