// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagetransfer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageTransferJob_basic(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testUpdatedDataSourceBucketName := acctest.RandString(t, 10)
	testUpdatedDataSinkBucketName := acctest.RandString(t, 10)
	testUpdatedTransferJobDescription := acctest.RandString(t, 10)
	testPubSubTopicName := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_omitNotificationConfig(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_omitSchedule(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testUpdatedDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testUpdatedDataSourceBucketName, testUpdatedDataSinkBucketName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testUpdatedDataSourceBucketName, testUpdatedDataSinkBucketName, testUpdatedTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageTransferJob_omitScheduleEndDate(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_omitScheduleEndDate(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageTransferJob_posixSource(t *testing.T) {
	t.Parallel()

	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testSourceAgentPoolName := fmt.Sprintf("tf-test-source-agent-pool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_posixSource(envvar.GetTestProjectFromEnv(), testDataSinkName, testTransferJobDescription, testSourceAgentPoolName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccStorageTransferJob_posixSink(t *testing.T) {
	t.Parallel()

	testDataSourceName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testSinkAgentPoolName := fmt.Sprintf("tf-test-sink-agent-pool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_posixSink(envvar.GetTestProjectFromEnv(), testDataSourceName, testTransferJobDescription, testSinkAgentPoolName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageTransferJob_transferOptions(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testOverwriteWhen := []string{"ALWAYS", "NEVER", "DIFFERENT"}
	testPubSubTopicName := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_transferOptions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, false, false, false, testOverwriteWhen[0], testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_transferOptions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, true, true, false, testOverwriteWhen[1], testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_transferOptions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, true, false, true, testOverwriteWhen[2], testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageTransferJob_objectConditions(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testPubSubTopicName := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_objectConditions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageTransferJob_notificationConfig(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	noneNotificationConfigPayloadFormat := "NONE"
	testPubSubTopicName := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	testPubSubTopicNameUpdate := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicName),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicNameUpdate),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_notificationPayloadFormat(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicNameUpdate, noneNotificationConfigPayloadFormat),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_notificationEventTypesUpdate(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicNameUpdate, noneNotificationConfigPayloadFormat),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_omitNotificationEventTypes(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testPubSubTopicNameUpdate),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStorageTransferJobDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_transfer_job" {
				continue
			}

			rs_attr := rs.Primary.Attributes
			name, ok := rs_attr["name"]
			if !ok {
				return fmt.Errorf("No name set")
			}

			project, err := acctest.GetTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			res, err := config.NewStorageTransferClient(config.UserAgent).TransferJobs.Get(name, project).Do()
			if err != nil {
				return fmt.Errorf("Transfer Job does not exist, should exist and be DELETED")
			}
			if res.Status != "DELETED" {
				return fmt.Errorf("Transfer Job not set to DELETED")
			}
		}

		return nil
	}
}

func testAccStorageTransferJob_omitSchedule(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, transferJobDescription, project)
}

func testAccStorageTransferJob_omitNotificationConfig(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, transferJobDescription, project)
}

func testAccStorageTransferJob_basic(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, pubsubTopicName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_topic_iam_member" "notification_config" {
  topic = google_pubsub_topic.topic.id
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  notification_config {
    pubsub_topic  = google_pubsub_topic.topic.id
    event_types   = [
      "TRANSFER_OPERATION_SUCCESS",
      "TRANSFER_OPERATION_FAILED"
    ]
    payload_format = "JSON"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
    google_pubsub_topic_iam_member.notification_config,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubsubTopicName, transferJobDescription, project)
}

func testAccStorageTransferJob_omitScheduleEndDate(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, transferJobDescription, project)
}

func testAccStorageTransferJob_posixSource(project string, dataSinkBucketName string, transferJobDescription string, sourceAgentPoolName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_project_iam_member" "pubsub" {
	project = data.google_storage_transfer_project_service_account.default.project
  role    = "roles/pubsub.admin"
  member  = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }

  depends_on = [google_project_iam_member.pubsub]
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    source_agent_pool_name = google_storage_transfer_agent_pool.foo.id
    posix_data_source {
    	root_directory = "/some/path"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
  }

  depends_on = [
    google_storage_bucket_iam_member.data_sink,
    google_project_iam_member.pubsub
  ]
}
`, project, dataSinkBucketName, project, sourceAgentPoolName, transferJobDescription, project)
}

func testAccStorageTransferJob_posixSink(project string, dataSourceBucketName string, transferJobDescription string, sinkAgentPoolName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_project_iam_member" "pubsub" {
	project = data.google_storage_transfer_project_service_account.default.project
  role    = "roles/pubsub.admin"
  member  = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }

  depends_on = [google_project_iam_member.pubsub]
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    sink_agent_pool_name  = google_storage_transfer_agent_pool.foo.id
    posix_data_sink {
    	root_directory = "/some/path"
    }
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_project_iam_member.pubsub
  ]
}
`, project, dataSourceBucketName, project, sinkAgentPoolName, transferJobDescription, project)
}

