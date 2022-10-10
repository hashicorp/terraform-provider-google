package google

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestBigQueryTableSchemaDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"empty schema": {
			Old:                "null",
			New:                "[]",
			ExpectDiffSuppress: true,
		},
		"empty schema -> non-empty": {
			Old: "null",
			New: `[
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: false,
		},
		"no change": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"finalKey\" : {} }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"finalKey\" : {} }]",
			ExpectDiffSuppress: true,
		},
		"remove key": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"finalKey\" : {} }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"finalKey\" : {} }]",
			ExpectDiffSuppress: false,
		},
		"empty description -> default description (empty)": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"description\": \"\"  }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\" }]",
			ExpectDiffSuppress: true,
		},
		"empty description -> other description": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"description\": \"\"  }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"description\": \"somethingRandom\"  }]",
			ExpectDiffSuppress: false,
		},
		"mode NULLABLE -> other mode": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"mode\": \"NULLABLE\"  }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"mode\": \"somethingRandom\"  }]",
			ExpectDiffSuppress: false,
		},
		"mode NULLABLE -> default mode (also NULLABLE)": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			New: `[
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"mode & type uppercase -> lowercase": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			New: `[
				{
					"mode": "nullable",
					"name": "PageNo",
					"type": "integer"
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"type INTEGER -> INT64": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"INTEGER\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"INT64\"  }]",
			ExpectDiffSuppress: true,
		},
		"type INTEGER -> other": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"INTEGER\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"somethingRandom\"  }]",
			ExpectDiffSuppress: false,
		},
		"type FLOAT -> FLOAT64": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"FLOAT\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"FLOAT64\"  }]",
			ExpectDiffSuppress: true,
		},
		"type FLOAT -> other": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"FLOAT\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"somethingRandom\" }]",
			ExpectDiffSuppress: false,
		},
		"type BOOLEAN -> BOOL": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOL\"  }]",
			ExpectDiffSuppress: true,
		},
		"type BOOLEAN -> other": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"somethingRandom\" }]",
			ExpectDiffSuppress: false,
		},
		// this is invalid but we need to make sure we don't cause a panic
		// if users provide an invalid schema
		"invalid - missing type for old": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
			ExpectDiffSuppress: false,
		},
		// this is invalid but we need to make sure we don't cause a panic
		// if users provide an invalid schema
		"invalid - missing type for new": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
			ExpectDiffSuppress: false,
		},
		"reordering fields": {
			Old: `[
				{
					"name": "PageNo",
					"type": "INTEGER"
				},
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				}
			]`,
			New: `[
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				},
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"reordering fields with value change": {
			Old: `[
				{
					"name": "PageNo",
					"type": "INTEGER",
					"description": "someVal"
				},
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				}
			]`,
			New: `[
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				},
				{
					"name": "PageNo",
					"type": "INTEGER",
					"description": "otherVal"
				}
			]`,
			ExpectDiffSuppress: false,
		},
		"nested field ordering changes": {
			Old: `[
				{
					"name": "someValue",
					"type": "INTEGER",
					"fields": [
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal"
						},
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				}
			]`,
			New: `[
				{
					"name": "someValue",
					"type": "INTEGER",
					"fields": [
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						},
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"policyTags": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "providerphone",
					"policyTags": {
						"names": [
							"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
						]
					},
					"type":"STRING"
				}
			]`,
			New: `[
			  {
			    "name": "providerphone",
			    "type": "STRING",
			    "policyTags": {
			          "names": ["projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"]
			        }
			  }
			]`,
			ExpectDiffSuppress: true,
		},
		"multiple levels of reordering with policyTags set": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "providerphone",
					"type":"STRING",
					"policyTags": {
						"names": [
							"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
						]
					},
					"fields": [
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						},
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				},
				{
					"name": "PageNo",
					"type": "INTEGER"
				},
				{
					"name": "IngestTime",
					"type": "TIMESTAMP",
					"fields": [
						{
							"name": "value3",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						},
						{
							"name": "value4",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				}
			]`,
			New: `[
				{
					"name": "IngestTime",
					"type": "TIMESTAMP",
					"fields": [
						{
							"name": "value4",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						},
						{
							"name": "value3",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						}
					]
				},
				{
					"mode": "NULLABLE",
					"name": "providerphone",
					"type":"STRING",
					"policyTags": {
						"names": [
							"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
						]
					},
					"fields": [
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						},
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				},
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		tn := tn
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var a, b interface{}
			if err := json.Unmarshal([]byte(tc.Old), &a); err != nil {
				t.Fatalf(fmt.Sprintf("unable to unmarshal old json - %v", err))
			}
			if err := json.Unmarshal([]byte(tc.New), &b); err != nil {
				t.Fatalf(fmt.Sprintf("unable to unmarshal new json - %v", err))
			}
			if bigQueryTableSchemaDiffSuppress("schema", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
				t.Fatalf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
			}
		})
	}
}

