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
					testAccCheckDataprocJobAttrMatch(
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
					testAccCheckDataprocJobAttrMatch(
						"google_dataproc_job.spark", "spark_config"),
				),
			},
		},
	})
}

func TestAccDataprocJob_Hadoop(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_hadoop(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJobAttrMatch(
						"google_dataproc_job.hadoop", "hadoop_config"),
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

func testAccCheckDataprocJobAttrMatch(n, jobType string) resource.TestCheckFunc {
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
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.args", job.PysparkJob.Args})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.jars", job.PysparkJob.JarFileUris})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.files", job.PysparkJob.PythonFileUris})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.archives", job.PysparkJob.ArchiveUris})
			clusterTests = append(clusterTests, jobTestField{"pyspark_config.0.properties", job.PysparkJob.Properties})
		}
		if jobType == "spark_config" {
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.main_class", job.SparkJob.MainClass})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.main_jar", job.SparkJob.MainJarFileUri})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.args", job.SparkJob.Args})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.jars", job.SparkJob.JarFileUris})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.files", job.SparkJob.FileUris})
			clusterTests = append(clusterTests, jobTestField{"spark_config.0.archives", job.SparkJob.ArchiveUris})
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
	name   = "dproc-job-test-%s"
	region = "us-central1"

    # Keep the costs down with smallest config we can get away with
    # Making use of the single node cluster feature (1 x master)
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }

	worker_config {}
	master_config {
		machine_type      = "n1-standard-2"
		boot_disk_size_gb = 10
	}


}

resource "google_dataproc_job" "pyspark" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    region       = "${google_dataproc_cluster.basic.region}"
    force_delete = true

    pyspark_config {
        main_python_file = "gs://dataproc-examples-2f10d78d114f6aaec76462e3c310f31f/src/pyspark/hello-world/hello-world.py"
        properties = {
            "spark.logConf" = "true"
        }
    }
}
`, rnd)
}

func testAccDataprocJob_spark(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name   = "dproc-job-test-%s"
	region = "us-central1"

    # Keep the costs down with smallest config we can get away with
    # Making use of the single node cluster feature (1 x master)
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }

    worker_config {}
	master_config {
		machine_type      = "n1-standard-2"
		boot_disk_size_gb = 10
	}
}

resource "google_dataproc_job" "spark" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    region       = "${google_dataproc_cluster.basic.region}"
    force_delete = true

    spark_config {
        main_class = "org.apache.spark.examples.SparkPi"
        jars       = ["file:///usr/lib/spark/examples/jars/spark-examples.jar"]
        args       = ["1000"]
        properties = {
            "spark.logConf" = "true"
        }
    }
}

output "spark_status" {
    value = "${google_dataproc_job.spark.status}"
}
`, rnd)
}

func testAccDataprocJob_hadoop(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name   = "dproc-job-test-%s"
	region = "us-central1"

    # Keep the costs down with smallest config we can get away with
    # Making use of the single node cluster feature (1 x master)
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }

    worker_config {}
	master_config {
		machine_type      = "n1-standard-2"
		boot_disk_size_gb = 10
	}
}

resource "google_dataproc_job" "hadoop" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    region       = "${google_dataproc_cluster.basic.region}"
    force_delete = true

    hadoop_config {
		main_jar   =  "file:///usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar"
		args       = [
		  "wordcount",
		  "file:///usr/lib/spark/NOTICE",
		  "gs://${google_dataproc_cluster.basic.bucket}/hadoopjob_output"
		]
    }
}
`, rnd)
}
