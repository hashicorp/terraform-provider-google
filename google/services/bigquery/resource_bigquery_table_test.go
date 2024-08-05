// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigQueryTable_Basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableTimePartitioning(datasetID, tableID, "DAY"),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableUpdated(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_DropColumns(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableTimePartitioningDropColumns(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableTimePartitioningDropColumnsUpdate(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_Kms(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	kms := acctest.BootstrapKMSKey(t)
	cryptoKeyName := kms.CryptoKey.Name

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableKms(cryptoKeyName, datasetID, tableID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_HourlyTimePartitioning(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableTimePartitioning(datasetID, tableID, "HOUR"),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableUpdated(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_MonthlyTimePartitioning(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableTimePartitioning(datasetID, tableID, "MONTH"),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableUpdated(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_YearlyTimePartitioning(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableTimePartitioning(datasetID, tableID, "YEAR"),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableUpdated(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_HivePartitioning(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableHivePartitioning(bucketName, datasetID, tableID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_HivePartitioningCustomSchema(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableHivePartitioningCustomSchema(bucketName, datasetID, tableID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_data_configuration.0.schema", "deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_AvroPartitioning(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	avroFilePath := "./test-fixtures/avro-generated.avro"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableAvroPartitioning(bucketName, avroFilePath, datasetID, tableID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_json(t *testing.T) {
	t.Parallel()
	bucketName := acctest.TestBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableJson(datasetID, tableID, bucketName, "UTF-8"),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_data_configuration.0.schema", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableJson(datasetID, tableID, bucketName, "UTF-16BE"),
			},
		},
	})
}

func TestAccBigQueryTable_RangePartitioning(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableRangePartitioning(datasetID, tableID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_PrimaryKey(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTablePrimaryKey(datasetID, tableID),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_ForeignKey(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID_pk := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID_fk := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableForeignKeys(projectID, datasetID, tableID_pk, tableID_fk),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_updateTableConstraints(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID_pk := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID_fk := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableForeignKeys(projectID, datasetID, tableID_pk, tableID_fk),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableTableConstraintsUpdate(projectID, datasetID, tableID_pk, tableID_fk),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_View(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithView(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_updateView(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithView(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableWithNewSqlView(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_WithViewAndSchema(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithViewAndSchema(datasetID, tableID, "table description1"),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableWithViewAndSchema(datasetID, tableID, "table description2"),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_MaterializedView_DailyTimePartioning_Basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	materialized_viewID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	query := fmt.Sprintf("SELECT count(some_string) as count, some_int, ts FROM `%s.%s` WHERE DATE(ts) = '2019-01-01' GROUP BY some_int, ts", datasetID, tableID)
	queryNew := strings.ReplaceAll(query, "2019", "2020")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithMatViewDailyTimePartitioning_basic(datasetID, tableID, materialized_viewID, query),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				ResourceName:            "google_bigquery_table.mv_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableWithMatViewDailyTimePartitioning_basic(datasetID, tableID, materialized_viewID, queryNew),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				ResourceName:            "google_bigquery_table.mv_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_MaterializedView_DailyTimePartioning_Update(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	materialized_viewID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	query := fmt.Sprintf("SELECT count(some_string) as count, some_int, ts FROM `%s.%s` WHERE DATE(ts) = '2019-01-01' GROUP BY some_int, ts", datasetID, tableID)

	enable_refresh := "false"
	refresh_interval_ms := "3600000"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithMatViewDailyTimePartitioning_basic(datasetID, tableID, materialized_viewID, query),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				ResourceName:            "google_bigquery_table.mv_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableWithMatViewDailyTimePartitioning(datasetID, tableID, materialized_viewID, enable_refresh, refresh_interval_ms, query),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				ResourceName:            "google_bigquery_table.mv_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_MaterializedView_NonIncremental_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	materialized_viewID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	query := fmt.Sprintf("SELECT count(some_string) as count, some_int, ts FROM `%s.%s` WHERE DATE(ts) = '2019-01-01' GROUP BY some_int, ts", datasetID, tableID)
	maxStaleness := "0-0 0 10:0:0"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithMatViewNonIncremental_basic(datasetID, tableID, materialized_viewID, query, maxStaleness),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection", "require_partition_filter", "time_partitioning.0.require_partition_filter"},
			},
			{
				ResourceName:            "google_bigquery_table.mv_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection", "require_partition_filter", "time_partitioning.0.require_partition_filter"},
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_deltaLake(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSDeltaLake(datasetID, tableID, bucketName),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_parquet(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.gz.parquet", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSParquet(datasetID, tableID, bucketName, objectName),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_parquetOptions(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.gz.parquet", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSParquetOptions(datasetID, tableID, bucketName, objectName, true, true),
			},
			{
				Config: testAccBigQueryTableFromGCSParquetOptions(datasetID, tableID, bucketName, objectName, false, false),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_iceberg(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSIceberg(datasetID, tableID, bucketName),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_parquetFileSetSpecType(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	parquetFileName := "test.parquet"
	manifestName := fmt.Sprintf("tf_test_%s.manifest.json", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSParquetManifest(datasetID, tableID, bucketName, manifestName, parquetFileName),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_queryAcceleration(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.gz.parquet", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	metadataCacheMode := "AUTOMATIC"
	// including an optional field. Should work without specifiying.
	// Has to follow google sql IntervalValue encoding
	maxStaleness := "0-0 0 10:0:0"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSParquetWithQueryAcceleration(connectionID, datasetID, tableID, bucketName, objectName, metadataCacheMode, maxStaleness),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_objectTable(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	// including an optional field. Should work without specifiying.
	// Has to follow google sql IntervalValue encoding
	maxStaleness := "0-0 0 10:0:0"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSObjectTable(connectionID, datasetID, tableID, bucketName, objectName, maxStaleness),
			},
			{
				Config: testAccBigQueryTableFromGCSObjectTableMetadata(connectionID, datasetID, tableID, bucketName, objectName, maxStaleness),
			},
			{
				Config: testAccBigQueryTableFromGCSObjectTable(connectionID, datasetID, tableID, bucketName, objectName, maxStaleness),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_connectionIdDiff_UseNameReference(t *testing.T) {
	t.Parallel()
	// Setup
	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	// Feature Under Test.
	location := "US"
	connection_id_reference := "google_bigquery_connection.test.name"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableExternalDataConfigurationConnectionID(location, connectionID, datasetID, tableID, bucketName, objectName, connection_id_reference),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_connectionIdDiff_UseIdReference(t *testing.T) {
	t.Parallel()
	// Setup
	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	// Feature Under Test.
	location := "US"
	connection_id_reference := "google_bigquery_connection.test.id"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableExternalDataConfigurationConnectionID(location, connectionID, datasetID, tableID, bucketName, objectName, connection_id_reference),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_connectionIdDiff_UseIdReference_UsCentral1LowerCase(t *testing.T) {
	t.Parallel()
	// Setup
	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	// Feature Under Test.
	location := "us-central1"
	connection_id_reference := "google_bigquery_connection.test.id"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableExternalDataConfigurationConnectionID(location, connectionID, datasetID, tableID, bucketName, objectName, connection_id_reference),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_connectionIdDiff_UseIdReference_UsEast1(t *testing.T) {
	t.Parallel()
	// Setup
	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	// Feature Under Test.
	location := "US-EAST1"
	connection_id_reference := "google_bigquery_connection.test.id"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableExternalDataConfigurationConnectionID(location, connectionID, datasetID, tableID, bucketName, objectName, connection_id_reference),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_connectionIdDiff_UseIdReference_EuropeWest8(t *testing.T) {
	t.Parallel()
	// Setup
	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	// Feature Under Test.
	location := "EUROPE-WEST8"
	connection_id_reference := "google_bigquery_connection.test.id"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableExternalDataConfigurationConnectionID(location, connectionID, datasetID, tableID, bucketName, objectName, connection_id_reference),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_CSV(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCS(datasetID, tableID, bucketName, objectName, TEST_CSV, "CSV", "\\\""),
				Check:  testAccCheckBigQueryExtData(t, "\""),
			},
			{
				Config: testAccBigQueryTableFromGCS(datasetID, tableID, bucketName, objectName, TEST_CSV, "CSV", ""),
				Check:  testAccCheckBigQueryExtData(t, ""),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_CSV_WithSchema_InvalidSchemas(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryTableFromGCSWithExternalDataConfigSchema(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_INVALID_SCHEMA_NOT_JSON),
				ExpectError: regexp.MustCompile("contains an invalid JSON"),
			},
			{
				Config:      testAccBigQueryTableFromGCSWithExternalDataConfigSchema(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_INVALID_SCHEMA_NOT_JSON_LIST),
				ExpectError: regexp.MustCompile("\"schema\" is not a JSON array"),
			},
			{
				Config:      testAccBigQueryTableFromGCSWithExternalDataConfigSchema(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_INVALID_SCHEMA_JSON_LIST_WITH_NULL_ELEMENT),
				ExpectError: regexp.MustCompile("\"schema\" contains a nil element"),
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_CSV_WithSchemaAndConnectionID_UpdateNoConnectionID(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSWithSchemaWithConnectionId(datasetID, tableID, connectionID, projectID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableFromGCSWithSchema(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_CSV_WithSchema_UpdateToConnectionID(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSWithSchema(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableFromGCSWithSchemaWithConnectionId(datasetID, tableID, connectionID, projectID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableFromGCSWithSchemaWithConnectionId2(datasetID, tableID, connectionID, projectID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryExternalDataTable_CSV_WithSchema_UpdateAllowQuotedNewlines(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", acctest.RandString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromGCSWithSchema(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTableFromGCSWithSchema_UpdatAllowQuotedNewlines(datasetID, tableID, bucketName, objectName, TEST_SIMPLE_CSV, TEST_SIMPLE_CSV_SCHEMA),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryDataTable_bigtable(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 8),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromBigtable(context),
			},
			{
				ResourceName:            "google_bigquery_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryDataTable_bigtable_options(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 8),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromBigtableOptions(context),
			},
			{
				ResourceName:            "google_bigquery_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableFromBigtable(context),
			},
		},
	})
}

func TestAccBigQueryDataTable_sheet(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableFromSheet(context),
			},
			{
				ResourceName:            "google_bigquery_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryDataTable_jsonEquivalency(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_jsonEq(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection", "labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryTable_jsonEqModeRemoved(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataTable_canReorderParameters(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// we don't run any checks because the resource will error out if
				// it attempts to destroy/tear down.
				Config: testAccBigQueryTable_jsonPreventDestroy(datasetID, tableID),
			},
			{
				Config: testAccBigQueryTable_jsonPreventDestroyOrderChanged(datasetID, tableID),
			},
			{
				Config: testAccBigQueryTable_jsonEq(datasetID, tableID),
			},
		},
	})
}

func TestAccBigQueryDataTable_expandArray(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_arrayInitial(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection", "labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryTable_arrayExpanded(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection", "labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryTable_allowDestroy(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_noAllowDestroy(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "labels", "terraform_labels"},
			},
			{
				Config:      testAccBigQueryTable_noAllowDestroy(datasetID, tableID),
				Destroy:     true,
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccBigQueryTable_noAllowDestroyUpdated(datasetID, tableID),
			},
		},
	})
}

func TestAccBigQueryTable_emptySchema(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_mimicCreateFromConsole(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTable_emptySchema(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_Update_SchemaWithoutPolicyTagsToWithPolicyTags(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableBasicSchema(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableBasicSchemaWithPolicyTags(datasetID, tableID, projectID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_Update_SchemaWithPolicyTagsToNoPolicyTag(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableBasicSchemaWithPolicyTags(datasetID, tableID, projectID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableBasicSchema(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_Update_SchemaWithPolicyTagsToEmptyPolicyTag(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableBasicSchemaWithPolicyTags(datasetID, tableID, projectID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableBasicSchemaWithEmptyPolicyTags(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_Update_SchemaWithPolicyTagsToEmptyPolicyTagNames(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	projectID := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableBasicSchemaWithPolicyTags(datasetID, tableID, projectID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
			{
				Config: testAccBigQueryTableBasicSchemaWithEmptyPolicyTagNames(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_invalidSchemas(t *testing.T) {
	t.Parallel()
	// Pending VCR support in https://github.com/hashicorp/terraform-provider-google/issues/15427.
	acctest.SkipIfVcr(t)

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryTableWithSchema(datasetID, tableID, TEST_INVALID_SCHEMA_NOT_JSON),
				ExpectError: regexp.MustCompile("contains an invalid JSON"),
			},
			{
				Config:      testAccBigQueryTableWithSchema(datasetID, tableID, TEST_INVALID_SCHEMA_NOT_JSON_LIST),
				ExpectError: regexp.MustCompile("\"schema\" is not a JSON array"),
			},
			{
				Config:      testAccBigQueryTableWithSchema(datasetID, tableID, TEST_INVALID_SCHEMA_JSON_LIST_WITH_NULL_ELEMENT),
				ExpectError: regexp.MustCompile("\"schema\" contains a nil element"),
			},
		},
	})
}

func TestAccBigQueryTable_schemaWithRequiredFieldAndView(t *testing.T) {
	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryTableWithSchemaWithRequiredFieldAndView(datasetID, tableID),
				ExpectError: regexp.MustCompile("Schema cannot contain required fields when creating a view"),
			},
		},
	})
}

func TestAccBigQueryTable_TableReplicationInfo_ConflictsWithView(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryTableWithReplicationInfoAndView(datasetID, tableID),
				ExpectError: regexp.MustCompile("Schema, view, or materialized view cannot be specified when table replication info is present"),
			},
		},
	})
}

func TestAccBigQueryTable_TableReplicationInfo_WithoutReplicationInterval(t *testing.T) {
	t.Parallel()

	projectID := envvar.GetTestProjectFromEnv()

	sourceDatasetID := fmt.Sprintf("tf_test_source_dataset_%s", acctest.RandString(t, 10))
	sourceTableID := fmt.Sprintf("tf_test_source_table_%s", acctest.RandString(t, 10))
	sourceMVID := fmt.Sprintf("tf_test_source_mv_%s", acctest.RandString(t, 10))
	replicaDatasetID := fmt.Sprintf("tf_test_replica_dataset_%s", acctest.RandString(t, 10))
	replicaMVID := fmt.Sprintf("tf_test_replica_mv_%s", acctest.RandString(t, 10))
	sourceMVJobID := fmt.Sprintf("tf_test_create_source_mv_job_%s", acctest.RandString(t, 10))
	dropMVJobID := fmt.Sprintf("tf_test_drop_source_mv_job_%s", acctest.RandString(t, 10))
	replicationIntervalExpr := ""
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithReplicationInfo(projectID, sourceDatasetID, sourceTableID, sourceMVID, replicaDatasetID, replicaMVID, sourceMVJobID, dropMVJobID, replicationIntervalExpr),
			},
			{
				ResourceName:            "google_bigquery_table.replica_mv",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_TableReplicationInfo_WithReplicationInterval(t *testing.T) {
	t.Parallel()

	projectID := envvar.GetTestProjectFromEnv()

	sourceDatasetID := fmt.Sprintf("tf_test_source_dataset_%s", acctest.RandString(t, 10))
	sourceTableID := fmt.Sprintf("tf_test_source_table_%s", acctest.RandString(t, 10))
	sourceMVID := fmt.Sprintf("tf_test_source_mv_%s", acctest.RandString(t, 10))
	replicaDatasetID := fmt.Sprintf("tf_test_replica_dataset_%s", acctest.RandString(t, 10))
	replicaMVID := fmt.Sprintf("tf_test_replica_mv_%s", acctest.RandString(t, 10))
	sourceMVJobID := fmt.Sprintf("tf_test_create_source_mv_job_%s", acctest.RandString(t, 10))
	dropMVJobID := fmt.Sprintf("tf_test_drop_source_mv_job_%s", acctest.RandString(t, 10))
	replicationIntervalExpr := "replication_interval_ms = 600000"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithReplicationInfo(projectID, sourceDatasetID, sourceTableID, sourceMVID, replicaDatasetID, replicaMVID, sourceMVJobID, dropMVJobID, replicationIntervalExpr),
			},
			{
				ResourceName:            "google_bigquery_table.replica_mv",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
			},
		},
	})
}

func TestAccBigQueryTable_ResourceTags(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":      envvar.GetTestProjectFromEnv(),
		"dataset_id":      fmt.Sprintf("tf_test_dataset_%s", acctest.RandString(t, 10)),
		"table_id":        fmt.Sprintf("tf_test_table_%s", acctest.RandString(t, 10)),
		"tag_key_name1":   fmt.Sprintf("tf_test_tag_key1_%s", acctest.RandString(t, 10)),
		"tag_value_name1": fmt.Sprintf("tf_test_tag_value1_%s", acctest.RandString(t, 10)),
		"tag_key_name2":   fmt.Sprintf("tf_test_tag_key2_%s", acctest.RandString(t, 10)),
		"tag_value_name2": fmt.Sprintf("tf_test_tag_value2_%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTableWithResourceTags(context),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "allow_resource_tags_on_deletion"},
			},
			{
				Config: testAccBigQueryTableWithResourceTagsUpdate(context),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "allow_resource_tags_on_deletion"},
			},
			// testAccBigQueryTableWithResourceTagsDestroy must be called at the end of this test to clear the resource tag bindings of the table before deletion.
			{
				Config: testAccBigQueryTableWithResourceTagsDestroy(context),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "allow_resource_tags_on_deletion"},
			},
		},
	})
}

func testAccCheckBigQueryExtData(t *testing.T, expectedQuoteChar string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_table" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			dataset := rs.Primary.Attributes["dataset_id"]
			table := rs.Primary.Attributes["table_id"]
			res, err := config.NewBigQueryClient(config.UserAgent).Tables.Get(config.Project, dataset, table).Do()
			if err != nil {
				return err
			}

			if res.Type != "EXTERNAL" {
				return fmt.Errorf("Table \"%s.%s\" is of type \"%s\", expected EXTERNAL.", dataset, table, res.Type)
			}
			edc := res.ExternalDataConfiguration
			cvsOpts := edc.CsvOptions
			if cvsOpts == nil || *cvsOpts.Quote != expectedQuoteChar {
				return fmt.Errorf("Table \"%s.%s\" quote should be '%s' but was '%s'", dataset, table, expectedQuoteChar, *cvsOpts.Quote)
			}
		}
		return nil
	}
}

func testAccCheckBigQueryTableDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_table" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			_, err := config.NewBigQueryClient(config.UserAgent).Tables.Get(config.Project, rs.Primary.Attributes["dataset_id"], rs.Primary.Attributes["table_id"]).Do()
			if err == nil {
				return fmt.Errorf("Table still present")
			}
		}

		return nil
	}
}

func testAccBigQueryTableBasicSchema(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema = <<EOH
[
  {
    "name": "id",
    "type": "INTEGER"
  }
]
EOH

}
`, datasetID, tableID)
}

func testAccBigQueryTableBasicSchemaWithPolicyTags(datasetID, tableID, projectID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema = <<EOH
[
  {
    "name": "id",
    "type": "INTEGER",
    "policyTags": {
      "names": [
        "projects/%s/locations/us/taxonomies/123/policyTags/1"
      ]
    }
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
      {
        "name": "id",
        "type": "INTEGER",
        "policyTags": {
          "names": [
            "projects/%s/locations/us/taxonomies/123/policyTags/1"
          ]
        }
      },
      {
        "name": "coord",
        "type": "RECORD",
        "fields": [
          {
            "name": "lon",
            "type": "FLOAT",
            "policyTags": {
              "names": [
                "projects/%s/locations/us/taxonomies/123/policyTags/1"
              ]
            }
          }
        ]
      }
    ]
  }
]
EOH

}
`, datasetID, tableID, projectID, projectID, projectID)
}

func testAccBigQueryTableBasicSchemaWithEmptyPolicyTags(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema = <<EOH
[
  {
    "name": "id",
    "type": "INTEGER",
    "policyTags": {}
  }
]
EOH

}
`, datasetID, tableID)
}

func testAccBigQueryTableBasicSchemaWithEmptyPolicyTagNames(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema = <<EOH
[
  {
    "name": "id",
    "type": "INTEGER",
    "policyTags": {
      "names": []
    }
  }
]
EOH

}
`, datasetID, tableID)
}

func testAccBigQueryTableTimePartitioning(datasetID, tableID, partitioningType string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type                     = "%s"
    field                    = "ts"
    expiration_ms            = 1000
  }
  require_partition_filter = true
  clustering = ["some_int", "some_string"]
  schema     = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "some_string",
    "type": "STRING"
  },
  {
    "name": "some_int",
    "type": "INTEGER"
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "coord",
    "type": "RECORD",
    "fields": [
    {
    "name": "lon",
    "type": "FLOAT"
    }
    ]
  }
    ]
  }
]
EOH

}
`, datasetID, tableID, partitioningType)
}

func testAccBigQueryTableTimePartitioningDropColumns(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema     = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "some_string",
    "type": "STRING"
  },
  {
    "name": "some_int",
    "type": "INTEGER"
  }
]
EOH

}
`, datasetID, tableID)
}

func testAccBigQueryTableTimePartitioningDropColumnsUpdate(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  schema     = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  }
]
EOH

}
`, datasetID, tableID)
}

func testAccBigQueryTableKms(cryptoKeyName, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
    dataset_id = "%s"
}

data "google_bigquery_default_service_account" "acct" {}

resource "google_kms_crypto_key_iam_member" "allow" {
  crypto_key_id = "%s"
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:${data.google_bigquery_default_service_account.acct.email}"
  depends_on = ["google_bigquery_dataset.test"]
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"

  time_partitioning {
    type = "DAY"
    field = "ts"
  }

  encryption_configuration {
    kms_key_name = "${google_kms_crypto_key_iam_member.allow.crypto_key_id}"
  }

  schema = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "coord",
    "type": "RECORD",
    "fields": [
    {
    "name": "lon",
    "type": "FLOAT"
    }
    ]
  }
    ]
  }
]
EOH
}
`, datasetID, cryptoKeyName, tableID)
}

func testAccBigQueryTableHivePartitioning(bucketName, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "key1=20200330/init.csv"
  content = ";"
  bucket  = google_storage_bucket.test.name
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  external_data_configuration {
    source_format = "CSV"
    autodetect = true
    source_uris= ["gs://${google_storage_bucket.test.name}/*"]

    hive_partitioning_options {
      mode = "AUTO"
      source_uri_prefix = "gs://${google_storage_bucket.test.name}/"
      require_partition_filter = true
    }

  }
  depends_on = ["google_storage_bucket_object.test"]
}
`, bucketName, datasetID, tableID)
}

func testAccBigQueryTableHivePartitioningCustomSchema(bucketName, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "key1=20200330/data.json"
  content = "{\"name\":\"test\", \"last_modification\":\"2020-04-01\"}"
  bucket  = google_storage_bucket.test.name
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  external_data_configuration {
    source_format = "NEWLINE_DELIMITED_JSON"
    autodetect = false
    source_uris= ["gs://${google_storage_bucket.test.name}/*"]

    hive_partitioning_options {
      mode = "CUSTOM"
      source_uri_prefix = "gs://${google_storage_bucket.test.name}/{key1:STRING}"
      require_partition_filter = true
    }

    schema = <<EOH
[
  {
    "name": "name",
    "type": "STRING"
  },
  {
    "name": "last_modification",
    "type": "DATE"
  }
]
EOH
        }
  depends_on = ["google_storage_bucket_object.test"]
}
`, bucketName, datasetID, tableID)
}

func testAccBigQueryTableAvroPartitioning(bucketName, avroFilePath, datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "key1=20200330/init.avro"
  source = "%s"
  bucket  = google_storage_bucket.test.name
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  external_data_configuration {
    source_format = "AVRO"
    autodetect = true
    source_uris= ["gs://${google_storage_bucket.test.name}/*"]

    avro_options {
      use_avro_logical_types = true
    }

  }
  depends_on = ["google_storage_bucket_object.test"]
}
`, bucketName, avroFilePath, datasetID, tableID)
}

func testAccBigQueryTableRangePartitioning(datasetID, tableID string) string {
	return fmt.Sprintf(`
  resource "google_bigquery_dataset" "test" {
    dataset_id = "%s"
  }

  resource "google_bigquery_table" "test" {
	  deletion_protection = false
    table_id   = "%s"
    dataset_id = google_bigquery_dataset.test.dataset_id

    range_partitioning {
      field = "id"
      range {
        start    = 0
        end      = 10000
        interval = 100
      }
    }

    schema = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "id",
    "type": "INTEGER"
  }
]
EOH
}
  `, datasetID, tableID)
}

func testAccBigQueryTableWithView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type = "DAY"
  }

  view {
    query          = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
    use_legacy_sql = true
  }
}
`, datasetID, tableID)
}

func testAccBigQueryTableWithViewAndSchema(datasetID, tableID, desc string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  description = "%s"

  time_partitioning {
    type = "DAY"
  }

  schema = jsonencode(
  [

	{
	"description":"desc1",
	"mode":"NULLABLE",
	"name":"col1",
	"type":"STRING"
	},
	{
	"description":"desc2",
	"mode":"NULLABLE",
	"name":"col2",
	"type":"STRING"
	}
  ]
  )

  view {
    query = <<SQL
select "val1" as col1, "val2" as col2
SQL
    use_legacy_sql = false
  }
}
`, datasetID, tableID, desc)
}

func testAccBigQueryTableWithNewSqlView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type = "DAY"
  }

  view {
    query          = "%s"
    use_legacy_sql = false
  }
}
`, datasetID, tableID, "SELECT state FROM `lookerdata.cdc.project_tycho_reports`")
}

