package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccStorageTransferJob_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageTransferJobDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccStorageTransferJob_basic(getTestProjectFromEnv()),
				Check: resource.ComposeTestCheckFunc(
					testAccStorageTransferJobExists("google_storage_transfer_job.transfer_job"),
					resource.TestCheckResourceAttrSet("google_storage_transfer_job.transfer_job", "name"),
					resource.TestCheckResourceAttrSet("google_storage_transfer_job.transfer_job", "description"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "project", getTestProjectFromEnv()),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.schedule_start_date.0.year", "2018"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.schedule_start_date.0.month", "10"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.schedule_start_date.0.day", "1"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.schedule_end_date.0.year", "2019"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.schedule_end_date.0.month", "10"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.schedule_end_date.0.day", "1"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.start_time_of_day.0.hours", "0"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.start_time_of_day.0.minutes", "30"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.start_time_of_day.0.seconds", "0"),
					resource.TestCheckResourceAttr("google_storage_transfer_job.transfer_job", "schedule.0.start_time_of_day.0.nanos", "0"),
				),
			},
		},
	})
}

func testAccStorageTransferJob_basic(project string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project       = "%s"
}

resource "google_storage_bucket" "data_source" {
  name          = "test-data-source-bucket-%s"
  project       = "%s"
  force_destroy = true
}

resource "google_storage_bucket_iam_member" "data_source" {
  bucket        = "${google_storage_bucket.data_source.name}"
  role          = "roles/storage.admin"
  member        = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"

  depends_on    = [
    "google_storage_bucket.data_source",
    "data.google_storage_transfer_project_service_account.default"
  ]
}

resource "google_storage_bucket" "data_sink" {
  name          = "test-data-sink-bucket-%s"
  project       = "%s"
  force_destroy = true
}

resource "google_storage_bucket_iam_member" "data_sink" {
  bucket        = "${google_storage_bucket.data_sink.name}"
  role          = "roles/storage.admin"
  member        = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"

  depends_on    = [
    "google_storage_bucket.data_sink",
    "data.google_storage_transfer_project_service_account.default"
  ]
}

resource "google_storage_transfer_job" "transfer_job" {
	description	= "transfer-job-%s"
	project     = "%s"

	transfer_spec {
		gcs_data_source {
			bucket_name = "${google_storage_bucket.data_source.name}"
		}
		gcs_data_sink {
			bucket_name = "${google_storage_bucket.data_sink.name}"
		}
	}

	schedule {
		schedule_start_date {
			year	= 2018
			month	= 10
			day		= 1
		}
		schedule_end_date {
			year	= 2019
			month	= 10
			day		= 1
		}
		start_time_of_day {
			hours	= 0
			minutes	= 30
			seconds	= 0
			nanos	= 0
		}
	}

	depends_on = [
		"google_storage_bucket_iam_member.data_source",
		"google_storage_bucket_iam_member.data_sink",
	]
}
`, project, acctest.RandString(10), project, acctest.RandString(10), project, acctest.RandString(10), project)
}

func testAccStorageTransferJobExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
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

		_, err = config.clientStorageTransfer.TransferJobs.Get(name).ProjectId(project).Do()
		if err != nil {
			return fmt.Errorf("Job does not exist")
		}

		return nil
	}
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
