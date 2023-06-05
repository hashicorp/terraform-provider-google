// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package dataproc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dataproc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceDataprocWorkflowTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocWorkflowTemplateCreate,
		Read:   resourceDataprocWorkflowTemplateRead,
		Delete: resourceDataprocWorkflowTemplateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDataprocWorkflowTemplateImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"jobs": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The Directed Acyclic Graph of Jobs to submit.",
				Elem:        DataprocWorkflowTemplateJobsSchema(),
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Output only. The resource name of the workflow template, as described in https://cloud.google.com/apis/design/resource_names. * For `projects.regions.workflowTemplates`, the resource name of the template has the following format: `projects/{project_id}/regions/{region}/workflowTemplates/{template_id}` * For `projects.locations.workflowTemplates`, the resource name of the template has the following format: `projects/{project_id}/locations/{location}/workflowTemplates/{template_id}`",
			},

			"placement": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. WorkflowTemplate scheduling information.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementSchema(),
			},

			"dag_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Timeout duration for the DAG of jobs, expressed in seconds (see [JSON representation of duration](https://developers.google.com/protocol-buffers/docs/proto3#json)). The timeout duration must be from 10 minutes (\"600s\") to 24 hours (\"86400s\"). The timer begins when the first job is submitted. If the workflow is running at the end of the timeout period, any remaining jobs are cancelled, the workflow is ended, and if the workflow was running on a [managed cluster](/dataproc/docs/concepts/workflows/using-workflows#configuring_or_selecting_a_cluster), the cluster is deleted.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The labels to associate with this template. These labels will be propagated to all jobs and clusters created by the workflow instance. Label **keys** must contain 1 to 63 characters, and must conform to [RFC 1035](https://www.ietf.org/rfc/rfc1035.txt). Label **values** may be empty, but, if present, must contain 1 to 63 characters, and must conform to [RFC 1035](https://www.ietf.org/rfc/rfc1035.txt). No more than 32 labels can be associated with a template.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Template parameters whose values are substituted into the template. Values for parameters must be provided when the template is instantiated.",
				Elem:        DataprocWorkflowTemplateParametersSchema(),
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Output only. The current version of this workflow template.",
				Deprecated:  "version is not useful as a configurable field, and will be removed in the future.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time template was created.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time template was last updated.",
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"step_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The step id. The id must be unique among all jobs within the template. The step id is used as prefix for job id, as job `goog-dataproc-workflow-step-id` label, and in prerequisiteStepIds field from other steps. The id must contain only letters (a-z, A-Z), numbers (0-9), underscores (_), and hyphens (-). Cannot begin or end with underscore or hyphen. Must consist of between 3 and 50 characters.",
			},

			"hadoop_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a Hadoop job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsHadoopJobSchema(),
			},

			"hive_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a Hive job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsHiveJobSchema(),
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The labels to associate with this job. Label keys must be between 1 and 63 characters long, and must conform to the following regular expression: p{Ll}p{Lo}{0,62} Label values must be between 1 and 63 characters long, and must conform to the following regular expression: [p{Ll}p{Lo}p{N}_-]{0,63} No more than 32 labels can be associated with a given job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"pig_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a Pig job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPigJobSchema(),
			},

			"prerequisite_step_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The optional list of prerequisite job step_ids. If not specified, the job will start at the beginning of workflow.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"presto_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a Presto job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPrestoJobSchema(),
			},

			"pyspark_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a PySpark job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPysparkJobSchema(),
			},

			"scheduling": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job scheduling configuration.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSchedulingSchema(),
			},

			"spark_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a Spark job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkJobSchema(),
			},

			"spark_r_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a SparkR job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkRJobSchema(),
			},

			"spark_sql_job": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Job is a SparkSql job.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkSqlJobSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplateJobsHadoopJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"archive_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of archives to be extracted in the working directory of Hadoop drivers and tasks. Supported file types: .jar, .tar, .tar.gz, .tgz, or .zip.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `-libjars` or `-Dfoo=bar`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS (Hadoop Compatible Filesystem) URIs of files to be copied to the working directory of Hadoop drivers and distributed tasks. Useful for naively parallel tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Jar file URIs to add to the CLASSPATHs of the Hadoop driver and tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsHadoopJobLoggingConfigSchema(),
			},

			"main_class": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the driver's main class. The jar file containing the class must be in the default CLASSPATH or specified in `jar_file_uris`.",
			},

			"main_jar_file_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCFS URI of the jar file containing the main class. Examples: 'gs://foo-bucket/analytics-binaries/extract-useful-metrics-mr.jar' 'hdfs:/tmp/test-samples/custom-wordcount.jar' 'file:///home/usr/lib/hadoop-mapreduce/hadoop-mapreduce-examples.jar'",
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values, used to configure Hadoop. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site and classes in user code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsHadoopJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsHiveJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"continue_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.",
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of jar files to add to the CLASSPATH of the Hive server and Hadoop MapReduce (MR) tasks. Can contain Hive SerDes and UDFs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names and values, used to configure Hive. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/hive/conf/hive-site.xml, and classes in user code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"query_file_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCFS URI of the script that contains Hive queries.",
			},

			"query_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "A list of queries.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsHiveJobQueryListSchema(),
			},

			"script_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Mapping of query variable names to values (equivalent to the Hive command: `SET name=\"value\";`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsHiveJobQueryListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queries": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPigJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"continue_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.",
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of jar files to add to the CLASSPATH of the Pig Client and Hadoop MapReduce (MR) tasks. Can contain Pig UDFs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPigJobLoggingConfigSchema(),
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values, used to configure Pig. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/pig/conf/pig.properties, and classes in user code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"query_file_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCFS URI of the script that contains the Pig queries.",
			},

			"query_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "A list of queries.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPigJobQueryListSchema(),
			},

			"script_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Mapping of query variable names to values (equivalent to the Pig command: `name=[value]`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPigJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPigJobQueryListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queries": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPrestoJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Presto client tags to attach to this query",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"continue_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Whether to continue executing queries if a query fails. The default value is `false`. Setting to `true` can be useful when executing independent parallel queries.",
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPrestoJobLoggingConfigSchema(),
			},

			"output_format": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The format in which query output will be displayed. See the Presto documentation for supported output formats",
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values. Used to set Presto [session properties](https://prestodb.io/docs/current/sql/set-session.html) Equivalent to using the --session flag in the Presto CLI",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"query_file_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCFS URI of the script that contains SQL queries.",
			},

			"query_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "A list of queries.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPrestoJobQueryListSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPrestoJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPrestoJobQueryListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queries": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPysparkJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"main_python_file_uri": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The HCFS URI of the main Python file to use as the driver. Must be a .py file.",
			},

			"archive_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of jar files to add to the CLASSPATHs of the Python driver and tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsPysparkJobLoggingConfigSchema(),
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values, used to configure PySpark. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"python_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS file URIs of Python files to pass to the PySpark framework. Supported file types: .py, .egg, and .zip.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsPysparkJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSchedulingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_failures_per_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Maximum number of times per hour a driver may be restarted as a result of driver exiting with non-zero code before job is reported failed. A job may be reported as thrashing if driver exits with non-zero code 4 times within 10 minute window. Maximum value is 10.",
			},

			"max_failures_total": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Maximum number of times in total a driver may be restarted as a result of driver exiting with non-zero code before job is reported failed. Maximum value is 240.",
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"archive_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkJobLoggingConfigSchema(),
			},

			"main_class": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the driver's main class. The jar file that contains the class must be in the default CLASSPATH or specified in `jar_file_uris`.",
			},

			"main_jar_file_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCFS URI of the jar file that contains the main class.",
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values, used to configure Spark. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkRJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"main_r_file_uri": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The HCFS URI of the main R file to use as the driver. Must be a .R file.",
			},

			"archive_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of archives to be extracted into the working directory of each executor. Supported file types: .jar, .tar, .tar.gz, .tgz, and .zip.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as `--conf`, that can be set as job properties, since a collision may occur that causes an incorrect job submission.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of files to be placed in the working directory of each executor. Useful for naively parallel tasks.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkRJobLoggingConfigSchema(),
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values, used to configure SparkR. Properties that conflict with values set by the Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkRJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkSqlJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. HCFS URIs of jar files to be added to the Spark CLASSPATH.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The runtime log config for job execution.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkSqlJobLoggingConfigSchema(),
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A mapping of property names to values, used to configure Spark SQL's SparkConf. Properties that conflict with values set by the Dataproc API may be overwritten.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"query_file_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCFS URI of the script that contains SQL queries.",
			},

			"query_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "A list of queries.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateJobsSparkSqlJobQueryListSchema(),
			},

			"script_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Mapping of query variable names to values (equivalent to the Spark SQL command: SET `name=\"value\";`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkSqlJobLoggingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The per-package log levels for the driver. This may include \"root\" package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateJobsSparkSqlJobQueryListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queries": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The queries to execute. You do not need to end a query expression with a semicolon. Multiple queries can be specified in one string by separating each with a semicolon. Here is an example of a Dataproc API snippet that uses a QueryList to specify a HiveJob: \"hiveJob\": { \"queryList\": { \"queries\": [ \"query1\", \"query2\", \"query3;query4\", ] } }",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_selector": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. A selector that chooses target cluster for jobs based on metadata. The selector is evaluated at the time each job is submitted.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementClusterSelectorSchema(),
			},

			"managed_cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "A cluster that is managed by the workflow.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementClusterSelectorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_labels": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The cluster labels. Cluster must have all labels to match.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The zone where workflow process executes. This parameter does not affect the selection of the cluster. If unspecified, the zone of the first cluster matching the selector is used.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The cluster name prefix. A unique cluster name will be formed by appending a random suffix. The name must contain only lower-case letters (a-z), numbers (0-9), and hyphens (-). Must begin with a letter. Cannot begin or end with hyphen. Must consist of between 2 and 35 characters.",
			},

			"config": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The cluster configuration.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSchema(),
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The labels to associate with this cluster. Label keys must be between 1 and 63 characters long, and must conform to the following PCRE regular expression: p{Ll}p{Lo}{0,62} Label values must be between 1 and 63 characters long, and must conform to the following PCRE regular expression: [p{Ll}p{Lo}p{N}_-]{0,63} No more than 32 labels can be associated with a given cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"autoscaling_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Autoscaling config for the policy associated with the cluster. Cluster does not autoscale if this field is unset.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSchema(),
			},

			"encryption_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Encryption settings for the cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSchema(),
			},

			"endpoint_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Port/endpoint configuration for this cluster",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSchema(),
			},

			"gce_cluster_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The shared Compute Engine config settings for all instances in a cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSchema(),
			},

			"initialization_actions": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Commands to execute on each node after config is completed. By default, executables are run on master and all worker nodes. You can test a node's `role` metadata to run an executable on a master or worker node, as shown below using `curl` (you can also use `wget`): ROLE=$(curl -H Metadata-Flavor:Google http://metadata/computeMetadata/v1/instance/attributes/dataproc-role) if [[ \"${ROLE}\" == 'Master' ]]; then ... master specific actions ... else ... worker specific actions ... fi",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSchema(),
			},

			"lifecycle_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Lifecycle setting for the cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSchema(),
			},

			"master_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine config settings for the master instance in a cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigSchema(),
			},

			"secondary_worker_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine config settings for additional worker instances in a cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSchema(),
			},

			"security_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Security settings for the cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSchema(),
			},

			"software_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The config settings for software inside the cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSchema(),
			},

			"staging_bucket": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. A Cloud Storage bucket used to stage job dependencies, config files, and job driver console output. If you do not specify a staging bucket, Cloud Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's staging bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket (see [Dataproc staging bucket](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/staging-bucket)). **This field requires a Cloud Storage bucket name, not a URI to a Cloud Storage bucket.**",
			},

			"temp_bucket": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. A Cloud Storage bucket used to store ephemeral cluster and jobs data, such as Spark and MapReduce history files. If you do not specify a temp bucket, Dataproc will determine a Cloud Storage location (US, ASIA, or EU) for your cluster's temp bucket according to the Compute Engine zone where your cluster is deployed, and then create and manage this project-level, per-location bucket. The default bucket has a TTL of 90 days, but you can use any TTL (or none) if you specify a bucket. **This field requires a Cloud Storage bucket name, not a URI to a Cloud Storage bucket.**",
			},

			"worker_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine config settings for worker instances in a cluster.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The autoscaling policy used by the cluster. Only resource names including projectid and location (region) are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/locations/[dataproc_region]/autoscalingPolicies/[policy_id]` * `projects/[project_id]/locations/[dataproc_region]/autoscalingPolicies/[policy_id]` Note that the policy must be in the same project and Dataproc region.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gce_pd_kms_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The Cloud KMS key name to use for PD disk encryption for all instances in the cluster.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_http_port_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. If true, enable http access to specific ports on the cluster from external sources. Defaults to false.",
			},

			"http_ports": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Output only. The map of port descriptions to URLs. Will only be populated if enable_http_port_access is true.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"internal_ip_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. If true, all instances in the cluster will only have internal IP addresses. By default, clusters are not restricted to internal IP addresses, and will have ephemeral external IP addresses assigned to each instance. This `internal_ip_only` restriction can only be enabled for subnetwork enabled networks, and all off-cluster dependencies must be configured to be accessible without external IP addresses.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The Compute Engine metadata entries to add to all instances (see [Project and instance metadata](https://cloud.google.com/compute/docs/storing-retrieving-metadata#project_and_instance_metadata)).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The Compute Engine network to be used for machine communications. Cannot be specified with subnetwork_uri. If neither `network_uri` nor `subnetwork_uri` is specified, the \"default\" network of the project is used, if it exists. Cannot be a \"Custom Subnet Network\" (see [Using Subnetworks](https://cloud.google.com/compute/docs/subnetworks) for more information). A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/regions/global/default` * `projects/[project_id]/regions/global/default` * `default`",
			},

			"node_group_affinity": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Node Group Affinity for sole-tenant clusters.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySchema(),
			},

			"private_ipv6_google_access": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The type of IPv6 access for a cluster. Possible values: PRIVATE_IPV6_GOOGLE_ACCESS_UNSPECIFIED, INHERIT_FROM_SUBNETWORK, OUTBOUND, BIDIRECTIONAL",
			},

			"reservation_affinity": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Reservation Affinity for consuming Zonal reservation.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySchema(),
			},

			"service_account": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The [Dataproc service account](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/service-accounts#service_accounts_in_dataproc) (also see [VM Data Plane identity](https://cloud.google.com/dataproc/docs/concepts/iam/dataproc-principals#vm_service_account_data_plane_identity)) used by Dataproc cluster VM instances to access Google Cloud Platform services. If not specified, the [Compute Engine default service account](https://cloud.google.com/compute/docs/access/service-accounts#default_service_account) is used.",
			},

			"service_account_scopes": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The URIs of service account scopes to be included in Compute Engine instances. The following base set of scopes is always included: * https://www.googleapis.com/auth/cloud.useraccounts.readonly * https://www.googleapis.com/auth/devstorage.read_write * https://www.googleapis.com/auth/logging.write If no scopes are specified, the following defaults are also provided: * https://www.googleapis.com/auth/bigquery * https://www.googleapis.com/auth/bigtable.admin.table * https://www.googleapis.com/auth/bigtable.data * https://www.googleapis.com/auth/devstorage.full_control",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"shielded_instance_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Shielded Instance Config for clusters using Compute Engine Shielded VMs.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSchema(),
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The Compute Engine subnetwork to be used for machine communications. Cannot be specified with network_uri. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/regions/us-east1/subnetworks/sub0` * `projects/[project_id]/regions/us-east1/subnetworks/sub0` * `sub0`",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "The Compute Engine tags to add to all instances (see [Tagging instances](https://cloud.google.com/compute/docs/label-or-tag-resources#tags)).",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},

			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The zone where the Compute Engine cluster will be located. On a create request, it is required in the \"global\" region. If omitted in a non-global Dataproc region, the service will pick a zone in the corresponding Compute Engine region. On a get request, zone will always be present. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/[zone]` * `projects/[project_id]/zones/[zone]` * `us-central1-f`",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_group": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. The URI of a sole-tenant [node group resource](https://cloud.google.com/compute/docs/reference/rest/v1/nodeGroups) that the cluster will be created on. A full URL, partial URI, or node group name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-central1-a/nodeGroups/node-group-1` * `projects/[project_id]/zones/us-central1-a/nodeGroups/node-group-1` * `node-group-1`",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"consume_reservation_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Type of reservation to consume Possible values: TYPE_UNSPECIFIED, NO_RESERVATION, ANY_RESERVATION, SPECIFIC_RESERVATION",
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Corresponds to the label key of reservation resource.",
			},

			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Corresponds to the label values of reservation resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_integrity_monitoring": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Defines whether instances have integrity monitoring enabled. Integrity monitoring compares the most recent boot measurements to the integrity policy baseline and returns a pair of pass/fail results depending on whether they match or not.",
			},

			"enable_secure_boot": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Defines whether the instances have Secure Boot enabled. Secure Boot helps ensure that the system only runs authentic software by verifying the digital signature of all boot components, and halting the boot process if signature verification fails.",
			},

			"enable_vtpm": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Defines whether the instance have the vTPM enabled. Virtual Trusted Platform Module protects objects like keys, certificates and enables Measured Boot by performing the measurements needed to create a known good boot baseline, called the integrity policy baseline.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"executable_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Required. Cloud Storage URI of executable file.",
			},

			"execution_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Amount of time executable has to complete. Default is 10 minutes (see JSON representation of [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)). Cluster creation fails with an explanatory error message (the name of the executable that caused the error and the exceeded timeout period) if the executable is not completed at end of the timeout period.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_delete_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The time when cluster will be auto-deleted (see JSON representation of [Timestamp](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
			},

			"auto_delete_ttl": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The lifetime duration of cluster. The cluster will be auto-deleted at the end of this period. Minimum value is 10 minutes; maximum value is 14 days (see JSON representation of [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
			},

			"idle_delete_ttl": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The duration to keep the cluster alive while idling (when no jobs are running). Passing this threshold will cause the cluster to be deleted. Minimum value is 5 minutes; maximum value is 14 days (see JSON representation of [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
			},

			"idle_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when cluster became idle (most recent job finished) and became eligible for deletion due to idleness (see JSON representation of [Timestamp](https://developers.google.com/protocol-buffers/docs/proto3#json)).",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerators": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine accelerator configuration for these instances.",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSchema(),
			},

			"disk_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Disk option config settings.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSchema(),
			},

			"image": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/[image-id]` * `projects/[project_id]/global/images/[image-id]` * `image-id` Image family examples. Dataproc will use the most recent image from the family: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/family/[custom-image-family-name]` * `projects/[project_id]/global/images/family/[custom-image-family-name]` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.",
			},

			"machine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `n1-standard-2` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.",
			},

			"min_cpu_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Specifies the minimum cpu platform for the Instance Group. See [Dataproc -> Minimum CPU Platform](https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).",
			},

			"num_instances": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The number of VM instances in the instance group. For [HA cluster](/dataproc/docs/concepts/configuring-clusters/high-availability) [master_config](#FIELDS.master_config) groups, **must be set to 3**. For standard cluster [master_config](#FIELDS.master_config) groups, **must be set to 1**.",
			},

			"preemptibility": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE",
			},

			"instance_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"is_preemptible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Specifies that this instance group contains preemptible instances.",
			},

			"managed_group_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerator_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of the accelerator cards of this type exposed to this instance.",
			},

			"accelerator_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See [Compute Engine AcceleratorTypes](https://cloud.google.com/compute/docs/reference/beta/acceleratorTypes). Examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `nvidia-tesla-k80` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"boot_disk_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Size in GB of the boot disk (default is 500GB).",
			},

			"boot_disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Type of the boot disk (default is \"pd-standard\"). Valid values: \"pd-balanced\" (Persistent Disk Balanced Solid State Drive), \"pd-ssd\" (Persistent Disk Solid State Drive), or \"pd-standard\" (Persistent Disk Hard Disk Drive). See [Disk types](https://cloud.google.com/compute/docs/disks#disk-types).",
			},

			"num_local_ssds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and [HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_group_manager_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Instance Group Manager for this group.",
			},

			"instance_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Instance Template used for the Managed Instance Group.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerators": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine accelerator configuration for these instances.",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSchema(),
			},

			"disk_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Disk option config settings.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSchema(),
			},

			"image": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/[image-id]` * `projects/[project_id]/global/images/[image-id]` * `image-id` Image family examples. Dataproc will use the most recent image from the family: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/family/[custom-image-family-name]` * `projects/[project_id]/global/images/family/[custom-image-family-name]` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.",
			},

			"machine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `n1-standard-2` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.",
			},

			"min_cpu_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Specifies the minimum cpu platform for the Instance Group. See [Dataproc -> Minimum CPU Platform](https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).",
			},

			"num_instances": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The number of VM instances in the instance group. For [HA cluster](/dataproc/docs/concepts/configuring-clusters/high-availability) [master_config](#FIELDS.master_config) groups, **must be set to 3**. For standard cluster [master_config](#FIELDS.master_config) groups, **must be set to 1**.",
			},

			"preemptibility": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE",
			},

			"instance_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"is_preemptible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Specifies that this instance group contains preemptible instances.",
			},

			"managed_group_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerator_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of the accelerator cards of this type exposed to this instance.",
			},

			"accelerator_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See [Compute Engine AcceleratorTypes](https://cloud.google.com/compute/docs/reference/beta/acceleratorTypes). Examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `nvidia-tesla-k80` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"boot_disk_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Size in GB of the boot disk (default is 500GB).",
			},

			"boot_disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Type of the boot disk (default is \"pd-standard\"). Valid values: \"pd-balanced\" (Persistent Disk Balanced Solid State Drive), \"pd-ssd\" (Persistent Disk Solid State Drive), or \"pd-standard\" (Persistent Disk Hard Disk Drive). See [Disk types](https://cloud.google.com/compute/docs/disks#disk-types).",
			},

			"num_local_ssds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and [HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_group_manager_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Instance Group Manager for this group.",
			},

			"instance_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Instance Template used for the Managed Instance Group.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kerberos_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Kerberos related configuration.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cross_realm_trust_admin_server": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The admin server (IP or hostname) for the remote trusted realm in a cross realm trust relationship.",
			},

			"cross_realm_trust_kdc": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The KDC (IP or hostname) for the remote trusted realm in a cross realm trust relationship.",
			},

			"cross_realm_trust_realm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The remote realm the Dataproc on-cluster KDC will trust, should the user enable cross realm trust.",
			},

			"cross_realm_trust_shared_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the shared password between the on-cluster Kerberos realm and the remote trusted realm, in a cross realm trust relationship.",
			},

			"enable_kerberos": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Flag to indicate whether to Kerberize the cluster (default: false). Set this field to true to enable Kerberos on a cluster.",
			},

			"kdc_db_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the master key of the KDC database.",
			},

			"key_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the password to the user provided key. For the self-signed certificate, this password is generated by Dataproc.",
			},

			"keystore": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of the keystore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.",
			},

			"keystore_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the password to the user provided keystore. For the self-signed certificate, this password is generated by Dataproc.",
			},

			"kms_key": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The uri of the KMS key used to encrypt various sensitive files.",
			},

			"realm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The name of the on-cluster Kerberos realm. If not specified, the uppercased domain of hostnames will be the realm.",
			},

			"root_principal_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the root principal password.",
			},

			"tgt_lifetime_hours": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The lifetime of the ticket granting ticket, in hours. If not specified, or user specifies 0, then default value 10 will be used.",
			},

			"truststore": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of the truststore file used for SSL encryption. If not provided, Dataproc will provide a self-signed certificate.",
			},

			"truststore_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Cloud Storage URI of a KMS encrypted file containing the password to the user provided truststore. For the self-signed certificate, this password is generated by Dataproc.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The version of software inside the cluster. It must be one of the supported [Dataproc Versions](https://cloud.google.com/dataproc/docs/concepts/versioning/dataproc-versions#supported_dataproc_versions), such as \"1.2\" (including a subminor version, such as \"1.2.29\"), or the [\"preview\" version](https://cloud.google.com/dataproc/docs/concepts/versioning/dataproc-versions#other_versions). If unspecified, it defaults to the latest Debian version.",
			},

			"optional_components": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The set of components to activate on the cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The properties to set on daemon config files. Property keys are specified in `prefix:property` format, for example `core:hadoop.tmp.dir`. The following are supported prefixes and their mappings: * capacity-scheduler: `capacity-scheduler.xml` * core: `core-site.xml` * distcp: `distcp-default.xml` * hdfs: `hdfs-site.xml` * hive: `hive-site.xml` * mapred: `mapred-site.xml` * pig: `pig.properties` * spark: `spark-defaults.conf` * yarn: `yarn-site.xml` For more information, see [Cluster properties](https://cloud.google.com/dataproc/docs/concepts/cluster-properties).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerators": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine accelerator configuration for these instances.",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSchema(),
			},

			"disk_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Disk option config settings.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSchema(),
			},

			"image": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The Compute Engine image resource used for cluster instances. The URI can represent an image or image family. Image examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/[image-id]` * `projects/[project_id]/global/images/[image-id]` * `image-id` Image family examples. Dataproc will use the most recent image from the family: * `https://www.googleapis.com/compute/beta/projects/[project_id]/global/images/family/[custom-image-family-name]` * `projects/[project_id]/global/images/family/[custom-image-family-name]` If the URI is unspecified, it will be inferred from `SoftwareConfig.image_version` or the system default.",
			},

			"machine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The Compute Engine machine type used for cluster instances. A full URL, partial URI, or short name are valid. Examples: * `https://www.googleapis.com/compute/v1/projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `projects/[project_id]/zones/us-east1-a/machineTypes/n1-standard-2` * `n1-standard-2` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the machine type resource, for example, `n1-standard-2`.",
			},

			"min_cpu_platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Specifies the minimum cpu platform for the Instance Group. See [Dataproc -> Minimum CPU Platform](https://cloud.google.com/dataproc/docs/concepts/compute/dataproc-min-cpu).",
			},

			"num_instances": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. The number of VM instances in the instance group. For [HA cluster](/dataproc/docs/concepts/configuring-clusters/high-availability) [master_config](#FIELDS.master_config) groups, **must be set to 3**. For standard cluster [master_config](#FIELDS.master_config) groups, **must be set to 1**.",
			},

			"preemptibility": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Specifies the preemptibility of the instance group. The default value for master and worker groups is `NON_PREEMPTIBLE`. This default cannot be changed. The default value for secondary instances is `PREEMPTIBLE`. Possible values: PREEMPTIBILITY_UNSPECIFIED, NON_PREEMPTIBLE, PREEMPTIBLE",
			},

			"instance_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The list of instance names. Dataproc derives the names from `cluster_name`, `num_instances`, and the instance group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"is_preemptible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Specifies that this instance group contains preemptible instances.",
			},

			"managed_group_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. The config for Compute Engine Instance Group Manager that manages this group. This is only used for preemptible instance groups.",
				Elem:        DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"accelerator_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of the accelerator cards of this type exposed to this instance.",
			},

			"accelerator_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Full URL, partial URI, or short name of the accelerator type resource to expose to this instance. See [Compute Engine AcceleratorTypes](https://cloud.google.com/compute/docs/reference/beta/acceleratorTypes). Examples: * `https://www.googleapis.com/compute/beta/projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `projects/[project_id]/zones/us-east1-a/acceleratorTypes/nvidia-tesla-k80` * `nvidia-tesla-k80` **Auto Zone Exception**: If you are using the Dataproc [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/configuring-clusters/auto-zone#using_auto_zone_placement) feature, you must use the short name of the accelerator type resource, for example, `nvidia-tesla-k80`.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"boot_disk_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Size in GB of the boot disk (default is 500GB).",
			},

			"boot_disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Type of the boot disk (default is \"pd-standard\"). Valid values: \"pd-balanced\" (Persistent Disk Balanced Solid State Drive), \"pd-ssd\" (Persistent Disk Solid State Drive), or \"pd-standard\" (Persistent Disk Hard Disk Drive). See [Disk types](https://cloud.google.com/compute/docs/disks#disk-types).",
			},

			"num_local_ssds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Number of attached SSDs, from 0 to 4 (default is 0). If SSDs are not attached, the boot disk is used to store runtime logs and [HDFS](https://hadoop.apache.org/docs/r1.2.1/hdfs_user_guide.html) data. If one or more SSDs are attached, this runtime bulk data is spread across them, and the boot disk contains only basic config and installed binaries.",
			},
		},
	}
}

func DataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_group_manager_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Instance Group Manager for this group.",
			},

			"instance_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Instance Template used for the Managed Instance Group.",
			},
		},
	}
}

func DataprocWorkflowTemplateParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fields": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Paths to all fields that the parameter replaces. A field is allowed to appear in at most one parameter's list of field paths. A field path is similar in syntax to a google.protobuf.FieldMask. For example, a field path that references the zone field of a workflow template's cluster selector would be specified as `placement.clusterSelector.zone`. Also, field paths can reference fields using the following syntax: * Values in maps can be referenced by key: * labels['key'] * placement.clusterSelector.clusterLabels['key'] * placement.managedCluster.labels['key'] * placement.clusterSelector.clusterLabels['key'] * jobs['step-id'].labels['key'] * Jobs in the jobs list can be referenced by step-id: * jobs['step-id'].hadoopJob.mainJarFileUri * jobs['step-id'].hiveJob.queryFileUri * jobs['step-id'].pySparkJob.mainPythonFileUri * jobs['step-id'].hadoopJob.jarFileUris[0] * jobs['step-id'].hadoopJob.archiveUris[0] * jobs['step-id'].hadoopJob.fileUris[0] * jobs['step-id'].pySparkJob.pythonFileUris[0] * Items in repeated fields can be referenced by a zero-based index: * jobs['step-id'].sparkJob.args[0] * Other examples: * jobs['step-id'].hadoopJob.properties['key'] * jobs['step-id'].hadoopJob.args[0] * jobs['step-id'].hiveJob.scriptVariables['key'] * jobs['step-id'].hadoopJob.mainJarFileUri * placement.clusterSelector.zone It may not be possible to parameterize maps and repeated fields in their entirety since only individual map values and individual items in repeated fields can be referenced. For example, the following field paths are invalid: - placement.clusterSelector.clusterLabels - jobs['step-id'].sparkJob.args",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Parameter name. The parameter name is used as the key, and paired with the parameter value, which are passed to the template when the template is instantiated. The name must contain only capital letters (A-Z), numbers (0-9), and underscores (_), and must not start with a number. The maximum length is 40 characters.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Brief description of the parameter. Must not exceed 1024 characters.",
			},

			"validation": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. Validation rules to be applied to this parameter's value.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateParametersValidationSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplateParametersValidationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"regex": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Validation based on regular expressions.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateParametersValidationRegexSchema(),
			},

			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Validation based on a list of allowed values.",
				MaxItems:    1,
				Elem:        DataprocWorkflowTemplateParametersValidationValuesSchema(),
			},
		},
	}
}

