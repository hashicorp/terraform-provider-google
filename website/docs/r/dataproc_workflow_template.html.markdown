---
subcategory: "Dataproc"
description: |-
  A Workflow Template is a reusable workflow configuration.
---

# google_dataproc_workflow_template

A Workflow Template is a reusable workflow configuration. It defines a graph of jobs with information on where to run those jobs.

## Example Usage

```hcl
resource "google_dataproc_workflow_template" "template" {
  name = "template-example"
  location = "us-central1"
  placement {
    managed_cluster {
      cluster_name = "my-cluster"
      config {
        gce_cluster_config {
          zone = "us-central1-a"
          tags = ["foo", "bar"]
        }
        master_config {
          num_instances = 1
          machine_type = "n1-standard-1"
          disk_config {
            boot_disk_type = "pd-ssd"
            boot_disk_size_gb = 15
          }
        }
        worker_config {
          num_instances = 3
          machine_type = "n1-standard-2"
          disk_config {
            boot_disk_size_gb = 10
            num_local_ssds = 2
          }
        }

        secondary_worker_config {
          num_instances = 2
        }
        software_config {
          image_version = "2.0.35-debian10"
        }
      }
    }
  }
  jobs {
    step_id = "someJob"
    spark_job {
      main_class = "SomeClass"
    }
  }
  jobs {
    step_id = "otherJob"
    prerequisite_step_ids = ["someJob"]
    presto_job {
      query_file_uri = "someuri"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `jobs` -
  (Required)
  Required. The Directed Acyclic Graph of Jobs to submit.

* `location` -
  (Required)
  The location for the resource

* `name` -
  (Required)
  Output only. The resource name of the workflow template, as described in https://cloud.google.com/apis/design/resource_names. * For `projects.regions.workflowTemplates`, the resource name of the template has the following format: `projects/{project_id}/regions/{region}/workflowTemplates/{template_id}` * For `projects.locations.workflowTemplates`, the resource name of the template has the following format: `projects/{project_id}/locations/{location}/workflowTemplates/{template_id}`

* `placement` -
  (Required)
  Required. WorkflowTemplate scheduling information.

The `jobs` block supports:

* `hadoop_job` -
  (Optional)
  Job is a Hadoop job.

* `hive_job` -
  (Optional)
  Job is a Hive job.

* `labels` -
  (Optional)
  The labels to associate with this job. Label keys must be between 1 and 63 characters long, and must conform to the following regular expression: {0,63} No more than 32 labels can be associated with a given job.

* `pig_job` -
  (Optional)
  Job is a Pig job.

* `prerequisite_step_ids` -
  (Optional)
  The optional list of prerequisite job step_ids. If not specified, the job will start at the beginning of workflow.

* `presto_job` -
  (Optional)
  Job is a Presto job.

* `pyspark_job` -
  (Optional)
  Job is a PySpark job.

* `scheduling` -
  (Optional)
  Job scheduling configuration.

* `spark_job` -
  (Optional)
  Job is a Spark job.

* `spark_r_job` -
  (Optional)
  Job is a SparkR job.

* `spark_sql_job` -
  (Optional)
  Job is a SparkSql job.

* `step_id` -
  (Required)
  Required. The step id. The id must be unique among all jobs within the template. The step id is used as prefix for job id, as job `goog-dataproc-workflow-step-id` label, and in field from other steps. The id must contain only letters (a-z, A-Z), numbers (0-9), underscores (_), and hyphens (-). Cannot begin or end with underscore or hyphen. Must consist of between 3 and 50 characters.


The `placement` block supports:

* `cluster_selector` -
  (Optional)
  A selector that chooses target cluster for jobs based on metadata. The selector is evaluated at the time each job is submitted.

* `managed_cluster` -
  (Optional)
  A cluster that is managed by the workflow.


The `config` block supports:

* `autoscaling_config` -
  (Optional)
  Autoscaling config for the policy associated with the cluster. Cluster does not autoscale if this field is unset.

* `encryption_config` -
  (Optional)
  Encryption settings for the cluster.

* `endpoint_config` -
  (Optional)
  Port/endpoint configuration for this cluster

* `gce_cluster_config` -
  (Optional)
  The shared Compute Engine config settings for all instances in a cluster.

* `gke_cluster_config` -
  (Optional)
  The Kubernetes Engine config for Dataproc clusters deployed to Kubernetes. Setting this is considered mutually exclusive with Compute Engine-based options such as `gce_cluster_config`, `master_config`, `worker_config`, `secondary_worker_config`, and `autoscaling_config`.

* `initialization_actions` -
  (Optional)
  Commands to execute on each node after config is completed. By default, executables are run on master and all worker nodes. You can test a node's `role` metadata to run an executable on a master or worker node, as shown below using `curl` (you can also use `wget`): ROLE=$(curl -H Metadata-Flavor:Google http://metadata/computeMetadata/v1/instance/attributes/dataproc-role) if ; then ... master specific actions ... else ... worker specific actions ... fi

* `lifecycle_config` -
  (Optional)
  Lifecycle setting for the cluster.

* `master_config` -
  (Optional)
  The Compute Engine config settings for additional worker instances in a cluster.

* `metastore_config` -
  (Optional)
  Metastore configuration.

* `secondary_worker_config` -
  (Optional)
  The Compute Engine config settings for additional worker instances in a cluster.

* `security_config` -
  (Optional)
  Security settings for the cluster.

* `software_config` -
  (Optional)
  The config settings for software inside the cluster.

* `staging_bucket` -
  (Optional)
  A Cloud Storage bucket used to stage job dependencies, config files, and job driver console output. If you do not specify a staging bucket, Cloud Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's staging bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket (see (https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/staging-bucket)).

* `temp_bucket` -
  (Optional)
  A Cloud Storage bucket used to store ephemeral cluster and jobs data, such as Spark and MapReduce history files. If you do not specify a temp bucket, Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's temp bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket. The default bucket has a TTL of 90 days, but you can use any TTL (or none) if you specify a bucket.

* `worker_config` -
  (Optional)
  The Compute Engine config settings for additional worker instances in a cluster.

- - -

* `dag_timeout` -
  (Optional)
  (Beta only) Optional. Timeout duration for the DAG of jobs. You can use "s", "m", "h", and "d" suffixes for second, minute, hour, and day duration values, respectively. The timeout duration must be from 10 minutes ("10m") to 24 hours ("24h" or "1d"). The timer begins when the first job is submitted. If the workflow is running at the end of the timeout period, any remaining jobs are cancelled, the workflow is ended, and if the workflow was running on a (/dataproc/docs/concepts/workflows/using-workflows#configuring_or_selecting_a_cluster), the cluster is deleted.

* `labels` -
  (Optional)
  The labels to associate with this template. These labels will be propagated to all jobs and clusters created by the workflow instance. Label **keys** must contain 1 to 63 characters, and must conform to (https://www.ietf.org/rfc/rfc1035.txt). No more than 32 labels can be associated with a template.

* `parameters` -
  (Optional)
  Template parameters whose values are substituted into the template. Values for parameters must be provided when the template is instantiated.

* `project` -
  (Optional)
  The project for the resource

* `version` -
  (Optional)
  Used to perform a consistent read-modify-write. This field should be left blank for a `CreateWorkflowTemplate` request. It is required for an `UpdateWorkflowTemplate` request, and must match the current server version. A typical update template flow would fetch the current template with a `GetWorkflowTemplate` request, which will return the current template with the `version` field filled in with the current server version. The user updates other fields in the template, then returns it as part of the `UpdateWorkflowTemplate` request.



The `hadoop_job` block supports:

* `archive_uris` -
  (Optional)
  HCFS URIs of archives to be extracted in the working directory of Hadoop drivers and tasks. Supported file types: .jar, .tar, .tar.gz, .tgz, or .zip.

* `args` -
  (Optional)
  The arguments to pass to the driver. Do not include arguments, such as `-libjars` or `-Dfoo=bar`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.

* `file_uris` -
  (Optional)
  HCFS (Hadoop Compatible Filesystem) URIs of files to be copied to the working directory of Hadoop drivers and distributed tasks. Useful for naively parallel tasks.

* `jar_file_uris` -
  (Optional)
  Jar file URIs to add to the CLASSPATHs of the Hadoop driver and tasks.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `main_class` -
  (Optional)
  The name of the driver's main class. The jar file containing the class must be in the default CLASSPATH or specified in `jar_file_uris`.

* `main_jar_file_uri` -
  (Optional)
  The HCFS URI of the jar file containing the main class. Examples: 'gs://foo-bucket/analytics-binaries/extract-useful-metrics-mr.jar' 'hdfs:/tmp/test-samples/custom-wordcount.jar' 'file:///home/usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar'

* `properties` -
  (Optional)
  A mapping of property names to values, used to configure Hadoop. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site and classes in user code.

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `hive_job` block supports:

* `continue_on_failure` -
  (Optional)
  Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.

* `jar_file_uris` -
  (Optional)
  HCFS URIs of jar files to add to the CLASSPATH of the Hive server and Hadoop MapReduce (MR) tasks. Can contain Hive SerDes and UDFs.

* `properties` -
  (Optional)
  A mapping of property names and values, used to configure Hive. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/hive/conf/hive-site.xml, and classes in user code.

* `query_file_uri` -
  (Optional)
  The HCFS URI of the script that contains Hive queries.

* `query_list` -
  (Optional)
  A list of queries.

* `script_variables` -
  (Optional)
  Mapping of query variable names to values (equivalent to the Hive command: `SET name="value";`).

The `query_list` block supports:

* `queries` -
  (Required)
  Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: "hiveJob": { "queryList": { "queries": } }

The `pig_job` block supports:

* `continue_on_failure` -
  (Optional)
  Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.

* `jar_file_uris` -
  (Optional)
  HCFS URIs of jar files to add to the CLASSPATH of the Pig Client and Hadoop MapReduce (MR) tasks. Can contain Pig UDFs.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `properties` -
  (Optional)
  A mapping of property names to values, used to configure Pig. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/pig/conf/pig.properties, and classes in user code.

* `query_file_uri` -
  (Optional)
  The HCFS URI of the script that contains the Pig queries.

* `query_list` -
  (Optional)
  A list of queries.

* `script_variables` -
  (Optional)
  Mapping of query variable names to values (equivalent to the Pig command: `name=`).

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `query_list` block supports:

* `queries` -
  (Required)
  Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: "hiveJob": { "queryList": { "queries": } }

The `presto_job` block supports:

* `client_tags` -
  (Optional)
  Presto client tags to attach to this query

* `continue_on_failure` -
  (Optional)
  Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `output_format` -
  (Optional)
  The format in which query output will be displayed. See the Presto documentation for supported output formats

* `properties` -
  (Optional)
  A mapping of property names to values. Used to set Presto (https://prestodb.io/docs/current/sql/set-session.html) Equivalent to using the --session flag in the Presto CLI

* `query_file_uri` -
  (Optional)
  The HCFS URI of the script that contains SQL queries.

* `query_list` -
  (Optional)
  A list of queries.

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `query_list` block supports:

* `queries` -
  (Required)
  Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: "hiveJob": { "queryList": { "queries": } }

The `pyspark_job` block supports:

* `archive_uris` -
  (Optional)
  HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.

* `args` -
  (Optional)
  The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.

* `file_uris` -
  (Optional)
  HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.

* `jar_file_uris` -
  (Optional)
  HCFS URIs of jar files to add to the CLASSPATHs of the Python driver and tasks.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `main_python_file_uri` -
  (Required)
  Required. The HCFS URI of the main Python file to use as the driver. Must be a .py file.

* `properties` -
  (Optional)
  A mapping of property names to values, used to configure PySpark. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.

* `python_file_uris` -
  (Optional)
  HCFS file URIs of Python files to pass to the PySpark framework. Supported file types: .py, .egg, and .zip.

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `scheduling` block supports:

* `max_failures_per_hour` -
  (Optional)
  Maximum number of times per hour a driver may be restarted as a result of driver exiting with non-zero code before job is reported failed. A job may be reported as thrashing if driver exits with non-zero code 4 times within 10 minute window. Maximum value is 10.

* `max_failures_total` -
  (Optional)
  Maximum number of times in total a driver may be restarted as a result of driver exiting with non-zero code before job is reported failed. Maximum value is 240

The `spark_job` block supports:

* `archive_uris` -
  (Optional)
  HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.

* `args` -
  (Optional)
  The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.

* `file_uris` -
  (Optional)
  HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.

* `jar_file_uris` -
  (Optional)
  HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `main_class` -
  (Optional)
  The name of the driver's main class. The jar file that contains the class must be in the default CLASSPATH or specified in `jar_file_uris`.

* `main_jar_file_uri` -
  (Optional)
  The HCFS URI of the jar file that contains the main class.

* `properties` -
  (Optional)
  A mapping of property names to values, used to configure Spark. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `spark_r_job` block supports:

* `archive_uris` -
  (Optional)
  HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.

* `args` -
  (Optional)
  The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.

* `file_uris` -
  (Optional)
  HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `main_r_file_uri` -
  (Required)
  Required. The HCFS URI of the main R file to use as the driver. Must be a .R file.

* `properties` -
  (Optional)
  A mapping of property names to values, used to configure SparkR. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `spark_sql_job` block supports:

* `jar_file_uris` -
  (Optional)
  HCFS URIs of jar files to be added to the Spark CLASSPATH.

* `logging_config` -
  (Optional)
  The runtime log config for job execution.

* `properties` -
  (Optional)
  A mapping of property names to values, used to configure Spark SQL's SparkConf. Properties that conflict with values set by the Dataproc API may be overwritten.

* `query_file_uri` -
  (Optional)
  The HCFS URI of the script that contains SQL queries.

* `query_list` -
  (Optional)
  A list of queries.

* `script_variables` -
  (Optional)
  Mapping of query variable names to values (equivalent to the Spark SQL command: SET `name="value";`).

The `logging_config` block supports:

* `driver_log_levels` -
  (Optional)
  The per-package log levels for the driver. This may include "root" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'

The `query_list` block supports:

* `queries` -
  (Required)
  Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: "hiveJob": { "queryList": { "queries": } }

The `parameters` block supports:

* `description` -
  (Optional)
  Brief description of the parameter. Must not exceed 1024 characters.

* `fields` -
  (Required)
  Required. Paths to all fields that the parameter replaces. A field is allowed to appear in at most one parameter's list of field paths. A field path is similar in syntax to a .sparkJob.args

* `name` -
  (Required)
  Required. Parameter name. The parameter name is used as the key, and paired with the parameter value, which are passed to the template when the template is instantiated. The name must contain only capital letters (A-Z), numbers (0-9), and underscores (_), and must not start with a number. The maximum length is 40 characters.

* `validation` -
  (Optional)
  Validation rules to be applied to this parameter's value.

The `validation` block supports:

* `regex` -
  (Optional)
  Validation based on regular expressions.

* `values` -
  (Optional)
  Validation based on a list of allowed values.

The `regex` block supports:

* `regexes` -
  (Required)
  Required. RE2 regular expressions used to validate the parameter's value. The value must match the regex in its entirety (substring matches are not sufficient).

The `values` block supports:

* `values` -
  (Required)
  Required. List of allowed values for the parameter.

The `cluster_selector` block supports:

* `cluster_labels` -
  (Required)
  Required. The cluster labels. Cluster must have all labels to match.

* `zone` -
  (Optional)
  The zone where workflow process executes. This parameter does not affect the selection of the cluster. If unspecified, the zone of the first cluster matching the selector is used.

The `managed_cluster` block supports:

* `cluster_name` -
  (Required)
  Required. The cluster name prefix. A unique cluster name will be formed by appending a random suffix. The name must contain only lower-case letters (a-z), numbers (0-9), and hyphens (-). Must begin with a letter. Cannot begin or end with hyphen. Must consist of between 2 and 35 characters.

* `config` -
  (Required)
  Required. The cluster configuration.

* `labels` -
  (Optional)
  The labels to associate with this cluster. Label keys must be between 1 and 63 characters long, and must conform to the following PCRE regular expression: {0,63} No more than 32 labels can be associated with a given cluster.

The `master_config` block supports:

* `accelerators` -
  (Optional)
  The Compute Engine accelerator configuration for these instances.

* `disk_config` -
  (Optional)
  Disk option config settings.

* `image` -
  (Optional)
  The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.

* `machine_type` -
  (Optional)
  The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/(https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.

* `min_cpu_platform` -
  (Optional)
  Specifies the minimum cpu platform for the Instance Group. See (https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).

* `num_instances` -
  (Optional)
  The number of VM instances in the instance group. For master instance groups, must be set to 1.

* `preemptibility` -
  (Optional)
  Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE

* `instance_names` -
  Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.

* `is_preemptible` -
  Output only. Specifies that this instance group contains preemptible instances.

* `managed_group_config` -
  Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.

The `accelerators` block supports:

* `accelerator_count` -
  (Optional)
  The number of the accelerator cards of this type exposed to this instance.

* `accelerator_type` -
  (Optional)
  Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See (https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.

The `disk_config` block supports:

* `boot_disk_size_gb` -
  (Optional)
  Size in GB of the boot disk (default is 500GB).

* `boot_disk_type` -
  (Optional)
  Type of the boot disk (default is "pd-standard"). Valid values: "pd-ssd" (Persistent Disk Solid State Drive) or "pd-standard" (Persistent Disk Hard Disk Drive).

* `num_local_ssds` -
  (Optional)
  Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and (https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.

The `autoscaling_config` block supports:

* `policy` -
  (Optional)
  The autoscaling policy used by the cluster. Only resource names including projectid and location (region) are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/` Note that the policy must be in the same project and Dataproc region.