func testAccBigQueryTableWithMatViewDailyTimePartitioning_basic(datasetID, tableID, mViewID, query string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type                     = "DAY"
    field                    = "ts"
  }
  require_partition_filter = true
  clustering = ["some_int", "some_string"]
  schema     = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "some_string",
    "type": "STRING"
  },
  {
    "name": "some_int",
    "type": "INTEGER"
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "coord",
    "type": "RECORD",
    "fields": [
    {
    "name": "lon",
    "type": "FLOAT"
    }
    ]
  }
    ]
  }
]
EOH

}

resource "google_bigquery_table" "mv_test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type    = "DAY"
    field   = "ts"
  }

  materialized_view {
    query          = "%s"
  }

  depends_on = [
    google_bigquery_table.test,
  ]
}
`, datasetID, tableID, mViewID, query)
}

func testAccBigQueryTableWithMatViewDailyTimePartitioning(datasetID, tableID, mViewID, enable_refresh, refresh_interval, query string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type                     = "DAY"
    field                    = "ts"
  }
  require_partition_filter = true
  clustering = ["some_int", "some_string"]
  schema     = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "some_string",
    "type": "STRING"
  },
  {
    "name": "some_int",
    "type": "INTEGER"
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "coord",
    "type": "RECORD",
    "fields": [
    {
    "name": "lon",
    "type": "FLOAT"
    }
    ]
  }
    ]
  }
]
EOH

}

resource "google_bigquery_table" "mv_test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type    = "DAY"
    field   = "ts"
  }

  materialized_view {
    enable_refresh = "%s"
    refresh_interval_ms = "%s"
    query          = "%s"
  }

  depends_on = [
    google_bigquery_table.test,
  ]
}
`, datasetID, tableID, mViewID, enable_refresh, refresh_interval, query)
}

