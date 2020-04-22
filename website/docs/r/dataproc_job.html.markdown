---
subcategory: "Dataproc"
layout: "google"
page_title: "Google: google_dataproc_job"
sidebar_current: "docs-google-dataproc-job"
description: |-
  Manages a job resource within a Dataproc cluster.
---

# google\_dataproc\_job

Manages a job resource within a Dataproc cluster within GCE. For more information see
[the official dataproc documentation](https://cloud.google.com/dataproc/).

!> **Note:** This resource does not support 'update' and changing any attributes will cause the resource to be recreated.

## Example usage

```hcl
resource "google_dataproc_cluster" "mycluster" {
  name   = "dproc-cluster-unique-name"
  region = "us-central1"
}

# Submit an example spark job to a dataproc cluster
resource "google_dataproc_job" "spark" {
  region       = google_dataproc_cluster.mycluster.region
  force_delete = true
  placement {
    cluster_name = google_dataproc_cluster.mycluster.name
  }

  spark_config {
    main_class    = "org.apache.spark.examples.SparkPi"
    jar_file_uris = ["file:///usr/lib/spark/examples/jars/spark-examples.jar"]
    args          = ["1000"]

    properties = {
      "spark.logConf" = "true"
    }

    logging_config {
      driver_log_levels = {
        "root" = "INFO"
      }
    }
  }
}

# Submit an example pyspark job to a dataproc cluster
resource "google_dataproc_job" "pyspark" {
  region       = google_dataproc_cluster.mycluster.region
  force_delete = true
  placement {
    cluster_name = google_dataproc_cluster.mycluster.name
  }

  pyspark_config {
    main_python_file_uri = "gs://dataproc-examples-2f10d78d114f6aaec76462e3c310f31f/src/pyspark/hello-world/hello-world.py"
    properties = {
      "spark.logConf" = "true"
    }
  }
}

# Check out current state of the jobs
output "spark_status" {
  value = google_dataproc_job.spark.status[0].state
}

output "pyspark_status" {
  value = google_dataproc_job.pyspark.status[0].state
}
```

## Argument Reference

* `placement.cluster_name` - (Required) The name of the cluster where the job 
   will be submitted.

* `xxx_config` - (Required) Exactly one of the specific job types to run on the
   cluster should be specified. If you want to submit multiple jobs, this will
   currently require the definition of multiple `google_dataproc_job` resources
   as shown in the example above, or by setting the `count` attribute.
   The following job configs are supported:

       * pyspark_config  - Submits a PySpark job to the cluster
       * spark_config    - Submits a Spark job to the cluster
       * hadoop_config   - Submits a Hadoop job to the cluster
       * hive_config     - Submits a Hive job to the cluster
       * hpig_config     - Submits a Pig job to the cluster
       * sparksql_config - Submits a Spark SQL job to the cluster

- - -

* `project` - (Optional) The project in which the `cluster` can be found and jobs
   subsequently run against. If it is not provided, the provider project is used.

* `region` - (Optional) The Cloud Dataproc region. This essentially determines which clusters are available
   for this job to be submitted to. If not specified, defaults to `global`.

* `force_delete` - (Optional) By default, you can only delete inactive jobs within
   Dataproc. Setting this to true, and calling destroy, will ensure that the
   job is first cancelled before issuing the delete.

* `labels` - (Optional) The list of labels (key/value pairs) to add to the job.

* `scheduling.max_failures_per_hour` - (Required) Maximum number of times per hour a driver may be restarted as a result of driver terminating with non-zero code before job is reported failed.

The `pyspark_config` block supports:

Submitting a pyspark job to the cluster. Below is an example configuration:

```hcl
# Submit a pyspark job to the cluster
resource "google_dataproc_job" "pyspark" {
  ...
  pyspark_config {
    main_python_file_uri = "gs://dataproc-examples-2f10d78d114f6aaec76462e3c310f31f/src/pyspark/hello-world/hello-world.py"
    properties = {
      "spark.logConf" = "true"
    }
  }
}
```

For configurations requiring Hadoop Compatible File System (HCFS) references, the options below
are generally applicable:

      - GCS files with the `gs://` prefix
      - HDFS files on the cluster with the `hdfs://` prefix
      - Local files on the cluster with the `file://` prefix

* `main_python_file_uri`- (Required) The HCFS URI of the main Python file to use as the driver. Must be a .py file.

* `args` - (Optional) The arguments to pass to the driver.

* `python_file_uris` - (Optional) HCFS file URIs of Python files to pass to the PySpark framework. Supported file types: .py, .egg, and .zip.

* `jar_file_uris` - (Optional) HCFS URIs of jar files to add to the CLASSPATHs of the Python driver and tasks.

* `file_uris` - (Optional) HCFS URIs of files to be copied to the working directory of Python drivers and distributed tasks. Useful for naively parallel tasks.

* `archive_uris` - (Optional) HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip.

* `properties` - (Optional) A mapping of property names to values, used to configure PySpark. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in `/etc/spark/conf/spark-defaults.conf` and classes in user code.

* `logging_config.driver_log_levels`- (Required) The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `spark_config` block supports:

```hcl
# Submit a spark job to the cluster
resource "google_dataproc_job" "spark" {
  ...
  spark_config {
    main_class    = "org.apache.spark.examples.SparkPi"
    jar_file_uris = ["file:///usr/lib/spark/examples/jars/spark-examples.jar"]
    args          = ["1000"]

    properties = {
      "spark.logConf" = "true"
    }

    logging_config {
      driver_log_levels = {
        "root" = "INFO"
      }
    }
  }
}
```

* `main_class`- (Optional) The class containing the main method of the driver. Must be in a
   provided jar or jar that is already on the classpath. Conflicts with `main_jar_file_uri`

* `main_jar_file_uri` - (Optional) The HCFS URI of jar file containing
   the driver jar. Conflicts with `main_class`

* `args` - (Optional) The arguments to pass to the driver.

* `jar_file_uris` - (Optional) HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.

* `file_uris` - (Optional) HCFS URIs of files to be copied to the working directory of Spark drivers and distributed tasks. Useful for naively parallel tasks.

* `archive_uris` - (Optional) HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip.

* `properties` - (Optional) A mapping of property names to values, used to configure Spark. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in `/etc/spark/conf/spark-defaults.conf` and classes in user code.

* `logging_config.driver_log_levels`- (Required) The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'


The `hadoop_config` block supports:

```hcl
# Submit a hadoop job to the cluster
resource "google_dataproc_job" "hadoop" {
  ...
  hadoop_config {
    main_jar_file_uri = "file:///usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar"
    args = [
      "wordcount",
      "file:///usr/lib/spark/NOTICE",
      "gs://${google_dataproc_cluster.basic.cluster_config[0].bucket}/hadoopjob_output",
    ]
  }
}
```

* `main_class`- (Optional) The name of the driver's main class. The jar file containing the class must be in the default CLASSPATH or specified in `jar_file_uris`. Conflicts with `main_jar_file_uri`

* `main_jar_file_uri` - (Optional) The HCFS URI of the jar file containing the main class. Examples: 'gs://foo-bucket/analytics-binaries/extract-useful-metrics-mr.jar' 'hdfs:/tmp/test-samples/custom-wordcount.jar' 'file:///home/usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar'. Conflicts with `main_class`

* `args` - (Optional) The arguments to pass to the driver. Do not include arguments, such as -libjars or -Dfoo=bar, that can be set as job properties, since a collision may occur that causes an incorrect job submission.

* `jar_file_uris` - (Optional) HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.

* `file_uris` - (Optional) HCFS URIs of files to be copied to the working directory of Hadoop drivers and distributed tasks. Useful for naively parallel tasks.

* `archive_uris` - (Optional) HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip.

* `properties` - (Optional) A mapping of property names to values, used to configure Hadoop. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in `/etc/hadoop/conf/*-site` and classes in user code..

* `logging_config.driver_log_levels`- (Required) The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `hive_config` block supports:

```hcl
# Submit a hive job to the cluster
resource "google_dataproc_job" "hive" {
  ...
  hive_config {
    query_list = [
      "DROP TABLE IF EXISTS dprocjob_test",
      "CREATE EXTERNAL TABLE dprocjob_test(bar int) LOCATION 'gs://${google_dataproc_cluster.basic.cluster_config[0].bucket}/hive_dprocjob_test/'",
      "SELECT * FROM dprocjob_test WHERE bar > 2",
    ]
  }
}
```

* `query_list`- (Optional) The list of Hive queries or statements to execute as part of the job.
   Conflicts with `query_file_uri`

* `query_file_uri` - (Optional) HCFS URI of file containing Hive script to execute as the job.
   Conflicts with `query_list`

* `continue_on_failure` - (Optional) Whether to continue executing queries if a query fails. The default value is false. Setting to true can be useful when executing independent parallel queries. Defaults to false.

* `script_variables` - (Optional) Mapping of query variable names to values (equivalent to the Hive command: `SET name="value";`).

* `properties` - (Optional)  A mapping of property names and values, used to configure Hive. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in `/etc/hadoop/conf/*-site.xml`, `/etc/hive/conf/hive-site.xml`, and classes in user code..

* `jar_file_uris` - (Optional) HCFS URIs of jar files to add to the CLASSPATH of the Hive server and Hadoop MapReduce (MR) tasks. Can contain Hive SerDes and UDFs.

The `pig_config` block supports:

```hcl
# Submit a pig job to the cluster
resource "google_dataproc_job" "pig" {
  ...
  pig_config {
    query_list = [
      "LNS = LOAD 'file:///usr/lib/pig/LICENSE.txt ' AS (line)",
      "WORDS = FOREACH LNS GENERATE FLATTEN(TOKENIZE(line)) AS word",
      "GROUPS = GROUP WORDS BY word",
      "WORD_COUNTS = FOREACH GROUPS GENERATE group, COUNT(WORDS)",
      "DUMP WORD_COUNTS",
    ]
  }
}
```

* `query_list`- (Optional) The list of Hive queries or statements to execute as part of the job.
   Conflicts with `query_file_uri`

* `query_file_uri` - (Optional) HCFS URI of file containing Hive script to execute as the job.
   Conflicts with `query_list`

* `continue_on_failure` - (Optional) Whether to continue executing queries if a query fails. The default value is false. Setting to true can be useful when executing independent parallel queries. Defaults to false.

* `script_variables` - (Optional) Mapping of query variable names to values (equivalent to the Pig command: `name=[value]`).

* `properties` - (Optional) A mapping of property names to values, used to configure Pig. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in `/etc/hadoop/conf/*-site.xml`, `/etc/pig/conf/pig.properties`, and classes in user code.

* `jar_file_uris` - (Optional) HCFS URIs of jar files to add to the CLASSPATH of the Pig Client and Hadoop MapReduce (MR) tasks. Can contain Pig UDFs.

* `logging_config.driver_log_levels`- (Required) The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'


The `sparksql_config` block supports:

```hcl
# Submit a spark SQL job to the cluster
resource "google_dataproc_job" "sparksql" {
  ...
  sparksql_config {
    query_list = [
      "DROP TABLE IF EXISTS dprocjob_test",
      "CREATE TABLE dprocjob_test(bar int)",
      "SELECT * FROM dprocjob_test WHERE bar > 2",
    ]
  }
}
```

* `query_list`- (Optional) The list of SQL queries or statements to execute as part of the job.
   Conflicts with `query_file_uri`

* `query_file_uri` - (Optional) The HCFS URI of the script that contains SQL queries.
   Conflicts with `query_list`

* `script_variables` - (Optional) Mapping of query variable names to values (equivalent to the Spark SQL command: `SET name="value";`).

* `properties` - (Optional) A mapping of property names to values, used to configure Spark SQL's SparkConf. Properties that conflict with values set by the Cloud Dataproc API may be overwritten.

* `jar_file_uris` - (Optional) HCFS URIs of jar files to be added to the Spark CLASSPATH.

* `logging_config.driver_log_levels`- (Required) The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `reference.0.cluster_uuid` - A cluster UUID generated by the Cloud Dataproc service when the job is submitted.

* `status.0.state` - A state message specifying the overall job state.

* `status.0.details` - Optional job state details, such as an error description if the state is ERROR.

* `status.0.state_start_time` - The time when this state was entered.

* `status.0.substate` - Additional state information, which includes status reported by the agent.

* `driver_output_resource_uri` - A URI pointing to the location of the stdout of the job's driver program.

* `driver_controls_files_uri` - If present, the location of miscellaneous control files which may be used as part of job setup and handling. If not present, control files may be placed in the same location as driver_output_uri.


## Timeouts

`google_dataproc_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for submitting a job to a dataproc cluster.
- `delete` - (Default `10 minutes`) Used for deleting a job from a dataproc cluster.