func testAccStorageTransferJob_transferOptions(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, overwriteObjectsAlreadyExistingInSink bool, deleteObjectsUniqueInSink bool, deleteObjectsFromSourceAfterTransfer bool, overwriteWhenVal string, pubSubTopicName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_topic_iam_member" "notification_config" {
  topic = google_pubsub_topic.topic.id
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
    transfer_options {
      overwrite_objects_already_existing_in_sink = %t
      delete_objects_unique_in_sink = %t
      delete_objects_from_source_after_transfer = %t
      overwrite_when = "%s"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  notification_config {
    pubsub_topic  = google_pubsub_topic.topic.id
    event_types   = [
      "TRANSFER_OPERATION_SUCCESS",
      "TRANSFER_OPERATION_FAILED"
    ]
    payload_format = "JSON"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
    google_pubsub_topic_iam_member.notification_config,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubSubTopicName, transferJobDescription, project, overwriteObjectsAlreadyExistingInSink, deleteObjectsUniqueInSink, deleteObjectsFromSourceAfterTransfer, overwriteWhenVal)
}

func testAccStorageTransferJob_objectConditions(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, pubSubTopicName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_topic_iam_member" "notification_config" {
  topic = google_pubsub_topic.topic.id
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
    object_conditions {
      last_modified_since = "2020-01-01T00:00:00Z"
      last_modified_before = "2020-01-01T00:00:00Z"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  notification_config {
    pubsub_topic  = google_pubsub_topic.topic.id
    event_types   = [
      "TRANSFER_OPERATION_SUCCESS",
      "TRANSFER_OPERATION_FAILED"
    ]
    payload_format = "JSON"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
    google_pubsub_topic_iam_member.notification_config,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubSubTopicName, transferJobDescription, project)
}

func testAccStorageTransferJob_notificationPayloadFormat(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, pubsubTopicName string, notificationPayloadFormat string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_topic_iam_member" "notification_config" {
  topic = google_pubsub_topic.topic.id
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  notification_config {
    pubsub_topic  = google_pubsub_topic.topic.id
    event_types   = [
      "TRANSFER_OPERATION_SUCCESS",
      "TRANSFER_OPERATION_FAILED"
    ]
    payload_format = "%s"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
    google_pubsub_topic_iam_member.notification_config,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubsubTopicName, transferJobDescription, project, notificationPayloadFormat)
}

func testAccStorageTransferJob_notificationEventTypesUpdate(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, pubsubTopicName string, notificationPayloadFormat string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_topic_iam_member" "notification_config" {
  topic = google_pubsub_topic.topic.id
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  notification_config {
    pubsub_topic  = google_pubsub_topic.topic.id
    event_types   = [
      "TRANSFER_OPERATION_ABORTED"
    ]
    payload_format = "%s"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
    google_pubsub_topic_iam_member.notification_config,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubsubTopicName, transferJobDescription, project, notificationPayloadFormat)
}

func testAccStorageTransferJob_omitNotificationEventTypes(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, pubsubTopicName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket = google_storage_bucket.data_source.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket = google_storage_bucket.data_sink.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "topic" {
  name = "%s"
}

resource "google_pubsub_topic_iam_member" "notification_config" {
  topic = google_pubsub_topic.topic.id
  role = "roles/pubsub.publisher"
  member = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    gcs_data_source {
      bucket_name = google_storage_bucket.data_source.name
      path  = "foo/bar/"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "foo/bar/"
    }
  }

  schedule {
    schedule_start_date {
      year  = 2018
      month = 10
      day   = 1
    }
    schedule_end_date {
      year  = 2019
      month = 10
      day   = 1
    }
    start_time_of_day {
      hours   = 0
      minutes = 30
      seconds = 0
      nanos   = 0
    }
	  repeat_interval = "604800s"
  }

  notification_config {
    pubsub_topic  = google_pubsub_topic.topic.id
    payload_format = "JSON"
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
    google_pubsub_topic_iam_member.notification_config,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubsubTopicName, transferJobDescription, project)
}
