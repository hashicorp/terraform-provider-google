---
subcategory: "Dataflow"
description: |-
  Creates a job in Dataflow according to a provided config file.
---

# google_dataflow_job

Creates a job on Dataflow, which is an implementation of Apache Beam running on Google Compute Engine. For more information see
the official documentation for
[Beam](https://beam.apache.org) and [Dataflow](https://cloud.google.com/dataflow/).

## Example Usage

```hcl
resource "google_dataflow_job" "big_data_job" {
  name              = "dataflow-job"
  template_gcs_path = "gs://my-bucket/templates/template_file"
  temp_gcs_location = "gs://my-bucket/tmp_dir"
  parameters = {
    foo = "bar"
    baz = "qux"
  }
}
```
## Example Usage - Streaming Job
```hcl
resource "google_pubsub_topic" "topic" {
	name     = "dataflow-job1"
}
resource "google_storage_bucket" "bucket1" {
	name          = "tf-test-bucket1"
	location      = "US"
	force_destroy = true
}
resource "google_storage_bucket" "bucket2" {
	name          = "tf-test-bucket2"
	location      = "US"
	force_destroy = true
}
resource "google_dataflow_job" "pubsub_stream" {
	name = "tf-test-dataflow-job1"
	template_gcs_path = "gs://my-bucket/templates/template_file"
	temp_gcs_location = "gs://my-bucket/tmp_dir"
	enable_streaming_engine = true
	parameters = {
	  inputFilePattern = "${google_storage_bucket.bucket1.url}/*.json"
	  outputTopic    = google_pubsub_topic.topic.id
	}
	transform_name_mapping = {
		name = "test_job"
		env = "test"
	}
	on_delete = "cancel"
}
```

## Note on "destroy" / "apply"
There are many types of Dataflow jobs.  Some Dataflow jobs run constantly, getting new data from (e.g.) a GCS bucket, and outputting data continuously.  Some jobs process a set amount of data then terminate.  All jobs can fail while running due to programming errors or other issues.  In this way, Dataflow jobs are different from most other Terraform / Google resources.

The Dataflow resource is considered 'existing' while it is in a nonterminal state.  If it reaches a terminal state (e.g. 'FAILED', 'COMPLETE', 'CANCELLED'), it will be recreated on the next 'apply'.  This is as expected for jobs which run continuously, but may surprise users who use this resource for other kinds of Dataflow jobs.

A Dataflow job which is 'destroyed' may be "cancelled" or "drained".  If "cancelled", the job terminates - any data written remains where it is, but no new data will be processed.  If "drained", no new data will enter the pipeline, but any data currently in the pipeline will finish being processed.  The default is "drain". When `on_delete` is set to `"drain"` in the configuration, you may experience a long wait for your `terraform destroy` to complete.

You can potentially short-circuit the wait by setting `skip_wait_on_job_termination` to `true`, but beware that unless you take active steps to ensure that the job `name` parameter changes between instances, the name will conflict and the launch of the new job will fail. One way to do this is with a [random_id](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/id) resource, for example:

```hcl
variable "big_data_job_subscription_id" {
  type    = string
  default = "projects/myproject/subscriptions/messages"
}

resource "random_id" "big_data_job_name_suffix" {
  byte_length = 4
  keepers = {
    region          = var.region
    subscription_id = var.big_data_job_subscription_id
  }
}
resource "google_dataflow_flex_template_job" "big_data_job" {
  provider                      = google-beta
  name                          = "dataflow-flextemplates-job-${random_id.big_data_job_name_suffix.dec}"
  region                        = var.region
  container_spec_gcs_path       = "gs://my-bucket/templates/template.json"
  skip_wait_on_job_termination = true
  parameters = {
    inputSubscription = var.big_data_job_subscription_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by Dataflow.
* `template_gcs_path` - (Required) The GCS path to the Dataflow job template.
* `temp_gcs_location` - (Required) A writeable location on GCS for the Dataflow job to dump its temporary data.

- - -

* `parameters` - **Template specific** Key/Value pairs to be forwarded to the pipeline's options; keys are
  case-sensitive based on the language on which the pipeline is coded, mostly Java.
  **Note**: do not configure Dataflow options here in parameters.
* `labels` - (Optional) User labels to be specified for the job. Keys and values should follow the restrictions
   specified in the [labeling restrictions](https://cloud.google.com/compute/docs/labeling-resources#restrictions) page.
   **Note**: This field is non-authoritative, and will only manage the labels present in your configuration. Please refer to the field `effective_labels` for all of the labels present on the resource.
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
* `transform_name_mapping` - (Optional) Only applicable when updating a pipeline. Map of transform name prefixes of the job to be replaced with the corresponding name prefixes of the new job. This field is not used outside of update.
* `max_workers` - (Optional) The number of workers permitted to work on the job.  More workers may improve processing speed at additional cost.
* `on_delete` - (Optional) One of "drain" or "cancel".  Specifies behavior of deletion during `terraform destroy`.  See above note.
* `skip_wait_on_job_termination` - (Optional)  If set to `true`, terraform will treat `DRAINING` and `CANCELLING` as terminal states when deleting the resource, and will remove the resource from terraform state and move on.  See above note.
* `project` - (Optional) The project in which the resource belongs. If it is not provided, the provider project is used.
* `zone` - (Optional) The zone in which the created job should run. If it is not provided, the provider zone is used.
* `region` - (Optional) The region in which the created job should run.
* `service_account_email` - (Optional) The Service Account email used to create the job. This should be just an email e.g. `myserviceaccount@myproject.iam.gserviceaccount.com`. Do not include any `serviceAccount:` or other prefix.
* `network` - (Optional) The network to which VMs will be assigned. If it is not provided, "default" will be used.
* `subnetwork` - (Optional) The subnetwork to which VMs will be assigned. Should be of the form "regions/REGION/subnetworks/SUBNETWORK". If the [subnetwork is located in a Shared VPC network](https://cloud.google.com/dataflow/docs/guides/specifying-networks#shared), you must use the complete URL. For example `"googleapis.com/compute/v1/projects/PROJECT_ID/regions/REGION/subnetworks/SUBNET_NAME"`
* `machine_type` - (Optional) The machine type to use for the job.
* `kms_key_name` - (Optional) The name for the Cloud KMS key for the job. Key format is: `projects/PROJECT_ID/locations/LOCATION/keyRings/KEY_RING/cryptoKeys/KEY`
* `ip_configuration` - (Optional) The configuration for VM IPs.  Options are `"WORKER_IP_PUBLIC"` or `"WORKER_IP_PRIVATE"`.
* `additional_experiments` - (Optional) List of experiments that should be used by the job. An example value is `["enable_stackdriver_agent_metrics"]`.
* `enable_streaming_engine` - (Optional) Enable/disable the use of [Streaming Engine](https://cloud.google.com/dataflow/docs/guides/deploying-a-pipeline#streaming-engine) for the job. Note that Streaming Engine is enabled by default for pipelines developed against the Beam SDK for Python v2.21.0 or later when using Python 3.

## Attributes Reference

* `job_id` - The unique ID of this job.
* `type` - The type of this job, selected from the [JobType enum](https://cloud.google.com/dataflow/docs/reference/rest/v1b3/projects.jobs#Job.JobType)
* `state` - The current state of the resource, selected from the [JobState enum](https://cloud.google.com/dataflow/docs/reference/rest/v1b3/projects.jobs#Job.JobState)

## Import

Dataflow jobs can be imported using the job `id` e.g.

* `{{id}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import dataflow jobs using one of the formats above. For example:

```tf
import {
  id = "{{id}}"
  to = google_dataflow_job.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), dataflow jobs can be imported using one of the formats above. For example:

```
$ terraform import google_dataflow_job.default {{id}}
```
