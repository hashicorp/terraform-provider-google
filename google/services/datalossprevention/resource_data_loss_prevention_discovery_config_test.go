// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datalossprevention_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataLossPreventionDiscoveryConfig_Update(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":      testAccDataLossPreventionDiscoveryConfig_BasicUpdate,
		"org":        testAccDataLossPreventionDiscoveryConfig_OrgUpdate,
		"actions":    testAccDataLossPreventionDiscoveryConfig_ActionsUpdate,
		"conditions": testAccDataLossPreventionDiscoveryConfig_ConditionsCadenceUpdate,
		"filter":     testAccDataLossPreventionDiscoveryConfig_FilterUpdate,
		"cloud_sql":  testAccDataLossPreventionDiscoveryConfig_CloudSqlUpdate,
		"bq_single":  testAccDataLossPreventionDiscoveryConfig_BqSingleTable,
		"sql_single": testAccDataLossPreventionDiscoveryConfig_SqlSingleTable,
		"secrets":    testAccDataLossPreventionDiscoveryConfig_SecretsUpdate,
		"gcs":        testAccDataLossPreventionDiscoveryConfig_CloudStorageUpdate,
		"gcs_single": testAccDataLossPreventionDiscoveryConfig_CloudStorageSingleBucket,
	}
	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccDataLossPreventionDiscoveryConfig_BasicUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_OrgUpdate(t *testing.T) {

	context := map[string]interface{}{
		"organization":  envvar.GetTestOrgFromEnv(t),
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgRunning(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgFolderPaused(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_ActionsUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigActions(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigActionsSensitivity(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_ConditionsCadenceUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigConditionsCadence(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_FilterUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigFilterRegexesAndConditions(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_CloudSqlUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStartCloudSql(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateCloudSql(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_BqSingleTable(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigBqSingleUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_SqlSingleTable(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStartCloudSql(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigCloudSqlSingleUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_CloudStorageUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStartCloudStorage(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateCloudStorage(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_CloudStorageSingleBucket(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStartCloudStorage(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigCloudStorageSingleUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_SecretsUpdate(t *testing.T) {

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"location":      envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataLossPreventionDiscoveryConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigSecretsStart(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
			{
				Config: testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigSecretsUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_discovery_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "parent", "last_run_time", "update_time", "errors"},
			},
		},
	})
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}/locations/%{location}"
	location = "%{location}"
	display_name = "display name"
	status = "RUNNING"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "custom_type" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		custom_info_types {
			info_type {
				name = "MY_CUSTOM_TYPE"
			}

			likelihood = "UNLIKELY"

			regex {
				pattern = "test*"
			}
		}
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}/locations/%{location}"
	location = "%{location}"
	status = "RUNNING"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
			conditions {
				or_conditions {
					min_row_count = 10
					min_age = "10800s"
				}
			}
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.custom_type.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigActions(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
	project_id = "%{project}"
}

resource "google_tags_tag_key" "tag_key" {
	parent = "projects/${data.google_project.project.number}"
	short_name = "environment"
}

resource "google_tags_tag_value" "tag_value" {
	parent = google_tags_tag_key.tag_key.id
	short_name = "prod"
}

resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_pubsub_topic" "basic" {
	name = "test-topic"
}

resource "google_project_iam_member" "tag_role" {
    project = "%{project}"
    role    = "roles/resourcemanager.tagUser"
    member = "serviceAccount:service-${data.google_project.project.number}@dlp-api.iam.gserviceaccount.com"
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}/locations/%{location}"
	location = "%{location}"
	status = "RUNNING"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	actions {
        export_data {
            profile_table {
                project_id = "%{project}"
                dataset_id = "dataset"
                table_id = "table"
            }
        }
    }
    actions { 
        pub_sub_notification {
			topic = "projects/%{project}/topics/${google_pubsub_topic.basic.name}"
			event = "NEW_PROFILE"
			pubsub_condition {
				expressions {
					logical_operator = "OR"
					conditions { 
						minimum_risk_score = "HIGH" 
					}
				}
			}
			detail_of_message = "TABLE_PROFILE"
		}
    }
	actions {
        tag_resources {
            tag_conditions {
                tag {
                    namespaced_value = "%{project}/environment/prod"
                }
                sensitivity_score {
                    score = "SENSITIVITY_HIGH"
                }
            }
            profile_generations_to_tag = ["PROFILE_GENERATION_NEW", "PROFILE_GENERATION_UPDATE"]
            lower_data_risk_to_low = true
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
	depends_on = [
		google_project_iam_member.tag_role,
		google_tags_tag_key.tag_key,
		google_tags_tag_value.tag_value,
	]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigActionsSensitivity(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_pubsub_topic" "basic" {
	name = "test-topic"
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}/locations/%{location}"
	location = "%{location}"
	status = "RUNNING"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	actions {
        export_data {
            profile_table {
                project_id = "project"
                dataset_id = "dataset"
                table_id = "table"
            }
        }
    }
    actions { 
        pub_sub_notification {
			topic = "projects/%{project}/topics/${google_pubsub_topic.basic.name}"
			event = "NEW_PROFILE"
			pubsub_condition {
				expressions {
					logical_operator = "OR"
					conditions { 
						minimum_sensitivity_score = "HIGH" 
					}
				}
			}
			detail_of_message = "TABLE_PROFILE"
		}
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgRunning(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "organizations/%{organization}/locations/%{location}"
	location = "%{location}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	org_config {
		project_id = "%{project}"
		location {
			organization_id = "%{organization}"
		}
	}
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
	status = "RUNNING"
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigOrgFolderPaused(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "organizations/%{organization}/locations/%{location}"
	location = "%{location}"

    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
	org_config {
		project_id = "%{project}"
		location {
			folder_id = 123
		}
	}
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
	status = "PAUSED"
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigConditionsCadence(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}/locations/%{location}"
	location = "%{location}"
	status = "RUNNING"

	targets {
		big_query_target {
			filter {
				other_tables {}
			}
			conditions {
				type_collection = "BIG_QUERY_COLLECTION_ALL_TYPES"
			}
			cadence {
				schema_modified_cadence {
					types = ["SCHEMA_NEW_COLUMNS"]
					frequency = "UPDATE_FREQUENCY_DAILY"
				}
				table_modified_cadence {
					types = ["TABLE_MODIFIED_TIMESTAMP"]
					frequency = "UPDATE_FREQUENCY_DAILY"
				}
				inspect_template_modified_cadence {
					frequency = "UPDATE_FREQUENCY_DAILY"
				}
			}
		}
	}
	inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigFilterRegexesAndConditions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}

resource "google_data_loss_prevention_discovery_config" "basic" {
	parent = "projects/%{project}/locations/%{location}"
	location = "%{location}"
	status = "RUNNING"

	targets {
        big_query_target {
            filter {
                tables {
                    include_regexes {
                        patterns {
                            project_id_regex = ".*"
                            dataset_id_regex = ".*"
                            table_id_regex = ".*"
                        }
                    }
                }
            }
            conditions {
                created_after = "2023-10-02T15:01:23Z"
                types {
                    types = ["BIG_QUERY_TABLE_TYPE_TABLE", "BIG_QUERY_TABLE_TYPE_EXTERNAL_BIG_LAKE"]
                }
                or_conditions {
                    min_row_count = 10
                    min_age = "21600s"
                }
            }
        }
    }
    targets {
        big_query_target {
            filter {
                other_tables {}
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStartCloudSql(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
    parent = "projects/%{project}"
    description = "Description"
    display_name = "Display"
    inspect_config {
        info_types {
            name = "EMAIL_ADDRESS"
        }
    }
}
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
        cloud_sql_target {
            filter {
                collection {
                    include_regexes {
                        patterns {
                            project_id_regex = ".*"
                            instance_regex = ".*"
                            database_regex = "do-not-scan.*"
                            database_resource_name_regex = ".*"
                        }
                    }
                }
            }
            conditions {
                database_engines = ["MYSQL", "POSTGRES"]
                types = ["DATABASE_RESOURCE_TYPE_ALL_SUPPORTED_TYPES"]
            }
            disabled {}
        }
    }
    targets {
        cloud_sql_target {
            filter {
                others {}
            }
            generation_cadence {
                schema_modified_cadence {
                    types = ["NEW_COLUMNS"]
                    frequency = "UPDATE_FREQUENCY_MONTHLY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_MONTHLY"
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateCloudSql(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
    parent = "projects/%{project}"
    description = "Description"
    display_name = "Display"
    inspect_config {
        info_types {
            name = "EMAIL_ADDRESS"  
        }
    }
}
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
        cloud_sql_target {
            filter {
                collection {
                    include_regexes {
                        patterns {
                            project_id_regex = ".*"
                            instance_regex = ".*"
                            database_regex = ".*"
                            database_resource_name_regex = "mytable.*"
                        }
                    }
                }
            }
            conditions {
                database_engines = ["ALL_SUPPORTED_DATABASE_ENGINES"]
                types = ["DATABASE_RESOURCE_TYPE_TABLE"]
            }
            generation_cadence {
                schema_modified_cadence {
                    types = ["NEW_COLUMNS", "REMOVED_COLUMNS"]
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_MONTHLY"
                inspect_template_modified_cadence {
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
            }
        }
    }
    targets {
        cloud_sql_target {
            filter {
                others {}
            }
            generation_cadence {
                schema_modified_cadence {
                    types = ["NEW_COLUMNS", "REMOVED_COLUMNS"]
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_DAILY"
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigBqSingleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
	parent = "projects/%{project}"
	description = "Description"
	display_name = "Display"

	inspect_config {
		info_types {
			name = "EMAIL_ADDRESS"
		}
	}
}
resource "google_bigquery_dataset" "default" {
    dataset_id                  = "tf_test_%{random_suffix}"
    friendly_name               = "terraform-test"
    description                 = "Description for the dataset created by terraform"
    location                    = "US"
    default_table_expiration_ms = 3600000

    labels = {
        env = "default"
    }
}
resource "google_bigquery_table" "default" {
    dataset_id          = google_bigquery_dataset.default.dataset_id
    table_id            = "tf_test_%{random_suffix}"
    deletion_protection = false

    labels = {
        env = "default"
    }

    schema = <<EOF
        [
        {
            "name": "quantity",
            "type": "NUMERIC",
            "mode": "NULLABLE",
            "description": "The quantity"
        },
        {
            "name": "name",
            "type": "STRING",
            "mode": "NULLABLE",
            "description": "Name of the object"
        }
        ]
    EOF
}

resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    display_name = "display name"
    status = "RUNNING"

    targets {
        big_query_target {
            filter {
                table_reference {
                    dataset_id = google_bigquery_dataset.default.dataset_id
                    table_id = google_bigquery_table.default.table_id
				}
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigCloudSqlSingleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
    parent = "projects/%{project}"
    description = "Description"
    display_name = "Display"
    inspect_config {
        info_types {
            name = "EMAIL_ADDRESS"  
        }
    }
}
resource "google_sql_database_instance" "instance" {
    name             = "tf-test-instance-%{random_suffix}"
    database_version = "POSTGRES_14"
    region           = "%{location}"

    settings {
        tier = "db-f1-micro"
    }
  
    deletion_protection = false
}
resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "database"
}
data "google_project" "project" {
}
resource "google_project_iam_member" "dlp_role" {
    project = "%{project}"
    role    = "roles/dlp.projectdriver"
    member = "serviceAccount:service-${data.google_project.project.number}@dlp-api.iam.gserviceaccount.com"
}
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
        cloud_sql_target {
            filter {
                database_resource_reference {
                    project_id = "%{project}"
                    instance = google_sql_database_instance.instance.name
                    database = "database"
                    database_resource = "resource"
                }
            }
            conditions {
                database_engines = ["ALL_SUPPORTED_DATABASE_ENGINES"]
                types = ["DATABASE_RESOURCE_TYPE_TABLE"]
            }
            generation_cadence {
                schema_modified_cadence {
                    types = ["NEW_COLUMNS", "REMOVED_COLUMNS"]
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_MONTHLY"
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigSecretsStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
       secrets_target {}
    }
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigSecretsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "PAUSED"
    targets {
       secrets_target {}
    }
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigStartCloudStorage(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
    parent = "projects/%{project}"
    description = "Description"
    display_name = "Display"
    inspect_config {
        info_types {
            name = "EMAIL_ADDRESS"
        }
    }
}
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
        cloud_storage_target {
            filter {
                others {}
            }
			generation_cadence {
                inspect_template_modified_cadence {
                    frequency = "UPDATE_FREQUENCY_MONTHLY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_MONTHLY"
            }
        }

    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigUpdateCloudStorage(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
    parent = "projects/%{project}"
    description = "Description"
    display_name = "Display"
    inspect_config {
        info_types {
            name = "EMAIL_ADDRESS"  
        }
    }
}
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
        cloud_storage_target {
            filter {
                collection {
                    include_regexes {
                        patterns {
                            cloud_storage_regex {
                                project_id_regex = "foo-project"
                                bucket_name_regex = "bucket"
                            }
                        }
                    }
                }
            }
            conditions {
                created_after = "2023-10-02T15:01:23Z"
                min_age = "10800s"
                cloud_storage_conditions {
                    included_object_attributes = ["ALL_SUPPORTED_OBJECTS"]
                    included_bucket_attributes = ["ALL_SUPPORTED_BUCKETS"]
                }
            }
            generation_cadence {
                inspect_template_modified_cadence {
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_MONTHLY"
            }
        }
    }
    targets {
        cloud_storage_target {
            filter {
                collection {
                    include_regexes {
                        patterns {
                            cloud_storage_regex {
                                project_id_regex = "foo-project"
                                bucket_name_regex = "do-not-scan"
                            }
                        }
                    }
                }
            }
            disabled {}
        }
    }
    targets {
        cloud_storage_target {
            filter {
                others {}
            }
            generation_cadence {
                inspect_template_modified_cadence {
                    frequency = "UPDATE_FREQUENCY_MONTHLY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_MONTHLY"
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}

func testAccDataLossPreventionDiscoveryConfig_dlpDiscoveryConfigCloudStorageSingleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_loss_prevention_inspect_template" "basic" {
    parent = "projects/%{project}"
    description = "Description"
    display_name = "Display"
    inspect_config {
        info_types {
            name = "EMAIL_ADDRESS"  
        }
    }
}
resource "google_storage_bucket" "testbucket" {
	name                        = "dlp_test_bucket%{random_suffix}"
	location                    = "%{location}"
	uniform_bucket_level_access = true
}
data "google_project" "project" {
}
resource "google_project_iam_member" "dlp_role" {
    project = "%{project}"
    role    = "roles/dlp.projectdriver"
    member = "serviceAccount:service-${data.google_project.project.number}@dlp-api.iam.gserviceaccount.com"
}
resource "google_data_loss_prevention_discovery_config" "basic" {
    parent = "projects/%{project}/locations/%{location}"
    location = "%{location}"
    status = "RUNNING"
    targets {
        cloud_storage_target {
            filter {
                cloud_storage_resource_reference {
                    project_id = "%{project}"
                    bucket_name = google_storage_bucket.testbucket.name
                }
            }
			generation_cadence {
                inspect_template_modified_cadence {
                    frequency = "UPDATE_FREQUENCY_DAILY"
                }
                refresh_frequency = "UPDATE_FREQUENCY_DAILY"
            }
        }
    }
    inspect_templates = ["projects/%{project}/inspectTemplates/${google_data_loss_prevention_inspect_template.basic.name}"]
}
`, context)
}
