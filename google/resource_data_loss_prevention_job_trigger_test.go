package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerBasic(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerUpdateExample2(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFields(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.identifying_fields",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
			{
				Config: testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFieldsUpdate(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.identifying_fields_update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func TestAccDataLossPreventionJobTrigger_dlpJobTriggerPubsub(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": getTestProjectFromEnv(),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataLossPreventionJobTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLossPreventionJobTrigger_publishToPubSub(context),
			},
			{
				ResourceName:            "google_data_loss_prevention_job_trigger.pubsub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent"},
			},
		},
	})
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "basic" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFields(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "identifying_fields" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			big_query_options {
				table_reference {
					project_id = "project"
					dataset_id = "dataset"
					table_id = "table_to_scan"
				}
				rows_limit = 1000
				sample_method = "RANDOM_START"
				identifying_fields {
					name = "field"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "basic" {
	parent = "projects/%{project}"
	description = "An updated description"
	display_name = "Different"

	triggers {
		schedule {
			recurrence_period_duration = "86500s"
		}
	}

	inspect_job {
		inspect_template_name = "other"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "different"
						dataset_id = "asdf"
					}
				}
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_dlpJobTriggerIdentifyingFieldsUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "identifying_fields_update" {
	parent = "projects/%{project}"
	description = "An updated description"
	display_name = "Different"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			save_findings {
				output_config {
					table {
						project_id = "project"
						dataset_id = "dataset123"
					}
				}
			}
		}
		storage_config {
			big_query_options {
				table_reference {
					project_id = "project"
					dataset_id = "dataset"
					table_id = "table_to_scan"
				}
				rows_limit = 1000
				sample_method = "RANDOM_START"
				identifying_fields {
					name = "different"
				}
			}
		}
	}
}
`, context)
}

func testAccDataLossPreventionJobTrigger_publishToPubSub(context map[string]interface{}) string {
	return Nprintf(`
resource "google_data_loss_prevention_job_trigger" "pubsub" {
	parent = "projects/%{project}"
	description = "Starting description"
	display_name = "display"

	triggers {
		schedule {
			recurrence_period_duration = "86400s"
		}
	}

	inspect_job {
		inspect_template_name = "fake"
		actions {
			pub_sub {
				topic = "projects/%{project}/topics/bar"
			}
		}
		storage_config {
			cloud_storage_options {
				file_set {
					url = "gs://mybucket/directory/"
				}
			}
		}
	}
}
`, context)
}
