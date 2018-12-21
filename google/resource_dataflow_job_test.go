package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

	parameters {
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

	parameters {
		inputFile = "gs://dataflow-samples/shakespeare/kinglear.txt"
		output    = "${google_storage_bucket.temp.url}/output"
	}
	region  = "us-central1"
	zone    = "us-central1-c"
	project = "%s"

	on_delete = "cancel"
}`, acctest.RandString(10), acctest.RandString(10), getTestProjectFromEnv())