func testAccBigQueryTableWithMatViewNonIncremental_basic(datasetID, tableID, mViewID, query, maxStaleness string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}
resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  time_partitioning {
    type                     = "DAY"
    field                    = "ts"
    require_partition_filter = true
  }
  clustering = ["some_int", "some_string"]
  schema     = <<EOH
[
  {
    "name": "ts",
    "type": "TIMESTAMP"
  },
  {
    "name": "some_string",
    "type": "STRING"
  },
  {
    "name": "some_int",
    "type": "INTEGER"
  },
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "coord",
    "type": "RECORD",
    "fields": [
    {
    "name": "lon",
    "type": "FLOAT"
    }
    ]
  }
    ]
  }
]
EOH
}
resource "google_bigquery_table" "mv_test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  time_partitioning {
    type    = "DAY"
    field   = "ts"
  }
  materialized_view {
    query          = "%s"
    allow_non_incremental_definition = true
  }
  depends_on = [
    google_bigquery_table.test,
  ]
  max_staleness = "%s"
}
`, datasetID, tableID, mViewID, query, maxStaleness)
}

func testAccBigQueryTableUpdated(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  time_partitioning {
    type = "DAY"
  }

  schema = <<EOH
[
  {
    "name": "city",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "coord",
    "type": "RECORD",
    "fields": [
    {
      "name": "lon",
      "type": "FLOAT"
    },
    {
      "name": "lat",
      "type": "FLOAT"
    }
    ]
  }
    ]
  },
  {
    "name": "country",
    "type": "RECORD",
    "fields": [
  {
    "name": "id",
    "type": "INTEGER"
  },
  {
    "name": "name",
    "type": "STRING"
  }
    ]
  }
]
EOH

}
`, datasetID, tableID)
}

