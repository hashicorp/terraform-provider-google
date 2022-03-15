package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStorageTransferJob_basic(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := randString(t, 10)
	testDataSinkName := randString(t, 10)
	testTransferJobDescription := randString(t, 10)
	testUpdatedDataSourceBucketName := randString(t, 10)
	testUpdatedDataSinkBucketName := randString(t, 10)
	testUpdatedTransferJobDescription := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_omitSchedule(getTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(getTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(getTestProjectFromEnv(), testUpdatedDataSourceBucketName, testDataSinkName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(getTestProjectFromEnv(), testUpdatedDataSourceBucketName, testUpdatedDataSinkBucketName, testTransferJobDescription),
			},
			{
				ResourceName:      "google_storage_transfer_job.transfer_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferJob_basic(getTestProjectFromEnv(), testUpdatedDataSourceBucketName, testUpdatedDataSinkBucketName, testUpdatedTransferJobDescription),
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

	testDataSourceBucketName := randString(t, 10)
	testDataSinkName := randString(t, 10)
	testTransferJobDescription := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_omitScheduleEndDate(getTestProjectFromEnv(), testDataSourceBucketName, testDataSinkName, testTransferJobDescription),
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

	testDataSinkName := randString(t, 10)
	testTransferJobDescription := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_posixSource(getTestProjectFromEnv(), testDataSinkName, testTransferJobDescription),
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

	testDataSourceName := randString(t, 10)
	testTransferJobDescription := randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferJob_posixSink(getTestProjectFromEnv(), testDataSourceName, testTransferJobDescription),
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
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_transfer_job" {
				continue
			}

			rs_attr := rs.Primary.Attributes
			name, ok := rs_attr["name"]
			if !ok {
				return fmt.Errorf("No name set")
			}

			project, err := getTestProject(rs.Primary, config)
			if err != nil {
				return err
			}

			res, err := config.NewStorageTransferClient(config.userAgent).TransferJobs.Get(name, project).Do()
			if res.Status != "DELETED" {
				return fmt.Errorf("Transfer Job not set to DELETED")
			}
			if err != nil {
				return fmt.Errorf("Transfer Job does not exist, should exist and be DELETED")
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

func testAccStorageTransferJob_basic(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
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
  }

  depends_on = [
    google_storage_bucket_iam_member.data_source,
    google_storage_bucket_iam_member.data_sink,
  ]
}
`, project, dataSourceBucketName, project, dataSinkBucketName, project, transferJobDescription, project)
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

func testAccStorageTransferJob_posixSource(project string, dataSinkBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_sink" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
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

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
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
`, project, dataSinkBucketName, project, transferJobDescription, project)
}

func testAccStorageTransferJob_posixSink(project string, dataSourceBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
  location      = "US"
  force_destroy = true
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

resource "google_storage_transfer_job" "transfer_job" {
  description = "%s"
  project     = "%s"

  transfer_spec {
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
`, project, dataSourceBucketName, project, transferJobDescription, project)
}