func DataprocWorkflowTemplateParametersValidationRegexSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"regexes": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. RE2 regular expressions used to validate the parameter's value. The value must match the regex in its entirety (substring matches are not sufficient).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func DataprocWorkflowTemplateParametersValidationValuesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Required. List of allowed values for the parameter.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDataprocWorkflowTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataproc.WorkflowTemplate{
		Jobs:       expandDataprocWorkflowTemplateJobsArray(d.Get("jobs")),
		Location:   dcl.String(d.Get("location").(string)),
		Name:       dcl.String(d.Get("name").(string)),
		Placement:  expandDataprocWorkflowTemplatePlacement(d.Get("placement")),
		DagTimeout: dcl.String(d.Get("dag_timeout").(string)),
		Labels:     tpgresource.CheckStringMap(d.Get("labels")),
		Parameters: expandDataprocWorkflowTemplateParametersArray(d.Get("parameters")),
		Project:    dcl.String(project),
		Version:    dcl.Int64OrNil(int64(d.Get("version").(int))),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataprocClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyWorkflowTemplate(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating WorkflowTemplate: %s", err)
	}

	log.Printf("[DEBUG] Finished creating WorkflowTemplate %q: %#v", d.Id(), res)

	return resourceDataprocWorkflowTemplateRead(d, meta)
}

func resourceDataprocWorkflowTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataproc.WorkflowTemplate{
		Jobs:       expandDataprocWorkflowTemplateJobsArray(d.Get("jobs")),
		Location:   dcl.String(d.Get("location").(string)),
		Name:       dcl.String(d.Get("name").(string)),
		Placement:  expandDataprocWorkflowTemplatePlacement(d.Get("placement")),
		DagTimeout: dcl.String(d.Get("dag_timeout").(string)),
		Labels:     tpgresource.CheckStringMap(d.Get("labels")),
		Parameters: expandDataprocWorkflowTemplateParametersArray(d.Get("parameters")),
		Project:    dcl.String(project),
		Version:    dcl.Int64OrNil(int64(d.Get("version").(int))),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataprocClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetWorkflowTemplate(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("DataprocWorkflowTemplate %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("jobs", flattenDataprocWorkflowTemplateJobsArray(res.Jobs)); err != nil {
		return fmt.Errorf("error setting jobs in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("placement", flattenDataprocWorkflowTemplatePlacement(res.Placement)); err != nil {
		return fmt.Errorf("error setting placement in state: %s", err)
	}
	if err = d.Set("dag_timeout", res.DagTimeout); err != nil {
		return fmt.Errorf("error setting dag_timeout in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("parameters", flattenDataprocWorkflowTemplateParametersArray(res.Parameters)); err != nil {
		return fmt.Errorf("error setting parameters in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("version", res.Version); err != nil {
		return fmt.Errorf("error setting version in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}

func resourceDataprocWorkflowTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataproc.WorkflowTemplate{
		Jobs:       expandDataprocWorkflowTemplateJobsArray(d.Get("jobs")),
		Location:   dcl.String(d.Get("location").(string)),
		Name:       dcl.String(d.Get("name").(string)),
		Placement:  expandDataprocWorkflowTemplatePlacement(d.Get("placement")),
		DagTimeout: dcl.String(d.Get("dag_timeout").(string)),
		Labels:     tpgresource.CheckStringMap(d.Get("labels")),
		Parameters: expandDataprocWorkflowTemplateParametersArray(d.Get("parameters")),
		Project:    dcl.String(project),
		Version:    dcl.Int64OrNil(int64(d.Get("version").(int))),
	}

	log.Printf("[DEBUG] Deleting WorkflowTemplate %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataprocClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteWorkflowTemplate(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting WorkflowTemplate: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting WorkflowTemplate %q", d.Id())
	return nil
}

func resourceDataprocWorkflowTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/workflowTemplates/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/workflowTemplates/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandDataprocWorkflowTemplateJobsArray(o interface{}) []dataproc.WorkflowTemplateJobs {
	if o == nil {
		return make([]dataproc.WorkflowTemplateJobs, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]dataproc.WorkflowTemplateJobs, 0)
	}

	items := make([]dataproc.WorkflowTemplateJobs, 0, len(objs))
	for _, item := range objs {
		i := expandDataprocWorkflowTemplateJobs(item)
		items = append(items, *i)
	}

	return items
}

func expandDataprocWorkflowTemplateJobs(o interface{}) *dataproc.WorkflowTemplateJobs {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobs
	}

	obj := o.(map[string]interface{})
	return &dataproc.WorkflowTemplateJobs{
		StepId:              dcl.String(obj["step_id"].(string)),
		HadoopJob:           expandDataprocWorkflowTemplateJobsHadoopJob(obj["hadoop_job"]),
		HiveJob:             expandDataprocWorkflowTemplateJobsHiveJob(obj["hive_job"]),
		Labels:              tpgresource.CheckStringMap(obj["labels"]),
		PigJob:              expandDataprocWorkflowTemplateJobsPigJob(obj["pig_job"]),
		PrerequisiteStepIds: tpgdclresource.ExpandStringArray(obj["prerequisite_step_ids"]),
		PrestoJob:           expandDataprocWorkflowTemplateJobsPrestoJob(obj["presto_job"]),
		PysparkJob:          expandDataprocWorkflowTemplateJobsPysparkJob(obj["pyspark_job"]),
		Scheduling:          expandDataprocWorkflowTemplateJobsScheduling(obj["scheduling"]),
		SparkJob:            expandDataprocWorkflowTemplateJobsSparkJob(obj["spark_job"]),
		SparkRJob:           expandDataprocWorkflowTemplateJobsSparkRJob(obj["spark_r_job"]),
		SparkSqlJob:         expandDataprocWorkflowTemplateJobsSparkSqlJob(obj["spark_sql_job"]),
	}
}

func flattenDataprocWorkflowTemplateJobsArray(objs []dataproc.WorkflowTemplateJobs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenDataprocWorkflowTemplateJobs(&item)
		items = append(items, i)
	}

	return items
}

func flattenDataprocWorkflowTemplateJobs(obj *dataproc.WorkflowTemplateJobs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"step_id":               obj.StepId,
		"hadoop_job":            flattenDataprocWorkflowTemplateJobsHadoopJob(obj.HadoopJob),
		"hive_job":              flattenDataprocWorkflowTemplateJobsHiveJob(obj.HiveJob),
		"labels":                obj.Labels,
		"pig_job":               flattenDataprocWorkflowTemplateJobsPigJob(obj.PigJob),
		"prerequisite_step_ids": obj.PrerequisiteStepIds,
		"presto_job":            flattenDataprocWorkflowTemplateJobsPrestoJob(obj.PrestoJob),
		"pyspark_job":           flattenDataprocWorkflowTemplateJobsPysparkJob(obj.PysparkJob),
		"scheduling":            flattenDataprocWorkflowTemplateJobsScheduling(obj.Scheduling),
		"spark_job":             flattenDataprocWorkflowTemplateJobsSparkJob(obj.SparkJob),
		"spark_r_job":           flattenDataprocWorkflowTemplateJobsSparkRJob(obj.SparkRJob),
		"spark_sql_job":         flattenDataprocWorkflowTemplateJobsSparkSqlJob(obj.SparkSqlJob),
	}

	return transformed

}

func expandDataprocWorkflowTemplateJobsHadoopJob(o interface{}) *dataproc.WorkflowTemplateJobsHadoopJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsHadoopJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsHadoopJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsHadoopJob{
		ArchiveUris:    tpgdclresource.ExpandStringArray(obj["archive_uris"]),
		Args:           tpgdclresource.ExpandStringArray(obj["args"]),
		FileUris:       tpgdclresource.ExpandStringArray(obj["file_uris"]),
		JarFileUris:    tpgdclresource.ExpandStringArray(obj["jar_file_uris"]),
		LoggingConfig:  expandDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(obj["logging_config"]),
		MainClass:      dcl.String(obj["main_class"].(string)),
		MainJarFileUri: dcl.String(obj["main_jar_file_uri"].(string)),
		Properties:     tpgresource.CheckStringMap(obj["properties"]),
	}
}

func flattenDataprocWorkflowTemplateJobsHadoopJob(obj *dataproc.WorkflowTemplateJobsHadoopJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"archive_uris":      obj.ArchiveUris,
		"args":              obj.Args,
		"file_uris":         obj.FileUris,
		"jar_file_uris":     obj.JarFileUris,
		"logging_config":    flattenDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(obj.LoggingConfig),
		"main_class":        obj.MainClass,
		"main_jar_file_uri": obj.MainJarFileUri,
		"properties":        obj.Properties,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsHadoopJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsHadoopJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsHadoopJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsHadoopJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsHadoopJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsHiveJob(o interface{}) *dataproc.WorkflowTemplateJobsHiveJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsHiveJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsHiveJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsHiveJob{
		ContinueOnFailure: dcl.Bool(obj["continue_on_failure"].(bool)),
		JarFileUris:       tpgdclresource.ExpandStringArray(obj["jar_file_uris"]),
		Properties:        tpgresource.CheckStringMap(obj["properties"]),
		QueryFileUri:      dcl.String(obj["query_file_uri"].(string)),
		QueryList:         expandDataprocWorkflowTemplateJobsHiveJobQueryList(obj["query_list"]),
		ScriptVariables:   tpgresource.CheckStringMap(obj["script_variables"]),
	}
}

func flattenDataprocWorkflowTemplateJobsHiveJob(obj *dataproc.WorkflowTemplateJobsHiveJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"continue_on_failure": obj.ContinueOnFailure,
		"jar_file_uris":       obj.JarFileUris,
		"properties":          obj.Properties,
		"query_file_uri":      obj.QueryFileUri,
		"query_list":          flattenDataprocWorkflowTemplateJobsHiveJobQueryList(obj.QueryList),
		"script_variables":    obj.ScriptVariables,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsHiveJobQueryList(o interface{}) *dataproc.WorkflowTemplateJobsHiveJobQueryList {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsHiveJobQueryList
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsHiveJobQueryList
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsHiveJobQueryList{
		Queries: tpgdclresource.ExpandStringArray(obj["queries"]),
	}
}

func flattenDataprocWorkflowTemplateJobsHiveJobQueryList(obj *dataproc.WorkflowTemplateJobsHiveJobQueryList) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"queries": obj.Queries,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPigJob(o interface{}) *dataproc.WorkflowTemplateJobsPigJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPigJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPigJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPigJob{
		ContinueOnFailure: dcl.Bool(obj["continue_on_failure"].(bool)),
		JarFileUris:       tpgdclresource.ExpandStringArray(obj["jar_file_uris"]),
		LoggingConfig:     expandDataprocWorkflowTemplateJobsPigJobLoggingConfig(obj["logging_config"]),
		Properties:        tpgresource.CheckStringMap(obj["properties"]),
		QueryFileUri:      dcl.String(obj["query_file_uri"].(string)),
		QueryList:         expandDataprocWorkflowTemplateJobsPigJobQueryList(obj["query_list"]),
		ScriptVariables:   tpgresource.CheckStringMap(obj["script_variables"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPigJob(obj *dataproc.WorkflowTemplateJobsPigJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"continue_on_failure": obj.ContinueOnFailure,
		"jar_file_uris":       obj.JarFileUris,
		"logging_config":      flattenDataprocWorkflowTemplateJobsPigJobLoggingConfig(obj.LoggingConfig),
		"properties":          obj.Properties,
		"query_file_uri":      obj.QueryFileUri,
		"query_list":          flattenDataprocWorkflowTemplateJobsPigJobQueryList(obj.QueryList),
		"script_variables":    obj.ScriptVariables,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPigJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsPigJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPigJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPigJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPigJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPigJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsPigJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPigJobQueryList(o interface{}) *dataproc.WorkflowTemplateJobsPigJobQueryList {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPigJobQueryList
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPigJobQueryList
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPigJobQueryList{
		Queries: tpgdclresource.ExpandStringArray(obj["queries"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPigJobQueryList(obj *dataproc.WorkflowTemplateJobsPigJobQueryList) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"queries": obj.Queries,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPrestoJob(o interface{}) *dataproc.WorkflowTemplateJobsPrestoJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPrestoJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPrestoJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPrestoJob{
		ClientTags:        tpgdclresource.ExpandStringArray(obj["client_tags"]),
		ContinueOnFailure: dcl.Bool(obj["continue_on_failure"].(bool)),
		LoggingConfig:     expandDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(obj["logging_config"]),
		OutputFormat:      dcl.String(obj["output_format"].(string)),
		Properties:        tpgresource.CheckStringMap(obj["properties"]),
		QueryFileUri:      dcl.String(obj["query_file_uri"].(string)),
		QueryList:         expandDataprocWorkflowTemplateJobsPrestoJobQueryList(obj["query_list"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPrestoJob(obj *dataproc.WorkflowTemplateJobsPrestoJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"client_tags":         obj.ClientTags,
		"continue_on_failure": obj.ContinueOnFailure,
		"logging_config":      flattenDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(obj.LoggingConfig),
		"output_format":       obj.OutputFormat,
		"properties":          obj.Properties,
		"query_file_uri":      obj.QueryFileUri,
		"query_list":          flattenDataprocWorkflowTemplateJobsPrestoJobQueryList(obj.QueryList),
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsPrestoJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPrestoJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPrestoJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPrestoJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsPrestoJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPrestoJobQueryList(o interface{}) *dataproc.WorkflowTemplateJobsPrestoJobQueryList {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPrestoJobQueryList
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPrestoJobQueryList
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPrestoJobQueryList{
		Queries: tpgdclresource.ExpandStringArray(obj["queries"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPrestoJobQueryList(obj *dataproc.WorkflowTemplateJobsPrestoJobQueryList) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"queries": obj.Queries,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPysparkJob(o interface{}) *dataproc.WorkflowTemplateJobsPysparkJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPysparkJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPysparkJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPysparkJob{
		MainPythonFileUri: dcl.String(obj["main_python_file_uri"].(string)),
		ArchiveUris:       tpgdclresource.ExpandStringArray(obj["archive_uris"]),
		Args:              tpgdclresource.ExpandStringArray(obj["args"]),
		FileUris:          tpgdclresource.ExpandStringArray(obj["file_uris"]),
		JarFileUris:       tpgdclresource.ExpandStringArray(obj["jar_file_uris"]),
		LoggingConfig:     expandDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(obj["logging_config"]),
		Properties:        tpgresource.CheckStringMap(obj["properties"]),
		PythonFileUris:    tpgdclresource.ExpandStringArray(obj["python_file_uris"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPysparkJob(obj *dataproc.WorkflowTemplateJobsPysparkJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"main_python_file_uri": obj.MainPythonFileUri,
		"archive_uris":         obj.ArchiveUris,
		"args":                 obj.Args,
		"file_uris":            obj.FileUris,
		"jar_file_uris":        obj.JarFileUris,
		"logging_config":       flattenDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(obj.LoggingConfig),
		"properties":           obj.Properties,
		"python_file_uris":     obj.PythonFileUris,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsPysparkJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsPysparkJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsPysparkJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsPysparkJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsPysparkJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsScheduling(o interface{}) *dataproc.WorkflowTemplateJobsScheduling {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsScheduling
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsScheduling
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsScheduling{
		MaxFailuresPerHour: dcl.Int64(int64(obj["max_failures_per_hour"].(int))),
		MaxFailuresTotal:   dcl.Int64(int64(obj["max_failures_total"].(int))),
	}
}

func flattenDataprocWorkflowTemplateJobsScheduling(obj *dataproc.WorkflowTemplateJobsScheduling) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"max_failures_per_hour": obj.MaxFailuresPerHour,
		"max_failures_total":    obj.MaxFailuresTotal,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkJob(o interface{}) *dataproc.WorkflowTemplateJobsSparkJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkJob{
		ArchiveUris:    tpgdclresource.ExpandStringArray(obj["archive_uris"]),
		Args:           tpgdclresource.ExpandStringArray(obj["args"]),
		FileUris:       tpgdclresource.ExpandStringArray(obj["file_uris"]),
		JarFileUris:    tpgdclresource.ExpandStringArray(obj["jar_file_uris"]),
		LoggingConfig:  expandDataprocWorkflowTemplateJobsSparkJobLoggingConfig(obj["logging_config"]),
		MainClass:      dcl.String(obj["main_class"].(string)),
		MainJarFileUri: dcl.String(obj["main_jar_file_uri"].(string)),
		Properties:     tpgresource.CheckStringMap(obj["properties"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkJob(obj *dataproc.WorkflowTemplateJobsSparkJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"archive_uris":      obj.ArchiveUris,
		"args":              obj.Args,
		"file_uris":         obj.FileUris,
		"jar_file_uris":     obj.JarFileUris,
		"logging_config":    flattenDataprocWorkflowTemplateJobsSparkJobLoggingConfig(obj.LoggingConfig),
		"main_class":        obj.MainClass,
		"main_jar_file_uri": obj.MainJarFileUri,
		"properties":        obj.Properties,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsSparkJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsSparkJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkRJob(o interface{}) *dataproc.WorkflowTemplateJobsSparkRJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkRJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkRJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkRJob{
		MainRFileUri:  dcl.String(obj["main_r_file_uri"].(string)),
		ArchiveUris:   tpgdclresource.ExpandStringArray(obj["archive_uris"]),
		Args:          tpgdclresource.ExpandStringArray(obj["args"]),
		FileUris:      tpgdclresource.ExpandStringArray(obj["file_uris"]),
		LoggingConfig: expandDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(obj["logging_config"]),
		Properties:    tpgresource.CheckStringMap(obj["properties"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkRJob(obj *dataproc.WorkflowTemplateJobsSparkRJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"main_r_file_uri": obj.MainRFileUri,
		"archive_uris":    obj.ArchiveUris,
		"args":            obj.Args,
		"file_uris":       obj.FileUris,
		"logging_config":  flattenDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(obj.LoggingConfig),
		"properties":      obj.Properties,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsSparkRJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkRJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkRJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkRJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsSparkRJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkSqlJob(o interface{}) *dataproc.WorkflowTemplateJobsSparkSqlJob {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkSqlJob
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkSqlJob
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkSqlJob{
		JarFileUris:     tpgdclresource.ExpandStringArray(obj["jar_file_uris"]),
		LoggingConfig:   expandDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(obj["logging_config"]),
		Properties:      tpgresource.CheckStringMap(obj["properties"]),
		QueryFileUri:    dcl.String(obj["query_file_uri"].(string)),
		QueryList:       expandDataprocWorkflowTemplateJobsSparkSqlJobQueryList(obj["query_list"]),
		ScriptVariables: tpgresource.CheckStringMap(obj["script_variables"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkSqlJob(obj *dataproc.WorkflowTemplateJobsSparkSqlJob) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"jar_file_uris":    obj.JarFileUris,
		"logging_config":   flattenDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(obj.LoggingConfig),
		"properties":       obj.Properties,
		"query_file_uri":   obj.QueryFileUri,
		"query_list":       flattenDataprocWorkflowTemplateJobsSparkSqlJobQueryList(obj.QueryList),
		"script_variables": obj.ScriptVariables,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(o interface{}) *dataproc.WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkSqlJobLoggingConfig{
		DriverLogLevels: tpgresource.CheckStringMap(obj["driver_log_levels"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(obj *dataproc.WorkflowTemplateJobsSparkSqlJobLoggingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"driver_log_levels": obj.DriverLogLevels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateJobsSparkSqlJobQueryList(o interface{}) *dataproc.WorkflowTemplateJobsSparkSqlJobQueryList {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkSqlJobQueryList
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateJobsSparkSqlJobQueryList
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateJobsSparkSqlJobQueryList{
		Queries: tpgdclresource.ExpandStringArray(obj["queries"]),
	}
}

func flattenDataprocWorkflowTemplateJobsSparkSqlJobQueryList(obj *dataproc.WorkflowTemplateJobsSparkSqlJobQueryList) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"queries": obj.Queries,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacement(o interface{}) *dataproc.WorkflowTemplatePlacement {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacement
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacement
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacement{
		ClusterSelector: expandDataprocWorkflowTemplatePlacementClusterSelector(obj["cluster_selector"]),
		ManagedCluster:  expandDataprocWorkflowTemplatePlacementManagedCluster(obj["managed_cluster"]),
	}
}

func flattenDataprocWorkflowTemplatePlacement(obj *dataproc.WorkflowTemplatePlacement) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster_selector": flattenDataprocWorkflowTemplatePlacementClusterSelector(obj.ClusterSelector),
		"managed_cluster":  flattenDataprocWorkflowTemplatePlacementManagedCluster(obj.ManagedCluster),
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementClusterSelector(o interface{}) *dataproc.WorkflowTemplatePlacementClusterSelector {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementClusterSelector
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementClusterSelector
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementClusterSelector{
		ClusterLabels: tpgresource.CheckStringMap(obj["cluster_labels"]),
		Zone:          dcl.StringOrNil(obj["zone"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementClusterSelector(obj *dataproc.WorkflowTemplatePlacementClusterSelector) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster_labels": obj.ClusterLabels,
		"zone":           obj.Zone,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedCluster(o interface{}) *dataproc.WorkflowTemplatePlacementManagedCluster {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedCluster
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedCluster
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedCluster{
		ClusterName: dcl.String(obj["cluster_name"].(string)),
		Config:      expandDataprocWorkflowTemplatePlacementManagedClusterConfig(obj["config"]),
		Labels:      tpgresource.CheckStringMap(obj["labels"]),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedCluster(obj *dataproc.WorkflowTemplatePlacementManagedCluster) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster_name": obj.ClusterName,
		"config":       flattenDataprocWorkflowTemplatePlacementManagedClusterConfig(obj.Config),
		"labels":       obj.Labels,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfig{
		AutoscalingConfig:     expandDataprocWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(obj["autoscaling_config"]),
		EncryptionConfig:      expandDataprocWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(obj["encryption_config"]),
		EndpointConfig:        expandDataprocWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(obj["endpoint_config"]),
		GceClusterConfig:      expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(obj["gce_cluster_config"]),
		InitializationActions: expandDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActionsArray(obj["initialization_actions"]),
		LifecycleConfig:       expandDataprocWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(obj["lifecycle_config"]),
		MasterConfig:          expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfig(obj["master_config"]),
		SecondaryWorkerConfig: expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(obj["secondary_worker_config"]),
		SecurityConfig:        expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(obj["security_config"]),
		SoftwareConfig:        expandDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(obj["software_config"]),
		StagingBucket:         dcl.String(obj["staging_bucket"].(string)),
		TempBucket:            dcl.String(obj["temp_bucket"].(string)),
		WorkerConfig:          expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(obj["worker_config"]),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"autoscaling_config":      flattenDataprocWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(obj.AutoscalingConfig),
		"encryption_config":       flattenDataprocWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(obj.EncryptionConfig),
		"endpoint_config":         flattenDataprocWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(obj.EndpointConfig),
		"gce_cluster_config":      flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(obj.GceClusterConfig),
		"initialization_actions":  flattenDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActionsArray(obj.InitializationActions),
		"lifecycle_config":        flattenDataprocWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(obj.LifecycleConfig),
		"master_config":           flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfig(obj.MasterConfig),
		"secondary_worker_config": flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(obj.SecondaryWorkerConfig),
		"security_config":         flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(obj.SecurityConfig),
		"software_config":         flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(obj.SoftwareConfig),
		"staging_bucket":          obj.StagingBucket,
		"temp_bucket":             obj.TempBucket,
		"worker_config":           flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(obj.WorkerConfig),
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{
		Policy: dcl.String(obj["policy"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"policy": obj.Policy,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{
		GcePdKmsKeyName: dcl.String(obj["gce_pd_kms_key_name"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"gce_pd_kms_key_name": obj.GcePdKmsKeyName,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigEndpointConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigEndpointConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{
		EnableHttpPortAccess: dcl.Bool(obj["enable_http_port_access"].(bool)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"enable_http_port_access": obj.EnableHttpPortAccess,
		"http_ports":              obj.HttpPorts,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{
		InternalIPOnly:          dcl.Bool(obj["internal_ip_only"].(bool)),
		Metadata:                tpgresource.CheckStringMap(obj["metadata"]),
		Network:                 dcl.String(obj["network"].(string)),
		NodeGroupAffinity:       expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(obj["node_group_affinity"]),
		PrivateIPv6GoogleAccess: dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef(obj["private_ipv6_google_access"].(string)),
		ReservationAffinity:     expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(obj["reservation_affinity"]),
		ServiceAccount:          dcl.String(obj["service_account"].(string)),
		ServiceAccountScopes:    tpgdclresource.ExpandStringArray(obj["service_account_scopes"]),
		ShieldedInstanceConfig:  expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(obj["shielded_instance_config"]),
		Subnetwork:              dcl.String(obj["subnetwork"].(string)),
		Tags:                    tpgdclresource.ExpandStringArray(obj["tags"]),
		Zone:                    dcl.StringOrNil(obj["zone"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"internal_ip_only":           obj.InternalIPOnly,
		"metadata":                   obj.Metadata,
		"network":                    obj.Network,
		"node_group_affinity":        flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(obj.NodeGroupAffinity),
		"private_ipv6_google_access": obj.PrivateIPv6GoogleAccess,
		"reservation_affinity":       flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(obj.ReservationAffinity),
		"service_account":            obj.ServiceAccount,
		"service_account_scopes":     obj.ServiceAccountScopes,
		"shielded_instance_config":   flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(obj.ShieldedInstanceConfig),
		"subnetwork":                 obj.Subnetwork,
		"tags":                       obj.Tags,
		"zone":                       obj.Zone,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{
		NodeGroup: dcl.String(obj["node_group"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"node_group": obj.NodeGroup,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{
		ConsumeReservationType: dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef(obj["consume_reservation_type"].(string)),
		Key:                    dcl.String(obj["key"].(string)),
		Values:                 tpgdclresource.ExpandStringArray(obj["values"]),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"consume_reservation_type": obj.ConsumeReservationType,
		"key":                      obj.Key,
		"values":                   obj.Values,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{
		EnableIntegrityMonitoring: dcl.Bool(obj["enable_integrity_monitoring"].(bool)),
		EnableSecureBoot:          dcl.Bool(obj["enable_secure_boot"].(bool)),
		EnableVtpm:                dcl.Bool(obj["enable_vtpm"].(bool)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"enable_integrity_monitoring": obj.EnableIntegrityMonitoring,
		"enable_secure_boot":          obj.EnableSecureBoot,
		"enable_vtpm":                 obj.EnableVtpm,
	}

	return []interface{}{transformed}

}
func expandDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActionsArray(o interface{}) []dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	if o == nil {
		return make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions, 0)
	}

	items := make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions, 0, len(objs))
	for _, item := range objs {
		i := expandDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActions(item)
		items = append(items, *i)
	}

	return items
}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActions(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigInitializationActions
	}

	obj := o.(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions{
		ExecutableFile:   dcl.String(obj["executable_file"].(string)),
		ExecutionTimeout: dcl.String(obj["execution_timeout"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActionsArray(objs []dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActions(&item)
		items = append(items, i)
	}

	return items
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigInitializationActions(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigInitializationActions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"executable_file":   obj.ExecutableFile,
		"execution_timeout": obj.ExecutionTimeout,
	}

	return transformed

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{
		AutoDeleteTime: dcl.String(obj["auto_delete_time"].(string)),
		AutoDeleteTtl:  dcl.String(obj["auto_delete_ttl"].(string)),
		IdleDeleteTtl:  dcl.String(obj["idle_delete_ttl"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"auto_delete_time": obj.AutoDeleteTime,
		"auto_delete_ttl":  obj.AutoDeleteTtl,
		"idle_delete_ttl":  obj.IdleDeleteTtl,
		"idle_start_time":  obj.IdleStartTime,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfig{
		Accelerators:   expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsArray(obj["accelerators"]),
		DiskConfig:     expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(obj["disk_config"]),
		Image:          dcl.String(obj["image"].(string)),
		MachineType:    dcl.String(obj["machine_type"].(string)),
		MinCpuPlatform: dcl.StringOrNil(obj["min_cpu_platform"].(string)),
		NumInstances:   dcl.Int64(int64(obj["num_instances"].(int))),
		Preemptibility: dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnumRef(obj["preemptibility"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"accelerators":         flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsArray(obj.Accelerators),
		"disk_config":          flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(obj.DiskConfig),
		"image":                obj.Image,
		"machine_type":         obj.MachineType,
		"min_cpu_platform":     obj.MinCpuPlatform,
		"num_instances":        obj.NumInstances,
		"preemptibility":       obj.Preemptibility,
		"instance_names":       obj.InstanceNames,
		"is_preemptible":       obj.IsPreemptible,
		"managed_group_config": flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(obj.ManagedGroupConfig),
	}

	return []interface{}{transformed}

}
func expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsArray(o interface{}) []dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	if o == nil {
		return nil
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, 0, len(objs))
	for _, item := range objs {
		i := expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(item)
		items = append(items, *i)
	}

	return items
}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{
		AcceleratorCount: dcl.Int64(int64(obj["accelerator_count"].(int))),
		AcceleratorType:  dcl.String(obj["accelerator_type"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsArray(objs []dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(&item)
		items = append(items, i)
	}

	return items
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"accelerator_count": obj.AcceleratorCount,
		"accelerator_type":  obj.AcceleratorType,
	}

	return transformed

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{
		BootDiskSizeGb: dcl.Int64(int64(obj["boot_disk_size_gb"].(int))),
		BootDiskType:   dcl.String(obj["boot_disk_type"].(string)),
		NumLocalSsds:   dcl.Int64OrNil(int64(obj["num_local_ssds"].(int))),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"boot_disk_size_gb": obj.BootDiskSizeGb,
		"boot_disk_type":    obj.BootDiskType,
		"num_local_ssds":    obj.NumLocalSsds,
	}

	return []interface{}{transformed}

}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"instance_group_manager_name": obj.InstanceGroupManagerName,
		"instance_template_name":      obj.InstanceTemplateName,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{
		Accelerators:   expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsArray(obj["accelerators"]),
		DiskConfig:     expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(obj["disk_config"]),
		Image:          dcl.String(obj["image"].(string)),
		MachineType:    dcl.String(obj["machine_type"].(string)),
		MinCpuPlatform: dcl.StringOrNil(obj["min_cpu_platform"].(string)),
		NumInstances:   dcl.Int64(int64(obj["num_instances"].(int))),
		Preemptibility: dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnumRef(obj["preemptibility"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"accelerators":         flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsArray(obj.Accelerators),
		"disk_config":          flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(obj.DiskConfig),
		"image":                obj.Image,
		"machine_type":         obj.MachineType,
		"min_cpu_platform":     obj.MinCpuPlatform,
		"num_instances":        obj.NumInstances,
		"preemptibility":       obj.Preemptibility,
		"instance_names":       obj.InstanceNames,
		"is_preemptible":       obj.IsPreemptible,
		"managed_group_config": flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(obj.ManagedGroupConfig),
	}

	return []interface{}{transformed}

}
func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsArray(o interface{}) []dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	if o == nil {
		return nil
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, 0, len(objs))
	for _, item := range objs {
		i := expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(item)
		items = append(items, *i)
	}

	return items
}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{
		AcceleratorCount: dcl.Int64(int64(obj["accelerator_count"].(int))),
		AcceleratorType:  dcl.String(obj["accelerator_type"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsArray(objs []dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(&item)
		items = append(items, i)
	}

	return items
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"accelerator_count": obj.AcceleratorCount,
		"accelerator_type":  obj.AcceleratorType,
	}

	return transformed

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{
		BootDiskSizeGb: dcl.Int64(int64(obj["boot_disk_size_gb"].(int))),
		BootDiskType:   dcl.String(obj["boot_disk_type"].(string)),
		NumLocalSsds:   dcl.Int64OrNil(int64(obj["num_local_ssds"].(int))),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"boot_disk_size_gb": obj.BootDiskSizeGb,
		"boot_disk_type":    obj.BootDiskType,
		"num_local_ssds":    obj.NumLocalSsds,
	}

	return []interface{}{transformed}

}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"instance_group_manager_name": obj.InstanceGroupManagerName,
		"instance_template_name":      obj.InstanceTemplateName,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{
		KerberosConfig: expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(obj["kerberos_config"]),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"kerberos_config": flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(obj.KerberosConfig),
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{
		CrossRealmTrustAdminServer:    dcl.String(obj["cross_realm_trust_admin_server"].(string)),
		CrossRealmTrustKdc:            dcl.String(obj["cross_realm_trust_kdc"].(string)),
		CrossRealmTrustRealm:          dcl.String(obj["cross_realm_trust_realm"].(string)),
		CrossRealmTrustSharedPassword: dcl.String(obj["cross_realm_trust_shared_password"].(string)),
		EnableKerberos:                dcl.Bool(obj["enable_kerberos"].(bool)),
		KdcDbKey:                      dcl.String(obj["kdc_db_key"].(string)),
		KeyPassword:                   dcl.String(obj["key_password"].(string)),
		Keystore:                      dcl.String(obj["keystore"].(string)),
		KeystorePassword:              dcl.String(obj["keystore_password"].(string)),
		KmsKey:                        dcl.String(obj["kms_key"].(string)),
		Realm:                         dcl.String(obj["realm"].(string)),
		RootPrincipalPassword:         dcl.String(obj["root_principal_password"].(string)),
		TgtLifetimeHours:              dcl.Int64(int64(obj["tgt_lifetime_hours"].(int))),
		Truststore:                    dcl.String(obj["truststore"].(string)),
		TruststorePassword:            dcl.String(obj["truststore_password"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cross_realm_trust_admin_server":    obj.CrossRealmTrustAdminServer,
		"cross_realm_trust_kdc":             obj.CrossRealmTrustKdc,
		"cross_realm_trust_realm":           obj.CrossRealmTrustRealm,
		"cross_realm_trust_shared_password": obj.CrossRealmTrustSharedPassword,
		"enable_kerberos":                   obj.EnableKerberos,
		"kdc_db_key":                        obj.KdcDbKey,
		"key_password":                      obj.KeyPassword,
		"keystore":                          obj.Keystore,
		"keystore_password":                 obj.KeystorePassword,
		"kms_key":                           obj.KmsKey,
		"realm":                             obj.Realm,
		"root_principal_password":           obj.RootPrincipalPassword,
		"tgt_lifetime_hours":                obj.TgtLifetimeHours,
		"truststore":                        obj.Truststore,
		"truststore_password":               obj.TruststorePassword,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	if o == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{
		ImageVersion:       dcl.String(obj["image_version"].(string)),
		OptionalComponents: expandDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsArray(obj["optional_components"]),
		Properties:         tpgresource.CheckStringMap(obj["properties"]),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"image_version":       obj.ImageVersion,
		"optional_components": flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsArray(obj.OptionalComponents),
		"properties":          obj.Properties,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{
		Accelerators:   expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsArray(obj["accelerators"]),
		DiskConfig:     expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(obj["disk_config"]),
		Image:          dcl.String(obj["image"].(string)),
		MachineType:    dcl.String(obj["machine_type"].(string)),
		MinCpuPlatform: dcl.StringOrNil(obj["min_cpu_platform"].(string)),
		NumInstances:   dcl.Int64(int64(obj["num_instances"].(int))),
		Preemptibility: dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnumRef(obj["preemptibility"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"accelerators":         flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsArray(obj.Accelerators),
		"disk_config":          flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(obj.DiskConfig),
		"image":                obj.Image,
		"machine_type":         obj.MachineType,
		"min_cpu_platform":     obj.MinCpuPlatform,
		"num_instances":        obj.NumInstances,
		"preemptibility":       obj.Preemptibility,
		"instance_names":       obj.InstanceNames,
		"is_preemptible":       obj.IsPreemptible,
		"managed_group_config": flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(obj.ManagedGroupConfig),
	}

	return []interface{}{transformed}

}
func expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsArray(o interface{}) []dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	if o == nil {
		return nil
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, 0, len(objs))
	for _, item := range objs {
		i := expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(item)
		items = append(items, *i)
	}

	return items
}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{
		AcceleratorCount: dcl.Int64(int64(obj["accelerator_count"].(int))),
		AcceleratorType:  dcl.String(obj["accelerator_type"].(string)),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsArray(objs []dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(&item)
		items = append(items, i)
	}

	return items
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"accelerator_count": obj.AcceleratorCount,
		"accelerator_type":  obj.AcceleratorType,
	}

	return transformed

}

func expandDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(o interface{}) *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{
		BootDiskSizeGb: dcl.Int64(int64(obj["boot_disk_size_gb"].(int))),
		BootDiskType:   dcl.String(obj["boot_disk_type"].(string)),
		NumLocalSsds:   dcl.Int64OrNil(int64(obj["num_local_ssds"].(int))),
	}
}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"boot_disk_size_gb": obj.BootDiskSizeGb,
		"boot_disk_type":    obj.BootDiskType,
		"num_local_ssds":    obj.NumLocalSsds,
	}

	return []interface{}{transformed}

}

func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(obj *dataproc.WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"instance_group_manager_name": obj.InstanceGroupManagerName,
		"instance_template_name":      obj.InstanceTemplateName,
	}

	return []interface{}{transformed}

}
func expandDataprocWorkflowTemplateParametersArray(o interface{}) []dataproc.WorkflowTemplateParameters {
	if o == nil {
		return make([]dataproc.WorkflowTemplateParameters, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]dataproc.WorkflowTemplateParameters, 0)
	}

	items := make([]dataproc.WorkflowTemplateParameters, 0, len(objs))
	for _, item := range objs {
		i := expandDataprocWorkflowTemplateParameters(item)
		items = append(items, *i)
	}

	return items
}

func expandDataprocWorkflowTemplateParameters(o interface{}) *dataproc.WorkflowTemplateParameters {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateParameters
	}

	obj := o.(map[string]interface{})
	return &dataproc.WorkflowTemplateParameters{
		Fields:      tpgdclresource.ExpandStringArray(obj["fields"]),
		Name:        dcl.String(obj["name"].(string)),
		Description: dcl.String(obj["description"].(string)),
		Validation:  expandDataprocWorkflowTemplateParametersValidation(obj["validation"]),
	}
}

func flattenDataprocWorkflowTemplateParametersArray(objs []dataproc.WorkflowTemplateParameters) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenDataprocWorkflowTemplateParameters(&item)
		items = append(items, i)
	}

	return items
}

func flattenDataprocWorkflowTemplateParameters(obj *dataproc.WorkflowTemplateParameters) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"fields":      obj.Fields,
		"name":        obj.Name,
		"description": obj.Description,
		"validation":  flattenDataprocWorkflowTemplateParametersValidation(obj.Validation),
	}

	return transformed

}

func expandDataprocWorkflowTemplateParametersValidation(o interface{}) *dataproc.WorkflowTemplateParametersValidation {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateParametersValidation
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateParametersValidation
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateParametersValidation{
		Regex:  expandDataprocWorkflowTemplateParametersValidationRegex(obj["regex"]),
		Values: expandDataprocWorkflowTemplateParametersValidationValues(obj["values"]),
	}
}

func flattenDataprocWorkflowTemplateParametersValidation(obj *dataproc.WorkflowTemplateParametersValidation) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"regex":  flattenDataprocWorkflowTemplateParametersValidationRegex(obj.Regex),
		"values": flattenDataprocWorkflowTemplateParametersValidationValues(obj.Values),
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateParametersValidationRegex(o interface{}) *dataproc.WorkflowTemplateParametersValidationRegex {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateParametersValidationRegex
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateParametersValidationRegex
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateParametersValidationRegex{
		Regexes: tpgdclresource.ExpandStringArray(obj["regexes"]),
	}
}

func flattenDataprocWorkflowTemplateParametersValidationRegex(obj *dataproc.WorkflowTemplateParametersValidationRegex) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"regexes": obj.Regexes,
	}

	return []interface{}{transformed}

}

func expandDataprocWorkflowTemplateParametersValidationValues(o interface{}) *dataproc.WorkflowTemplateParametersValidationValues {
	if o == nil {
		return dataproc.EmptyWorkflowTemplateParametersValidationValues
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataproc.EmptyWorkflowTemplateParametersValidationValues
	}
	obj := objArr[0].(map[string]interface{})
	return &dataproc.WorkflowTemplateParametersValidationValues{
		Values: tpgdclresource.ExpandStringArray(obj["values"]),
	}
}

func flattenDataprocWorkflowTemplateParametersValidationValues(obj *dataproc.WorkflowTemplateParametersValidationValues) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"values": obj.Values,
	}

	return []interface{}{transformed}

}
func flattenDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsArray(obj []dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
func expandDataprocWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsArray(o interface{}) []dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum {
	objs := o.([]interface{})
	items := make([]dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum, 0, len(objs))
	for _, item := range objs {
		i := dataproc.WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumRef(item.(string))
		items = append(items, *i)
	}
	return items
}