func testAccBigQueryTableFromGCS(datasetID, tableID, bucketName, objectName, content, format, quoteChar string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  content = <<EOF
%s
EOF

  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = true
    source_format = "%s"
    csv_options {
      encoding = "UTF-8"
      quote    = "%s"
    }

    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
    ]
  }
}
`, datasetID, bucketName, objectName, content, tableID, format, quoteChar)
}

func testAccBigQueryTableFromGCSParquetWithQueryAcceleration(connectionID, datasetID, tableID, bucketName, objectName, metadataCacheMode, maxStaleness string) string {
	return fmt.Sprintf(`
resource "google_bigquery_connection" "test" {
	connection_id = "%s"
	location = "US"
	cloud_resource {}
}

locals {
	connection_id_split = split("/", google_bigquery_connection.test.name)
	connection_id_reformatted = "${local.connection_id_split[1]}.${local.connection_id_split[3]}.${local.connection_id_split[5]}"
 }

 data "google_project" "project" {}

 resource "google_project_iam_member" "test" {
	role = "roles/storage.objectViewer"
	project = data.google_project.project.id
	member = "serviceAccount:${google_bigquery_connection.test.cloud_resource[0].service_account_id}"
 }

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
	connection_id   = local.connection_id_reformatted
    autodetect    = false
    source_format = "PARQUET"

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
	metadata_cache_mode = "%s"
	hive_partitioning_options {
		source_uri_prefix = "gs://${google_storage_bucket.test.name}/"
	}
  }

  max_staleness = "%s"

  depends_on = [
	google_project_iam_member.test
  ]
}
`, connectionID, datasetID, bucketName, objectName, tableID, metadataCacheMode, maxStaleness)
}

func testAccBigQueryTableFromGCSDeltaLake(datasetID, tableID, bucketName string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

# Setup Empty Delta Lake table in Bucket.

// Upload Metadata File.
resource "google_storage_bucket_object" "metadata" {
	name    = "_delta_log/00000000000000000000.json"
	source = "./test-fixtures/simple/metadata/00000000000000000000.json"
	bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = true
    source_format = "DELTA_LAKE"
	reference_file_schema_uri = "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.metadata.name}"

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
  }
}
`, datasetID, bucketName, tableID)
}