func TestAccBigQueryTable_Basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryTable_Kms(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	kms := BootstrapKMSKey(t)
	cryptoKeyName := kms.CryptoKey.Name

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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
	bucketName := testBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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
	bucketName := testBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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
	bucketName := testBucketName(t)
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	avroFilePath := "./test-fixtures/bigquerytable/avro-generated.avro"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryTable_RangePartitioning(t *testing.T) {
	t.Parallel()
	resourceName := "google_bigquery_table.test"
	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryTable_View(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	materialized_viewID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	query := fmt.Sprintf("SELECT count(some_string) as count, some_int, ts FROM `%s.%s` WHERE DATE(ts) = '2019-01-01' GROUP BY some_int, ts", datasetID, tableID)
	queryNew := strings.ReplaceAll(query, "2019", "2020")

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	materialized_viewID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	query := fmt.Sprintf("SELECT count(some_string) as count, some_int, ts FROM `%s.%s` WHERE DATE(ts) = '2019-01-01' GROUP BY some_int, ts", datasetID, tableID)

	enable_refresh := "false"
	refresh_interval_ms := "3600000"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryExternalDataTable_CSV(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", randString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryExternalDataTable_CSV_WithSchemaAndConnectionID_UpdateNoConnectionID(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", randString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	projectID := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", randString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	connectionID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	projectID := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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
		},
	})
}

func TestAccBigQueryExternalDataTable_CSV_WithSchema_UpdateAllowQuotedNewlines(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	objectName := fmt.Sprintf("tf_test_%s.csv", randString(t, 10))

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 8),
		"project":       getTestProjectFromEnv(),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