The `encryption_config` block supports:

* `gce_pd_kms_key_name` -
  (Optional)
  The Cloud KMS key name to use for PD disk encryption for all instances in the cluster.

The `endpoint_config` block supports:

* `enable_http_port_access` -
  (Optional)
  If true, enable http access to specific ports on the cluster from external sources. Defaults to false.

* `http_ports` -
  Output only. The map of port descriptions to URLs. Will only be populated if enable_http_port_access is true.

The `gce_cluster_config` block supports:

* `internal_ip_only` -
  (Optional)
  If true, all instances in the cluster will only have internal IP addresses. By default, clusters are not restricted to internal IP addresses, and will have ephemeral external IP addresses assigned to each instance. This `internal_ip_only` restriction can only be enabled for subnetwork enabled networks, and all off-cluster dependencies must be configured to be accessible without external IP addresses.

* `metadata` -
  (Optional)
  The Compute Engine metadata entries to add to all instances (see (https://cloud.google.com/compute/docs/storing-retrieving-metadata#project_and_instance_metadata)).

* `network` -
  (Optional)
  The Compute Engine network to be used for machine communications. Cannot be specified with subnetwork_uri. If neither `network_uri` nor `subnetwork_uri` is specified, the "default" network of the project is used, if it exists. Cannot be a "Custom Subnet Network" (see /regions/global/default` * `default`

* `node_group_affinity` -
  (Optional)
  Node Group Affinity for sole-tenant clusters.

* `private_ipv6_google_access` -
  (Optional)
  The type of IPv6 access for a cluster. Possible values: PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED, INHERIT_FROM_SUBNETWORK, OUTBOUND, BIDIRECTIONAL

* `reservation_affinity` -
  (Optional)
  Reservation Affinity for consuming Zonal reservation.

* `service_account` -
  (Optional)
  The (https://cloud.google.com/compute/docs/access/service-accounts#default_service_account) is used.

* `service_account_scopes` -
  (Optional)
  The URIs of service account scopes to be included in Compute Engine instances. The following base set of scopes is always included: * https://www.googleapis.com/auth/cloud.useraccounts.readonly * https://www.googleapis.com/auth/devstorage.read_write * https://www.googleapis.com/auth/logging.write If no scopes are specified, the following defaults are also provided: * https://www.googleapis.com/auth/bigquery * https://www.googleapis.com/auth/bigtable.admin.table * https://www.googleapis.com/auth/bigtable.data * https://www.googleapis.com/auth/devstorage.full_control

* `shielded_instance_config` -
  (Optional)
  Shielded Instance Config for clusters using [Compute Engine Shielded VMs](https://cloud.google.com/security/shielded-cloud/shielded-vm). Structure [defined below](#nested_shielded_instance_config).

* `subnetwork` -
  (Optional)
  The Compute Engine subnetwork to be used for machine communications. Cannot be specified with network_uri. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects//regions/us-east1/subnetworks/sub0` * `sub0`

* `tags` -
  (Optional)
  The Compute Engine tags to add to all instances (see (https://cloud.google.com/compute/docs/label-or-tag-resources#tags)).

* `zone` -
  (Optional)
  The zone where the Compute Engine cluster will be located. On a create request, it is required in the "global" region. If omitted in a non-global Dataproc region, the service will pick a zone in the corresponding Compute Engine region. On a get request, zone will always be present. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/` * `us-central1-f`

The `node_group_affinity` block supports:

* `node_group` -
  (Required)
  Required. The URI of a sole-tenant /zones/us-central1-a/nodeGroups/node-group-1` * `node-group-1`

The `reservation_affinity` block supports:

* `consume_reservation_type` -
  (Optional)
  Type of reservation to consume Possible values: TYPE_UNSPECIFIED, NO_RESERVATION, ANY_RESERVATION, SPECIFIC_RESERVATION

* `key` -
  (Optional)
  Corresponds to the label key of reservation resource.

* `values` -
  (Optional)
  Corresponds to the label values of reservation resource.

<a name="nested_shielded_instance_config"></a>The `shielded_instance_config` block supports:

```hcl
cluster_config {
  gce_cluster_config {
    shielded_instance_config {
      enable_secure_boot          = true
      enable_vtpm                 = true
      enable_integrity_monitoring = true
    }
  }
}
```

* `enable_secure_boot` -
  (Optional)
  Defines whether instances have [Secure Boot](https://cloud.google.com/compute/shielded-vm/docs/shielded-vm#secure-boot) enabled.

* `enable_vtpm` -
  (Optional)
  Defines whether instances have the [vTPM](https://cloud.google.com/compute/shielded-vm/docs/shielded-vm#vtpm) enabled.

* `enable_integrity_monitoring` -
  (Optional)
  Defines whether instances have [Integrity Monitoring](https://cloud.google.com/compute/shielded-vm/docs/shielded-vm#integrity-monitoring) enabled.

The `gke_cluster_config` block supports:

* `namespaced_gke_deployment_target` -
  (Optional)
  A target for the deployment.

The `namespaced_gke_deployment_target` block supports:

* `cluster_namespace` -
  (Optional)
  A namespace within the GKE cluster to deploy into.

* `target_gke_cluster` -
  (Optional)
  The target GKE cluster to deploy to. Format: 'projects/{project}/locations/{location}/clusters/{cluster_id}'

The `initialization_actions` block supports:

* `executable_file` -
  (Optional)
  Required. Cloud Storage URI of executable file.

* `execution_timeout` -
  (Optional)
  Amount of time executable has to complete. Default is 10 minutes (see JSON representation of (https://developers.google.com/protocol-buffers/docs/proto3#json)). Cluster creation fails with an explanatory error message (the name of the executable that caused the error and the exceeded timeout period) if the executable is not completed at end of the timeout period.

The `lifecycle_config` block supports:

* `auto_delete_time` -
  (Optional)
  The time when cluster will be auto-deleted (see JSON representation of (https://developers.google.com/protocol-buffers/docs/proto3#json)).

* `auto_delete_ttl` -
  (Optional)
  The lifetime duration of cluster. The cluster will be auto-deleted at the end of this period. Minimum value is 10 minutes; maximum value is 14 days (see JSON representation of (https://developers.google.com/protocol-buffers/docs/proto3#json)).

* `idle_delete_ttl` -
  (Optional)
  The duration to keep the cluster alive while idling (when no jobs are running). Passing this threshold will cause the cluster to be deleted. Minimum value is 5 minutes; maximum value is 14 days (see JSON representation of (https://developers.google.com/protocol-buffers/docs/proto3#json).

* `idle_start_time` -
  Output only. The time when cluster became idle (most recent job finished) and became eligible for deletion due to idleness (see JSON representation of (https://developers.google.com/protocol-buffers/docs/proto3#json)).

The `metastore_config` block supports:

* `dataproc_metastore_service` -
  (Required)
  Required. Resource name of an existing Dataproc Metastore service. Example: * `projects/`

The `security_config` block supports:

* `kerberos_config` -
  (Optional)
  Kerberos related configuration.

The `kerberos_config` block supports:

* `cross_realm_trust_admin_server` -
  (Optional)
  The admin server (IP or hostname) for the remote trusted realm in a cross realm trust relationship.

* `cross_realm_trust_kdc` -
  (Optional)
  The KDC (IP or hostname) for the remote trusted realm in a cross realm trust relationship.

* `cross_realm_trust_realm` -
  (Optional)
  The remote realm the Dataproc on-cluster KDC will trust, should the user enable cross realm trust.

* `cross_realm_trust_shared_password` -
  (Optional)
  The Cloud Storage URI of a KMS encrypted file containing the shared password between the on-cluster Kerberos realm and the remote trusted realm, in a cross realm trust relationship.

* `enable_kerberos` -
  (Optional)
  Flag to indicate whether to Kerberize the cluster (default: false). Set this field to true to enable Kerberos on a cluster.

* `kdc_db_key` -
  (Optional)
  The Cloud Storage URI of a KMS encrypted file containing the master key of the KDC database.

* `key_password` -
  (Optional)
  The Cloud Storage URI of a KMS encrypted file containing the password to the user provided key. For the self-signed certificate, this password is generated by Dataproc.

* `keystore` -
  (Optional)
  The Cloud Storage URI of the keystore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.

* `keystore_password` -
  (Optional)
  The Cloud Storage URI of a KMS encrypted file containing the password to the user provided keystore. For the self-signed certificate, this password is generated by Dataproc.

* `kms_key` -
  (Optional)
  The uri of the KMS key used to encrypt various sensitive files.

* `realm` -
  (Optional)
  The name of the on-cluster Kerberos realm. If not specified, the uppercased domain of hostnames will be the realm.

* `root_principal_password` -
  (Optional)
  The Cloud Storage URI of a KMS encrypted file containing the root principal password.

* `tgt_lifetime_hours` -
  (Optional)
  The lifetime of the ticket granting ticket, in hours. If not specified, or user specifies 0, then default value 10 will be used.

* `truststore` -
  (Optional)
  The Cloud Storage URI of the truststore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.

* `truststore_password` -
  (Optional)
  The Cloud Storage URI of a KMS encrypted file containing the password to the user provided truststore. For the self-signed certificate, this password is generated by Dataproc.

The `software_config` block supports:

* `image_version` -
  (Optional)
  The version of software inside the cluster. It must be one of the supported [Dataproc Versions](https://cloud.google.com/dataproc/docs/concepts/versioning/dataproc-versions#supported_dataproc_versions), such as "1.2" (including a subminor version, such as "1.2.29"), or the ["preview" version](https://cloud.google.com/dataproc/docs/concepts/versioning/dataproc-versions#other_versions). If unspecified, it defaults to the latest Debian version.

* `optional_components` -
  (Optional)
  The set of components to activate on the cluster.

* `properties` -
  (Optional)
  The properties to set on daemon config files.

  Property keys are specified in `prefix:property` format, for example `core:hadoop.tmp.dir`. The following are supported prefixes and their mappings:

  * capacity-scheduler: `capacity-scheduler.xml`
  * core: `core-site.xml`
  * distcp: `distcp-default.xml`
  * hdfs: `hdfs-site.xml`
  * hive: `hive-site.xml`
  * mapred: `mapred-site.xml`
  * pig: `pig.properties`
  * spark: `spark-defaults.conf`
  * yarn: `yarn-site.xml`

  
  For more information, see [Cluster properties](https://cloud.google.com/dataproc/docs/concepts/cluster-properties).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}`

* `create_time` -
  Output only. The time template was created.

* `update_time` -
  Output only. The time template was last updated.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

WorkflowTemplate can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}`
* `{{project}}/{{location}}/{{name}}`
* `{{location}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import WorkflowTemplate using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}"
  to = google_dataproc_workflow_template.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), WorkflowTemplate can be imported using one of the formats above. For example:

```
$ terraform import google_dataproc_workflow_template.default projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}
$ terraform import google_dataproc_workflow_template.default {{project}}/{{location}}/{{name}}
$ terraform import google_dataproc_workflow_template.default {{location}}/{{name}}
```