func testAccBigQueryTableFromGCSParquet(datasetID, tableID, bucketName, objectName string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = false
    source_format = "PARQUET"
	reference_file_schema_uri = "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}"

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
  }
}
`, datasetID, bucketName, objectName, tableID)
}

func testAccBigQueryTableFromGCSParquetOptions(datasetID, tableID, bucketName, objectName string, enum, list bool) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = false
    source_format = "PARQUET"
    reference_file_schema_uri = "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}"

    parquet_options {
      enum_as_string        = "%t"
      enable_list_inference = "%t"
    }

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
  }
}
`, datasetID, bucketName, objectName, tableID, enum, list)
}

func testAccBigQueryTableFromGCSIceberg(datasetID, tableID, bucketName string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

# Setup Empty Iceberg table in Bucket.
// .
// ├── data
// └── metadata
//     ├── 00000-1114da6b-bb88-4b5a-94bd-370f286c858a.metadata.json
// Upload Data Files
resource "google_storage_bucket_object" "empty_data_folder" {
	name   = "data/"
	content = " "
	bucket = google_storage_bucket.test.name
}
// Upload Metadata File.
resource "google_storage_bucket_object" "metadata" {
	name    = "simple/metadata/00000-1114da6b-bb88-4b5a-94bd-370f286c858a.metadata.json"
	source = "./test-fixtures/simple/metadata/00000-1114da6b-bb88-4b5a-94bd-370f286c858a.metadata.json"
	bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = false
    source_format = "ICEBERG"
	# Point to metadata.json.
    source_uris = [
      "gs://${google_storage_bucket.test.name}/simple/metadata/00000-1114da6b-bb88-4b5a-94bd-370f286c858a.metadata.json",
    ]
  }
  # Depends on Iceberg Table Files
  depends_on = [
	google_storage_bucket_object.empty_data_folder,
	google_storage_bucket_object.metadata, 
  ]
}
`, datasetID, bucketName, tableID)
}

func testAccBigQueryTableFromGCSParquetManifest(datasetID, tableID, bucketName, manifestName, parquetFileName string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

# Upload Data File.
resource "google_storage_bucket_object" "datafile" {
	name = "%s"
	source = "./test-fixtures/simple/data/00000-0-4e4a11ad-368c-496b-97ae-e3ac28051a4d-00001.parquet"
	bucket = google_storage_bucket.test.name
}

# Upload Metadata file
resource "google_storage_bucket_object" "manifest" {
	name = "%s" 
	content = "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.datafile.name}"
	bucket = google_storage_bucket.test.name
}


resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = false
    source_format = "PARQUET"
	# Specify URI is a manifest.
	file_set_spec_type = "FILE_SET_SPEC_TYPE_NEW_LINE_DELIMITED_MANIFEST"
	# Point to metadata.json.
    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.manifest.name}",
    ]
  }
}
`, datasetID, bucketName, manifestName, parquetFileName, tableID)
}

func testAccBigQueryTableExternalDataConfigurationConnectionID(location, connectionID, datasetID, tableID, bucketName, objectName, connectionIdReference string) string {
	return fmt.Sprintf(`
resource "google_bigquery_connection" "test" {
   connection_id = "%s"
   location = "%s"
   cloud_resource {}
}

data "google_project" "project" {}

resource "google_project_iam_member" "test" {
   role = "roles/storage.objectViewer"
   project = data.google_project.project.id
   member = "serviceAccount:${google_bigquery_connection.test.cloud_resource[0].service_account_id}"
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
  location = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "%s"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {

	# Feature Under Test
	connection_id   = %s

    autodetect      = false
	object_metadata = "SIMPLE"
	metadata_cache_mode = "MANUAL"

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
  }
}
`, connectionID, location, datasetID, location, bucketName, location, objectName, tableID, connectionIdReference)
}

func testAccBigQueryTableFromGCSObjectTable(connectionID, datasetID, tableID, bucketName, objectName, maxStaleness string) string {
	return fmt.Sprintf(`
resource "google_bigquery_connection" "test" {
   connection_id = "%s"
   location = "US"
   cloud_resource {}
}

locals {
   connection_id_split = split("/", google_bigquery_connection.test.name)
   connection_id_reformatted = "${local.connection_id_split[1]}.${local.connection_id_split[3]}.${local.connection_id_split[5]}"
}

data "google_project" "project" {}

resource "google_project_iam_member" "test" {
   role = "roles/storage.objectViewer"
   project = data.google_project.project.id
   member = "serviceAccount:${google_bigquery_connection.test.cloud_resource[0].service_account_id}"
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
	connection_id   = local.connection_id_reformatted
    autodetect      = false
	object_metadata = "SIMPLE"
	metadata_cache_mode = "MANUAL"

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
  }
  max_staleness = "%s"
}
`, connectionID, datasetID, bucketName, objectName, tableID, maxStaleness)
}

func testAccBigQueryTableFromGCSObjectTableMetadata(connectionID, datasetID, tableID, bucketName, objectName, maxStaleness string) string {
	return fmt.Sprintf(`
resource "google_bigquery_connection" "test" {
   connection_id = "%s"
   location = "US"
   cloud_resource {}
}

locals {
   connection_id_split = split("/", google_bigquery_connection.test.name)
   connection_id_reformatted = "${local.connection_id_split[1]}.${local.connection_id_split[3]}.${local.connection_id_split[5]}"
}

data "google_project" "project" {}

resource "google_project_iam_member" "test" {
   role = "roles/storage.objectViewer"
   project = data.google_project.project.id
   member = "serviceAccount:${google_bigquery_connection.test.cloud_resource[0].service_account_id}"
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "%s"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
	connection_id       = local.connection_id_reformatted
    autodetect          = false
	object_metadata     = "SIMPLE"
	metadata_cache_mode = "MANUAL"

    source_uris = [
      "gs://${google_storage_bucket.test.name}/*",
    ]
  }
  max_staleness = "%s"
  depends_on = [google_project_iam_member.test]
}
`, connectionID, datasetID, bucketName, objectName, tableID, maxStaleness)
}