func TestAccBigQueryDataTable_sheet(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_jsonEq(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTable_jsonEqModeRemoved(datasetID, tableID),
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

func TestAccBigQueryDataTable_canReorderParameters(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_arrayInitial(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "last_modified_time", "deletion_protection"},
			},
			{
				Config: testAccBigQueryTable_arrayExpanded(datasetID, tableID),
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

func TestAccBigQueryTable_allowDestroy(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryTable_noAllowDestroy(datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_table.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection"},
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

	datasetID := fmt.Sprintf("tf_test_%s", randString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigQueryTableDestroyProducer(t),
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

type testUnitBigQueryDataTableJSONChangeableTestCase struct {
	name       string
	jsonOld    string
	jsonNew    string
	changeable bool
}

func (testcase *testUnitBigQueryDataTableJSONChangeableTestCase) check(t *testing.T) {
	var old, new interface{}
	if err := json.Unmarshal([]byte(testcase.jsonOld), &old); err != nil {
		t.Fatalf("unable to unmarshal json - %v", err)
	}
	if err := json.Unmarshal([]byte(testcase.jsonNew), &new); err != nil {
		t.Fatalf("unable to unmarshal json - %v", err)
	}
	changeable, err := resourceBigQueryTableSchemaIsChangeable(old, new)
	if err != nil {
		t.Errorf("%s failed unexpectedly: %s", testcase.name, err)
	}
	if changeable != testcase.changeable {
		t.Errorf("expected changeable result of %v but got %v for testcase %s", testcase.changeable, changeable, testcase.name)
	}

	d := &ResourceDiffMock{
		Before: map[string]interface{}{},
		After:  map[string]interface{}{},
	}

	d.Before["schema"] = testcase.jsonOld
	d.After["schema"] = testcase.jsonNew

	err = resourceBigQueryTableSchemaCustomizeDiffFunc(d)
	if err != nil {
		t.Errorf("error on testcase %s - %v", testcase.name, err)
	}
	if !testcase.changeable != d.IsForceNew {
		t.Errorf("%s: expected d.IsForceNew to be %v, but was %v", testcase.name, !testcase.changeable, d.IsForceNew)
	}
}

var testUnitBigQueryDataTableIsChangableTestCases = []testUnitBigQueryDataTableJSONChangeableTestCase{
	{
		name:       "defaultEquality",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: true,
	},
	{
		name:       "arraySizeIncreases",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"asomeValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: true,
	},
	{
		name:       "arraySizeDecreases",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"asomeValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: false,
	},
	{
		name:       "descriptionChanges",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeInteger",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INT64\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeFloat",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"FLOAT\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"FLOAT64\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeBool",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOL\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeChangeIncompatible",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"DATETIME\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: false,
	},
	// this is invalid but we need to make sure we don't cause a panic
	// if users provide an invalid schema
	{
		name:       "typeChangeIgnoreNewMissingType",
		jsonOld:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
		changeable: true,
	},
	// this is invalid but we need to make sure we don't cause a panic
	// if users provide an invalid schema
	{
		name:       "typeChangeIgnoreOldMissingType",
		jsonOld:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
		changeable: true,
	},
	{
		name:       "typeModeReqToNull",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REQUIRED\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeModeIncompatible",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REQUIRED\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REPEATED\", \"description\" : \"some new value\" }]",
		changeable: false,
	},
	{
		name:       "modeToDefaultNullable",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REQUIRED\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "orderOfArrayChangesAndDescriptionChanges",
		jsonOld:    "[{\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"value2\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"value2\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"newVal\" },  {\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: true,
	},
	{
		name:       "orderOfArrayChangesAndNameChanges",
		jsonOld:    "[{\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"value2\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"value3\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"newVal\" },  {\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: false,
	},
	{
		name: "policyTags",
		jsonOld: `[
			{
				"mode": "NULLABLE",
				"name": "providerphone",
				"policyTags": {
					"names": ["projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"]
				},
				"type":"STRING"
			}
		]`,
		jsonNew: `[
			{
				"name": "providerphone",
				"type": "STRING",
				"policyTags": {
					"names": ["projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"]
				}
			}
		]`,
		changeable: true,
	},
}

func TestUnitBigQueryDataTable_schemaIsChangable(t *testing.T) {
	t.Parallel()
	for _, testcase := range testUnitBigQueryDataTableIsChangableTestCases {
		testcase.check(t)
		testcaseNested := &testUnitBigQueryDataTableJSONChangeableTestCase{
			testcase.name + "Nested",
			fmt.Sprintf("[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"fields\" : %s }]", testcase.jsonOld),
			fmt.Sprintf("[{\"name\": \"someValue\", \"type\" : \"INT64\", \"fields\" : %s }]", testcase.jsonNew),
			testcase.changeable,
		}
		testcaseNested.check(t)
	}
}

func testAccCheckBigQueryExtData(t *testing.T, expectedQuoteChar string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigquery_table" {
				continue
			}

			config := googleProviderConfig(t)
			dataset := rs.Primary.Attributes["dataset_id"]
			table := rs.Primary.Attributes["table_id"]
			res, err := config.NewBigQueryClient(config.userAgent).Tables.Get(config.Project, dataset, table).Do()
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

			config := googleProviderConfig(t)
			_, err := config.NewBigQueryClient(config.userAgent).Tables.Get(config.Project, rs.Primary.Attributes["dataset_id"], rs.Primary.Attributes["table_id"]).Do()
			if err == nil {
				return fmt.Errorf("Table still present")
			}
		}

		return nil
	}
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
`, datasetID, tableID, partitioningType)
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
}
`, datasetID, bucketName, objectName, content, connectionID, projectID, tableID, schema)
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
	return Nprintf(`
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

func testAccBigQueryTableFromSheet(context map[string]interface{}) string {
	return Nprintf(`
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
    "terrafrom_managed" = "true"
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
