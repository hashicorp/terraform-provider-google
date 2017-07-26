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

func TestAccDataprocJob_Hive(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_hive(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJobAttrMatch(
						"google_dataproc_job.hive", "hive_config"),
				),
			},
		},
	})
}

func TestAccDataprocJob_Pig(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_pig(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJobAttrMatch(
						"google_dataproc_job.pig", "pig_config"),
				),
			},
		},
	})
}

func TestAccDataprocJob_SparkSql(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocJob_sparksql(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocJobAttrMatch(
						"google_dataproc_job.sparksql", "sparksql_config"),
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

		jobTests := []jobTestField{

			{"cluster", job.Placement.ClusterName},
			{"labels", job.Labels},
		}

		if jobType == "pyspark_config" {
			jobTests = append(jobTests, jobTestField{"pyspark_config.0.main_python_file", job.PysparkJob.MainPythonFileUri})
			jobTests = append(jobTests, jobTestField{"pyspark_config.0.args", job.PysparkJob.Args})
			jobTests = append(jobTests, jobTestField{"pyspark_config.0.jars", job.PysparkJob.JarFileUris})
			jobTests = append(jobTests, jobTestField{"pyspark_config.0.files", job.PysparkJob.PythonFileUris})
			jobTests = append(jobTests, jobTestField{"pyspark_config.0.archives", job.PysparkJob.ArchiveUris})
			jobTests = append(jobTests, jobTestField{"pyspark_config.0.properties", job.PysparkJob.Properties})
		}
		if jobType == "spark_config" {
			jobTests = append(jobTests, jobTestField{"spark_config.0.main_class", job.SparkJob.MainClass})
			jobTests = append(jobTests, jobTestField{"spark_config.0.main_jar", job.SparkJob.MainJarFileUri})
			jobTests = append(jobTests, jobTestField{"spark_config.0.args", job.SparkJob.Args})
			jobTests = append(jobTests, jobTestField{"spark_config.0.jars", job.SparkJob.JarFileUris})
			jobTests = append(jobTests, jobTestField{"spark_config.0.files", job.SparkJob.FileUris})
			jobTests = append(jobTests, jobTestField{"spark_config.0.archives", job.SparkJob.ArchiveUris})
			jobTests = append(jobTests, jobTestField{"spark_config.0.properties", job.SparkJob.Properties})
		}
		if jobType == "hadoop_config" {
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.main_class", job.HadoopJob.MainClass})
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.main_jar", job.HadoopJob.MainJarFileUri})
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.args", job.HadoopJob.Args})
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.jars", job.HadoopJob.JarFileUris})
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.files", job.HadoopJob.FileUris})
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.archives", job.HadoopJob.ArchiveUris})
			jobTests = append(jobTests, jobTestField{"hadoop_config.0.properties", job.HadoopJob.Properties})
		}
		if jobType == "hive_config" {
			queries := []string{}
			if job.HiveJob.QueryList != nil {
				queries = job.HiveJob.QueryList.Queries
			}
			jobTests = append(jobTests, jobTestField{"hive_config.0.execution_queries", queries})
			jobTests = append(jobTests, jobTestField{"hive_config.0.execution_file", job.HiveJob.QueryFileUri})
			jobTests = append(jobTests, jobTestField{"hive_config.0.params", job.HiveJob.ScriptVariables})
			jobTests = append(jobTests, jobTestField{"hive_config.0.jars", job.HiveJob.JarFileUris})
			jobTests = append(jobTests, jobTestField{"hive_config.0.properties", job.HiveJob.Properties})
		}
		if jobType == "pig_config" {
			queries := []string{}
			if job.PigJob.QueryList != nil {
				queries = job.PigJob.QueryList.Queries
			}
			jobTests = append(jobTests, jobTestField{"pig_config.0.execution_queries", queries})
			jobTests = append(jobTests, jobTestField{"pig_config.0.execution_file", job.PigJob.QueryFileUri})
			jobTests = append(jobTests, jobTestField{"pig_config.0.params", job.PigJob.ScriptVariables})
			jobTests = append(jobTests, jobTestField{"pig_config.0.jars", job.PigJob.JarFileUris})
			jobTests = append(jobTests, jobTestField{"pig_config.0.properties", job.PigJob.Properties})
		}
		if jobType == "sparksql_config" {
			queries := []string{}
			if job.SparkSqlJob.QueryList != nil {
				queries = job.SparkSqlJob.QueryList.Queries
			}
			jobTests = append(jobTests, jobTestField{"sparksql_config.0.execution_queries", queries})
			jobTests = append(jobTests, jobTestField{"sparksql_config.0.execution_file", job.SparkSqlJob.QueryFileUri})
			jobTests = append(jobTests, jobTestField{"sparksql_config.0.params", job.SparkSqlJob.ScriptVariables})
			jobTests = append(jobTests, jobTestField{"sparksql_config.0.jars", job.SparkSqlJob.JarFileUris})
			jobTests = append(jobTests, jobTestField{"sparksql_config.0.properties", job.SparkSqlJob.Properties})
		}

		for _, attrs := range jobTests {
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

func testAccDataprocJob_hive(rnd string) string {
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

resource "google_dataproc_job" "hive" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    region       = "${google_dataproc_cluster.basic.region}"
    force_delete = true

    hive_config {
        execution_queries       = [
            "DROP TABLE IF EXISTS dprocjob_test",
            "CREATE EXTERNAL TABLE dprocjob_test(bar int) LOCATION 'gs://${google_dataproc_cluster.basic.bucket}/hive_dprocjob_test/'",
            "SELECT * FROM dprocjob_test WHERE bar > 2",
        ]
    }
}
`, rnd)
}

func testAccDataprocJob_pig(rnd string) string {
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

resource "google_dataproc_job" "pig" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    region       = "${google_dataproc_cluster.basic.region}"
    force_delete = true

    pig_config {
        execution_queries       = [
            "LNS = LOAD 'file:///usr/lib/pig/LICENSE.txt ' AS (line)",
            "WORDS = FOREACH LNS GENERATE FLATTEN(TOKENIZE(line)) AS word",
            "GROUPS = GROUP WORDS BY word",
            "WORD_COUNTS = FOREACH GROUPS GENERATE group, COUNT(WORDS)",
            "DUMP WORD_COUNTS"
        ]
    }
}
`, rnd)
}

func testAccDataprocJob_sparksql(rnd string) string {
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

resource "google_dataproc_job" "sparksql" {
    cluster      = "${google_dataproc_cluster.basic.name}"
    region       = "${google_dataproc_cluster.basic.region}"
    force_delete = true

    sparksql_config {
        execution_queries       = [
            "DROP TABLE IF EXISTS dprocjob_test",
            "CREATE TABLE dprocjob_test(bar int)",
            "SELECT * FROM dprocjob_test WHERE bar > 2",
        ]
    }
}
`, rnd)
}