func testAccBigQueryTableFromGCSWithSchemaWithConnectionId(datasetID, tableID, connectionID, projectID, bucketName, objectName, content, schema string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
resource "google_storage_bucket_object" "test" {
  name    = "%s"
  content = <<EOF
%s
EOF
  bucket = google_storage_bucket.test.name
}
resource "google_bigquery_connection" "test" {
   connection_id = "%s"
   location = "US"
   cloud_resource {}
}
locals {
   connection_id_split = split("/", google_bigquery_connection.test.name)
   connection_id_reformatted = "${local.connection_id_split[1]}.${local.connection_id_split[3]}.${local.connection_id_split[5]}"
}
resource "google_project_iam_member" "test" {
   role = "roles/storage.objectViewer"
   project = "%s"
   member = "serviceAccount:${google_bigquery_connection.test.cloud_resource[0].service_account_id}"
}
resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  schema = <<EOF
  %s
  EOF
  external_data_configuration {
    autodetect    = false
    connection_id = local.connection_id_reformatted
    source_format = "CSV"
    csv_options {
      encoding = "UTF-8"
      quote = ""
    }
    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
    ]
  }
  depends_on = [google_project_iam_member.test]
}
`, datasetID, bucketName, objectName, content, connectionID, projectID, tableID, schema)
}

func testAccBigQueryTableFromGCSWithSchemaWithConnectionId2(datasetID, tableID, connectionID, projectID, bucketName, objectName, content, schema string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
resource "google_storage_bucket_object" "test" {
  name    = "%s"
  content = <<EOF
%s
EOF
  bucket = google_storage_bucket.test.name
}
resource "google_bigquery_connection" "test" {
   connection_id = "%s"
   location = "US"
   cloud_resource {}
}
locals {
   connection_id_reformatted = google_bigquery_connection.test.name
}
resource "google_project_iam_member" "test" {
   role = "roles/storage.objectViewer"
   project = "%s"
   member = "serviceAccount:${google_bigquery_connection.test.cloud_resource[0].service_account_id}"
}
resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  schema = <<EOF
  %s
  EOF
  external_data_configuration {
    autodetect    = false
    connection_id = local.connection_id_reformatted
    source_format = "CSV"
    csv_options {
      encoding = "UTF-8"
      quote = ""
    }
    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
    ]
  }
  depends_on = [google_project_iam_member.test]
}
`, datasetID, bucketName, objectName, content, connectionID, projectID, tableID, schema)
}

func testAccBigQueryTableJson(bucketName, datasetID, tableID, encoding string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "test" {
  name    = "key1=20200330/data.json"
  content = "{\"name\":\"test\", \"last_modification\":\"2020-04-01\"}"
  bucket  = google_storage_bucket.test.name
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  external_data_configuration {
    source_format = "NEWLINE_DELIMITED_JSON"
    autodetect = false
    source_uris= ["gs://${google_storage_bucket.test.name}/*"]

    json_options {
      encoding = "%s"
    }

	json_extension = "GEOJSON"

    hive_partitioning_options {
      mode = "CUSTOM"
      source_uri_prefix = "gs://${google_storage_bucket.test.name}/{key1:STRING}"
      require_partition_filter = true
    }

    schema = <<EOH
[
  {
    "name": "name",
    "type": "STRING"
  },
  {
    "name": "last_modification",
    "type": "DATE"
  }
]
EOH
  }
  depends_on = ["google_storage_bucket_object.test"]
}
`, datasetID, bucketName, tableID, encoding)
}

func testAccBigQueryTableFromGCSWithSchema(datasetID, tableID, bucketName, objectName, content, schema string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
resource "google_storage_bucket_object" "test" {
  name    = "%s"
  content = <<EOF
%s
EOF
  bucket = google_storage_bucket.test.name
}
resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  schema = <<EOF
  %s
  EOF
  external_data_configuration {
    autodetect    = false
    source_format = "CSV"
    csv_options {
      encoding = "UTF-8"
      quote = ""
    }
    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
    ]
  }
}
`, datasetID, bucketName, objectName, content, tableID, schema)
}

func testAccBigQueryTableFromGCSWithExternalDataConfigSchema(datasetID, tableID, bucketName, objectName, content, schema string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
resource "google_storage_bucket_object" "test" {
  name    = "%s"
  content = <<EOF
%s
EOF
  bucket = google_storage_bucket.test.name
}
resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  external_data_configuration {
    autodetect    = false
    source_format = "CSV"
    csv_options {
      encoding = "UTF-8"
      quote = ""
    }
    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
    ]
    schema = <<EOF
    %s
    EOF
  }
}
`, datasetID, bucketName, objectName, content, tableID, schema)
}

func testAccBigQueryTableFromGCSWithSchema_UpdatAllowQuotedNewlines(datasetID, tableID, bucketName, objectName, content, schema string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}
resource "google_storage_bucket" "test" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
resource "google_storage_bucket_object" "test" {
  name    = "%s"
  content = <<EOF
%s
EOF
  bucket = google_storage_bucket.test.name
}
resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  schema = <<EOF
  %s
  EOF
  external_data_configuration {
    autodetect    = false
    source_format = "CSV"
    csv_options {
      encoding = "UTF-8"
      quote = ""
	  allow_quoted_newlines = "false"
      allow_jagged_rows     = "false"
    }
    source_uris = [
      "gs://${google_storage_bucket.test.name}/${google_storage_bucket_object.test.name}",
    ]
  }
}
`, datasetID, bucketName, objectName, content, tableID, schema)
}

func testAccBigQueryTableFromBigtable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigtable_instance" "instance" {
  name = "tf-test-bigtable-inst-%{random_suffix}"
  cluster {
    cluster_id = "tf-test-bigtable-%{random_suffix}"
    zone       = "us-central1-b"
  }
  instance_type = "DEVELOPMENT"
  deletion_protection = false
}
resource "google_bigtable_table" "table" {
  name          = "%{random_suffix}"
  instance_name = google_bigtable_instance.instance.name
  column_family {
    family = "cf-%{random_suffix}-first"
  }
  column_family {
    family = "cf-%{random_suffix}-second"
  }
}
resource "google_bigquery_table" "table" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  table_id   = "tf_test_bigtable_%{random_suffix}"
  external_data_configuration {
    autodetect            = true
    source_format         = "BIGTABLE"
    ignore_unknown_values = true
    source_uris = [
    "https://googleapis.com/bigtable/${google_bigtable_table.table.id}",
    ]
  }
}
resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "tf_test_ds_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
  default_table_expiration_ms = 3600000
  labels = {
    env = "default"
  }
}
`, context)
}

func testAccBigQueryTableFromBigtableOptions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigtable_instance" "instance" {
  name = "tf-test-bigtable-inst-%{random_suffix}"
  cluster {
    cluster_id = "tf-test-bigtable-%{random_suffix}"
    zone       = "us-central1-b"
  }
  instance_type = "DEVELOPMENT"
  deletion_protection = false
}
resource "google_bigtable_table" "table" {
  name          = "%{random_suffix}"
  instance_name = google_bigtable_instance.instance.name
  column_family {
    family = "cf-%{random_suffix}-first"
  }
  column_family {
    family = "cf-%{random_suffix}-second"
  }
}
resource "google_bigquery_table" "table" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  table_id   = "tf_test_bigtable_%{random_suffix}"
  external_data_configuration {
    autodetect            = true
    source_format         = "BIGTABLE"
    ignore_unknown_values = true
    source_uris = [
    "https://googleapis.com/bigtable/${google_bigtable_table.table.id}",
    ]
	bigtable_options {
      column_family {
        family_id        = "cf-%{random_suffix}-first"
		column {
			field_name       = "cf-%{random_suffix}-first"
			type             = "STRING"
			encoding         = "TEXT"
			only_read_latest = true
		  }
		type             = "STRING"
		encoding         = "TEXT"
		only_read_latest = true
	  }
      column_family {
        family_id        = "cf-%{random_suffix}-second"
		type             = "STRING"
		encoding         = "TEXT"
		only_read_latest = false
	  }
      ignore_unspecified_column_families = true
      read_rowkey_as_string              = true
      output_column_families_as_json     = true
	}
  }
}
resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "tf_test_ds_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
  delete_contents_on_destroy  = true
  default_table_expiration_ms = 3600000
  labels = {
    env = "default"
  }
}
`, context)
}

