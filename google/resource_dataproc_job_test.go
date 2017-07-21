package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/googleapi"
	"regexp"
)

func TestAccDataprocJob_failForMissingJobConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataprocJob_missingJobConf(),
				ExpectError: regexp.MustCompile("You must define and configure exactly one xxx_config block"),
			},
		},
	})
}

func TestAccDataprocJob_PySpark(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_pySpark(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJob(
						"google_dataproc_job.pyspark", "pyspark_config"),
				),
			},
		},
	})
}

func TestAccDataprocJob_Spark(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_spark(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJob(
						"google_dataproc_job.spark", "spark_config"),
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

		if rs.Primary.ID == "" {
			return fmt.Errorf("Unable to verify delete of dataproc job ID is empty")
		}
		attributes := rs.Primary.Attributes

		_, err := config.clientDataproc.Projects.Regions.Jobs.Get(
			config.Project, attributes["region"], rs.Primary.ID).Do()

		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
				return nil
			} else if ok {
				return fmt.Errorf("Error make GCP platform call: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
			}
			return fmt.Errorf("Error make GCP platform call: %s", err.Error())
		}
		return fmt.Errorf("Dataproc job still exists")
	}

	return nil
}

func testAccCheckDataprocJob(n, jobType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		jobId := s.RootModule().Resources[n].Primary.ID
		config := testAccProvider.Meta().(*Config)
		job, err := config.clientDataproc.Projects.Regions.Jobs.Get(
			config.Project, attributes["region"], jobId).Do()
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
		}

		if jobType == "pyspark_config" {
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.main_python_file", job.PysparkJob.MainPythonFileUri})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.additional_python_files", job.PysparkJob.PythonFileUris})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.jar_files", job.PysparkJob.JarFileUris})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.args", job.PysparkJob.Args})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.properties", job.PysparkJob.Properties})
		}
		if jobType == "spark_config" {
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.main_class", job.SparkJob.MainClass})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.main_jar", job.SparkJob.MainJarFileUri})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.jar_files", job.SparkJob.JarFileUris})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.args", job.SparkJob.Args})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.properties", job.SparkJob.Properties})
		}

		for _, attrs := range clusterTests {
			if c := checkMatch(attributes, attrs.tf_attr, attrs.gcp_attr); c != "" {
				return fmt.Errorf(c)
			}
		}

		return nil
	}
}

func testAccDataprocJob_missingJobConf() string {
	return `
resource "google_dataproc_job" "missing_config" {
    cluster      = "na"
    force_delete = true
}`
}

func testAccDataprocJob_pySpark(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name = "dproc-job-test-%s-url"
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
`, rnd)
}

func testAccDataprocJob_spark(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name = "dproc-job-test-%s-url"
	zone = "us-central1-a"

	worker_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}

resource "google_dataproc_job" "spark" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    force_delete = true

    spark_config {
        main_class = "org.apache.spark.examples.SparkPi"
        jar_files  = ["file:///usr/lib/spark/examples/jars/spark-examples.jar"]
        args       = ["1000"]
    }
}
`, rnd)
}
