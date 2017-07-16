package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataprocJob_PySpark(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_pySpark,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJob(
						"google_dataproc_job.pyspark"),
				),
			},
		},
	})
}

func testAccCheckDataprocJobDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_dataproc_job" {
			continue
		}

		attributes := rs.Primary.Attributes
		_, err := config.clientDataproc.Projects.Regions.Jobs.Get(
			config.Project, attributes["region"], attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("Job still exists")
		}
	}

	return nil
}

func testAccCheckDataprocJob(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)
		job, err := config.clientDataproc.Projects.Regions.Jobs.Get(
			config.Project, attributes["region"], attributes["id"]).Do()
		if err != nil {
			return err
		}

		type jobTestField struct {
			tf_attr  string
			gcp_attr interface{}
		}

		clusterTests := []jobTestField{

			{"cluster", job.Placement.ClusterName},
			{"labels", job.Labels},

			{"pyspark_config.0.main_python_file", job.PysparkJob.MainPythonFileUri},
			{"pyspark_config.0.additional_python_files", job.PysparkJob.PythonFileUris},
			{"pyspark_config.0.jar_files", job.PysparkJob.JarFileUris},
			{"pyspark_config.0.args", job.PysparkJob.Args},
			{"pyspark_config.0.properties", job.PysparkJob.Properties},
			{"pyspark_config.0.main_python_file", job.PysparkJob.MainPythonFileUri},
		}

		for _, attrs := range clusterTests {
			if c := checkMatch(attributes, attrs.tf_attr, attrs.gcp_attr); c != "" {
				return fmt.Errorf(c)
			}
		}

		return nil
	}
}

var testAccDataprocJob_pySpark = fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name = "cluster-test-%s"
	zone = "us-central1-a"

	worker_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}

resource "google_dataproc_job" "pyspark" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    force_delete = true

    pyspark_config {
        main_python_file = "gs://dataproc-examples-2f10d78d114f6aaec76462e3c310f31f/src/pyspark/hello-world/hello-world.py"
    }
}
`, acctest.RandString(10))