func testAccBigQueryTableFromSheet(context map[string]interface{}) string {
	return acctest.Nprintf(`
  resource "google_bigquery_table" "table" {
	  deletion_protection = false
    dataset_id = google_bigquery_dataset.dataset.dataset_id
    table_id   = "tf_test_sheet_%{random_suffix}"

    external_data_configuration {
      autodetect            = true
      source_format         = "GOOGLE_SHEETS"
      ignore_unknown_values = true

      google_sheets_options {
      skip_leading_rows = 1
      }

      source_uris = [
      "https://drive.google.com/open?id=xxxx",
      ]
    }

    schema = <<EOF
    [
    {
      "name": "permalink",
      "type": "STRING",
      "mode": "NULLABLE",
      "description": "The Permalink"
    },
    {
      "name": "state",
      "type": "STRING",
      "mode": "NULLABLE",
      "description": "State where the head office is located"
    }
    ]
    EOF
    }

    resource "google_bigquery_dataset" "dataset" {
    dataset_id                  = "tf_test_ds_%{random_suffix}"
    friendly_name               = "test"
    description                 = "This is a test description"
    location                    = "EU"
    default_table_expiration_ms = 3600000

    labels = {
      env = "default"
    }
    }
`, context)
}

func testAccBigQueryTable_jsonEq(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "nullable"
        type        = "integer"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_jsonEqModeRemoved(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "false"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        type        = "INTEGER"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_jsonPreventDestroy(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id
	lifecycle {
		prevent_destroy = true
	}

  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "NULLABLE"
        type        = "INTEGER"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_jsonPreventDestroyOrderChanged(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id
	lifecycle {
		prevent_destroy = true
	}

  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
			},
			{
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "NULLABLE"
        type        = "INTEGER"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_noAllowDestroy(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id
  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "NULLABLE"
        type        = "INTEGER"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_noAllowDestroyUpdated(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
	dataset_id = google_bigquery_dataset.test.dataset_id
  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "NULLABLE"
        type        = "INTEGER"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_arrayInitial(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "NULLABLE"
        type        = "INTEGER"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_arrayExpanded(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id

  friendly_name = "bigquerytest"
  labels = {
    "terrafrom_managed" = "true"
  }

  schema = jsonencode(
    [
      {
        description = "Time snapshot was taken, in Epoch milliseconds. Same across all rows and all tables in the snapshot, and uniquely defines a particular snapshot."
        name        = "snapshot_timestamp"
        mode        = "NULLABLE"
        type        = "INTEGER"
      },
      {
        description = "Timestamp of dataset creation"
        name        = "creation_time"
        type        = "TIMESTAMP"
      },
			{
        description = "some new value"
        name        = "a_new_value"
        type        = "TIMESTAMP"
      },
    ])
}
`, datasetID, tableID)
}

func testAccBigQueryTable_mimicCreateFromConsole(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  schema = <<EOF
  [
  ]
  EOF
}
`, datasetID, tableID)
}

func testAccBigQueryTable_emptySchema(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
}
`, datasetID, tableID)
}

func testAccBigQueryTablePrimaryKey(datasetID, tableID string) string {
	return fmt.Sprintf(`
  resource "google_bigquery_dataset" "foo" {
    dataset_id = "%s"
  }

  resource "google_bigquery_table" "test" {
    deletion_protection = false
    table_id   = "%s"
    dataset_id = google_bigquery_dataset.foo.dataset_id

    table_constraints {
      primary_key {
        columns = ["id"]
      }
    }

    schema = <<EOH
    [
      {
        "name": "ts",
        "type": "TIMESTAMP"
      },
      {
        "name": "id",
        "type": "INTEGER"
      }
    ]
    EOH
  }
  `, datasetID, tableID)
}

func testAccBigQueryTableForeignKeys(projectID, datasetID, tableID_pk, tableID_fk string) string {
	return fmt.Sprintf(`
  resource "google_bigquery_dataset" "foo" {
    dataset_id = "%s"
  }

  resource "google_bigquery_table" "table_pk" {
    deletion_protection = false
    table_id   	= "%s"
    dataset_id 	= google_bigquery_dataset.foo.dataset_id

    table_constraints {
      primary_key {
        columns = ["id"]
      }
    }

    schema = <<EOH
    [
      {
        "name": "ts",
        "type": "TIMESTAMP"
      },
      {
        "name": "id",
        "type": "INTEGER"
      },
      {
        "name": "str",
        "type": "STRING"
      }
    ]
    EOH
  }

  resource "google_bigquery_table" "test" {
    deletion_protection = false
    table_id   	= "%s"
    dataset_id 	= google_bigquery_dataset.foo.dataset_id

    table_constraints {
      foreign_keys {
        name = "test_fk"
        referenced_table {
          project_id  	= "%s"
          dataset_id 	= google_bigquery_dataset.foo.dataset_id
          table_id   	= google_bigquery_table.table_pk.table_id
        }
        column_references {
          referencing_column 	= "id2"
          referenced_column 	= "id"
        }
      }
    }

    schema = <<EOH
    [
      {
        "name": "ts2",
        "type": "TIMESTAMP"
      },
      {
        "name": "id2",
        "type": "INTEGER"
      }
    ]
    EOH
  }
  `, datasetID, tableID_pk, tableID_fk, projectID)
}

func testAccBigQueryTableTableConstraintsUpdate(projectID, datasetID, tableID_pk, tableID_fk string) string {
	return fmt.Sprintf(`
  resource "google_bigquery_dataset" "foo" {
    dataset_id 	= "%s"
  }

  resource "google_bigquery_table" "table_pk" {
	deletion_protection = false
    table_id   	= "%s"
    dataset_id 	= google_bigquery_dataset.foo.dataset_id

    table_constraints {
      primary_key {
        columns = ["str"]
      }
      foreign_keys {
        name = "test_fk"
        referenced_table {
          project_id  	= "%s"
          dataset_id 	= google_bigquery_dataset.foo.dataset_id
          table_id   	= google_bigquery_table.test.table_id
        }
        column_references {
        referencing_column = "id"
        referenced_column = "id2"
        }
      }
      foreign_keys {
        name = "test_fk2"
        referenced_table {
          project_id  	= "%s"
          dataset_id 	= google_bigquery_dataset.foo.dataset_id
          table_id   	= google_bigquery_table.test.table_id
        }
        column_references {
          referencing_column 	= "ts"
          referenced_column 	= "ts2"
        }
      }
    }

    schema = <<EOH
    [
      {
        "name": "ts",
        "type": "TIMESTAMP"
      },
      {
        "name": "id",
        "type": "INTEGER"
      },
      {
        "name": "str",
        "type": "INTEGER"
      }
    ]
    EOH
  }

  resource "google_bigquery_table" "test" {
    deletion_protection = false
    table_id   	= "%s"
    dataset_id 	= google_bigquery_dataset.foo.dataset_id

    table_constraints {
      primary_key {
        columns = ["id2", "ts2"]
      }
    }

    schema = <<EOH
    [
      {
        "name": "ts2",
        "type": "TIMESTAMP"
      },
      {
        "name": "id2",
        "type": "INTEGER"
      }
    ]
    EOH
  }
  `, datasetID, tableID_pk, projectID, projectID, tableID_fk)
}

func testAccBigQueryTableWithSchema(datasetID, tableID, schema string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
	schema = <<EOF
  %s
  EOF
}
`, datasetID, tableID, schema)
}

func testAccBigQueryTableWithReplicationInfoAndView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  view {
    query          = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
    use_legacy_sql = true
  }
  table_replication_info {
    source_project_id = "source_project_id"
    source_dataset_id = "source_dataset_id"
    source_table_id = "source_table_id"
  }
}
`, datasetID, tableID)
}

func testAccBigQueryTableWithSchemaWithRequiredFieldAndView(datasetID, tableID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.test.dataset_id
  schema = <<EOF
  [
    {
      "name": "requiredField",
      "type": "STRING",
      "mode": "REQUIRED",
      "description": "requiredField"
    },
    {
      "name": "optionalField",
      "type": "STRING",
      "mode": "NULLABLE",
      "description": "optionalField"
    }
  ]
  EOF
  view {
    query = <<EOF
      SELECT 'a' AS requiredField, 'b' AS optionalField
    EOF
    use_legacy_sql = false
  }
}
`, datasetID, tableID)
}

