package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccStorageTransferJob_basic(t *testing.T) {
	t.Parallel()

	testDataSourceBucketName := acctest.RandString(10)
	testDataSinkName := acctest.RandString(10)
	testTransferJobDescription := acctest.RandString(10)
	testUpdatedDataSourceBucketName := acctest.RandString(10)
	testUpdatedDataSinkBucketName := acctest.RandString(10)
	testUpdatedTransferJobDescription := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroy,
		Steps: []resource.TestStep{
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

	testDataSourceBucketName := acctest.RandString(10)
	testDataSinkName := acctest.RandString(10)
	testTransferJobDescription := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroy,
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

func testAccStorageTransferJobDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

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

		res, err := config.clientStorageTransfer.TransferJobs.Get(name).ProjectId(project).Do()
		if res.Status != "DELETED" {
			return fmt.Errorf("Transfer Job not set to DELETED")
		}
		if err != nil {
			return fmt.Errorf("Transfer Job does not exist, should exist and be DELETED")
		}
	}

	return nil
}

func testAccStorageTransferJob_basic(project string, dataSourceBucketName string, dataSinkBucketName string, transferJobDescription string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "%s"
  project       = "%s"
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
