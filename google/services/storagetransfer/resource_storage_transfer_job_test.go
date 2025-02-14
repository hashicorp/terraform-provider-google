// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagetransfer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccStorageTransferReplicationJob_basic(t *testing.T) {
	t.Parallel()

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gs-project-accounts.iam.gserviceaccount.com",
			Role:   "roles/pubsub.publisher",
		},
		{
			Member: "serviceAccount:project-{project_number}@storage-transfer-service.iam.gserviceaccount.com",
			Role:   "roles/storagetransfer.serviceAgent",
		},
	})

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferReplicationJobDescription := acctest.RandString(t, 10)
	testUpdatedTransferReplicationJobDescription := acctest.RandString(t, 10)
	testOverwriteWhen := []string{"ALWAYS", "NEVER", "DIFFERENT"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferReplicationJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferReplicationJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferReplicationJob_basic(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testUpdatedTransferReplicationJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferReplicationJob_with_transferOptions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testUpdatedTransferReplicationJobDescription, true, false, testOverwriteWhen[0]),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferReplicationJob_with_transferOptions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testUpdatedTransferReplicationJobDescription, false, false, testOverwriteWhen[1]),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferReplicationJob_with_transferOptions(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testUpdatedTransferReplicationJobDescription, false, false, testOverwriteWhen[2]),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}})
}

func TestAccStorageTransferJob_transferJobName(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testTransferJobName := fmt.Sprintf("tf-test-transfer-job-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_transferJobName(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription, testTransferJobName),
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

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:project-{project_number}@storage-transfer-service.iam.gserviceaccount.com",
			Role:   "roles/pubsub.admin",
		},
	})

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

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:project-{project_number}@storage-transfer-service.iam.gserviceaccount.com",
			Role:   "roles/pubsub.admin",
		},
	})

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

func TestAccStorageTransferJob_eventStream(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(t, 10)
	testDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testPubSubTopicName := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	testEventStreamPubSubTopicName := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	testPubSubSubscriptionName := fmt.Sprintf("tf-test-subscription-%s", acctest.RandString(t, 10))
	eventStreamStart := []string{"2014-10-02T15:01:23Z", "2019-10-02T15:01:23Z"}
	eventStreamEnd := []string{"2022-10-02T15:01:23Z", "2032-10-02T15:01:23Z"}

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
				Config: testAccStorageTransferJob_eventStream(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testEventStreamPubSubTopicName, testPubSubSubscriptionName, testTransferJobDescription, eventStreamStart[0], eventStreamEnd[0]),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_eventStream(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testEventStreamPubSubTopicName, testPubSubSubscriptionName, testTransferJobDescription, eventStreamStart[1], eventStreamEnd[0]),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_eventStream(envvar.GetTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testEventStreamPubSubTopicName, testPubSubSubscriptionName, testTransferJobDescription, eventStreamStart[1], eventStreamEnd[1]),
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

func TestAccStorageTransferJob_hdfsSource(t *testing.T) {
	t.Parallel()

	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:project-{project_number}@storage-transfer-service.iam.gserviceaccount.com",
			Role:   "roles/pubsub.admin",
		},
	})

	testDataSinkName := acctest.RandString(t, 10)
	otherDataSinkName := acctest.RandString(t, 10)
	testTransferJobDescription := acctest.RandString(t, 10)
	testSourceAgentPoolName := fmt.Sprintf("tf-test-source-agent-pool-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_hdfsSource(envvar.GetTestProjectFromEnv(), testDataSinkName, testTransferJobDescription, testSourceAgentPoolName, "/root/", ""),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_hdfsSource(envvar.GetTestProjectFromEnv(), otherDataSinkName, testTransferJobDescription, testSourceAgentPoolName, "/root/dir/", "object/"),
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

func testAccStorageTransferJob_eventStream(project string, dataSourceBucketName string, dataSinkBucketName string, pubsubTopicName string, pubsubSubscriptionName string, transferJobDescription string, eventStreamStart string, eventStreamEnd string) string {
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

resource "google_pubsub_subscription_iam_member" "editor" {
  subscription = google_pubsub_subscription.example.name
  role         = "roles/editor"
  member       = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_pubsub_topic" "example" {
  name = "%s"
}

resource "google_pubsub_subscription" "example" {
  name  = "%s"
  topic = google_pubsub_topic.example.name

  ack_deadline_seconds = 20

  labels = {
    foo = "bar"
  }

  push_config {
    push_endpoint = "https://example.com/push"

    attributes = {
      x-goog-version = "v1"
    }
  }
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  event_stream {
    name = google_pubsub_subscription.example.id
    event_stream_start_time = "%s"
    event_stream_expiration_time = "%s"
  }

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
    google_pubsub_subscription_iam_member.editor,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, pubsubTopicName, pubsubSubscriptionName, transferJobDescription, project, eventStreamStart, eventStreamEnd)
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

func testAccStorageTransferJob_transferJobName(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, testTransferJobName string) string {
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
    name        = "transferJobs/%s"
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
  `, project, dataSourceBucketName, project, dataSinkBucketName, project, testTransferJobName, transferJobDescription, project)
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

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }
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

  depends_on = [google_storage_bucket_iam_member.data_sink]
}
`, project, dataSinkBucketName, project, sourceAgentPoolName, transferJobDescription, project)
}

func testAccStorageTransferJob_hdfsSource(project string, dataSinkBucketName string, transferJobDescription string, sourceAgentPoolName string, hdfsPath string, gcsPath string) string {
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

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }
}

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
    source_agent_pool_name = google_storage_transfer_agent_pool.foo.id
    hdfs_data_source {
    	path = "%s"
    }
    gcs_data_sink {
      bucket_name = google_storage_bucket.data_sink.name
      path  = "%s"
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

  depends_on = [google_storage_bucket_iam_member.data_sink]
}
`, project, dataSinkBucketName, project, sourceAgentPoolName, transferJobDescription, project, hdfsPath, gcsPath)
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

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }
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

  depends_on = [google_storage_bucket_iam_member.data_source]
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

func testAccStorageTransferReplicationJob_basic(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string) string {
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

  replication_spec {
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
    google_storage_bucket_iam_member.data_sink
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, transferJobDescription, project)
}

func testAccStorageTransferReplicationJob_with_transferOptions(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string, overwriteObjectsAlreadyExistingInSink bool, deleteObjectsUniqueInSink bool, overwriteWhenVal string) string {
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

  replication_spec {
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
      overwrite_when = "%s"
      delete_objects_from_source_after_transfer = false
    }
    object_conditions {
      last_modified_since = "2020-01-01T00:00:00Z"
      last_modified_before = "2020-01-01T00:00:00Z"
      exclude_prefixes = [
        "a/b/c", 
      ]
      include_prefixes = [
        "a/b"
      ]
      max_time_elapsed_since_last_modification="300s"
      min_time_elapsed_since_last_modification="3s"
    }
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, transferJobDescription, project, overwriteObjectsAlreadyExistingInSink, deleteObjectsUniqueInSink, overwriteWhenVal)
}