func testAccBigQueryTableWithReplicationInfo(projectID, sourceDatasetID, sourceTableID, sourceMVID, replicaDatasetID, replicaMVID, sourceMVJobID, dropMVJobID, replicationIntervalExpr string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "source" {
  dataset_id = "%s"
  location = "aws-us-east-1"
}

resource "google_bigquery_table" "source_table" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.source.dataset_id
  external_data_configuration {
    connection_id   = "bigquerytestdefault.aws-us-east-1.e2e_ccmv_tf_test"
    autodetect      = true
		metadata_cache_mode = "AUTOMATIC"
    source_format = "PARQUET"
    source_uris = [
      "s3://bq-testing/parquet/all_types.parquet",
    ]
  }
  max_staleness = "0-0 0 10:0:0"
}

resource "google_bigquery_job" "source_mv_job" {
  job_id = "%s"

  location = "aws-us-east-1"
  query {
    query = "CREATE MATERIALIZED VIEW %s.%s OPTIONS(max_staleness=INTERVAL \"10:00:0\" HOUR TO SECOND) AS SELECT * FROM %s.%s"
    use_legacy_sql = false
    create_disposition = ""
    write_disposition = ""
  }

  depends_on = [google_bigquery_table.source_table]
}

resource "time_sleep" "wait_10_seconds" {
  depends_on = [google_bigquery_job.source_mv_job]
  create_duration = "10s"
}

resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.source.dataset_id
  view {
    project_id = "%s"
    dataset_id = google_bigquery_dataset.source.dataset_id
    table_id   = "%s"
  }

  depends_on = [time_sleep.wait_10_seconds]
}

resource "google_bigquery_dataset" "replica" {
  dataset_id = "%s"
  location = "us"
}

resource "google_bigquery_table" "replica_mv" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.replica.dataset_id
  table_id   = "%s"
  table_replication_info {
    source_project_id = "%s"
    source_dataset_id = google_bigquery_dataset.source.dataset_id
    source_table_id = "%s"
    %s
  }

  depends_on = [google_bigquery_dataset_access.access]
}

resource "google_bigquery_job" "drop_source_mv_job" {
  job_id = "%s"

  location = "aws-us-east-1"
  query {
    query = "DROP MATERIALIZED VIEW %s.%s"
    use_legacy_sql = false
    create_disposition = ""
    write_disposition = ""
  }

  depends_on = [google_bigquery_table.replica_mv]
}

resource "time_sleep" "wait_10_seconds_last" {
  depends_on = [google_bigquery_job.drop_source_mv_job]
  create_duration = "10s"
}
`, sourceDatasetID, sourceTableID, sourceMVJobID, sourceDatasetID, sourceMVID, sourceDatasetID, sourceTableID, projectID, sourceMVID, replicaDatasetID, replicaMVID, projectID, sourceMVID, replicationIntervalExpr, dropMVJobID, sourceDatasetID, sourceMVID)
}

func testAccBigQueryTableWithResourceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key1" {
  parent = "projects/%{project_id}"
  short_name = "%{tag_key_name1}"
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "%{tag_value_name1}"
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%{dataset_id}"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  allow_resource_tags_on_deletion = true
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"
  table_id   = "%{table_id}"
  resource_tags = {
    "%{project_id}/${google_tags_tag_key.key1.short_name}" = "${google_tags_tag_value.value1.short_name}"
  }
}
`, context)
}

func testAccBigQueryTableWithResourceTagsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key1" {
  parent = "projects/%{project_id}"
  short_name = "%{tag_key_name1}"
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "%{tag_value_name1}"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%{project_id}"
  short_name = "%{tag_key_name2}"
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "%{tag_value_name2}"
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%{dataset_id}"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  allow_resource_tags_on_deletion = true
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"
  table_id   = "%{table_id}"
  resource_tags = {
    "%{project_id}/${google_tags_tag_key.key1.short_name}" = "${google_tags_tag_value.value1.short_name}"
    "%{project_id}/${google_tags_tag_key.key2.short_name}" = "${google_tags_tag_value.value2.short_name}"
  }
}
`, context)
}

func testAccBigQueryTableWithResourceTagsDestroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_tags_tag_key" "key1" {
  parent = "projects/%{project_id}"
  short_name = "%{tag_key_name1}"
}

resource "google_tags_tag_value" "value1" {
  parent = "tagKeys/${google_tags_tag_key.key1.name}"
  short_name = "%{tag_value_name1}"
}

resource "google_tags_tag_key" "key2" {
  parent = "projects/%{project_id}"
  short_name = "%{tag_key_name2}"
}

resource "google_tags_tag_value" "value2" {
  parent = "tagKeys/${google_tags_tag_key.key2.name}"
  short_name = "%{tag_value_name2}"
}

resource "google_bigquery_dataset" "test" {
  dataset_id = "%{dataset_id}"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  allow_resource_tags_on_deletion = true
  dataset_id = "${google_bigquery_dataset.test.dataset_id}"
  table_id   = "%{table_id}"
  resource_tags = {}
}
`, context)
}

var TEST_CSV = `lifelock,LifeLock,,web,Tempe,AZ,1-May-07,6850000,USD,b
lifelock,LifeLock,,web,Tempe,AZ,1-Oct-06,6000000,USD,a
lifelock,LifeLock,,web,Tempe,AZ,1-Jan-08,25000000,USD,c
mycityfaces,MyCityFaces,7,web,Scottsdale,AZ,1-Jan-08,50000,USD,seed
flypaper,Flypaper,,web,Phoenix,AZ,1-Feb-08,3000000,USD,a
infusionsoft,Infusionsoft,105,software,Gilbert,AZ,1-Oct-07,9000000,USD,a
gauto,gAuto,4,web,Scottsdale,AZ,1-Jan-08,250000,USD,seed
chosenlist-com,ChosenList.com,5,web,Scottsdale,AZ,1-Oct-06,140000,USD,seed
chosenlist-com,ChosenList.com,5,web,Scottsdale,AZ,25-Jan-08,233750,USD,angel
`
var TEST_SIMPLE_CSV = `US,phone,100
JP,tablet,300
UK,laptop,200
`
var TEST_SIMPLE_CSV_SCHEMA = `[
    {
      "name": "country",
      "type": "STRING"
    },
    {
      "name": "product",
      "type": "STRING"
    },
    {
      "name": "price",
      "type": "INT64"
    }
  ]`
var TEST_INVALID_SCHEMA_NOT_JSON = `
	not a valid table schema
	`
var TEST_INVALID_SCHEMA_NOT_JSON_LIST = `
    {
      "name": "country",
      "type": "STRING"
    }`
var TEST_INVALID_SCHEMA_JSON_LIST_WITH_NULL_ELEMENT = `[
    {
      "name": "country",
      "type": "STRING"
    },
    null,
    {
      "name": "price",
      "type": "INT64"
    }
  ]`
