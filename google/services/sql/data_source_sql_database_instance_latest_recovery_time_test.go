// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	resourceName := "data.google_sql_database_instance_latest_recovery_time.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "instance"),
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "latest_recovery_time"),
				),
			},
		},
	})
}

func testAccDataSourceSqlDatabaseInstanceLatestRecoveryTime_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    tier = "db-g1-small"
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "20:55"
      transaction_log_retention_days = "3"
    }
  }

  deletion_protection = false
}

data "google_sql_database_instance_latest_recovery_time" "default" {
  instance = resource.google_sql_database_instance.main.name
}
`, context)
}
