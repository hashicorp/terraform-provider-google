---
layout: "google"
page_title: "Google: google_dataproc_job"
sidebar_current: "docs-google-dataproc-job"
description: |-
  Manages a job resource within a Dataproc cluster.
---

# google\_dataproc\_job

Manages a job resource within a Dataproc cluster within GCE. For more information see
[the official dataproc documentation](https://cloud.google.com/dataproc/).

!> **Note:** This resource does not really support 'update' functionality. Once created
   (aka submitted to the cluster) there is not much point in changing anything. As a result
   changing any attributes will essentially cause the creation (submission) of a whole new job.

## Example usage

```hcl
resource "google_dataproc_cluster" "mycluster" {
    name   = "dproc-cluster-unique-name"
    region = "us-central1"
}

# Submit an example spark job to a dataproc cluster
resource "google_dataproc_job" "spark" {
    cluster      = "${google_dataproc_cluster.mycluster.name}"
    region       = "${google_dataproc_cluster.mycluster.region}"
    force_delete = true

    spark_config {
        main_class = "org.apache.spark.examples.SparkPi"
        jars       = ["file:///usr/lib/spark/examples/jars/spark-examples.jar"]
        args       = ["1000"]
    }
}

# Submit an example pyspark job to a dataproc cluster
resource "google_dataproc_job" "pyspark" {
    cluster      = "${google_dataproc_cluster.mycluster.name}"
    region       = "${google_dataproc_cluster.mycluster.region}"
    force_delete = true

    pyspark_config {
        main_python_file = "gs://dataproc-examples-2f10d78d114f6aaec76462e3c310f31f/src/pyspark/hello-world/hello-world.py"
    }
}

# Check out current state of running jobs
output "spark_status" {
    value = "${google_dataproc_job.spark.status}"
}

output "pyspark_status" {
    value = "${google_dataproc_job.pyspark.status}"
}
```

## Argument Reference

* `cluster` - (Required) The Dataproc cluster to submit the job to. Note: the list
   of available clusters to choose from is determined by the `region` value.


* `xxx_config` - (Required) Exactly one of the specific job types to run on the
   cluster should be specified. If you want to submit multiple jobs, this will
   currently require the definition of multiple `google_dataproc_job` resources
   as shown in the example above, or by setting the `count` attribute.
   The following job configs are supported:

       * pyspark_config - Submits a PySpark job to the cluster
       * spark_config   - Submits a Spark job to the cluster

   These job configs are not yet implemented:

       * hadoop
       * hive
       * pig
       * spark-sql

- - -

* `region` - (Optional) The Cloud Dataproc region. This essentially determines which clusters are available
   for this job to be submitted to. If not specified, defaults to `global`.

* `force_delete` - (Optional) By default, you can only delete inactive jobs within
   Dataproc. Setting this to true, and calling destroy, will ensure that the
   job is first cancelled before issuing the delete.

* `labels` - (Optional) The list of labels (key/value pairs) to add to the job.

The **pyspark_config** supports:

Submitting a pyspark job to the cluster. Below is an example configuration:

```hcl

# Submit a pyspark job to the cluster
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
```

For configurations requiring Hadoop Compatible File System (HCFS) references, the options below
are generally applicable:

      - GCS files with the `gs://` prefix
      - HDFS files on the cluster with the `hdfs://` prefix
      - Local files on the cluster with the `file://` prefix

* `main_python_file`- (Required) The HCFS URI of the main Python file (.py) to use as the driver.

* `files` - (Optional) A list of HCFS file URIs of Python files to pass to the
   PySpark framework. These are copied to the working directory of Python drivers
   and distributed tasks

* `jars` - (Optional) A list of HCFS jar files URIs to add to the
   CLASSPATHs of the Python driver and tasks.

* `archives` - (Optional) A list of HCFS archive URIs to be extracted in the
   working directory, typically of .jar, .tar, .tar.gz, .tgz, and .zip extentions.

* `args` - (Optional) The arguments to pass to the driver.

* `properties` - (Optional) A list of key value pairs to configure PySpark.


The **spark_config** supports:


```hcl

# Submit a spark job to the cluster
resource "google_dataproc_job" "pyspark" {
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
```

* `main_class`- (Optional) The class containing the main method of the driver. Must be in a
   provided jar or jar that is already on the classpath. Conflicts with `main_jar`

* `main_jar` - (Optional) The HCFS URI of jar file containing
   the driver jar. Conflicts with `main_class`

* `args` - (Optional) The arguments to pass to the main class.

* `jars` - (Optional) A list of HCFS jar files URIs to be provided to the executor and driver classpaths.

* `files` - (Optional) A list of HCFS files URIs to be provided to the job.

* `archives` - (Optional) A list of HCFS archive files URIs to to be provided to the job. must be one
   of the following file formats: .zip, .tar, .tar.gz, or .tgz.

* `properties` - (Optional) A list of key value pairs to configure Spark.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `status` - The current status of the job.

* `outputUri` - A URI pointing to the location of the stdout of the job's driver program.

<a id="timeouts"></a>
## Timeouts

`google_dataproc_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for submitting a job to a dataproc cluster.
- `delete` - (Default `10 minutes`) Used for deleting a job from a dataproc cluster.
