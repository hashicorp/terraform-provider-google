package google

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccDataflowJobCreate(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJob,
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(
						"google_dataflow_job.big_data"),
				),
			},
		},
	})
}

func TestAccDataflowJobRegionCreate(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJobRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobRegionExists(
						"google_dataflow_job.big_data"),
				),
			},
		},
	})
}

func TestAccDataflowJobCreateWithServiceAccount(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowJobWithServiceAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(
						"google_dataflow_job.big_data"),
					testAccDataflowJobHasServiceAccount(
						"google_dataflow_job.big_data"),
				),
			},
		},
	})
}

func testAccCheckDataflowJobDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_dataflow_job" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		job, err := config.clientDataflow.Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
		if job != nil {
			if _, ok := dataflowTerminalStatesMap[job.CurrentState]; !ok {
				return fmt.Errorf("Job still present")
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckDataflowJobRegionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_dataflow_job" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		job, err := config.clientDataflow.Projects.Locations.Jobs.Get(config.Project, "us-central1", rs.Primary.ID).Do()
		if job != nil {
			if _, ok := dataflowTerminalStatesMap[job.CurrentState]; !ok {
				return fmt.Errorf("Job still present")
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func testAccDataflowJobExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientDataflow.Projects.Jobs.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Job does not exist")
		}

		return nil
	}
}

func testAccDataflowJobHasServiceAccount(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		// Check that the service account was applied to the Dataflow job's
		// generated instance template.
		if serviceAccountEmail, ok := rs.Primary.Attributes["service_account_email"]; ok {
			filter := fmt.Sprintf("properties.labels.dataflow_job_id = %s", rs.Primary.ID)
			var serviceAccounts []*compute.ServiceAccount

			// Wait for instance template generation.
			err := resource.Retry(1*time.Minute, func() *resource.RetryError {
				var err error
				instanceTemplates, err :=
					config.clientCompute.InstanceTemplates.List(config.Project).Filter(filter).MaxResults(2).Fields("items/properties/serviceAccounts/email").Do()
				if err != nil {
					return resource.NonRetryableError(err)
				}
				if len(instanceTemplates.Items) == 0 {
					return resource.RetryableError(fmt.Errorf("no instance template found for dataflow job"))
				}
				if len(instanceTemplates.Items) > 1 {
					return resource.NonRetryableError(fmt.Errorf("Wrong number of matching instance templates for dataflow job: %s, %d", rs.Primary.ID, len(instanceTemplates.Items)))
				}
				serviceAccounts = instanceTemplates.Items[0].Properties.ServiceAccounts
				return nil
			})

			if err != nil {
				return fmt.Errorf("Error getting service account from instance template: %s", err)
			}

			if len(serviceAccounts) > 1 {
				return fmt.Errorf("Found multiple service accounts for dataflow job: %s, %d", rs.Primary.ID, len(serviceAccounts))
			}
			if serviceAccountEmail != serviceAccounts[0].Email {
				return fmt.Errorf("Service account mismatch: %s != %s", serviceAccountEmail, serviceAccounts[0].Email)
			}
		}

		return nil
	}
}

func testAccDataflowJobRegionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientDataflow.Projects.Locations.Jobs.Get(config.Project, "us-central1", rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Job does not exist")
		}

		return nil
	}
}

var testAccDataflowJob = fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
	name = "dfjob-test-%s-temp"

	force_destroy = true
}

resource "google_dataflow_job" "big_data" {
	name = "dfjob-test-%s"

	template_gcs_path = "gs://dataflow-templates/wordcount/template_file"
	temp_gcs_location = "${google_storage_bucket.temp.url}"

	parameters = {
		inputFile = "gs://dataflow-samples/shakespeare/kinglear.txt"
		output    = "${google_storage_bucket.temp.url}/output"
	}
	zone = "us-central1-f"
	project = "%s"

	on_delete = "cancel"
}`, acctest.RandString(10), acctest.RandString(10), getTestProjectFromEnv())

var testAccDataflowJobRegion = fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
	name = "dfjob-test-%s-temp"

	force_destroy = true
}

resource "google_dataflow_job" "big_data" {
	name = "dfjob-test-%s"

	template_gcs_path = "gs://dataflow-templates/wordcount/template_file"
	temp_gcs_location = "${google_storage_bucket.temp.url}"

	parameters = {
		inputFile = "gs://dataflow-samples/shakespeare/kinglear.txt"
		output    = "${google_storage_bucket.temp.url}/output"
	}
	region  = "us-central1"
	zone    = "us-central1-c"
	project = "%s"

	on_delete = "cancel"
}`, acctest.RandString(10), acctest.RandString(10), getTestProjectFromEnv())

var testAccDataflowJobWithServiceAccount = fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
	name = "dfjob-test-%s-temp"

	force_destroy = true
}

resource "google_service_account" "dataflow-sa" {
  account_id   = "dataflow-sa"
  display_name = "DataFlow Service Account"
}

resource "google_storage_bucket_iam_member" "dataflow-gcs" {
  bucket = "${google_storage_bucket.temp.name}"
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "google_project_iam_member" "dataflow-worker" {
  role   = "roles/dataflow.worker"
  member = "serviceAccount:${google_service_account.dataflow-sa.email}"
}

resource "google_dataflow_job" "big_data" {
	name = "dfjob-test-%s"

	template_gcs_path = "gs://dataflow-templates/wordcount/template_file"
	temp_gcs_location = "${google_storage_bucket.temp.url}"

	parameters = {
		inputFile = "gs://dataflow-samples/shakespeare/kinglear.txt"
		output    = "${google_storage_bucket.temp.url}/output"
	}
	zone = "us-central1-f"
	project = "%s"
	service_account_email = "${google_service_account.dataflow-sa.email}"

	on_delete = "cancel"
}`, acctest.RandString(10), acctest.RandString(10), getTestProjectFromEnv())
